package main

import (
	"fmt"
	//. "nec-fogflow/common/ngsi"
)

type publishContextFunc func(ctxObj map[string]interface{})

func fogfunction(ctxObj map[string]interface{}, publish publishContextFunc) error {
	fmt.Println(ctxObj)
	//==============Implemet losic ==============
	publish(ctxObj)
	return nil
}
