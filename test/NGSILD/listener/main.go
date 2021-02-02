package main

import (
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

	mux "github.com/gufranmirza/go-router"
)

var previous_num = int(0)
var current_num = int(0)
var ticker *time.Ticker

func main() {
	myPort := flag.Int("p", 8066, "the port of this agent")

	flag.Parse()

	router := mux.NewRouter()
	router.POST("/notifyContext", onNotify)
	go http.ListenAndServe(":"+strconv.Itoa(*myPort), router)

	fmt.Println("the subscriber is listening on port " + strconv.Itoa(*myPort))

	// start a timer to do something periodically
	ticker = time.NewTicker(time.Second)
	go func() {
		for {
			<-ticker.C
			onTimer()
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	<-c
}

func onTimer() {
	//fmt.Println("timer")
	if current_num != previous_num {
		fmt.Printf("total =  %d, throughput = %d \r\n", current_num, current_num-previous_num)
	}

	previous_num = current_num
}

// to deal with the received notification messages

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
		//var context []interface{}
		//context = append(context, DEFAULT_CONTEXT)
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
