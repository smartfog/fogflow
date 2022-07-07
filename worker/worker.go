package main

import (
	"encoding/json"
	"runtime"
	"sync"
	"time"

	. "fogflow/common/communicator"
	. "fogflow/common/config"
	. "fogflow/common/datamodel"
	. "fogflow/common/ngsi"
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

	profile WorkerProfile
}

func (w *Worker) Start(config *Config) bool {
	w.cfg = config

	w.profile.WID = w.id
	w.profile.Capacity = config.Worker.Capacity
	w.profile.PLocation = config.Location
	w.profile.Workload = 0
	w.profile.OSType = runtime.GOOS
	w.profile.HWType = runtime.GOARCH

	w.allTasks = make(map[string]*ScheduledTaskInstance)

	cfg := MessageBusConfig{}
	cfg.Broker = w.cfg.GetMessageBus()
	cfg.Exchange = "fogflow"
	cfg.ExchangeType = "topic"
	cfg.DefaultQueue = w.id
	cfg.BindingKeys = []string{w.id + ".*"}

	w.selectedBrokerURL = config.GetBrokerURL()

	// start the executor to interact with docker
	w.executor = &Executor{}
	if w.executor.Init(w.cfg, w.selectedBrokerURL, w) == false {
		ERROR.Println("Failed to initialize the underlying container engine: ", config.Worker.ContainerManagement)
		return false
	}

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
	w.ticker = time.NewTicker(time.Second * 10)
	go func() {
		for {
			<-w.ticker.C
			w.onTimer()
		}
	}()

	w.publishMyself()

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

func (w *Worker) Process(msg *RecvMessage) error {
	var err error
	INFO.Println(msg.Type)

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
		var dockerImage DockerImage
		err = json.Unmarshal(msg.PayLoad, &dockerImage)
		if err == nil {
			w.onPrefetchImage(&dockerImage)
		}
	}

	return err
}

func (w *Worker) onTimer() {
	w.profile.Workload = w.getNumTasks()
	w.heartbeat()
}

func (w *Worker) getNumTasks() int {
	w.taskList_lock.RLock()
	defer w.taskList_lock.RUnlock()

	return len(w.allTasks)
}

func (w *Worker) publishMyself() {
	msg := SendMessage{Type: "WORKER_JOIN", RoutingKey: "heartbeat.", From: w.id, PayLoad: w.profile}
	INFO.Println(msg)
	w.communicator.Publish(&msg)
}

func (w *Worker) unpublishMyself() {
	msg := SendMessage{Type: "WORKER_LEAVE", RoutingKey: "heartbeat.", From: w.id, PayLoad: w.profile}
	INFO.Println(msg)
	w.communicator.Publish(&msg)
}

func (w *Worker) heartbeat() {
	msg := SendMessage{Type: "WORKER_HEARTBEAT", RoutingKey: "heartbeat.", From: w.id, PayLoad: w.profile}
	DEBUG.Println(msg)
	w.communicator.Publish(&msg)
}

func (w *Worker) onAddInput(from string, flow *FlowInfo) {
	w.executor.onAddInput(flow)
}

func (w *Worker) onRemoveInput(from string, flow *FlowInfo) {
	w.executor.onRemoveInput(flow)
}

func (w *Worker) TaskUpdate(topologyName string, taskName string, taskID string, serviceIntentID string, state string) {
	tp := TaskUpdate{}
	tp.TopologyName = topologyName
	tp.TaskName = taskName
	tp.TaskID = taskID
	tp.ServiceIntentID = serviceIntentID

	tp.Status = state

	taskUpdateMsg := SendMessage{Type: "TASK_UPDATE", RoutingKey: "task.", From: w.id, PayLoad: tp}

	go w.communicator.Publish(&taskUpdateMsg)
}

func (w *Worker) TaskInfo(topologyName string, taskName string, taskID string, serviceIntentID string, info string) {
	tInfo := TaskInfo{}
	tInfo.TopologyName = topologyName
	tInfo.TaskName = taskName
	tInfo.TaskID = taskID

	tInfo.Info = info

	taskUpdateMsg := SendMessage{Type: "TASK_INFO", RoutingKey: "task.", From: w.id, PayLoad: tInfo}

	go w.communicator.Publish(&taskUpdateMsg)
}

func (w *Worker) onScheduledTask(from string, task *ScheduledTaskInstance) {
	INFO.Println("execute task ", task.ID, " with operation", task.DockerImage)
	INFO.Printf("task configuration %+v\n", *task)

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

			w.TaskUpdate(existTask.TopologyName, existTask.TaskName, existTask.ID, existTask.ServiceIntentID, "paused")
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
		w.TaskUpdate(task.TopologyName, task.TaskName, task.ID, task.ServiceIntentID, "running")
	} else {
		// add the new task into the local task list
		task.Status = "paused"
		w.allTasks[task.ID] = task

		// send ACK back to the master
		w.TaskUpdate(task.TopologyName, task.TaskName, task.ID, task.ServiceIntentID, "paused")
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

				w.TaskUpdate(task.TopologyName, task.TaskName, task.ID, task.ServiceIntentID, "running")
			}
		}
	} else {
		// resume other tasks
		for _, task := range w.allTasks {
			if task.Status != "running" {
				// restart this task temporarily
				go w.executor.LaunchTask(task)
				task.Status = "running"

				w.TaskUpdate(task.TopologyName, task.TaskName, task.ID, task.ServiceIntentID, "running")
			}
		}
	}
}

func (w *Worker) onPrefetchImage(dockerImage *DockerImage) {
	INFO.Println("I am going to fetch the docker image", dockerImage.ImageName)
	go w.executor.PullImage(dockerImage.ImageName, dockerImage.ImageTag)
}
