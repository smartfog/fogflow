package main

import (
	"fmt"

	"github.com/satori/go.uuid"

	. "github.com/smartfog/fogflow/common/config"
	. "github.com/smartfog/fogflow/common/datamodel"
	. "github.com/smartfog/fogflow/common/ngsi"
)

type Constraint struct {
	key   string
	value interface{}
}

func GenerateExcutionPlan(rootTask *TaskNode, streamObjects []*ContextObject) []*TaskInstance {
	restrictions := make([]Constraint, 0)
	rootInstances := generateTaskInstances(rootTask, streamObjects, restrictions)

	return rootInstances
}

func SubtractExecutionPlan(currentPlan []*TaskInstance, previousPlan []*TaskInstance) []*TaskInstance {
	newTaskInstanceList := make([]*TaskInstance, 0)

	for _, instance := range currentPlan {
		var isNew = true
		for _, existInstance := range previousPlan {
			if instance.Equal(existInstance) == true {
				isNew = false
				break
			}
		}

		if isNew == true {
			newTaskInstanceList = append(newTaskInstanceList, instance)
		}
	}

	return newTaskInstanceList
}

func generateTaskInstances(root *TaskNode, streams []*ContextObject, restrictions []Constraint) []*TaskInstance {
	values := searchUniqueValuesInScope(root.Task.Granularity, restrictions, streams)

	INFO.Printf("restrictions: %v\n", restrictions)
	INFO.Printf("unique values: %v\n", values)

	// construct a new task instance for each unique value
	instances := make([]*TaskInstance, 0)
	for _, value := range values {
		instance := TaskInstance{}

		uid, err := uuid.NewV4()
		if err != nil {
			fmt.Printf("Something went wrong: %s", err)
			return instances
		}

		instance.ID = root.Task.Name + "." + uid.String()
		instance.TaskNode = root

		// update the set of restrictions for all sub tasks
		newRestrictions := make([]Constraint, len(restrictions))
		copy(newRestrictions, restrictions)

		newConstraint := Constraint{}
		newConstraint.key = root.Task.Granularity
		newConstraint.value = value

		newRestrictions = append(newRestrictions, newConstraint)

		// go through all child tasks with their updated restrictions
		for _, childtask := range root.Children {
			instanceList := generateTaskInstances(childtask, streams, newRestrictions)
			for _, subtaskInstance := range instanceList {
				instance.Children = append(instance.Children, subtaskInstance)
			}
		}

		// configure the inputs and outputs of the new task instance
		configurateTask(&instance, newRestrictions, streams)

		instances = append(instances, &instance)
	}

	return instances
}

func configurateTask(instance *TaskInstance, restrictions []Constraint, streams []*ContextObject) {
	configurateOutputs(instance, restrictions)
	configurateInputs(instance, restrictions, streams)
}

//
// configure the input streams for the current task instance
//
func configurateInputs(instance *TaskInstance, restrictions []Constraint, streams []*ContextObject) {
	instance.Inputs = make([]InputStream, 0)

	for _, stream := range instance.TaskNode.Task.InputStreams {
		onetype := InputStream{}
		onetype.Type = stream.Topic
		onetype.Streams = make([]string, 0)
		onetype.URLs = make(map[string]string)

		if instance.Children == nil || len(instance.Children) == 0 { // at the lowest layer, without any children
			conditions := make([]Constraint, 0)

			if stream.Shuffling != "broadcast" {
				conditions = append(conditions, restrictions...)
			}

			newConstraint := Constraint{}
			newConstraint.key = "type"
			newConstraint.value = stream.Topic

			conditions = append(conditions, newConstraint)

			streamSet := getMatchedStreams(conditions, streams)
			INFO.Printf("condition to select streams : %+v\n", conditions)
			INFO.Printf("returned streams : %+v\n", streamSet)

			for _, stream := range streamSet {
				INFO.Printf("==returned stream ID : %+v\n", stream.ID)
				INFO.Printf("==returned stream URL : %+v\n", stream.URL)
				onetype.Streams = append(onetype.Streams, stream.ID)

				if stream.StreamType == "PULL" {
					onetype.URLs[stream.ID] = stream.URL
				}
			}
		} else { // at the upper layer with some child tasks
			for _, subtaskInstance := range instance.Children {
				for _, output := range subtaskInstance.Outputs {
					if output.Type == stream.Topic {
						onetype.Streams = append(onetype.Streams, output.StreamID)
					}
				}
			}
		}

		INFO.Printf("add new type of input streams %+v", onetype)

		instance.Inputs = append(instance.Inputs, onetype)
	}
}

