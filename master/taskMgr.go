package main

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"strconv"
	"strings"
	"sync"

	. "fogflow/common/datamodel"
	. "fogflow/common/ngsi"
)

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

type Selector struct {
	Name               string      `json:"name"`
	Conditions         []Condition `json:"conditionList,omitempty"`
	SelectedAttributes []string    `json:"selectedAttributeList,omitempty"`
	GroupBy            []string    `json:"groupedAttributeList,omitempty"`
}

type Condition struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type TaskConfig struct {
	TaskID   string
	Name     string
	Operator string

	WorkerID string

	Status string

	Inputs  []InputEntity     // the ID list of its input context entities
	Outputs []*ContextElement // the list context elements to report its generated results
}

func (taskCfg *TaskConfig) removeInput(entityID string) {
	for i := 0; i < len(taskCfg.Inputs); i++ {
		if taskCfg.Inputs[i].ID == entityID {
			taskCfg.Inputs = append(taskCfg.Inputs[:i], taskCfg.Inputs[i+1:]...)
			i--
		}
	}
}

type InputEntity struct {
	ID            string
	Type          string
	AttributeList []string
	Location      Point
}

type InputSubscription struct {
	InputSelector               InputStreamConfig
	SubID                       string
	ReceivedEntityRegistrations map[string]*EntityRegistration
}

type DeploymentAction struct {
	ActionType string      // "ADD_TASK", "REMOVE_TASK", "ADD_INPUT", "REMOVE_INPUT"
	ActionInfo interface{} //  can be either "ScheduledTaskInstance" or "FlowInfo"
}

type KVPair struct {
	Key   string
	Value interface{}
}

type GroupInfo map[string]interface{}

func (gf *GroupInfo) Set(group *GroupInfo) {
	for key, value := range *group {
		(*gf)[key] = value
	}
}

func (gf *GroupInfo) ByID() bool {
	for key, _ := range *gf {
		pair := strings.Split(key, "-")
		if pair[1] == "EntityID" {
			return true
		}
	}

	return false
}

// to generate a unique hash code from its values in the order of sorted keys
func (gf *GroupInfo) GetHash() string {
	sortedpairs := make([]*KVPair, 0)

	for k, v := range *gf {
		DEBUG.Printf("group k: %s, v: %+v\r\n", k, v)

		kvpair := KVPair{}
		kvpair.Key = k
		kvpair.Value = v

		//add it to the end
		sortedpairs = append(sortedpairs, &kvpair)

		//sort the list
		for i := len(sortedpairs) - 1; i > 0; i-- {
			if sortedpairs[i].Key < sortedpairs[i-1].Key {
				tmp := sortedpairs[i]
				sortedpairs[i] = sortedpairs[i-1]
				sortedpairs[i-1] = tmp
			}
		}
	}

	// generate the has code
	text := ""
	for _, pair := range sortedpairs {
		temp, _ := json.Marshal(pair.Value)
		text += string(temp)
	}

	hashID := fmt.Sprintf("%08d", hash(text))

	return hashID
}

type FogFlow struct {
	Intent *TaskIntent

	//to keep the unique values of all grouped keys
	UniqueKeys map[string][]interface{}

	Subscriptions map[string]*InputSubscription

	ExecutionPlan  map[string]*TaskConfig            // represent the derived execution plan
	DeploymentPlan map[string]*ScheduledTaskInstance // represent the derived deployment plan
}

func (flow *FogFlow) Init() {
	flow.UniqueKeys = make(map[string][]interface{})
	flow.Subscriptions = make(map[string]*InputSubscription)
	flow.ExecutionPlan = make(map[string]*TaskConfig)
	flow.DeploymentPlan = make(map[string]*ScheduledTaskInstance)
}

