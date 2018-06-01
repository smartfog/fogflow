package main

import (
	"encoding/json"
	"math"

	. "fogflow/common/config"
	. "fogflow/common/datamodel"
	. "fogflow/common/ngsi"
)

//
//	assign configured task instances to workers
//
func GenerateDeploymentPlan(workerObjects []*ContextObject, streamObjects []*ContextObject, executionPlan []*TaskInstance, requirement *Requirement) []*ScheduledTaskInstance {
	// perform the task assignment
	switch requirement.ScheduleMethod {
	case "random":
		taskAssignmentRandom(workerObjects, streamObjects, executionPlan)
	case "closest_first":
		taskAssignmentClosestFirst(workerObjects, streamObjects, executionPlan)
	default:
		taskAssignmentRandom(workerObjects, streamObjects, executionPlan)
	}

	// prepare the deployment actions to carry out the generated task assignment
	scheduledTasks := prepareDeploymentActions(executionPlan, requirement.Topology)

	return scheduledTasks
}

//
// prepare the deployment actions to be sent out
//
func prepareDeploymentActions(taskInstances []*TaskInstance, topology *Topology) []*ScheduledTaskInstance {
	scheduledTasks := make([]*ScheduledTaskInstance, 0)

	for _, instance := range taskInstances {
		scheduledTask := ScheduledTaskInstance{}
		scheduledTask.ID = instance.ID
		scheduledTask.WorkerID = instance.WorkerID
		scheduledTask.TaskName = instance.TaskNode.Task.Name
		scheduledTask.DockerImage = instance.TaskNode.Task.Operator
		scheduledTask.Inputs = instance.Inputs
		scheduledTask.Outputs = instance.Outputs

		INFO.Printf(" scheduled task instance : %+v\n", instance)

		scheduledTasks = append(scheduledTasks, &scheduledTask)
	}

	for _, instance := range taskInstances {
		travelTaskInstanceTree(instance.Children, &scheduledTasks)
	}

	// put the topology name and priority into each schedule task
	for _, scheduledTask := range scheduledTasks {
		scheduledTask.ServiceName = topology.Name
		scheduledTask.IsExclusive = topology.Priority.IsExclusive
		scheduledTask.PriorityLevel = topology.Priority.Level
	}

	return scheduledTasks
}

func travelTaskInstanceTree(taskInstances []*TaskInstance, scheduledTasks *[]*ScheduledTaskInstance) {
	for _, instance := range taskInstances {
		scheduledTask := ScheduledTaskInstance{}
		scheduledTask.ID = instance.ID
		scheduledTask.WorkerID = instance.WorkerID
		scheduledTask.TaskName = instance.TaskNode.Task.Name
		scheduledTask.DockerImage = instance.TaskNode.Task.Operator
		scheduledTask.Inputs = instance.Inputs
		scheduledTask.Outputs = instance.Outputs

		INFO.Printf(" scheduled task instance : %+v\n", instance)

		*scheduledTasks = append(*scheduledTasks, &scheduledTask)
	}

	// go through the sub tasks iteratively
	for _, instance := range taskInstances {
		travelTaskInstanceTree(instance.Children, scheduledTasks)
	}
}

func extractWorkerProfile(workerCtxObjects []*ContextObject) map[string]*WorkerProfile {
	workerProfileList := make(map[string]*WorkerProfile)

	for _, worker := range workerCtxObjects {
		wProfile := WorkerProfile{}

		wProfile.WID = worker.Entity.ID

		jsonPLText, _ := json.Marshal(worker.Attributes["physical_location"].Value)
		json.Unmarshal(jsonPLText, &wProfile.PLocation)

		jsonLLText, _ := json.Marshal(worker.Attributes["logical_location"].Value)
		json.Unmarshal(jsonLLText, &wProfile.LLocation)

		workerProfileList[wProfile.WID] = &wProfile
	}

	return workerProfileList
}

//
// ================= assign a task instance to a random worker in the list ================
//
func taskAssignmentRandom(workerObjects []*ContextObject, streamObjects []*ContextObject, executionPlan []*TaskInstance) {
	workerList := extractWorkerProfile(workerObjects)

	pos := 0
	for _, taskInstance := range executionPlan {
		INFO.Println("schedule the tasks at highest level")
		randomFirst(taskInstance, workerList, &pos)
	}
}

func randomFirst(taskInstance *TaskInstance, workers map[string]*WorkerProfile, index *int) {
	n := 0
	for _, w := range workers {
		if n == (*index) {
			taskInstance.WorkerID = w.WID
		}
		n++
	}

	*index = *index + 1
	if *index == len(workers) {
		*index = 0
	}

	INFO.Println("schedule task ", taskInstance.ID, " on worker ", taskInstance.WorkerID)

	for _, subTask := range taskInstance.Children {
		randomFirst(subTask, workers, index)
	}
}

