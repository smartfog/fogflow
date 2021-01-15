package main

import (
	"fmt"
)

type publishContextFunc func(ctxObj map[string]interface{})

// publish update on FogFlow broker

func fogfunction(ctxObj map[string]interface{}, publish publishContextFunc) error {
	fmt.Println(ctxObj)
	//============== publish data  ==============
	publish(ctxObj)
	return nil
}
