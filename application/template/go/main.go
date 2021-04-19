package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	. "github.com/smartfog/fogflow/common/ngsi"
)

var ctxUpdateBuffer []*ContextObject

var isConfigured = false
var brokerURL = ""
var myReferenceURL = ""

var inputEntityId = ""
var inputEntityType = ""

func readConfig(fileName string) []ConfigCommand {
	config, e := ioutil.ReadFile(fileName)
	if e != nil {
		fmt.Printf("File Error: [%v]\n", e)
		os.Exit(1)
	}

	dec := json.NewDecoder(bytes.NewReader(config))
	var commands []ConfigCommand
	dec.Decode(&commands)

	return commands
}

func startApp() {
	fmt.Println("start to receive input data streams via a listening port")
}

func stopApp() {
	fmt.Println("clean up the app")
}

// handle the commands received from the engine
func handleAdmin(commands []ConfigCommand) {
	fmt.Println("=============configuration commands=============")
	fmt.Println(commands)

	handleCmds(commands)

	isConfigured = true
}

func onNotify(notifyCtxReq *NotifyContextRequest) {
	for _, ctxResponse := range notifyCtxReq.ContextResponses {
		ctxObject := CtxElement2Object(&ctxResponse.ContextElement)
		fogfunction(ctxObject, publish)
	}
}

func notify2execution() {
	// apply the configuration
	adminCfg := os.Getenv("adminCfg")
	fmt.Println("handle the initial admin configuration " + adminCfg)
	var commands []ConfigCommand
	json.Unmarshal([]byte(adminCfg), &commands)
	handleCmds(commands)

	// get the listening port number from the environment variables given by the FogFlow edge worker
	myport, err := strconv.Atoi(os.Getenv("myport"))
	if err != nil {
		fmt.Println("myport is not set up properly for receiving notification")
		return
	}

	// start the NGSI agent
	agent := NGSIAgent{Port: myport}
	agent.Start()
	agent.SetContextNotifyHandler(onNotify)
	agent.SetAdminHandler(handleAdmin)

	startApp()

	// wait for the signal to stop the main thread
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	<-c

	stopApp()
}

func query2execution() {
	query := QueryContextRequest{}

	query.Entities = make([]EntityId, 0)

	entity := EntityId{}
	entity.Type = inputEntityType
	entity.IsPattern = true
	query.Entities = append(query.Entities, entity)

	client := NGSI10Client{IoTBrokerURL: brokerURL}
	ctxObjects, err := client.QueryContext(&query)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, ctxObj := range ctxObjects {
		fogfunction(ctxObj, publish)
	}
}

func handleCmds(commands []ConfigCommand) {
	for _, cmd := range commands {
		handleCmd(cmd)
	}

	// send the updates in the buffer
	sendUpdateWithinBuffer()
}

func sendUpdateWithinBuffer() {
	for _, ctxUpdate := range ctxUpdateBuffer {
		ngsi10client := NGSI10Client{IoTBrokerURL: brokerURL}
		err := ngsi10client.UpdateContextObject(ctxUpdate)
		if err != nil {
			fmt.Println(err)
		}
	}

	ctxUpdateBuffer = nil
}

func handleCmd(cmd ConfigCommand) {
	switch cmd.CommandType {
	case "CONNECT_BROKER":
		connectBroker(cmd)
	case "SET_INPUTS":
		setInputs(cmd)
	case "SET_OUTPUTS":
		setOutputs(cmd)
	case "SET_REFERENCE":
		setReferenceURL(cmd)
	}
}

func connectBroker(cmd ConfigCommand) {
	brokerURL = cmd.BrokerURL
	fmt.Println("set brokerURL = " + brokerURL)
}

func setInputs(cmd ConfigCommand) {
	inputEntityId = cmd.InputEntityId
	inputEntityType = cmd.InputEntityType
	fmt.Println("input has been set to (Id: " + inputEntityId + ", Type : " + inputEntityType + ")")
}

func setOutputs(cmd ConfigCommand) {

}

func setReferenceURL(cmd ConfigCommand) {
	myReferenceURL = cmd.ReferenceURL
	fmt.Println("your application can subscribe addtional inputs under the reference URL: " + myReferenceURL)
}

//
// publish context entities:
//
func publish(ctxUpdate *ContextObject) {
	ctxUpdateBuffer = append(ctxUpdateBuffer, ctxUpdate)

	if brokerURL == "" {
		fmt.Println("=== broker is not configured for your update")
		return
	}

	for _, ctxUpdate := range ctxUpdateBuffer {
		ngsi10client := NGSI10Client{IoTBrokerURL: brokerURL}
		err := ngsi10client.UpdateContextObject(ctxUpdate)
		if err != nil {
			fmt.Println(err)
		}
	}

	ctxUpdateBuffer = nil
}

func runInTestMode(runOnce bool) {
	fmt.Println("=== TEST MODE ====")

	// load the configuration
	commands := readConfig("config.json")
	fmt.Println(commands)
	handleCmds(commands)

	// query the required inputs and trigger the data processing function
	query2execution()
}

func runInOperationMode() {
	fmt.Println("=== OPERATION MODE ====")

	syncMode := os.Getenv("sync")
	if syncMode == "yes" {
		// query the required inputs and trigger the data processing function
		query2execution()
	} else {
		// trigger the data processing function to handle the received notification
		notify2execution()
	}
}

func main() {
	if len(os.Args) == 2 && os.Args[1] == "-o" {
		runInOperationMode()
	} else {
		runInTestMode(true)
	}
}
