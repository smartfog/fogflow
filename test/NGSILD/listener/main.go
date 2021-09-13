package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	mux "github.com/gufranmirza/go-router"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
	"os"
	"os/signal"
	"syscall"
	 "sync"
)

var previous_num = int(0)
var current_num = int(0)
var startTime = time.Now()
var ticker *time.Ticker
var flagv = int(0)
var count = int(0)
var mutex   sync.Mutex
func main() {
	myPort := flag.Int("p", 8066, "the port of this agent")
	flag.Parse()
	router := mux.NewRouter()
	router.POST("/notifyContext", onNotify)
	go http.ListenAndServe(":"+strconv.Itoa(*myPort), router)
	fmt.Println("the subscriber is listening on port " + strconv.Itoa(*myPort))
	c := make(chan os.Signal, 1)
        signal.Notify(c, os.Interrupt)
        signal.Notify(c, syscall.SIGTERM)
        <-c

	// start a timer to do something periodically
	//ticker = time.NewTicker(time.Second)
}

/*func onTimer() {
	fmt.Println("timer")
	if current_num != previous_num {
		fmt.Printf("total =  %d, throughput = %d \r\n", current_num, current_num-previous_num)
	}
	previous_num = current_num
}*/
type publishContextFunc func(ctxObj map[string]interface{})

// publish update on FogFlow broker
func fogfunction(ctxObj map[string]interface{}, publish publishContextFunc) error {
	//fmt.Println(ctxObj)
	//============== publish data  ==============
	publish(ctxObj)
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
		//var context []interface{}
		//context = append(context, DEFAULT_CONTEXT)
		 mutex.Lock()
	         count = count + 1
		 mutex.Unlock()

		if flagv == 0 {
			ticker = time.NewTicker(1 * time.Second)
			flagv = 1
		}
		for _ = range ticker.C {
			// execute in one sec
			fmt.Println(count)
			mutex.Lock()
			count = 0
			mutex.Unlock()
		}
		w.WriteHeader(200)
	}
}