//
// configure the output streams for the current task instance
//
func configurateOutputs(instance *TaskInstance, restrictions []Constraint) {
	instance.Outputs = make([]OutputStream, 0)

	prefix := ""
	for _, constraint := range restrictions {
		prefix += fmt.Sprintf(".%v", constraint.value)
	}

	for _, item := range instance.TaskNode.Task.OutputStreams {
		out := OutputStream{}
		out.Type = item.Topic
		out.StreamID = "Stream." + item.Topic + prefix

		instance.Outputs = append(instance.Outputs, out)
	}
}

func searchUniqueValuesInScope(granularity string, restrictions []Constraint, streams []*ContextObject) []interface{} {
	INFO.Println("************RESTRICTION: ", restrictions)
	INFO.Println(granularity)

	uniqueValues := make([]interface{}, 0)

	if granularity == "*" || granularity == "all" {
		uniqueValues = append(uniqueValues, "all")
		return uniqueValues
	}

	// find out all streams that fit the restrictions
	fits := make([]*ContextObject, 0)
	for _, item := range streams {
		if IsMatchedWithRestrictions(item, restrictions) {
			fits = append(fits, item)
		}
	}

	// find out all unique values of the new scope attribute
	for _, item := range fits {
		fmt.Printf("stream object %v\n", item)

		if _, hasKey := item.Metadata[granularity]; hasKey == false {
			continue
		}

		v := item.Metadata[granularity].Value

		var exist = false
		for _, value := range uniqueValues {
			if value == v {
				exist = true
				break
			}
		}

		if exist == false {
			uniqueValues = append(uniqueValues, v)
		}
	}

	INFO.Println("************ # of unique values: ", len(uniqueValues))

	return uniqueValues
}

func getMatchedStreams(restrictions []Constraint, streams []*ContextObject) []*StreamProfile {
	// find out all streams that fit the restrictions
	INFO.Printf("============= restriction %+v===", restrictions)

	fits := make([]*StreamProfile, 0)
	for _, item := range streams {
		if IsMatchedWithRestrictions(item, restrictions) {
			sProfile := StreamProfile{}

			sProfile.ID = item.Entity.ID
			sProfile.Category = item.Entity.Type

			if _, pullBased := item.Attributes["URL"]; pullBased {
				sProfile.URL = item.Attributes["URL"].Value.(string)
			} else {
				sProfile.URL = ""
			}

			fits = append(fits, &sProfile)
		}
	}

	return fits
}

func IsMatchedWithRestrictions(stream *ContextObject, restrictions []Constraint) bool {
	fmt.Printf("====stream object %v \n====", stream)
	fmt.Printf("====restrictions %v \n====", restrictions)

	for _, constraint := range restrictions {
		k := constraint.key
		v := constraint.value

		if k == "*" || k == "all" {
			continue
		}

		if k == "type" {
			if stream.Entity.Type == v {
				continue
			}

			fmt.Printf("key = %v, value = %s, type = %s\n", k, v, stream.Entity.Type)
			fmt.Println("====not matched====")
			return false
		}

		value := stream.Metadata[k].Value
		if value != v {
			fmt.Printf("key = %v, value = %v, stream value = %v\n", k, v, stream.Metadata[k].Value)
			fmt.Println("====not matched====")
			return false
		}
	}

	fmt.Println("====matched====")

	return true
}

func getStreamProfileList(streamObjects []*ContextObject) map[string]*StreamProfile {
	streamProfileList := make(map[string]*StreamProfile)

	for _, ctxObject := range streamObjects {
		sProfile := StreamProfile{}

		sProfile.ID = ctxObject.Entity.ID
		sProfile.Category = ctxObject.Entity.Type
		sProfile.StreamObject = ctxObject

		streamProfileList[sProfile.ID] = &sProfile
	}

	return streamProfileList
}
