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
	ServiceName string
	TaskName    string
	TaskID      string
	Status      string
}

// =========== messages used as the interfaces between different components ====================

type PhysicalLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
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

type OptPreference struct {
	Minimize []string
	Maximize []string
}

type OptConstraint struct {
	Subject   string
	Relation  string
	Objective string
}

type QoS struct {
	Preference  OptPreference
	Constraints []OptConstraint
}

type ServiceIntent struct {
	ID             string         `json:"id"`
	QoS            string         `json:"qos"`
	GeoScope       OperationScope `json:"geoscope"`
	Priority       Priority       `json:"priority"`
	TopologyName   string         `json:"topology"`
	TopologyObject *Topology
}

type TaskIntent struct {
	ID          string         `json:"id"`
	QoS         string         `json:"qos"`
	GeoScope    OperationScope `json:"geoscope"`
	Priority    Priority       `json:"priority"`
	ServiceName string         `json:"service"`
	TaskObject  Task           `json:"task"`
}

type InputStreamConfig struct {
	EntityType         string   `json:"selected_type"`
	SelectedAttributes []string `json:"selected_attributes"`
	GroupBy            string   `json:"groupby"`
	Scoped             bool     `json:"scoped"`
}

type OutputStreamConfig struct {
	EntityType string `json:"entity_type"`
}

type Parameter struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

type Operator struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Parameters  []Parameter `json:"parameters"`
}

type Task struct {
	Name          string               `json:"name"`
	Operator      string               `json:"operator"`
	InputStreams  []InputStreamConfig  `json:"input_streams"`
	OutputStreams []OutputStreamConfig `json:"output_streams"`
}

func (task *Task) CanBeDivided() bool {
	var flag = true

	for _, inputStream := range task.InputStreams {
		if inputStream.GroupBy != "EntityID" {
			flag = false
			break
		}
	}

	return flag
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
	Id          string `json:"id"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Tasks       []Task `json:"tasks"`
}

type FogFunction struct {
	Id       string        `json:"id"`
	Name     string        `json:"name"`
	Topology Topology      `json:"topology"`
	Intent   ServiceIntent `json:"intent"`
}

type DockerImage struct {
	OperatorName   string `json:"operater"`
	ImageName      string `json:"name"`
	ImageTag       string `json:"tag"`
	TargetedHWType string `json:"hwType"`
	TargetedOSType string `json:"osType"`
	Prefetched     bool   `json:"prefetched"`
}

type InputStream struct {
	Type          string
	ID            string
	AttributeList []string
}

func (myInputStream *InputStream) Equal(otherInputStream *InputStream) bool {
	if myInputStream.Type == otherInputStream.Type && myInputStream.ID == myInputStream.ID {
		return true
	} else {
		return false
	}
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
			if myInputStream.Equal(&otherInputStream) == true {
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

type FlowInfo struct {
	InputStream    InputStream
	TaskInstanceID string
	WorkerID       string
}

type ScheduledTaskInstance struct {
	ID          string
	ServiceName string
	TaskName    string

	OperatorName string

	TaskType     string
	FunctionCode string
	DockerImage  string
	Parameters   []Parameter

	WorkerID string

	IsExclusive   bool
	PriorityLevel int

	Status string

	Inputs  []InputStream
	Outputs []OutputStream
}

type WorkerProfile struct {
	WID          string
	PLocation    PhysicalLocation
	GeohashID    string
	OSType       string
	HWType       string
	Capacity     int
	Workload     int
	CAdvisorPort int
	EdgeAddress  string
}

func (worker *WorkerProfile) IsOverloaded() bool {
	if worker.Workload >= worker.Capacity {
		return true
	} else {
		return false
	}
}

type StreamProfile struct {
	ID         string
	StreamType string
	URL        string
	Category   string
	Location   PhysicalLocation

	StreamObject *ContextObject
}

type PrometheusConfig struct {
	Targets []string `json:"targets"`
	Labels  struct {
		Job string `json:"job"`
	} `json:"labels"`
}