//
// to update the execution plan based on the changes of registered context availability
//
func (flow *FogFlow) MetadataDrivenTaskOrchestration(subID string, entityAction string, registredEntity *EntityRegistration, workerSelection ProximityWorkerSelectionFn) []*DeploymentAction {
	if _, exist := flow.Subscriptions[subID]; exist == false {
		DEBUG.Println(subID, "subscription does not exist any more")
		return nil
	}

	inputSubscription := flow.Subscriptions[subID]
	entityID := registredEntity.ID
	DEBUG.Println(entityAction, " entity ", entityID, "from Subscription ", subID)
	switch entityAction {
	case "CREATE", "UPDATE":
		//update context availability
		if _, exist := inputSubscription.ReceivedEntityRegistrations[entityID]; exist {
			existEntityRegistration := inputSubscription.ReceivedEntityRegistrations[entityID]
			existEntityRegistration.Update(registredEntity)
		} else {
			inputSubscription.ReceivedEntityRegistrations[entityID] = registredEntity
		}

		//update the group key-value table for orchestration
		flow.updateGroupedKeyValueTable(inputSubscription, entityID)

		//check what needs to be instantiated when all required inputs are available
		if flow.checkInputAvailability() == true {
			return flow.expandExecutionPlan(entityID, inputSubscription, workerSelection)
		}

	case "DELETE":
		_, exist := inputSubscription.ReceivedEntityRegistrations[entityID]
		if exist == false {
			INFO.Println("entity registration has not arrived yet")
			return nil
		}

		if flow.checkInputAvailability() == true {
			return flow.removeExecutionPlan(entityID, inputSubscription)
		}

		delete(inputSubscription.ReceivedEntityRegistrations, entityID)
	}

	return nil
}

//
// to check if we already received some context registration
// for all required and subscribed context availability
//
func (flow *FogFlow) checkInputAvailability() bool {
	for _, inputSubscription := range flow.Subscriptions {
		if len(inputSubscription.ReceivedEntityRegistrations) == 0 {
			return false
		}
	}

	return true
}

//
// check the available of all required input stream for a specific task instance
//
func (flow *FogFlow) checkInputsOfTaskInstance(taskCfg *TaskConfig) bool {
	INFO.Println(taskCfg)
	INFO.Println(flow.Intent.TaskObject)

	for _, inputstream := range flow.Intent.TaskObject.InputStreams {
		entityType := inputstream.EntityType

		var exist = false
		for _, input := range taskCfg.Inputs {
			if input.Type == entityType {
				exist = true
				break
			}
		}

		if exist == false {
			return false
		}
	}

	return true
}

