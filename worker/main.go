package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	. "fogflow/common/config"
	. "fogflow/common/ngsi"
)

func generateID(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func main() {
	configurationFile := flag.String("f", "config.json", "A configuration file")
	id := flag.String("i", "0", "its ID in the current site")

	flag.Parse()
	config, err := LoadConfig(*configurationFile)
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("%s\n", err.Error()))
		INFO.Println("please specify the configuration file, for example, \r\n\t./worker -f config.json")
		os.Exit(-1)
	}

	// force to use only http for the communication between worker and broker
	config.HTTPS.Enabled = false

	// construct the unique id for this worker
	myID := "Worker." + config.SiteID
	if *id != "0" {
		myID = myID + "." + *id
	}

	// start the worker to deal with tasks
	var worker = &Worker{id: myID}
	worker.Start(&config)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	<-c

	worker.Quit()
}
