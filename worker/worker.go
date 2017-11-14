package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	. "fogflow/common/communicator"
	. "fogflow/common/datamodel"
	. "fogflow/common/ngsi"
)

type Worker struct {
	id           string
	communicator *Communicator
	ticker       *time.Ticker
	executor     *Executor
	allTasks     map[string]*ScheduledTaskInstance
	cfg          *Config
	profile      WorkerProfile
}

func (w *Worker) Start(config *Config) bool {
	w.cfg = config

	// construct a unique ID for this edge node based on its logical location
	w.id = "SysComponent.Worker." + strconv.Itoa(w.cfg.LLocation.SiteNo) + strconv.Itoa(w.cfg.LLocation.NodeNo)

	w.profile.WID = w.id
	w.profile.Capacity = 10
	w.profile.PLocation = config.PLocation
	w.profile.LLocation = config.LLocation

	w.allTasks = make(map[string]*ScheduledTaskInstance)

	cfg := MessageBusConfig{}
	cfg.Broker = w.cfg.MessageBus
	cfg.Exchange = "fogflow"
	cfg.ExchangeType = "topic"
	cfg.DefaultQueue = w.id
	cfg.BindingKeys = []string{w.id + ".*"}

	// find a nearby IoT Broker
	for {
		nearby := NearBy{}
		nearby.Latitude = w.cfg.PLocation.Latitude
		nearby.Longitude = w.cfg.PLocation.Longitude
		nearby.Limit = 1

		client := NGSI9Client{IoTDiscoveryURL: w.cfg.DiscoveryURL}
		selectedBroker, err := client.DiscoveryNearbyIoTBroker(nearby)

		if err == nil && selectedBroker != "" {
			w.cfg.BrokerURL = selectedBroker
			break
		} else {
			if err != nil {
				ERROR.Println(err)
			}

			INFO.Println("continue to look up a nearby IoT broker")
			time.Sleep(5 * time.Second)
		}
	}

	w.publishMyself()

	// start the executor to interact with docker
	w.executor = &Executor{}
	w.executor.Init(w.cfg)

	// create the communicator with the broker info and topics
	w.communicator = NewCommunicator(&cfg)

	// start the message consumer
	go func() {
		for {
			retry, err := w.communicator.StartConsuming(w.id, w)
			if retry {
				fmt.Printf("Going to retry launching the edge node. Error: %v", err)
			} else {
				break
			}
		}
	}()

	// start a timer to do something periodically
	w.ticker = time.NewTicker(time.Second * 5)
	go func() {
		for {
			<-w.ticker.C
			w.onTimer()
		}
	}()

	return true
}

func (w *Worker) Quit() {
	w.unpublishMyself()
	w.communicator.StopConsuming()
	w.ticker.Stop()
	w.executor.Shutdown()
	fmt.Println("stop consuming the messages")
}

func (w *Worker) publishMyself() {
	ctxObj := ContextObject{}

	ctxObj.Entity.ID = w.id
	ctxObj.Entity.Type = "Worker"
	ctxObj.Entity.IsPattern = false

	ctxObj.Attributes = make(map[string]ValueObject)
	ctxObj.Attributes["id"] = ValueObject{Type: "string", Value: w.id}
	ctxObj.Attributes["capacity"] = ValueObject{Type: "integer", Value: 2}
	ctxObj.Attributes["physical_location"] = ValueObject{Type: "object", Value: w.cfg.PLocation}
	ctxObj.Attributes["logical_location"] = ValueObject{Type: "object", Value: w.cfg.LLocation}

	ctxObj.Metadata = make(map[string]ValueObject)
	mylocation := Point{}
	mylocation.Latitude = w.cfg.PLocation.Latitude
	mylocation.Longitude = w.cfg.PLocation.Longitude
	ctxObj.Metadata["location"] = ValueObject{Type: "point", Value: mylocation}
	ctxObj.Metadata["role"] = ValueObject{Type: "string", Value: w.cfg.MyRole}

	client := NGSI10Client{IoTBrokerURL: w.cfg.BrokerURL}
	err := client.UpdateContext(&ctxObj)
	if err != nil {
		fmt.Println(err)
	}
}

func (w *Worker) unpublishMyself() {
	entity := EntityId{}
	entity.ID = w.id
	entity.Type = "Worker"
	entity.IsPattern = false

	client := NGSI10Client{IoTBrokerURL: w.cfg.BrokerURL}
	err := client.DeleteContext(&entity)
	if err != nil {
		fmt.Println(err)
	}
}

func (w *Worker) Process(msg *RecvMessage) error {
	var err error

	switch msg.Type {
	case "ADD_TASK":
		task := ScheduledTaskInstance{}
		err = json.Unmarshal(msg.PayLoad, &task)
		if err == nil {
			w.onScheduledTask(msg.From, &task)
		}
	case "REMOVE_TASK":
		task := ScheduledTaskInstance{}
		err = json.Unmarshal(msg.PayLoad, &task)
		if err == nil {
			w.onTerminateTask(msg.From, &task)
		}

	case "ADD_INPUT":
		flow := FlowInfo{}
		err = json.Unmarshal(msg.PayLoad, &flow)
		if err == nil {
			w.onAddInput(msg.From, &flow)
		}
	case "REMOVE_INPUT":
		flow := FlowInfo{}
		err = json.Unmarshal(msg.PayLoad, &flow)
		if err == nil {
			w.onRemoveInput(msg.From, &flow)
		}

	case "prefetch_image":
		imageList := make([]DockerImage, 0)
		err = json.Unmarshal(msg.PayLoad, &imageList)
		if err == nil {
			w.onPrefetchImage(imageList)
		}
	}

	return err
}

