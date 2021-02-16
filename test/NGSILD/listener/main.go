package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
	//"string"

	mux "github.com/gufranmirza/go-router"
	. "github.com/smartfog/fogflow/common/ngsi"
)

var previous_num = int(0)
var current_num = int(0)
var ticker *time.Ticker
var startTime = time.Now()
var myport = "8066"
var no_of_update = 2
var my_ip, updateBrokerURL, subscriberBrokerURL, id, Etype string
var counter = int64(0)

// main handler

func main() {
	myPort := flag.Int("p", 8066, "the port of this agent")

	flag.Parse()

	router := mux.NewRouter()
	router.POST("/notifyContext", onNotify)
	go http.ListenAndServe(":"+strconv.Itoa(*myPort), router)

	fmt.Println("the subscriber is listening on port " + strconv.Itoa(*myPort))


	configType := string(os.Args[1])
	var file string
	if configType == "Fog" {
	    file = "configFog.json"
	} else if configType == "orion" {
	    file = "configOrion.json"
	} else if configType == "scorpio" {
	    file = "configScorpio.json"
	}

	fmt.Println(file)
	setConfig(file)

	startTime = time.Now()
	fmt.Println(startTime)
	// start timer to get the latency

	fmt.Println(startTime)
	sid, _ := subscriber()
	fmt.Println(sid)
	for i := 0; i < no_of_update; i = i+1 {
	    update(i)
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	<-c
}

/*
   update context request
*/

func update(i int) {
	ctxEle := make(map[string]interface{})
	ctxEle["id"] = id + strconv.Itoa(i)
	//ctxEle["id"] = id
	ctxEle["type"] = Etype
	newEle1 := make(map[string]string)
	newEle1["type"] = "property"
	newEle1["value"] = "BMW"
	newEle2 := make(map[string]string)
	newEle2["type"] = "relationship"
	newEle2["object"] = "urn:ngsi-ld:Car:A111"
	ctxEle["brand"] = newEle1
	ctxEle["isparked"] = newEle2
        ctxElements := make([]map[string]interface{},0)
	ctxElements = append(ctxElements,ctxEle)
	err := UpdateLdContext(ctxElements,updateBrokerURL)
	if err != nil {
	    fmt.Println(err)
	}
}

/*
 subscription creation
*/

func subscriber() (string, error) {
	LdSubscription := LDSubscriptionRequest{}
	newEntity := EntityId{}
	newEntity.Type = Etype
	newEntity.IdPattern = id + ".*"
	LdSubscription.Entities = make([]EntityId, 0)
	LdSubscription.Entities = append(LdSubscription.Entities, newEntity)
	LdSubscription.Type = "Subscription"
	LdSubscription.Notification.Format = "normalized"
	LdSubscription.Notification.Endpoint.URI = my_ip + ":" + myport+ "/notifyContext"
	brokerURL := subscriberBrokerURL
	sid, err := SubscribeContextRequestForNGSILD(LdSubscription, brokerURL)
	if err != nil {
		ERROR.Println(err)
		return "", err
	} else {
		return sid, nil
	}
}

/*
   set basic config for testing the latency
*/

func setConfig(file string) {
	config1, e  := ioutil.ReadFile(file)
	if e != nil {
		fmt.Printf("File Error: [%v]\n", e)
		os.Exit(1)
	}

	dec := json.NewDecoder(bytes.NewReader(config1))
	var commands map[string]interface{}
	dec.Decode(&commands)
	my_ip  = commands["my_ip"].(string)
	updateBrokerURL = commands["update_broker_url"].(string)
	subscriberBrokerURL = commands["subscribe_broker_url"].(string)
	id = commands["id"].(string)
	Etype  = commands["type"].(string)
	fmt.Println(commands)
}

func fogfunction(ctxObj map[string]interface{}) error {
	for k, v := range ctxObj {
		fmt.Printf("%s\t%+v\n", k, v)
	}

	return nil
}

/*
	convert request into map for processing
*/

func getStringInterfaceMap(r *http.Request) (map[string]interface{}, error) {
	// Get bite array of request body
	reqBytes, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return nil, err
	}
	// Unmarshal using a generic interface
	var req interface{}
	err = json.Unmarshal(reqBytes, &req)
	if err != nil {
		fmt.Println("Invalid Request.")
		return nil, err
	}
	// Parse the JSON object into a map with string keys
	itemsMap := req.(map[string]interface{})

	if len(itemsMap) != 0 {
		return itemsMap, nil
	} else {
		return nil, errors.New("EmptyPayload!")
	}
}

/*
	handler for receiving notification
*/

func onNotify(w http.ResponseWriter, r *http.Request) {
	if ctype := r.Header.Get("Content-Type"); ctype == "application/json" || ctype == "application/ld+json" {
		fmt.Println(time.Since(startTime))
	        startTime = time.Now()
		fmt.Println(startTime)
		notifyElement, _ := getStringInterfaceMap(r)
		fmt.Println("+v\n", notifyElement)
		notifyElemtData := notifyElement["data"]
		notifyEleDatamap := notifyElemtData.([]interface{})
		w.WriteHeader(200)

		fmt.Println("=======================\r\n")

		for _, data := range notifyEleDatamap {
			notifyData := data.(map[string]interface{})
			fogfunction(notifyData)
		}
	}
}