func (flow *FogFlow) expandExecutionPlan(entityID string, inputSubscription *InputSubscription, workerSelection ProximityWorkerSelectionFn) []*DeploymentAction {
	groups := flow.getRelevantGroups(inputSubscription, entityID)
	deploymentActions := make([]*DeploymentAction, 0)

	DEBUG.Println("[related groups]: ", groups)

	for _, group := range groups {
		INFO.Println("# hash =", group.GetHash())
		hashID := group.GetHash()
		// check if the associated task instance is already created
		if task, exist := flow.ExecutionPlan[hashID]; exist {
			entitiesList := flow.searchRelevantEntities(&group, entityID)
			for _, entity := range entitiesList {
				newInput := true
				for _, input := range task.Inputs {
					if input.ID == entity.ID {
						newInput = false
						break
					}
				}

				if newInput == true {
					DEBUG.Printf("new input %+v to task %+v\r\n", entity, task)

					inputEntity := InputEntity{}
					inputEntity.ID = entity.ID
					inputEntity.Type = entity.Type
					inputEntity.Location = entity.Location
					inputEntity.AttributeList = inputSubscription.InputSelector.SelectedAttributes

					task.Inputs = append(task.Inputs, inputEntity)

					//generate a deployment action
					flowInfo := FlowInfo{}

					flowInfo.InputStream.ID = inputEntity.ID
					flowInfo.InputStream.Type = inputEntity.Type
					flowInfo.InputStream.AttributeList = inputEntity.AttributeList

					flowInfo.TaskInstanceID = task.TaskID
					flowInfo.WorkerID = flow.DeploymentPlan[task.TaskID].WorkerID

					deploymentAction := DeploymentAction{}
					deploymentAction.ActionType = "ADD_INPUT"
					deploymentAction.ActionInfo = flowInfo

					deploymentActions = append(deploymentActions, &deploymentAction)
				} else {
					// check if the location in this input entity is changed
					locationChanged := false
					for i := 0; i < len(task.Inputs); i++ {
						if task.Inputs[i].ID == entity.ID && task.Inputs[i].Location.IsEqual(&entity.Location) == false {
							locationChanged = true
							DEBUG.Println("[location changed] entity: ", entity.ID)
							// update the input entities with the new location
							task.Inputs[i].Location = entity.Location
							break
						}
					}

					// if the location is changed, calculate the new optimal worker assignment
					newOptimalWorkerID := task.WorkerID
					if locationChanged == true {
						locations := make([]Point, 0)
						for _, input := range task.Inputs {
							locations = append(locations, input.Location)
						}

						selectedWorkerID := workerSelection(locations)
						if newOptimalWorkerID != selectedWorkerID {
							newOptimalWorkerID = selectedWorkerID
						}
					}

					// prepare the new deployment actions for this task migration
					if newOptimalWorkerID != task.WorkerID {
						INFO.Println("[Task migration]:", task.TaskID, " migrated from ", task.WorkerID, " to ", newOptimalWorkerID)

						removeTaskAction := flow.removeExistingTask(task)
						deploymentActions = append(deploymentActions, removeTaskAction)

						task.WorkerID = newOptimalWorkerID

						addTaskAction := flow.addNewTask(task)
						deploymentActions = append(deploymentActions, addTaskAction)
					}
				}
			}
		} else {
			task := TaskConfig{}

			task.TaskID = "Task." + flow.Intent.TopologyName + "." + flow.Intent.TaskObject.Name + "." + hashID
			task.Operator = flow.Intent.TaskObject.Operator
			task.Name = flow.Intent.TaskObject.Name
			task.Status = "scheduled"
			task.Inputs = flow.searchRelevantEntities(&group, entityID)
			task.Outputs = flow.generateOutputs(&group)

			flow.ExecutionPlan[hashID] = &task

			DEBUG.Printf("new task %+v, hashID %s, taskID %s\r\n", task, hashID, task.TaskID)

			//generate a deployment action
			deploymentAction := flow.addNewTask(&task)
			deploymentActions = append(deploymentActions, deploymentAction)
		}
	}

	return deploymentActions
}

func (flow *FogFlow) addNewTask(task *TaskConfig) *DeploymentAction {
	taskInstance := ScheduledTaskInstance{}

	taskInstance.ID = task.TaskID

	taskInstance.TopologyName = flow.Intent.TopologyName
	taskInstance.OperatorName = task.Operator
	taskInstance.TaskName = task.Name
	taskInstance.ServiceIntentID = flow.Intent.ServiceIntentID

	taskInstance.IsExclusive = flow.Intent.Priority.IsExclusive
	taskInstance.PriorityLevel = flow.Intent.Priority.Level
	taskInstance.Status = "scheduled"

	// set up its input streams
	taskInstance.Inputs = make([]InputStream, 0)
	for _, inputEntity := range task.Inputs {
		instream := InputStream{}
		instream.Type = inputEntity.Type
		instream.ID = inputEntity.ID
		instream.AttributeList = inputEntity.AttributeList

		taskInstance.Inputs = append(taskInstance.Inputs, instream)
	}

	// set up its output streams
	taskInstance.Outputs = make([]OutputStream, 0)
	for _, ctxElem := range task.Outputs {
		outstream := OutputStream{}
		outstream.Type = ctxElem.Entity.Type
		outstream.StreamID = ctxElem.Entity.ID
		outstream.Annotations = ctxElem.Attributes

		taskInstance.Outputs = append(taskInstance.Outputs, outstream)
	}

	// create a deployment action
	deploymentAction := DeploymentAction{}
	deploymentAction.ActionType = "ADD_TASK"
	deploymentAction.ActionInfo = taskInstance

	return &deploymentAction
}