//
// ================= assign each task instance close to the input data sources ================
//

func extractStreamProfile(streamObjects []*ContextObject) map[string]*StreamProfile {
	streamProfileList := make(map[string]*StreamProfile)

	for _, ctxObject := range streamObjects {
		sProfile := StreamProfile{}

		sProfile.ID = ctxObject.Entity.ID
		sProfile.Category = ctxObject.Entity.Type

		// check the location information for the stream object
		if location, exist := ctxObject.Metadata["location"]; exist {
			if location.Type == "point" {
				point := location.Value.(Point)
				sProfile.Location.Latitude = point.Latitude
				sProfile.Location.Longitude = point.Longitude
			}
		}

		sProfile.StreamObject = ctxObject

		streamProfileList[sProfile.ID] = &sProfile
	}

	return streamProfileList
}

func taskAssignmentClosestFirst(workerObjects []*ContextObject, streamObjects []*ContextObject, executionPlan []*TaskInstance) {
	// preparing the informaiton for assignment algorithm
	workerList := extractWorkerProfile(workerObjects)
	streamList := extractStreamProfile(streamObjects)

	// run the assignment algorithm
	for _, taskInstance := range executionPlan {
		closeEdgeNodeFirst(taskInstance, workerList, streamList)
	}
}

func closeEdgeNodeFirst(taskInstance *TaskInstance, workers map[string]*WorkerProfile, streams map[string]*StreamProfile) {
	// if it is the task instance at lowest layer
	if taskInstance.Children == nil {
		taskInstance.WorkerID = searchNearbyWorker(taskInstance, workers, streams)
		INFO.Println(" ALLOCATE LEAF TASK INSTANCE ", taskInstance.ID, " , on Worker ", taskInstance.WorkerID)
	} else {
		// allocate tasks at lower layer fisrt
		for _, subTask := range taskInstance.Children {
			closeEdgeNodeFirst(subTask, workers, streams)
		}

		// decide where to allocate the current task, one layer above the first child task instance is located
		parentSiteNo := workers[taskInstance.Children[0].WorkerID].LLocation.ParentSiteNo

		// find some worker at the parent site
		workerID := ""
		for _, node := range workers {
			if node.LLocation.SiteNo == parentSiteNo {
				workerID = node.WID
				break
			}
		}

		if workerID != "" {
			taskInstance.WorkerID = workerID
		} else {
			// then just assign it to the same worker node
			taskInstance.WorkerID = taskInstance.Children[0].WorkerID
		}

		INFO.Println(" ALLOCATE TASK INSTANCE ", taskInstance.ID, " , on Worker ", taskInstance.WorkerID)
	}
}

//
// decide where to allocate a leaf task instance, based on where the first input data source comes from
//
func searchNearbyWorker(taskInstance *TaskInstance, workers map[string]*WorkerProfile, streams map[string]*StreamProfile) string {
	// the locations of all input streams
	locations := make([]Point, 0)
	for _, inputstream := range taskInstance.Inputs {
		for _, streamID := range inputstream.Streams {
			if stream, exist := streams[streamID]; exist {
				point := Point{}
				point.Longitude = stream.Location.Longitude
				point.Latitude = stream.Location.Latitude

				if point.Latitude != 0 && point.Longitude != 0 {
					locations = append(locations, point)
				}
			}
		}
	}

	INFO.Printf("=========================\r\n")

	// which worker is the closest one to all input streams, each of which represents one IoT device or one input data source
	closestWorkerID := ""
	closestTotalDistance := uint64(18446744073709551615)
	for _, worker := range workers {
		INFO.Printf("check worker %+v\r\n", worker)

		wp := Point{}
		wp.Latitude = worker.PLocation.Latitude
		wp.Longitude = worker.PLocation.Longitude

		totalDistance := uint64(0)

		for _, location := range locations {
			distance := Distance(wp, location)
			totalDistance += distance
			INFO.Printf("distance = %d between %+v, %+v\r\n", distance, wp, location)
		}

		if totalDistance < closestTotalDistance {
			closestWorkerID = worker.WID
			closestTotalDistance = totalDistance
		}

		INFO.Println("closest worker ", closestWorkerID, " with the closest distance ", closestTotalDistance)
	}

	INFO.Printf("=============end============\r\n")

	return closestWorkerID
}

func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

func Distance(p1 Point, p2 Point) uint64 {
	// convert to radians
	// must cast radius as float to multiply later
	var la1, lo1, la2, lo2, r float64
	la1 = p1.Latitude * math.Pi / 180
	lo1 = p1.Longitude * math.Pi / 180
	la2 = p2.Latitude * math.Pi / 180
	lo2 = p2.Longitude * math.Pi / 180

	r = 6378100 // Earth radius in METERS

	// calculate
	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)

	return uint64(2 * r * math.Asin(math.Sqrt(h)))
}
