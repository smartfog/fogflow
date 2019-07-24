package main

import (
	"encoding/json"
	"runtime"
	"strconv"
	"sync"
	"time"

	. "github.com/smartfog/fogflow/common/communicator"
	. "github.com/smartfog/fogflow/common/config"
	. "github.com/smartfog/fogflow/common/datamodel"
	. "github.com/smartfog/fogflow/common/ngsi"
)

type Worker struct {
	id           string
	communicator *Communicator
	ticker       *time.Ticker
	executor     *Executor

	allTasks      map[string]*ScheduledTaskInstance
	taskList_lock sync.RWMutex

	cfg               *Config
	selectedBrokerURL string
	profile           WorkerProfile
}

func (w *Worker) Start(config *Config) bool {
	w.cfg = config

	w.profile.WID = w.id
	w.profile.Capacity = 10
	w.profile.PLocation = config.PLocation

	w.profile.OSType = runtime.GOOS
	w.profile.HWType = runtime.GOARCH

	w.allTasks = make(map[string]*ScheduledTaskInstance)

	cfg := MessageBusConfig{}
	cfg.Broker = w.cfg.GetMessageBus()
	cfg.Exchange = "fogflow"
	cfg.ExchangeType = "topic"
	cfg.DefaultQueue = w.id
	cfg.BindingKeys = []string{w.id + ".*"}

	// if no broker is configured in the configuration file, the worker needs to find a nearby IoT Broker
	// otherwise, just use the configured broker
	if config.Broker.Port != 0 {
		w.selectedBrokerURL = "http://" + config.InternalIP + ":" + strconv.Itoa(config.Broker.Port) + "/ngsi10"
	} else {
		// find a nearby IoT Broker
		for {
			nearby := NearBy{}
			nearby.Latitude = w.cfg.PLocation.Latitude
			nearby.Longitude = w.cfg.PLocation.Longitude
			nearby.Limit = 1

			client := NGSI9Client{IoTDiscoveryURL: w.cfg.GetDiscoveryURL()}
			selectedBroker, err := client.DiscoveryNearbyIoTBroker(nearby)
			if err == nil && selectedBroker != "" {
				w.selectedBrokerURL = selectedBroker
				INFO.Println("find out a nearby broker ", selectedBroker)
				break
			} else {
				if err != nil {
					ERROR.Println(err)
				}

				INFO.Println("continue to look up a nearby IoT broker")
				time.Sleep(5 * time.Second)
			}
		}
	}

	INFO.Println("communicating with the broker ", w.selectedBrokerURL)

	for {
		err := w.publishMyself()
		if err != nil {
			INFO.Println("wait for the assigned broker to be ready")
			time.Sleep(5 * time.Second)
		} else {
			INFO.Println("annouce myself to the nearby broker")
			break
		}
	}

	// start the executor to interact with docker
	w.executor = &Executor{}
	w.executor.Init(w.cfg, w.selectedBrokerURL)

	// create the communicator with the broker info and topics
	w.communicator = NewCommunicator(&cfg)

	// start the message consumer
	go func() {
		for {
			retry, err := w.communicator.StartConsuming(w.id, w)
			if retry {
				INFO.Printf("Going to retry launching the edge node. Error: %v", err)
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
	INFO.Println("unregister myself")
	w.unpublishMyself()

	INFO.Println("stop the timer")
	w.ticker.Stop()

	INFO.Println("stop consuming the messages")
	w.communicator.StopConsuming()

	INFO.Println("to stop the worker")
	w.executor.Shutdown()
}

func (w *Worker) publishMyself() error {
	ctxObj := ContextObject{}

	ctxObj.Entity.ID = w.id
	ctxObj.Entity.Type = "Worker"
	ctxObj.Entity.IsPattern = false

	ctxObj.Attributes = make(map[string]ValueObject)
	ctxObj.Attributes["id"] = ValueObject{Type: "string", Value: w.id}
	ctxObj.Attributes["capacity"] = ValueObject{Type: "integer", Value: w.profile.Capacity}
	ctxObj.Attributes["location"] = ValueObject{Type: "object", Value: w.cfg.PLocation}

	ctxObj.Metadata = make(map[string]ValueObject)
	mylocation := Point{}
	mylocation.Latitude = w.cfg.PLocation.Latitude
	mylocation.Longitude = w.cfg.PLocation.Longitude
	ctxObj.Metadata["location"] = ValueObject{Type: "point", Value: mylocation}

	client := NGSI10Client{IoTBrokerURL: w.selectedBrokerURL}
	err := client.UpdateContextObject(&ctxObj)
	return err
}

func (w *Worker) unpublishMyself() {
	entity := EntityId{}
	entity.ID = w.id
	entity.Type = "Worker"
	entity.IsPattern = false

	client := NGSI10Client{IoTBrokerURL: w.selectedBrokerURL}
	err := client.DeleteContext(&entity)
	if err != nil {
		ERROR.Println(err)
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

	case "PREFETCH_IMAGE":
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

func (w *Worker) TaskUpdate(masterID string, task *ScheduledTaskInstance, state string) {
	tp := TaskUpdate{}
	tp.ServiceName = task.ServiceName
	tp.TaskName = task.TaskName
	tp.TaskID = task.ID
	tp.Status = state
	taskUpdateMsg := SendMessage{Type: "task_update", RoutingKey: "master." + masterID + ".", From: w.id, PayLoad: tp}

	go w.communicator.Publish(&taskUpdateMsg)
}

func (w *Worker) onScheduledTask(from string, task *ScheduledTaskInstance) {
	INFO.Println("execute task ", task.ID, " with operation", task.DockerImage)
	INFO.Printf("task configuration %+v\n", (*task))

	w.taskList_lock.Lock()
	defer w.taskList_lock.Unlock()

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

			w.TaskUpdate(from, existTask, "paused")
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
		w.TaskUpdate(from, task, "running")
	} else {
		// add the new task into the local task list
		task.Status = "paused"
		w.allTasks[task.ID] = task

		// send ACK back to the master
		w.TaskUpdate(from, task, "paused")
	}
}

func (w *Worker) onTerminateTask(from string, task *ScheduledTaskInstance) {
	w.taskList_lock.Lock()
	defer w.taskList_lock.Unlock()

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

				w.TaskUpdate(from, task, "running")
			}
		}
	} else {
		// resume other tasks
		for _, task := range w.allTasks {
			if task.Status != "running" {
				// restart this task temporarily
				go w.executor.LaunchTask(task)
				task.Status = "running"

				w.TaskUpdate(from, task, "running")
			}
		}
	}
}

func (w *Worker) onPrefetchImage(imageList []DockerImage) {
	for _, dockerImage := range imageList {
		INFO.Println("I am going to fetch the docker image", dockerImage.ImageName)
		go w.executor.PullImage(dockerImage.ImageName, dockerImage.ImageTag)
	}
}