func (flow *FogFlow) removeExistingTask(task *TaskConfig) *DeploymentAction {
	if _, exist := flow.DeploymentPlan[task.TaskID]; exist {
		//generate a deployment action
		taskInstance := ScheduledTaskInstance{}
		taskInstance.ID = task.TaskID
		taskInstance.WorkerID = flow.DeploymentPlan[task.TaskID].WorkerID

		// create a deployment action
		deploymentAction := DeploymentAction{}
		deploymentAction.ActionType = "REMOVE_TASK"
		deploymentAction.ActionInfo = taskInstance

		return &deploymentAction
	}

	return nil
}

func (flow *FogFlow) removeExecutionPlan(entityID string, inputSubscription *InputSubscription) []*DeploymentAction {
	groups := flow.getRelevantGroups(inputSubscription, entityID)

	DEBUG.Printf("removing groups = %+v\r\n", groups)

	deploymentActions := make([]*DeploymentAction, 0)

	for _, group := range groups {
		INFO.Printf("Hash of relevant groups : %s\r\n", group.GetHash())
		hashID := group.GetHash()
		// check if the associated task instance is already created
		if task, exist := flow.ExecutionPlan[hashID]; exist {
			INFO.Printf("inputs: %+v", task.Inputs)

			// remove it from the task inputs
			task.removeInput(entityID)

			//if any of the input streams is delete, the task will be terminated
			if flow.checkInputsOfTaskInstance(task) == false {
				// remove this task
				DEBUG.Printf("removing an existing task %+v\r\n", task)

				// add it into the deployment action list
				deploymentAction := flow.removeExistingTask(task)
				if deploymentAction != nil {
					deploymentActions = append(deploymentActions, deploymentAction)
				}

				// remove the group key from the table
				DEBUG.Printf(" GROUP KEY %+v\r\n", group)
				DEBUG.Printf(" table %+v\r\n", flow.UniqueKeys)

				// remove this task from the execution plan
				delete(flow.ExecutionPlan, hashID)
			} else {
				// remove only the specific input
				DEBUG.Printf("remove an existing input %+v to task %+v\r\n", entityID, task)

				//generate a deployment action
				flowInfo := FlowInfo{}

				flowInfo.InputStream.ID = entityID
				flowInfo.TaskInstanceID = task.TaskID
				flowInfo.WorkerID = flow.DeploymentPlan[task.TaskID].WorkerID

				deploymentAction := DeploymentAction{}
				deploymentAction.ActionType = "REMOVE_INPUT"
				deploymentAction.ActionInfo = flowInfo

				// add it into the deployment action list
				deploymentActions = append(deploymentActions, &deploymentAction)
			}
		}
	}

	return deploymentActions
}

func (flow *FogFlow) getLocationOfInputs(hashID string) []Point {
	locations := make([]Point, 0)

	task := flow.ExecutionPlan[hashID]

	for _, input := range task.Inputs {
		locations = append(locations, input.Location)
	}

	return locations
}

func (flow *FogFlow) updateDeploymentPlan(scheduledTask *ScheduledTaskInstance) {
	flow.DeploymentPlan[scheduledTask.ID] = scheduledTask
}

func (flow *FogFlow) updateExecutionPlan(scheduledTask *ScheduledTaskInstance) {
	for _, task := range flow.ExecutionPlan {
		if task.TaskID == scheduledTask.ID {
			task.WorkerID = scheduledTask.WorkerID
		}
	}
}

func (flow *FogFlow) removeGroupKeyFromTable(groupInfo *GroupInfo) {

}

