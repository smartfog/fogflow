package datamodel

import (
	"encoding/json"

	. "github.com/smartfog/fogflow/common/ngsi"
)

// Message represents a single task invocation
type SendMessage struct {
	Type       string
	RoutingKey string
	From       string
	PayLoad    interface{}
}

type RecvMessage struct {
	Type       string
	RoutingKey string
	From       string
	PayLoad    json.RawMessage
}

type TaskUpdate struct {
	TaskID   string
	Topology string
	Status   string
}

// =========== messages used as the interfaces between different components ====================

type PhysicalLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Section   string  `json:"section"`
	District  string  `json:"district"`
	City      string  `json:"city"`
}

type LogicalLocation struct {
	LayerNo      int `json:"layer_no"`
	SiteNo       int `json:"site_no"`
	NodeNo       int `json:"node_no"`
	ParentSiteNo int `json:"parent_site_no"`
}

type ProfileInfo struct {
	StreamType string  `json:"type"`
	URL        string  `json:"url"`
	Category   string  `json:"category"`
	ProducerID string  `json:"producerID"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	Section    string  `json:"section"`
	District   string  `json:"district"`
	City       string  `json:"city"`
}

type Requirement struct {
	ID             string
	Output         string
	ScheduleMethod string

	Restriction *Restriction
	Topology    *Topology
}

type InputStreamConfig struct {
	Topic     string `json:"type"`
	Shuffling string `json:"shuffling"`
	Scoped    bool   `json:"scoped"`
}

type OutputStreamConfig struct {
	Topic string `json:"type"`
}

type Task struct {
	Name          string               `json:"name"`
	Operator      string               `json:"operator"`
	Granularity   string               `json:"groupBy"`
	InputStreams  []InputStreamConfig  `json:"input_streams"`
	OutputStreams []OutputStreamConfig `json:"output_streams"`
}

type TaskOrchestration struct {
	Task      *Task
	Topology  *Topology
	Instances []ScheduledTaskInstance
}

type TaskNode struct {
	Task     *Task
	Children []*TaskNode
}

type Priority struct {
	IsExclusive bool `json:"exclusive"`
	Level       int  `json:"level"`
}

type Topology struct {
	Description string   `json:"description"`
	Name        string   `json:"name"`
	Priority    Priority `json:"priority"`
	Trigger     string   `json:"trigger"`
	Tasks       []Task   `json:"tasks"`
}

type DockerImage struct {
	OperatorName   string
	ImageName      string
	ImageTag       string
	TargetedHWType string
	TargetedOSType string
	Prefetched     bool
}

type InputStream struct {
	Type    string
	Streams []string // a set of stream IDs
	URLs    map[string]string
}

type OutputStream struct {
	Type        string
	StreamID    string
	Annotations []ContextAttribute
}

type TaskInstance struct {
	ID       string
	Children []*TaskInstance
	TaskNode *TaskNode
	Inputs   []InputStream
	Outputs  []OutputStream
	WorkerID string
}

func compareStreamSet(setA []string, setB []string) bool {
	if len(setA) != len(setB) {
		return false
	}

	for _, idA := range setA {
		var exist = false
		for _, idB := range setB {
			if idB == idA {
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

func (myInstance *TaskInstance) Equal(otherInstance *TaskInstance) bool {
	// check the task name
	if myInstance.TaskNode.Task.Name != otherInstance.TaskNode.Task.Name {
		return false
	}

	// check the input streams
	if len(myInstance.Inputs) != len(otherInstance.Inputs) {
		return false
	}
	for _, myInputStream := range myInstance.Inputs {
		var exist = false
		for _, otherInputStream := range otherInstance.Inputs {
			if myInputStream.Type == otherInputStream.Type {
				if compareStreamSet(myInputStream.Streams, otherInputStream.Streams) == true {
					exist = true
					break
				}
			}
		}

		if exist == false {
			return false
		}
	}

	return true
}

type FlowInfo struct {
	EntityID       string
	EntityType     string
	TaskInstanceID string
	WorkerID       string
}

type ScheduledTaskInstance struct {
	ID           string
	ServiceName  string
	TaskType     string
	TaskName     string
	FunctionCode string
	DockerImage  string

	WorkerID string

	IsExclusive   bool
	PriorityLevel int

	Status string

	Inputs  []InputStream
	Outputs []OutputStream
}

type WorkerProfile struct {
	WID       string
	PLocation PhysicalLocation
	LLocation LogicalLocation
	Capacity  int
	OSType    string
	HWType    string
}

type StreamProfile struct {
	ID         string
	StreamType string
	URL        string
	Category   string
	Location   PhysicalLocation

	StreamObject *ContextObject
}