func (w *Worker) onTimer() {
	w.heartbeat()
}

func (w *Worker) heartbeat() {
	taskUpdateMsg := SendMessage{Type: "heart_beat", RoutingKey: "heartbeat.", From: w.id, PayLoad: w.profile}
	w.communicator.Publish(&taskUpdateMsg)
}

func (w *Worker) onAddInput(from string, flow *FlowInfo) {
	w.executor.onAddInput(flow)
}

func (w *Worker) onRemoveInput(from string, flow *FlowInfo) {
	w.executor.onRemoveInput(flow)
}

func (w *Worker) onScheduledTask(from string, task *ScheduledTaskInstance) {
	INFO.Println("execute task ", task.ID, " with operation", task.DockerImage)
	INFO.Printf("task configuration %+v\n", (*task))

	Runnable := true
	for _, existTask := range w.allTasks {
		// judge if the incoming new task can not run against the existing task
		if task.IsExclusive == false && existTask.IsExclusive == true {
			Runnable = false
		}
		if task.IsExclusive == true && existTask.IsExclusive == true && existTask.PriorityLevel > task.PriorityLevel {
			Runnable = false
		}

		// judge if the existing task should be overtaken by the incoming new task
		overTaken := false
		if task.IsExclusive == true && existTask.IsExclusive == false {
			overTaken = true
		}
		if task.IsExclusive == true && existTask.IsExclusive == true && existTask.PriorityLevel < task.PriorityLevel {
			overTaken = true
		}

		if overTaken == true {
			// pause this task temporarily
			go w.executor.TerminateTask(existTask.ID, true)
			existTask.Status = "paused"

			tp := TaskUpdate{}
			tp.TaskID = existTask.ID
			tp.Status = "paused"
			taskUpdateMsg := SendMessage{Type: "task_update", RoutingKey: "master." + from + ".", From: w.id, PayLoad: tp}
			w.communicator.Publish(&taskUpdateMsg)
		}
	}

	INFO.Printf("task runnable =  %+v\n", Runnable)

	// if the incoming new task is able to run
	if Runnable == true {
		go w.executor.LaunchTask(task)

		// add the new task into the local task list
		task.Status = "running"
		w.allTasks[task.ID] = task

		// send ACK back to the master
		tp := TaskUpdate{}
		tp.TaskID = task.ID
		tp.Status = "running"
		taskUpdateMsg := SendMessage{Type: "task_update", RoutingKey: "master." + from + ".", From: w.id, PayLoad: tp}
		w.communicator.Publish(&taskUpdateMsg)
	} else {
		// add the new task into the local task list
		task.Status = "pause"
		w.allTasks[task.ID] = task

		// send ACK back to the master
		tp := TaskUpdate{}
		tp.TaskID = task.ID
		tp.Status = "pause"
		taskUpdateMsg := SendMessage{Type: "task_update", RoutingKey: "master." + from + ".", From: w.id, PayLoad: tp}
		w.communicator.Publish(&taskUpdateMsg)
	}
}

func (w *Worker) onTerminateTask(from string, task *ScheduledTaskInstance) {
	INFO.Println("terminate task ", task.ID, " with operation", task.DockerImage)

	myTask := w.allTasks[task.ID]

	if myTask == nil {
		return
	}

	if myTask.Status == "running" {
		go w.executor.TerminateTask(task.ID, false)
	}

	// remove it from the local task list
	delete(w.allTasks, myTask.ID)

	if myTask.IsExclusive == false {
		return
	}

	// check the other paused tasks
	hasActiveUrgentTask := false
	for _, t := range w.allTasks {
		if t.IsExclusive == true && t.Status == "running" {
			hasActiveUrgentTask = true
			break
		}
	}

	if hasActiveUrgentTask == true {
		return
	}

	/// no more running exclusive tasks, then check if there is any other paused exclusive task
	hasExclusiveTask := false
	priority := 0
	for _, task := range w.allTasks {
		if task.IsExclusive == true && task.Status != "running" && task.PriorityLevel > priority {
			hasExclusiveTask = true
			priority = task.PriorityLevel
		}
	}

	if hasExclusiveTask == true {
		for _, task := range w.allTasks {
			if task.PriorityLevel == priority {
				// restart this task temporarily
				go w.executor.LaunchTask(task)
				task.Status = "running"

				tp := TaskUpdate{}
				tp.TaskID = task.ID
				tp.Status = "running"
				taskUpdateMsg := SendMessage{Type: "task_update", RoutingKey: "master." + from + ".", From: w.id, PayLoad: tp}
				w.communicator.Publish(&taskUpdateMsg)
			}
		}
	} else {
		// resume other tasks
		for _, task := range w.allTasks {
			if task.Status != "running" {
				// restart this task temporarily
				go w.executor.LaunchTask(task)
				task.Status = "running"

				tp := TaskUpdate{}
				tp.TaskID = task.ID
				tp.Status = "running"
				taskUpdateMsg := SendMessage{Type: "task_update", RoutingKey: "master." + from + ".", From: w.id, PayLoad: tp}
				w.communicator.Publish(&taskUpdateMsg)
			}
		}
	}
}

func (w *Worker) onPrefetchImage(imageList []DockerImage) {
	for _, dockerImage := range imageList {
		INFO.Println("I am going to fetch the docker image", dockerImage.ImageName)
		imageURL := w.cfg.Registry.ServerAddress + "/" + dockerImage.ImageName
		w.executor.PullImage(imageURL, dockerImage.ImageTag)
	}
}