func (flow *FogFlow) updateGroupedKeyValueTable(sub *InputSubscription, entityID string) {
	selector := sub.InputSelector
	name := selector.EntityType
	groupKey := selector.GroupBy

	if groupKey == "ALL" {
		key := name + "-" + groupKey
		_, exist := flow.UniqueKeys[key]
		if exist == false {
			flow.UniqueKeys[key] = make([]interface{}, 0)
			flow.UniqueKeys[key] = append(flow.UniqueKeys[key], "ALL")
		}
	} else {
		key := name + "-" + groupKey
		entity := sub.ReceivedEntityRegistrations[entityID]

		var value interface{}

		switch groupKey {
		case "EntityID":
			value = entity.ID
		case "EntityType":
			value = entity.Type
		default:
			value = entity.MetadataList[groupKey]
		}

		if _, exist := flow.UniqueKeys[key]; exist { // add this value for the existing key
			inList := false
			items := flow.UniqueKeys[key]
			for _, item := range items {
				if item == value {
					inList = true
					break
				}
			}

			if inList == false {
				flow.UniqueKeys[key] = append(flow.UniqueKeys[key], value)
			}
		} else { // create a new key
			flow.UniqueKeys[key] = make([]interface{}, 0)
			flow.UniqueKeys[key] = append(flow.UniqueKeys[key], value)
		}
	}

	DEBUG.Println("[Unique Key Table]:", flow.UniqueKeys)
}

func (flow *FogFlow) getRelevantGroups(sub *InputSubscription, entityID string) []GroupInfo {
	// group set for the current selector
	groups := make([]GroupInfo, 0)
	selector := sub.InputSelector
	name := selector.EntityType

	entity := sub.ReceivedEntityRegistrations[entityID]

	myKeySet := make(map[string]bool)
	info := make(GroupInfo)

	groupKey := selector.GroupBy

	key := name + "-" + groupKey
	if groupKey == "ALL" {
		info[key] = "ALL"
	} else {
		var value interface{}
		switch groupKey {
		case "EntityID":
			value = entity.ID
		case "EntityType":
			value = entity.Type
		default:
			value = entity.MetadataList[groupKey]
		}
		info[key] = value
	}
	myKeySet[key] = true

	groups = append(groups, info)

	// multiple with all other keys
	for uniqueKey, uniqueValueItemList := range flow.UniqueKeys {
		if _, exist := myKeySet[uniqueKey]; exist == false {
			oldgroups := groups
			groups = make([]GroupInfo, 0)

			for _, uniqueValue := range uniqueValueItemList {
				for _, info := range oldgroups {
					newInfo := make(GroupInfo)
					newInfo.Set(&info)
					newInfo[uniqueKey] = uniqueValue

					groups = append(groups, newInfo)
				}
			}
		}
	}

	return groups
}

func (flow *FogFlow) searchRelevantEntities(group *GroupInfo, updatedEntityID string) []InputEntity {
	entities := make([]InputEntity, 0)

	for _, inputSub := range flow.Subscriptions {
		selector := inputSub.InputSelector

		// optimization for this specific case
		/*		if group.ByID() == true {
				entityRegistration := inputSub.ReceivedEntityRegistrations[updatedEntityID]

				DEBUG.Printf("REGISTR %+v\r\n", entityRegistration)

				inputEntity := InputEntity{}
				inputEntity.ID = entityRegistration.ID
				inputEntity.Type = entityRegistration.Type

				inputEntity.AttributeList = selector.SelectedAttributes

				//the location metadata will be used later to decide where to deploy the fog function instance
				inputEntity.Location = entityRegistration.getLocation()

				entities = append(entities, inputEntity)
			} else {  */
		// construct the restriction
		restrictions := make(map[string]interface{})
		key := selector.GroupBy
		groupKey := selector.EntityType + "-" + key
		if v, exist := (*group)[groupKey]; exist {
			restrictions[key] = v
		}

		// filtering
		for _, entityRegistration := range inputSub.ReceivedEntityRegistrations {
			if entityRegistration.IsMatched(restrictions) == true {
				inputEntity := InputEntity{}
				inputEntity.ID = entityRegistration.ID
				inputEntity.Type = entityRegistration.Type

				inputEntity.AttributeList = selector.SelectedAttributes

				//the location metadata will be used later to decide where to deploy the fog function instance
				inputEntity.Location = entityRegistration.GetLocation()

				entities = append(entities, inputEntity)
			}
		}
		//}
	}

	return entities
}

