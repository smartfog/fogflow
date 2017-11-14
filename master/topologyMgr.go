package main

import (
	"encoding/json"
	"sync"

	. "fogflow/common/datamodel"
	. "fogflow/common/ngsi"
)

type ProcessingPlane struct {
	Requirement *Requirement // requirement issued by external applications

	ExecutionPlan  []*TaskInstance          // represent the derived execution plan
	DeploymentPlan []*ScheduledTaskInstance // represent the derived deployment plan
}

type TopologyMgr struct {
	master *Master

	//list of all submitted topologies
	topologyList      map[string]*Topology
	topologyList_lock sync.RWMutex

	//for topology-based processing flows
	processList map[string]*ProcessingPlane
}

func NewTopologyMgr(myMaster *Master) *TopologyMgr {
	return &TopologyMgr{master: myMaster}
}

func (tMgr *TopologyMgr) Init() {
	tMgr.topologyList = make(map[string]*Topology)
	tMgr.processList = make(map[string]*ProcessingPlane)
}

//
// update the execution plan and deployment plan according to the system changes
//

func (tMgr *TopologyMgr) handleTopologyUpdate(responses []ContextElementResponse, sid string) {
	INFO.Println("handle topology update")
	topologyCtxObj := CtxElement2Object(&(responses[0].ContextElement))

	DEBUG.Printf("%+v\r\n", topologyCtxObj)

	// handle the incoming new requirement to trigger data processing tasks
	if topologyCtxObj.Attributes["status"].Value == "enabled" {
		tMgr.enableTopology(topologyCtxObj)
	} else if topologyCtxObj.Attributes["status"].Value == "disabled" {
		tMgr.disableTopology(topologyCtxObj)
	}
}

func (tMgr *TopologyMgr) enableTopology(topologyCtxObj *ContextObject) {
	// find out the requested topology
	topologyID := topologyCtxObj.Entity.ID
	topology := tMgr.getTopology(topologyID)
	if topology == nil {
		ERROR.Println("failed to read the topology structure")
	} else {
		INFO.Println("topology description :", topology.Description)
	}
}

func (tMgr *TopologyMgr) disableTopology(topologyCtxObj *ContextObject) {
	// find out the requested topology
	topologyID := topologyCtxObj.Entity.ID
	topology := tMgr.getTopology(topologyID)
	if topology == nil {
		ERROR.Println("failed to read the topology structure")
	} else {
		INFO.Println("topology description :", topology.Description)
	}
}

func (tMgr *TopologyMgr) handleRequirementUpdate(responses []ContextElementResponse, sid string) {
	INFO.Println("=================handle requirement update=================")

	requirementCtxObj := CtxElement2Object(&(responses[0].ContextElement))
	if requirementCtxObj.IsEmpty() == true { // the requirement is deleted
		tMgr.cancelRequirement(requirementCtxObj.Entity.ID)
		return
	}

	INFO.Println("read the requirement entity")

	// extract the issued requirement
	requirement := Requirement{}

	requirement.ID = requirementCtxObj.Entity.ID

	requirement.Output = requirementCtxObj.Attributes["output"].Value.(string)
	requirement.ScheduleMethod = requirementCtxObj.Attributes["scheduler"].Value.(string)

	if requirementCtxObj.Attributes["restriction"].Value == nil {
		requirement.Restriction = nil
	} else {
		restriction := Restriction{}
		jsondata, _ := json.Marshal(requirementCtxObj.Attributes["restriction"].Value.(map[string]interface{}))
		err := json.Unmarshal(jsondata, &restriction)
		if err != nil {
			ERROR.Println("failed to read the given restriction")
			requirement.Restriction = nil
		} else {
			requirement.Restriction = &restriction
		}
	}

	if requirementCtxObj.Metadata["topology"].Value == nil {
		ERROR.Println("the topology ID is not specified in the requirement")
		return
	}

	INFO.Println("read the topology entity")

	topologyID := requirementCtxObj.Metadata["topology"].Value.(string)
	topology := tMgr.getTopology(topologyID)
	if topology == nil {
		ERROR.Println("the topology is not submitted yet for this requirement")
	} else {
		requirement.Topology = topology
		INFO.Printf("requirement: output stream %s in topology %+v\n", requirement.Output, requirement.Topology)
		tMgr.onRequirement(&requirement)
	}
}

func (tMgr *TopologyMgr) onRequirement(requirement *Requirement) {
	if processingPlane, exist := tMgr.processList[requirement.ID]; exist {
		INFO.Printf("update requirement: %+v\r\n", requirement)
		tMgr.updateExistRequirement(processingPlane, requirement)
	} else {
		INFO.Printf("new requirement: %+v\r\n", requirement)
		tMgr.createNewRequirement(requirement)
	}
}

