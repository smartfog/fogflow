package main

import (
	"flag"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func generateRandomNum() string {
	return strconv.Itoa(rand.Intn(100))
}

func main() {
	configurationFile := flag.String("f", "config.json", "A configuration file")
	myID := flag.String("i", "", "the id of this worker node")

	flag.Parse()

	config := CreateConfig(*configurationFile)

	if *myID == "" {
		*myID = generateRandomNum()
	}

	// overwrite the configuration with environment variables
	if value, exist := os.LookupEnv("myip"); exist {
		config.MyIP = value
	}
	if value, exist := os.LookupEnv("discoveryURL"); exist {
		config.IoTDiscoveryURL = value
	}
	if value, exist := os.LookupEnv("rabbitmq"); exist {
		config.MessageBus = value
	}

	master := Master{myID: *myID}
	master.Start(&config)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	<-c

	master.Quit()
}