func (flow *FogFlow) generateOutputs(group *GroupInfo) []*ContextElement {
	outEntities := make([]*ContextElement, 0)

	for i, outputStream := range flow.Intent.TaskObject.OutputStreams {
		ctxElem := ContextElement{}

		ctxElem.Entity.ID = outputStream.EntityType + "." + group.GetHash() + "." + strconv.Itoa(i+1)
		ctxElem.Entity.Type = outputStream.EntityType

		outEntities = append(outEntities, &ctxElem)
	}

	return outEntities
}

type TaskMgr struct {
	master *Master

	//list of all task intents
	taskIntentList      map[string]*TaskIntent
	taskIntentList_lock sync.RWMutex

	//for function-based processing flows
	fogFlows      map[string]*FogFlow
	fogFlows_lock sync.RWMutex

	//mapping from availability subscription to function
	subID2FogFunc      map[string]string
	subID2FogFunc_lock sync.RWMutex
}

func NewTaskMgr(myMaster *Master) *TaskMgr {
	return &TaskMgr{master: myMaster}
}

func (tMgr *TaskMgr) Init() {
	tMgr.taskIntentList = make(map[string]*TaskIntent)
	tMgr.fogFlows = make(map[string]*FogFlow)
	tMgr.subID2FogFunc = make(map[string]string)
}

//
// deal with received task intents
//
func (tMgr *TaskMgr) handleTaskIntentUpdate(intentCtxObj *ContextObject) {
	INFO.Println("handle taskintent update")
	INFO.Println(intentCtxObj)

	taskIntent := TaskIntent{}
	jsonText, _ := json.Marshal(intentCtxObj.Attributes["intent"].Value.(map[string]interface{}))
	err := json.Unmarshal(jsonText, &taskIntent)
	if err == nil {
		INFO.Println(taskIntent)
	} else {
		INFO.Println(err)
	}

	tMgr.handleTaskIntent(&taskIntent)
}

func (tMgr *TaskMgr) handleTaskIntent(taskIntent *TaskIntent) {
	tMgr.handleASynchronousTaskIntent(taskIntent)
	// if taskIntent.SType == "SYN" {
	// 	tMgr.handleSynchronousTaskIntent(taskIntent)
	// } else if taskIntent.SType == "ASYN" {
	// 	tMgr.handleASynchronousTaskIntent(taskIntent)
	// }
}

func (tMgr *TaskMgr) handleSynchronousTaskIntent(taskIntent *TaskIntent) {
	INFO.Println("[SYNC]orchestrating task intent: %+v", taskIntent)

	fogflow := FogFlow{}

	fogflow.Init()
	fogflow.Intent = taskIntent

	fID := taskIntent.TopologyName + "." + taskIntent.TaskObject.Name

	task := taskIntent.TaskObject

	for _, inputStreamConfig := range task.InputStreams {
		INFO.Println(inputStreamConfig)
		subID := tMgr.selector2Subscription(&inputStreamConfig, taskIntent.GeoScope)

		if subID == "" {
			ERROR.Printf("failed to issue a subscription for this type of input, %+v\r\n", inputStreamConfig)
			continue
		}

		subscription := InputSubscription{}
		subscription.InputSelector = inputStreamConfig
		subscription.SubID = subID
		subscription.ReceivedEntityRegistrations = make(map[string]*EntityRegistration)

		fogflow.Subscriptions[subID] = &subscription

		// link this subscriptionId with the fog function name
		tMgr.subID2FogFunc_lock.Lock()
		tMgr.subID2FogFunc[subID] = fID
		tMgr.subID2FogFunc_lock.Unlock()
	}

	// add this fog function into the function map
	tMgr.fogFlows_lock.Lock()
	tMgr.fogFlows[fID] = &fogflow
	tMgr.fogFlows_lock.Unlock()
}

