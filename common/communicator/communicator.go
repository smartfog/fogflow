package communicator

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"

	. "fogflow/common/datamodel"
)

type TaskProcessor interface {
	Process(msg *RecvMessage) error
}

// Config holds all configuration for our program
type MessageBusConfig struct {
	Broker       string
	Exchange     string
	ExchangeType string
	DefaultQueue string
	BindingKeys  []string
}

// Communicator represents an AMQP broker
type Communicator struct {
	config *MessageBusConfig

	retry     bool
	retryFunc func()

	conn *amqp.Connection

	consuming bool
	stopChan  chan int
}

var RetryClosure = func() {
	retryIn := 2 // retry after 2 seconds
	log.Printf("Retrying to connect RabbitMQ in %v seconds", retryIn)
	time.Sleep(time.Second)
}

func NewCommunicator(cnf *MessageBusConfig) *Communicator {
	return &Communicator{config: cnf, retry: true}
}

// StartConsuming enters a loop and waits for incoming messages
func (communicator *Communicator) StartConsuming(consumerTag string, taskProcessor TaskProcessor) (bool, error) {
	if communicator.retryFunc == nil {
		communicator.retryFunc = RetryClosure
	}

	channel, queue, err := communicator.openSubscriber()

	if channel != nil {
		defer channel.Close()
	}

	if err != nil {
		communicator.retryFunc()
		communicator.consuming = false
		return communicator.retry, err
	}

	communicator.stopChan = make(chan int)

	if err := channel.Qos(
		3,     // prefetch count
		0,     // prefetch size
		false, // global
	); err != nil {
		return communicator.retry, fmt.Errorf("Channel Qos: %s", err)
	}

	deliveries, err := channel.Consume(
		queue.Name,  // queue
		consumerTag, // consumer tag
		false,       // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		return communicator.retry, fmt.Errorf("Queue Consume: %s", err)
	}

	log.Print("[*] Waiting for messages. To exit press CTRL+C")
	communicator.consuming = true

	if err := communicator.consume(deliveries, taskProcessor); err != nil {
		return communicator.retry, err // retry true
	}

	return communicator.retry, nil
}

// StopConsuming quits the loop
func (communicator *Communicator) StopConsuming() {
	// Do not retry from now on
	communicator.retry = false

	// Notifying the stop channel stops consuming of messages
	if communicator.consuming == true {
		communicator.stopChan <- 1
	}
}

// Publish places a new message on the default queue
func (communicator *Communicator) Publish(msg *SendMessage) error {
	channel, confirmsChan, err := communicator.openPublisher()
	if err != nil {
		return err
	}
	defer channel.Close()

	message, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("JSON Encode Message: %v", err)
	}

	if err := channel.Publish(
		communicator.config.Exchange, // exchange
		msg.RoutingKey,               // routing key
		false,                        // mandatory
		false,                        // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         message,
			DeliveryMode: amqp.Transient,
			//DeliveryMode: amqp.Persistent,
		},
	); err != nil {
		return err
	}

	confirmed := <-confirmsChan

	if confirmed.Ack {
		return nil
	}

	return fmt.Errorf("Failed delivery of delivery tag: %v", confirmed.DeliveryTag)
}

// Consume a single message
func (communicator *Communicator) consumeOne(d amqp.Delivery, taskProcessor TaskProcessor, errorsChan chan error) {
	if len(d.Body) == 0 {
		d.Nack(false, false)                                   // multiple, requeue
		errorsChan <- errors.New("Received an empty message.") // RabbitMQ down?
		return
	}

	//log.Printf("Received new message: %s", d.Body)

	msg := RecvMessage{}
	if err := json.Unmarshal(d.Body, &msg); err != nil {
		d.Nack(false, false) // multiple, requeue
		errorsChan <- err
		return
	}

	if err := taskProcessor.Process(&msg); err != nil {
		errorsChan <- err
	}

	d.Ack(false) // multiple
}

