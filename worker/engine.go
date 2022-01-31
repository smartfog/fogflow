package main

import (
	. "fogflow/common/config"
	. "fogflow/common/datamodel"
)

// define the interface to interact with the underlying docker management service,
// such as docker-engine, kubernetes, and MEC controller
type Engine interface {
	Init(cfg *Config) bool
	PullImage(dockerImage string, tag string) (string, error)
	StartTask(task *ScheduledTaskInstance, brokerURL string) (string, string, error)
	StopTask(ContainerID string)
}