func (tMgr *TaskMgr) handleASynchronousTaskIntent(taskIntent *TaskIntent) {
	INFO.Println("[Task Intent]: ", taskIntent)

	fogflow := FogFlow{}

	fogflow.Init()
	fogflow.Intent = taskIntent

	fID := taskIntent.TopologyName + "." + taskIntent.TaskObject.Name
	task := taskIntent.TaskObject

	for _, inputStreamConfig := range task.InputStreams {
		INFO.Println(inputStreamConfig)
		subID := tMgr.selector2Subscription(&inputStreamConfig, taskIntent.GeoScope)
		if subID == "" {
			ERROR.Printf("failed to issue a subscription for this type of input, %+v\r\n", inputStreamConfig)
			continue
		}

		subscription := InputSubscription{}
		subscription.InputSelector = inputStreamConfig
		subscription.SubID = subID
		subscription.ReceivedEntityRegistrations = make(map[string]*EntityRegistration)

		fogflow.Subscriptions[subID] = &subscription

		// link this subscriptionId with the fog function name
		tMgr.subID2FogFunc_lock.Lock()
		tMgr.subID2FogFunc[subID] = fID
		tMgr.subID2FogFunc_lock.Unlock()
	}

	// add this fog function into the function map
	tMgr.fogFlows_lock.Lock()
	tMgr.fogFlows[fID] = &fogflow
	tMgr.fogFlows_lock.Unlock()
}

func (tMgr *TaskMgr) removeTaskIntent(taskIntent *TaskIntent) {
	fID := taskIntent.TopologyName + "." + taskIntent.TaskObject.Name

	// remove all related subscriptions to IoT Discovery
	sidList := make([]string, 0)
	tMgr.subID2FogFunc_lock.Lock()
	for subscriptionID, functionID := range tMgr.subID2FogFunc {
		if functionID == fID {
			sidList = append(sidList, subscriptionID)
		}
	}
	tMgr.subID2FogFunc_lock.Unlock()

	// issue unscriptions
	for _, sid := range sidList {
		tMgr.master.unsubscribeContextAvailability(sid)
	}

	// send commands to terminate all existing task instances
	tMgr.fogFlows_lock.Lock()
	var fogflow = tMgr.fogFlows[fID]
	for _, scheduledTaskInstance := range fogflow.DeploymentPlan {
		tMgr.master.TerminateTask(scheduledTaskInstance)
	}

	delete(tMgr.fogFlows, fID)
	tMgr.fogFlows_lock.Unlock()
}

func (tMgr *TaskMgr) selector2Subscription(inputSelector *InputStreamConfig, geoscope OperationScope) string {
	availabilitySubscription := SubscribeContextAvailabilityRequest{}

	// define the specified restrictions

	// apply the required entity type
	newEntity := EntityId{}
	newEntity.Type = inputSelector.EntityType
	newEntity.IsPattern = true
	availabilitySubscription.Entities = make([]EntityId, 0)
	availabilitySubscription.Entities = append(availabilitySubscription.Entities, newEntity)

	// apply the required attributes
	availabilitySubscription.Attributes = make([]string, 0)
	for _, attribute := range inputSelector.SelectedAttributes {
		if strings.EqualFold(attribute, "all") == false {
			availabilitySubscription.Attributes = append(availabilitySubscription.Attributes, attribute)
		}
	}

	// apply the required geoscope
	if inputSelector.Scoped == true {
		availabilitySubscription.Restriction.Scopes = append(availabilitySubscription.Restriction.Scopes, geoscope)
	}

	INFO.Println("[NGSI9 subscription]: ", availabilitySubscription)

	// issue the constructed subscription to IoT Discovery
	subscriptionId := tMgr.master.subscribeContextAvailability(&availabilitySubscription)
	return subscriptionId
}