func (tMgr *TopologyMgr) createNewRequirement(requirement *Requirement) {
	// STEP 1: preparation

	// find out the trigger processing logic from the service topology
	rootTask := tMgr.getProcessingLogic(requirement.Output, requirement.Topology)
	if rootTask == nil {
		ERROR.Println("failed to extract the requested processing logic from the service topology")
		return
	}
	INFO.Printf("# of root task in the processing logic: %s\n", rootTask.Task.Name)

	// query input streams with regards to their scopes
	inputTypes := make([]InputStreamConfig, 0)
	findInputTypes(rootTask, &inputTypes)

	streams := tMgr.queryStreams(requirement.Restriction, inputTypes)
	INFO.Printf("# of streams: %d, \n", len(streams))
	if len(streams) == 0 {
		ERROR.Println("no input streams found!!!")
		return
	}

	// query all edge nodes available
	workers := tMgr.queryEdgeNodes()
	INFO.Printf("# of workers: %d\n", len(workers))
	if len(workers) == 0 {
		ERROR.Println("no worker found!!!")
		return
	}

	// STEP 2:  derive execution plan
	executionPlan := GenerateExcutionPlan(rootTask, streams)
	if executionPlan == nil {
		ERROR.Println("failed to derive the execution plan")
		return
	}

	// STEP 3:  derive deployment plan
	deploymentPlan := GenerateDeploymentPlan(workers, streams, executionPlan, requirement)
	if deploymentPlan == nil {
		ERROR.Println("failed to derive the deployment plan")
		return
	}

	// STEP 4:  carry out the generated deployment plan by sending out scheduled tasks
	tMgr.master.DeployTasks(deploymentPlan)

	// STEP 5:  record the processing plane
	processingPlane := ProcessingPlane{}

	processingPlane.Requirement = requirement
	processingPlane.ExecutionPlan = executionPlan
	processingPlane.DeploymentPlan = deploymentPlan

	tMgr.processList[requirement.ID] = &processingPlane
}

func (tMgr *TopologyMgr) updateExistRequirement(processingPlane *ProcessingPlane, requirement *Requirement) {
	// STEP 1: preparation

	// find out the trigger processing logic from the service topology
	rootTask := tMgr.getProcessingLogic(requirement.Output, requirement.Topology)
	if rootTask == nil {
		ERROR.Println("failed to extract the requested processing logic from the service topology")
		return
	}
	INFO.Printf("# of root task in the processing logic: %s\n", rootTask.Task.Name)

	// query input streams with regards to their scopes
	inputTypes := make([]InputStreamConfig, 0)
	findInputTypes(rootTask, &inputTypes)

	streams := tMgr.queryStreams(requirement.Restriction, inputTypes)
	INFO.Printf("# of streams: %d, \n", len(streams))
	if len(streams) == 0 {
		ERROR.Println("no input streams found!!!")
		return
	}

	// query all edge nodes available
	workers := tMgr.queryEdgeNodes()
	INFO.Printf("# of workers: %d\n", len(workers))
	if len(workers) == 0 {
		ERROR.Println("no worker found!!!")
		return
	}

	// STEP 2:  derive execution plan
	executionPlan := GenerateExcutionPlan(rootTask, streams)
	if executionPlan == nil {
		ERROR.Println("failed to derive the execution plan")
		return
	}

	INFO.Println("************************BEFORE*****************")
	printExcutionPlan(executionPlan)

	// STEP 3:  calculate the delta between the new execution plan and the old execution plan
	deltaExecutionPlan := SubtractExecutionPlan(executionPlan, processingPlane.ExecutionPlan)

	INFO.Println("************************AFTER*****************")
	printExcutionPlan(deltaExecutionPlan)
	INFO.Println("************************END*****************")

	// STEP 4:  figure out the deployment plan for the new task instances
	deploymentPlan := GenerateDeploymentPlan(workers, streams, deltaExecutionPlan, requirement)
	if deploymentPlan == nil {
		ERROR.Println("failed to derive the deployment plan")
		return
	}

	// STEP 5:  carry out the generated deployment plan by sending out scheduled tasks
	tMgr.master.DeployTasks(deploymentPlan)

	// STEP 6: update the processing plane
	processingPlane.ExecutionPlan = executionPlan

	processingPlane.DeploymentPlan = append(processingPlane.DeploymentPlan, deploymentPlan...)

	processingPlane.Requirement = requirement
}

func (tMgr *TopologyMgr) cancelRequirement(requirementID string) {
	// find out tasks that have been scheduled for this topology
	processingPlane := tMgr.processList[requirementID]

	if processingPlane != nil {
		// terminate all associated task instances
		scheduledTasks := processingPlane.DeploymentPlan
		tMgr.master.TerminateTasks(scheduledTasks)
	}
}