// Consumes messages...
func (communicator *Communicator) consume(deliveries <-chan amqp.Delivery, taskProcessor TaskProcessor) error {
	errorsChan := make(chan error)
	for {
		select {
		case <-communicator.stopChan:
			return nil
		case err := <-errorsChan:
			return err
		case d := <-deliveries:
			// Consume the task inside a gotourine so multiple tasks
			// can be processed concurrently
			go func() {
				communicator.consumeOne(d, taskProcessor, errorsChan)
			}()
		}
	}
}

// Connects to the message queue, opens a channel, declares a queue
func (communicator *Communicator) openSubscriber() (*amqp.Channel, amqp.Queue, error) {
	var (
		conn    *amqp.Connection
		channel *amqp.Channel
		queue   amqp.Queue
		err     error
	)

	// Connect, reuse the same connection if there is an established connection with the broker
	if communicator.conn != nil {
		conn = communicator.conn
	} else {
		conn, err = amqp.Dial(communicator.config.Broker)
		if err != nil {
			return channel, queue, fmt.Errorf("Dial: %s\r\n", err)
		}
		communicator.conn = conn
	}

	// Open a channel
	channel, err = conn.Channel()
	if err != nil {
		return channel, queue, fmt.Errorf("Channel: %s\r\n", err)
	}

	// Declare an exchange
	if err := channel.ExchangeDeclare(
		communicator.config.Exchange,     // name of the exchange
		communicator.config.ExchangeType, // type
		true,                             // durable
		true,                             // delete when complete
		false,                            // internal
		false,                            // noWait
		nil,                              // arguments
	); err != nil {
		return channel, queue, fmt.Errorf("Exchange Declare: %s\r\n", err)
	}

	// Declare a queue
	queue, err = channel.QueueDeclare(
		communicator.config.DefaultQueue, // name
		true,                             // durable
		true,                             // delete when unused
		false,                            // exclusive
		false,                            // no-wait
		nil,                              // arguments
	)
	if err != nil {
		return channel, queue, fmt.Errorf("Queue Declare: %s\r\n", err)
	}

	// Bind topics with the queue
	for _, key := range communicator.config.BindingKeys {
		if err := channel.QueueBind(
			queue.Name,                   // name of the queue
			key,                          // binding topic
			communicator.config.Exchange, // source exchange
			false,                        // noWait
			nil,                          // arguments
		); err != nil {
			return channel, queue, fmt.Errorf("Queue Bind: %s\r\n", err)
		}
	}

	return channel, queue, nil
}

// Connects to the message queue, opens a channel, declares a queue
func (communicator *Communicator) openPublisher() (*amqp.Channel, <-chan amqp.Confirmation, error) {
	var (
		conn    *amqp.Connection
		channel *amqp.Channel
		err     error
	)

	// Connect
	if communicator.conn != nil {
		conn = communicator.conn
	} else {
		conn, err = amqp.Dial(communicator.config.Broker)
		if err != nil {
			return channel, nil, fmt.Errorf("Dial: %s", err)
		}
		communicator.conn = conn
	}

	// Open a channel
	channel, err = conn.Channel()
	if err != nil {
		return channel, nil, fmt.Errorf("Channel: %s", err)
	}

	// Declare an exchange
	if err := channel.ExchangeDeclare(
		communicator.config.Exchange,     // name of the exchange
		communicator.config.ExchangeType, // type
		true,                             // durable
		true,                             // delete when complete
		false,                            // internal
		false,                            // noWait
		nil,                              // arguments
	); err != nil {
		return channel, nil, fmt.Errorf("Exchange Declare: %s", err)
	}

	// Enable publish confirmations
	if err := channel.Confirm(false); err != nil {
		return channel, nil, fmt.Errorf("Channel could not be put into confirm mode: %s", err)
	}

	return channel, channel.NotifyPublish(make(chan amqp.Confirmation, 1)), nil
}
