package main

import (
	. "fogflow/common/config"
)

type EdgeController struct {
	workerCfg *Config
}

func (mec *EdgeController) Init(cfg *Config) {

}

func (mec *EdgeController) PullImage(dockerImage string, tag string) (string, error) {
	return "test", nil
}

func (mec *EdgeController) StartTask() {

}

func (mec *EdgeController) StopTask(ContainerID string) {

}
