/*
	steps to execute
	1- go get
	2- go build

	To test subscription excute :
		./testperformance noOfthread subscription
	To test create request excute :
		./testperformance noOfthread create
	To test get Entity by ID excute :
                ./testperformance noOfthread gbid
	To test get Entity by Type excute :
                ./testperformance noOfthread gbtype
	 
*/
package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
)

var updateBrokerURL = "http://180.179.214.202:8070"  // broker URL
var id = "urn:ngsi-ld:latency:A102"
var Etype = "latency"
var my_ip = "http://180.179.214.202:8888"   // listner URL

var wg = &sync.WaitGroup{}

func main() {
	testType := string(os.Args[1])

	// set no of thred to be executed
	n := os.Args[2]
	noOfThread ,_:= strconv.Atoi(n)
	if testType == "subscription" {
	    subscribeContext(noOfThread)
	} else if testType == "create" {
	    updateContext(noOfThread)
	} else if testType == "gbid" {
	    getEntityByID(noOfThread)
	} else if testType == "gbtype" {
	    fmt.Println("test by type")
	}

}

func updateContext(n int) {
	for i := 0; i < n ; i = i+1 {
		wg.Add(1)
		go func() {
			update(i)
			wg.Done()
		}()
	}
	wg.Wait()
}


func subscribeContext(n int) {
        for i := 0; i < n ; i = i+1 {
		wg.Add(1)
		go func() {
			fmt.Println("Subscribe")
			subscribe(i)
			wg.Done()
		}()
        }
	wg.Wait()
}

func getEntityByID(n int) {
	for i := 0; i < n ; i = i+1 {
                wg.Add(1)
                go func() {
                        fmt.Println("get Entity ById")
			Eid := id + strconv.Itoa(i)
                        queryContext(Eid,updateBrokerURL)
                        wg.Done()
                }()
        }
        wg.Wait()
}

/*
   update context request
*/

func update(i int) {
	fmt.Println("This is update part")
	ctxEle := make(map[string]interface{})
	ctxEle["id"] = id + strconv.Itoa(i)
	ctxEle["type"] = Etype
	newEle1 := make(map[string]string)
	newEle1["type"] = "Property"
	newEle1["value"] = "BMW"
	newEle2 := make(map[string]string)
	newEle2["type"] = "Relationship"
	newEle2["object"] = "urn:ngsi-ld:Car:A111"
	ctxEle["brand"] = newEle1
	ctxEle["isparked"] = newEle2
        ctxElements := make([]map[string]interface{},0)
	ctxElements = append(ctxElements,ctxEle)
	fmt.Println(ctxElements)
	err := UpdateLdContext(ctxElements,updateBrokerURL)
	if err != nil {
	    fmt.Println(err)
	}
}

func subscribe(i int) {
	LdSubscription := make(map[string]interface{})
	newEntities := make([]map[string]interface{},0)
	newEntity := make(map[string]interface{})
        newEntity["type"] = Etype
	newEntity["id"] = id + strconv.Itoa(i)
        newEntities  = append(newEntities,newEntity)
	LdSubscription["entities"] = newEntities
	notification := make(map[string]interface{})
	notification["format"] = "normalized"
	endpoint := make(map[string]interface{})
	endpoint["uri"] = my_ip + "/notifyContext"
	endpoint["accept"] = "application/json"
	notification["endpoint"] = endpoint
	LdSubscription["notification"] = notification
	LdSubscription["type"] = "Subscription"
	brokerURL := updateBrokerURL
	sid , _ := SubscribeContextRequestForNGSILD(LdSubscription, brokerURL)
	fmt.Println(sid)
}
