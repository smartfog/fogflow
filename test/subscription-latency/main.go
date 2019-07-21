package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	. "github.com/smartfog/nec-fogflow/common/ngsi"
)

var startTime = time.Now()
var started = false
var counter = int64(0)

func main() {
	configurationFile := flag.String("f", "config.json", "A configuration file")
	myPort := flag.Int("p", 8050, "the port of this agent")
	num := flag.Int("n", 2, "number of updates")

	flag.Parse()

	config := CreateConfig(*configurationFile)

	config.MyPort = *myPort

	startAgent(&config)
	sid := subscribe(&config)

	time.Sleep(10 * time.Second)

	startTime = time.Now()
	started = true
	counter = 0
	for i := 1; i < *num; i++ {
		update(&config, i)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	<-c

	unsubscribe(&config, sid)

	time.Sleep(5 * time.Second)
}

func startAgent(config *Config) {
	agent := NGSIAgent{Port: config.MyPort}
	agent.Start()
	agent.SetContextNotifyHandler(HandleNotifyContext)
}

func HandleNotifyContext(notifyCtxReq *NotifyContextRequest) {
	//INFO.Println("===========RECEIVE NOTIFY CONTEXT=========")
	//INFO.Printf("<< %+v >>\r\n", notifyCtxReq)

	if started {
		counter = counter + 1

		now := time.Now()
		delta := now.Sub(startTime)
		throughput := int64(float64(counter) / delta.Seconds())

		fmt.Printf("throughput %d, %d \r\n", counter, throughput)
	}

	/*
		for _, v := range notifyCtxReq.ContextResponses {
			ctxObj := CtxElement2Object(&(v.ContextElement))
			currentTime := int64(time.Now().UnixNano() / 1000000)
			latency := currentTime - ctxObj.Attributes["time"].Value.(int64)
			num := ctxObj.Attributes["no"].Value.(int64)
			fmt.Printf("No. %d, latency: %d \r\n", num, latency)
		} */

}

func subscribe(config *Config) string {
	subscription := SubscribeContextRequest{}

	newEntity := EntityId{}
	newEntity.Type = "Car"
	newEntity.IsPattern = true
	subscription.Entities = make([]EntityId, 0)
	subscription.Entities = append(subscription.Entities, newEntity)

	subscription.Reference = "http://" + config.MyIP + ":" + strconv.Itoa(config.MyPort)

	client := NGSI10Client{IoTBrokerURL: config.SubscribeBrokerURL}
	sid, err := client.SubscribeContext(&subscription, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(sid)

	return sid
}

func unsubscribe(config *Config, sid string) {
	// unsubscribe context
	fmt.Println("=============unsubscribe===================")
	client := NGSI10Client{IoTBrokerURL: config.SubscribeBrokerURL}
	err := client.UnsubscribeContext(sid)
	if err != nil {
		fmt.Println(err)
	}
}

func update(config *Config, i int) {
	ctxObj := ContextObject{}

	ctxObj.Entity.ID = "Car." + strconv.Itoa(i)
	ctxObj.Entity.Type = "Car"
	ctxObj.Entity.IsPattern = false

	ctxObj.Attributes = make(map[string]ValueObject)
	ctxObj.Attributes["no"] = ValueObject{Type: "integer", Value: i}

	currentTime := int64(time.Now().UnixNano() / 1000000)
	ctxObj.Attributes["time"] = ValueObject{Type: "integer", Value: currentTime}

	ctxObj.Metadata = make(map[string]ValueObject)
	ctxObj.Metadata["no"] = ValueObject{Type: "integer", Value: i}
	ctxObj.Metadata["time"] = ValueObject{Type: "integer", Value: currentTime}
	ctxObj.Metadata["level"] = ValueObject{Type: "integer", Value: i}
	ctxObj.Metadata["level3"] = ValueObject{Type: "integer", Value: i}
	ctxObj.Metadata["level4"] = ValueObject{Type: "integer", Value: i}

	client := NGSI10Client{IoTBrokerURL: config.UpdateBrokerURL}
	err := client.UpdateContextObject(&ctxObj)
	if err != nil {
		fmt.Println(err)
	}

	//INFO.Println("send update ", i)
}
