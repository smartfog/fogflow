package main

import (
	"fmt"

	cadvisor "github.com/google/cadvisor/client"
)

type Monitor struct {
	client *cadvisor.Client
}

func (monitor *Monitor) init() {
	client, err := cadvisor.NewClient("http://localhost:8080/")
	if err != nil {
		fmt.Println("error !")
		return
	}

	monitor.client = client
}
