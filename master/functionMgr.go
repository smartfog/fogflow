package main

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"sync"

	. "fogflow/common/datamodel"
	. "fogflow/common/ngsi"
)

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

// to support serverless fog computing over cloud and edges
type FogFunction struct {
	Name        string `json:"name"`
	User        string `json:"user"`
	Type        string `json:"type"`
	Code        string `json:"code"`
	DockerImage string `json:"dockerImage"`

	InputTriggers    []Selector  `json:"inputTriggers,omitempty"`
	OutputAnnotators []Annotator `json:"outputAnnotators,omitempty"`
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

type Annotator struct {
	EntityType     string `json:"entityType"`
	InputInherited bool   `json:"groupInherited,omitempty"`
}

type FunctionTask struct {
	TaskID string

	FunctionType string
	FunctionName string
	FunctionCode string
	DockerImage  string

	WorkerID string

	Status string

	Inputs  []InputEntity     // the ID list of its input context entities
	Outputs []*ContextElement // the list context elements to report its generated results
}

type InputEntity struct {
	ID   string
	Type string
}

type EntityRegistration struct {
	ID                   string
	Type                 string
	AttributesList       map[string]ContextRegistrationAttribute
	MetadataList         map[string]ContextMetadata
	ProvidingApplication string
}

func (registredEntity *EntityRegistration) IsMatched(restrictions map[string]interface{}) bool {
	DEBUG.Printf(" ====restriction = %+v\r\n", restrictions)
	DEBUG.Printf(" ====registration = %+v\r\n", registredEntity)

	matched := true

	for key, value := range restrictions {
		if key == "all" && value == "all" {
			continue
		}

		switch key {
		case "id":
			if registredEntity.ID != value {
				matched = false
				break
			}
		case "type":
			if registredEntity.Type != value {
				matched = false
				break
			}
		default:
			if registredEntity.MetadataList[key] != value {
				matched = false
				break
			}
		}
	}

	DEBUG.Printf(" ====matched = %+v\r\n", matched)

	return matched
}

func (registredEntity *EntityRegistration) Update(newUpdates *EntityRegistration) {
	if newUpdates.Type != "" {
		registredEntity.Type = newUpdates.Type
	}

	if newUpdates.ProvidingApplication != "" {
		registredEntity.ProvidingApplication = newUpdates.ProvidingApplication
	}

	for _, attribute := range newUpdates.AttributesList {
		registredEntity.AttributesList[attribute.Name] = attribute
	}

	for _, meta := range newUpdates.MetadataList {
		registredEntity.MetadataList[meta.Name] = meta
	}
}

type InputSubscription struct {
	InputSelector               Selector
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

// to generate a unique hash code from its values in the order of sorted keys
func (gInfo *GroupInfo) GetHash() string {
	sortedpairs := make([]*KVPair, 0)

	for k, v := range *gInfo {
		DEBUG.Printf("group k: %s, v: %+v\r\n", k, v)

		kvpair := KVPair{}
		kvpair.Key = k
		kvpair.Value = v

		//add it to the end
		sortedpairs = append(sortedpairs, &kvpair)

		//sort the list
		for i := len(sortedpairs) - 1; i > 0; i++ {
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
	Function *FogFunction

	//to keep the unique values of all grouped keys
	UniqueKeys map[string][]interface{}

	Subscriptions    map[string]*InputSubscription
	OutputAnnotators []Annotator

	ExecutionPlan  map[string]*FunctionTask          // represent the derived execution plan
	DeploymentPlan map[string]*ScheduledTaskInstance // represent the derived deployment plan
}

func (flow *FogFlow) Init() {
	flow.UniqueKeys = make(map[string][]interface{})
	flow.OutputAnnotators = make([]Annotator, 0)
	flow.Subscriptions = make(map[string]*InputSubscription)
	flow.ExecutionPlan = make(map[string]*FunctionTask)
	flow.DeploymentPlan = make(map[string]*ScheduledTaskInstance)
}

//
// to update the execution plan based on the changes of registered context availability
//
func (flow *FogFlow) MetadataDrivenTaskOrchestration(subID string, entityAction string, registredEntity *EntityRegistration) []*DeploymentAction {
	if _, exist := flow.Subscriptions[subID]; exist == false {
		return nil
	}

	inputSubscription := flow.Subscriptions[subID]
	entityID := registredEntity.ID

	switch entityAction {
	case "CREATE", "UPDATE":
		//update context availability
		if _, exist := inputSubscription.ReceivedEntityRegistrations[entityID]; exist {
			DEBUG.Println("update an existing entity")
			//update context availability
			existEntityRegistration := inputSubscription.ReceivedEntityRegistrations[entityID]
			existEntityRegistration.Update(registredEntity)
		} else {
			inputSubscription.ReceivedEntityRegistrations[entityID] = registredEntity
			DEBUG.Println("create new entity")
		}

		//check what needs to be instantiated when all required inputs are available
		if flow.checkInputAvailability() == true {
			INFO.Println("input available")
			return flow.expandExecutionPlan(entityID, inputSubscription)
		}

	case "DELETE":
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

func (flow *FogFlow) expandExecutionPlan(entityID string, inputSubscription *InputSubscription) []*DeploymentAction {
	flow.updateGroupedKeyValueTable(inputSubscription, entityID)

	groups := flow.getRelevantGroups(inputSubscription, entityID)

	DEBUG.Printf("groups = %+v\r\n", groups)

	deploymentActions := make([]*DeploymentAction, 0)

	for _, group := range groups {
		INFO.Println("# hash =", group.GetHash())
		hashID := group.GetHash()
		// check if the associated task instance is already created
		if task, exist := flow.ExecutionPlan[hashID]; exist {
			INFO.Printf("inputs: %+v", task.Inputs)

			entitiesList := flow.searchRelevantEntities(&group)
			for _, entity := range entitiesList {
				DEBUG.Printf("input entity : %+v\r\n", entity)
				newInput := true
				for _, input := range task.Inputs {
					if input.ID == entity.ID {
						newInput = false
						break
					}
				}

				if newInput == true {
					DEBUG.Printf("new input %+v to task %+v\r\n", entity, task)

					//generate a deployment action
					flowInfo := FlowInfo{}

					flowInfo.EntityID = entity.ID
					flowInfo.EntityType = entity.Type
					flowInfo.TaskInstanceID = task.TaskID
					flowInfo.WorkerID = flow.DeploymentPlan[task.TaskID].WorkerID

					deploymentAction := DeploymentAction{}
					deploymentAction.ActionType = "ADD_INPUT"
					deploymentAction.ActionInfo = flowInfo

					deploymentActions = append(deploymentActions, &deploymentAction)
				}
			}
		} else {
			task := FunctionTask{}
			task.TaskID = flow.Function.Name + "." + hashID
			task.FunctionType = flow.Function.Type
			task.FunctionName = flow.Function.Name
			task.FunctionCode = flow.Function.Code
			task.DockerImage = flow.Function.DockerImage
			task.Status = "scheduled"

			task.Inputs = flow.searchRelevantEntities(&group)
			task.Outputs = flow.generateOutputs(&group)

			flow.ExecutionPlan[hashID] = &task

			//generate a deployment action
			DEBUG.Printf("new task %+v\r\n", task)

			taskInstance := ScheduledTaskInstance{}

			taskInstance.ID = task.TaskID
			taskInstance.ServiceName = "system"
			taskInstance.TaskType = task.FunctionType
			taskInstance.TaskName = task.FunctionName
			taskInstance.FunctionCode = task.FunctionCode
			taskInstance.DockerImage = task.DockerImage
			taskInstance.IsExclusive = false
			taskInstance.PriorityLevel = 100
			taskInstance.Status = "scheduled"

			// set up its input streams
			taskInstance.Inputs = make([]InputStream, 0)
			for _, inputEntity := range task.Inputs {
				instream := InputStream{}
				instream.Type = inputEntity.Type
				instream.Streams = []string{inputEntity.ID}

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

			deploymentActions = append(deploymentActions, &deploymentAction)
		}
	}

	return deploymentActions
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

			if len(task.Inputs) <= 1 {
				// remove this task
				DEBUG.Printf("removing an existing task %+v\r\n", task)

				//generate a deployment action
				taskInstance := ScheduledTaskInstance{}
				taskInstance.ID = task.TaskID
				taskInstance.WorkerID = flow.DeploymentPlan[task.TaskID].WorkerID

				// create a deployment action
				deploymentAction := DeploymentAction{}
				deploymentAction.ActionType = "REMOVE_TASK"
				deploymentAction.ActionInfo = taskInstance

				// add it into the deployment action list
				deploymentActions = append(deploymentActions, &deploymentAction)

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

				flowInfo.EntityID = entityID
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

func (flow *FogFlow) updateDeploymentPlan(scheduledTask *ScheduledTaskInstance) {
	flow.DeploymentPlan[scheduledTask.ID] = scheduledTask
}

func (flow *FogFlow) removeGroupKeyFromTable(groupInfo *GroupInfo) {

}

func (flow *FogFlow) updateGroupedKeyValueTable(sub *InputSubscription, entityID string) {
	selector := sub.InputSelector
	name := selector.Name
	for _, groupKey := range selector.GroupBy {
		if groupKey == "all" {
			key := name + "-" + groupKey
			_, exist := flow.UniqueKeys[key]
			if exist == false {
				flow.UniqueKeys[key] = make([]interface{}, 0)
				flow.UniqueKeys[key] = append(flow.UniqueKeys[key], "all")
			}
		} else {
			key := name + "-" + groupKey
			entity := sub.ReceivedEntityRegistrations[entityID]

			var value interface{}

			switch groupKey {
			case "id":
				value = entity.ID
			case "type":
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
	}

	DEBUG.Printf("unique key table %+v\r\n", flow.UniqueKeys)
}

func (flow *FogFlow) getRelevantGroups(sub *InputSubscription, entityID string) []GroupInfo {
	// group set for the current selector
	groups := make([]GroupInfo, 0)
	selector := sub.InputSelector
	name := selector.Name

	entity := sub.ReceivedEntityRegistrations[entityID]

	myKeySet := make(map[string]bool)
	info := make(GroupInfo)
	for _, groupKey := range selector.GroupBy {
		DEBUG.Printf("group key = %+v\r\n", groupKey)

		key := name + "-" + groupKey
		if groupKey == "all" {
			info[key] = "all"
		} else {
			var value interface{}
			switch groupKey {
			case "id":
				value = entity.ID
			case "type":
				value = entity.Type
			default:
				value = entity.MetadataList[groupKey]
			}
			info[key] = value
		}
		myKeySet[key] = true
	}

	DEBUG.Printf("info %+v\r\n", info)

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

func (flow *FogFlow) searchRelevantEntities(group *GroupInfo) []InputEntity {
	entities := make([]InputEntity, 0)

	for _, inputSub := range flow.Subscriptions {
		selector := inputSub.InputSelector

		DEBUG.Printf("SELECTOR %+v\r\n", selector)

		//restriction
		restrictions := make(map[string]interface{})
		for _, key := range selector.GroupBy {
			groupKey := selector.Name + "-" + key
			if v, exist := (*group)[groupKey]; exist {
				restrictions[key] = v
			}
		}

		DEBUG.Printf("restriction %+v\r\n", restrictions)

		// filtering
		for _, entityRegistration := range inputSub.ReceivedEntityRegistrations {
			if entityRegistration.IsMatched(restrictions) == true {
				inputEntity := InputEntity{}
				inputEntity.ID = entityRegistration.ID
				inputEntity.Type = entityRegistration.Type

				entities = append(entities, inputEntity)
			}
		}
	}

	return entities
}

func (flow *FogFlow) generateOutputs(group *GroupInfo) []*ContextElement {
	outEntities := make([]*ContextElement, 0)

	DEBUG.Println("<<<< length of output annotators : ", len(flow.OutputAnnotators))

	for _, annotator := range flow.OutputAnnotators {
		ctxElem := ContextElement{}

		ctxElem.Entity.ID = "Stream." + annotator.EntityType + ".01"
		ctxElem.Entity.Type = annotator.EntityType

		if annotator.InputInherited == true {
			for key, value := range *group {
				DEBUG.Printf("key = %s, value = %s\r\n", key, value)
			}
		}

		outEntities = append(outEntities, &ctxElem)
	}

	return outEntities
}

type FunctionMgr struct {
	master *Master

	//list of all fog functions
	fogFuncList      map[string]*FogFunction
	fogFuncList_lock sync.RWMutex

	//for function-based processing flows
	functionFlows      map[string]*FogFlow
	functionFlows_lock sync.RWMutex

	//mapping from availability subscription to function
	subID2FogFunc      map[string]string
	subID2FogFunc_lock sync.RWMutex
}

func NewFogFunctionMgr(myMaster *Master) *FunctionMgr {
	return &FunctionMgr{master: myMaster}
}

func (fMgr *FunctionMgr) Init() {
	fMgr.fogFuncList = make(map[string]*FogFunction)
	fMgr.functionFlows = make(map[string]*FogFlow)
	fMgr.subID2FogFunc = make(map[string]string)
}

//
// deal with the updates of fog functions
//
func (fMgr *FunctionMgr) handleFogFunctionUpdate(responses []ContextElementResponse, sid string) {
	INFO.Println("handle any update of fog functions")
	ctxObj := CtxElement2Object(&(responses[0].ContextElement))

	INFO.Printf("%+v\r\n", ctxObj)

	// handle the incoming new requirement to trigger data processing tasks
	if ctxObj.Attributes["status"].Value == "enabled" {
		fogFunc := fMgr.getFogFunction(ctxObj.Entity.ID)
		fMgr.enableFogFunction(fogFunc)
	} else if ctxObj.Attributes["status"].Value == "disabled" {
		fogFunc := fMgr.getFogFunction(ctxObj.Entity.ID)
		fMgr.disableFogFunction(fogFunc)
	}
}

func (fMgr *FunctionMgr) getFogFunction(entityID string) *FogFunction {
	//check if it is already exist in the topology list
	fMgr.fogFuncList_lock.RLock()
	if fogfunc, ok := fMgr.fogFuncList[entityID]; ok {
		fMgr.fogFuncList_lock.RUnlock()
		return fogfunc
	}
	fMgr.fogFuncList_lock.RUnlock()

	fogfunctionEntity := fMgr.master.RetrieveContextEntity(entityID)
	if fogfunctionEntity.Attributes["fogfunction"].Value != nil {
		fogfunction := FogFunction{}

		valueData, _ := json.Marshal(fogfunctionEntity.Attributes["fogfunction"].Value.(map[string]interface{}))
		err := json.Unmarshal(valueData, &fogfunction)
		if err == nil {
			fMgr.fogFuncList_lock.Lock()
			fMgr.fogFuncList[entityID] = &fogfunction
			fMgr.fogFuncList_lock.Unlock()

			return &fogfunction
		} else {
			ERROR.Println("=======error happens when loading fog function structure=============")
			ERROR.Println(err)
			return nil
		}
	} else {
		return nil
	}
}

func (fMgr *FunctionMgr) enableFogFunction(f *FogFunction) {
	INFO.Printf("enable fog function %s\r\n", f.Name)
	INFO.Printf("function code: %s\r\n", f.Code)

	fogflow := FogFlow{}

	fogflow.Init()

	fogflow.Function = f

	for _, annotator := range f.OutputAnnotators {
		fogflow.OutputAnnotators = append(fogflow.OutputAnnotators, annotator)
	}

	for _, inputSelector := range f.InputTriggers {
		INFO.Printf("selector: %+v\r\n", inputSelector)
		subID := fMgr.selector2Subscription(&inputSelector)

		if subID == "" {
			ERROR.Printf("failed to issue a subscription for this type of input, %+v\r\n", inputSelector)
			continue
		}

		subscription := InputSubscription{}
		subscription.InputSelector = inputSelector
		subscription.SubID = subID
		subscription.ReceivedEntityRegistrations = make(map[string]*EntityRegistration)

		fogflow.Subscriptions[subID] = &subscription

		// link this subscriptionId with the fog function name
		fMgr.subID2FogFunc_lock.Lock()
		fMgr.subID2FogFunc[subID] = f.Name
		fMgr.subID2FogFunc_lock.Unlock()

	}

	// add this fog function into the function map
	fMgr.functionFlows_lock.Lock()
	fMgr.functionFlows[f.Name] = &fogflow
	fMgr.functionFlows_lock.Unlock()
}

func (fMgr *FunctionMgr) disableFogFunction(f *FogFunction) {
	INFO.Printf("disable fog function %s\r\n", f.Name)

	// remove this fog function from the function map
	fMgr.functionFlows_lock.Lock()
	delete(fMgr.functionFlows, f.Name)
	fMgr.functionFlows_lock.Unlock()
}

func (fMgr *FunctionMgr) selector2Subscription(inputSelector *Selector) string {
	availabilitySubscription := SubscribeContextAvailabilityRequest{}

	// define the selected attributes
	availabilitySubscription.Attributes = make([]string, 0)
	for _, attribute := range inputSelector.SelectedAttributes {
		if attribute != "all" {
			availabilitySubscription.Attributes = append(availabilitySubscription.Attributes, attribute)
		}
	}

	// define the specified restrictions
	for _, condition := range inputSelector.Conditions {
		INFO.Printf("condition: %+v\r\n", condition)

		switch condition.Type {
		case "EntityId":
			newEntity := EntityId{}
			newEntity.ID = condition.Value
			newEntity.IsPattern = false
			availabilitySubscription.Entities = make([]EntityId, 0)
			availabilitySubscription.Entities = append(availabilitySubscription.Entities, newEntity)

		case "EntityType":
			newEntity := EntityId{}
			newEntity.Type = condition.Value
			newEntity.IsPattern = true
			availabilitySubscription.Entities = make([]EntityId, 0)
			availabilitySubscription.Entities = append(availabilitySubscription.Entities, newEntity)

		case "GeoScope(Nearby)":
			scope := OperationScope{}
			scope.Type = "nearby"
			scope.Value = condition.Value
			availabilitySubscription.Restriction.Scopes = append(availabilitySubscription.Restriction.Scopes, scope)

		case "GeoScope(InCircle)":
			scope := OperationScope{}
			scope.Type = "circle"
			scope.Value = condition.Value
			availabilitySubscription.Restriction.Scopes = append(availabilitySubscription.Restriction.Scopes, scope)

		case "GeoScope(InPolygon)":
			scope := OperationScope{}
			scope.Type = "polygon"
			scope.Value = condition.Value
			availabilitySubscription.Restriction.Scopes = append(availabilitySubscription.Restriction.Scopes, scope)

		case "TimeScope":
			// to be supported

		case "StringQuery":
			scope := OperationScope{}
			scope.Type = "stringquery"
			scope.Value = condition.Value
			availabilitySubscription.Restriction.Scopes = append(availabilitySubscription.Restriction.Scopes, scope)
		}
	}

	DEBUG.Printf("issue NGSI9 subscription: %+v\r\n", availabilitySubscription)

	// issue the constructed subscription to IoT Discovery
	subscriptionId := fMgr.master.subscribeContextAvailability(&availabilitySubscription)
	return subscriptionId
}

func (fMgr *FunctionMgr) HandleContextAvailabilityUpdate(subID string, entityAction string, entityRegistration *EntityRegistration) {
	INFO.Println("handle the change of stream availability")
	INFO.Println(subID, entityAction, entityRegistration.ID)

	fMgr.subID2FogFunc_lock.RLock()
	if _, exist := fMgr.subID2FogFunc[subID]; exist == false {
		INFO.Println("this subscripption is not issued by me")
		fMgr.subID2FogFunc_lock.RUnlock()
		return
	}

	funcName := fMgr.subID2FogFunc[subID]

	fMgr.subID2FogFunc_lock.RUnlock()

	// update the received context availability information
	fMgr.functionFlows_lock.Lock()
	defer fMgr.functionFlows_lock.Unlock()

	fogflow := fMgr.functionFlows[funcName]

	deploymentActions := fogflow.MetadataDrivenTaskOrchestration(subID, entityAction, entityRegistration)

	if deploymentActions == nil || len(deploymentActions) == 0 {
		DEBUG.Println("nothing is triggered!!!")
		return
	}

	for _, deploymentAction := range deploymentActions {
		switch deploymentAction.ActionType {
		case "ADD_TASK":
			//figure out where to deploy the new task
			INFO.Printf("add task %+v\r\n", deploymentAction.ActionInfo)

			scheduledTaskInstance := deploymentAction.ActionInfo.(ScheduledTaskInstance)

			// query the location information of all available inputs
			locations := make([]Point, 0)

			// to do

			// determine where to assign the deployment action
			scheduledTaskInstance.WorkerID = fMgr.master.SelectWorker(locations)

			if scheduledTaskInstance.WorkerID != "" {
				fMgr.master.DeployTask(&scheduledTaskInstance)
			}

			// update where the task has been assigned in the deployment plan
			fogflow.updateDeploymentPlan(&scheduledTaskInstance)

		case "REMOVE_TASK":
			INFO.Printf("remove task %+v\r\n", deploymentAction.ActionInfo)

			scheduledTaskInstance := deploymentAction.ActionInfo.(ScheduledTaskInstance)
			if scheduledTaskInstance.WorkerID != "" {
				fMgr.master.TerminateTask(&scheduledTaskInstance)
			}

		case "ADD_INPUT":
			INFO.Printf("add input %+v\r\n", deploymentAction.ActionInfo)

			flowInfo := deploymentAction.ActionInfo.(FlowInfo)
			fMgr.master.AddInputEntity(flowInfo)

		case "REMOVE_INPUT":
			INFO.Printf("remove input %+v\r\n", deploymentAction.ActionInfo)

			flowInfo := deploymentAction.ActionInfo.(FlowInfo)
			fMgr.master.RemoveInputEntity(flowInfo)
		}
	}
}