func (tMgr *TopologyMgr) getTopology(topologyID string) *Topology {
	//check if it is already exist in the topology list
	tMgr.topologyList_lock.RLock()
	if topology, ok := tMgr.topologyList[topologyID]; ok {
		tMgr.topologyList_lock.RUnlock()
		return topology
	}
	tMgr.topologyList_lock.RUnlock()

	topologyEntity := tMgr.master.RetrieveContextEntity(topologyID)
	if topologyEntity.Attributes["template"].Value != nil {
		topology := Topology{}
		valueData, _ := json.Marshal(topologyEntity.Attributes["template"].Value.(map[string]interface{}))
		err := json.Unmarshal(valueData, &topology)
		if err == nil {
			tMgr.topologyList_lock.Lock()
			tMgr.topologyList[topologyID] = &topology
			tMgr.topologyList_lock.Unlock()

			return &topology
		} else {
			ERROR.Println("=======loading topology structure=============")
			ERROR.Println(err)
			return nil
		}
	} else {
		return nil
	}
}

func (tMgr *TopologyMgr) queryEdgeNodes() []*ContextObject {
	query := QueryContextRequest{}

	query.Entities = make([]EntityId, 0)

	entity := EntityId{}
	entity.Type = "Worker"
	entity.IsPattern = true
	query.Entities = append(query.Entities, entity)

	restriction := Restriction{}
	restriction.Scopes = make([]OperationScope, 0)

	scope := OperationScope{}
	scope.Type = "stringQuery"
	scope.Value = "role=EdgeNode"
	restriction.Scopes = append(restriction.Scopes, scope)

	query.Restriction = restriction

	client := NGSI10Client{IoTBrokerURL: tMgr.master.BrokerURL}
	ctxObjects, err := client.QueryContext(&query, nil)
	if err != nil {
		ERROR.Println(err)
	}

	return ctxObjects
}

func (tMgr *TopologyMgr) queryStreams(restriction *Restriction, streamTypes []InputStreamConfig) []*ContextObject {
	streamObjectList := make([]*ContextObject, 0)

	for _, streamType := range streamTypes {
		query := QueryContextRequest{}

		query.Entities = make([]EntityId, 0)

		entity := EntityId{}
		entity.ID = "Stream." + streamType.Topic + ".*"
		entity.IsPattern = true
		query.Entities = append(query.Entities, entity)

		if restriction != nil && streamType.Scoped == true {
			query.Restriction = *restriction
		}

		client := NGSI10Client{IoTBrokerURL: tMgr.master.BrokerURL}
		ctxObjects, err := client.QueryContext(&query, nil)
		if err != nil {
			ERROR.Println(err)
		} else {
			streamObjectList = append(streamObjectList, ctxObjects...)
		}
	}

	return streamObjectList
}

//
// find out the processing logic in the topology
//
func (tMgr *TopologyMgr) getProcessingLogic(topic string, tp *Topology) *TaskNode {
	for _, tk := range tp.Tasks { // trigger only part of the service topology
		if isGeneratedByTask(topic, &tk) == true {
			rootTask := TaskNode{}
			rootTask.Task = &tk
			rootTask.Children = make([]*TaskNode, 0)

			// create a sub-tree for each input stream topic
			for _, input := range tk.InputStreams {
				subTreeRoot := findChildTaskTree(input.Topic, tp)
				if subTreeRoot != nil {
					rootTask.Children = append(rootTask.Children, subTreeRoot)
				}
			}

			return &rootTask
		}
	}

	return nil
}

//
// look for a sub tree producing the required topic
//
func findChildTaskTree(topic string, tp *Topology) *TaskNode {
	for _, item := range tp.Tasks {
		if isGeneratedByTask(topic, &item) {
			taskTree := TaskNode{}
			taskTree.Task = &item
			taskTree.Children = make([]*TaskNode, 0)
			// create a sub-tree for each input stream topic
			for _, input := range item.InputStreams {
				node := findChildTaskTree(input.Topic, tp)
				if node != nil {
					taskTree.Children = append(taskTree.Children, node)
				}
			}
			return &taskTree
		}
	}
	return nil
}

//
// check if this task can provide an output stream with the required topic
//
func isGeneratedByTask(topic string, task *Task) bool {
	for _, output := range task.OutputStreams {
		if output.Topic == topic {
			return true
		}
	}

	return false
}

//
//	check if a task is the root task
//
func isRootTask(task *Task, tp *Topology) bool {
	for _, output := range task.OutputStreams {
		for _, tk := range tp.Tasks {
			for _, input := range tk.InputStreams {
				if input.Topic == output.Topic {
					return false
				}
			}
		}
	}

	return true
}

//
//  find out all input stream types for the given data processing logic
//
func findInputTypes(rt *TaskNode, inputTypeList *[]InputStreamConfig) {
	if rt == nil {
		return
	}

	var isLeaf = true

	for _, child := range rt.Children {
		findInputTypes(child, inputTypeList)
		isLeaf = false
	}

	if isLeaf == true {
		for _, inputstream := range rt.Task.InputStreams {
			*inputTypeList = append(*inputTypeList, inputstream)
		}
	}
}

//
//    print out the execution plan
//
func printExcutionPlan(instances []*TaskInstance) {
	for _, instance := range instances {
		INFO.Printf("task instance %+v\n", instance)
	}

	for _, instance := range instances {
		printExcutionPlan(instance.Children)
	}
}