//
// the main function to deal with data-driven and context aware task orchestration
//
func (tMgr *TaskMgr) HandleContextAvailabilityUpdate(subID string, entityAction string, entityRegistration *EntityRegistration) {
	INFO.Println("[Registration update]: ", subID, entityAction, entityRegistration.ID)

	tMgr.subID2FogFunc_lock.RLock()
	funcName, fogFunctionExist := tMgr.subID2FogFunc[subID]
	if fogFunctionExist == false {
		INFO.Println("this subscripption is not issued by me")
		tMgr.subID2FogFunc_lock.RUnlock()
		return
	}
	tMgr.subID2FogFunc_lock.RUnlock()

	// update the received context availability information
	tMgr.fogFlows_lock.Lock()
	defer tMgr.fogFlows_lock.Unlock()

	fogflow, fogFlowExist := tMgr.fogFlows[funcName]
	if fogFlowExist == false {
		INFO.Println("no flow established for this function: ", funcName)
		return
	}

	// derive the deployment actions according to the received registration
	deploymentActions := fogflow.MetadataDrivenTaskOrchestration(subID, entityAction, entityRegistration, tMgr.master.SelectWorker)
	if deploymentActions == nil || len(deploymentActions) == 0 {
		DEBUG.Println("nothing is triggered!!!")
		return
	}

	// schedule and send out the deployment actions
	for _, deploymentAction := range deploymentActions {
		DEBUG.Println("[Orchestration Action]: ", deploymentAction.ActionType, deploymentAction.ActionInfo)

		switch deploymentAction.ActionType {
		case "ADD_TASK":
			scheduledTaskInstance := deploymentAction.ActionInfo.(ScheduledTaskInstance)

			// figure out where to deploy this task instance
			itemList := strings.Split(scheduledTaskInstance.ID, ".")
			hashID := itemList[len(itemList)-1]

			// find out the worker close to the available inputs
			locations := fogflow.getLocationOfInputs(hashID)
			selectedWorkerID := tMgr.master.SelectWorker(locations)
			if selectedWorkerID == "" {
				ERROR.Println("===NOT ABLE TO FIND A WORKER FOR THIS TASK===")
				return
			}

			scheduledTaskInstance.WorkerID = selectedWorkerID

			// find out which implementation image to be used by the assigned worker
			operator := scheduledTaskInstance.OperatorName
			workerID := scheduledTaskInstance.WorkerID
			scheduledTaskInstance.DockerImage = tMgr.master.DetermineDockerImage(operator, workerID)

			// carry the paramemters associated with this operator
			scheduledTaskInstance.Parameters = tMgr.master.GetOperatorParamters(operator)

			if scheduledTaskInstance.WorkerID != "" {
				tMgr.master.DeployTask(&scheduledTaskInstance)
			}

			// update the deployment plan
			fogflow.updateDeploymentPlan(&scheduledTaskInstance)

			// reflect the worker assignment into the execution plan
			fogflow.updateExecutionPlan(&scheduledTaskInstance)

		case "REMOVE_TASK":
			scheduledTaskInstance := deploymentAction.ActionInfo.(ScheduledTaskInstance)
			if scheduledTaskInstance.WorkerID != "" {
				tMgr.master.TerminateTask(&scheduledTaskInstance)
			}

		case "ADD_INPUT":
			flowInfo := deploymentAction.ActionInfo.(FlowInfo)
			tMgr.master.AddInputEntity(flowInfo)

		case "REMOVE_INPUT":
			flowInfo := deploymentAction.ActionInfo.(FlowInfo)
			tMgr.master.RemoveInputEntity(flowInfo)
		}
	}
}
