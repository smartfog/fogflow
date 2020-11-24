package main

import (
	"fmt"

	. "nec-fogflow/common/ngsi"
)

type publishContextFunc func(ctxObj *ContextObject)

func fogfunction(ctxObj *ContextObject, publish publishContextFunc) error {
	fmt.Println(ctxObj)

	return nil
}
