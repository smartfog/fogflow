package main

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"os/signal"
	"syscall"
)

func generateID(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func main() {
	config := LoadConfig()

	// overwrite the configuration with environment variables
	if value, exist := os.LookupEnv("myip"); exist {
		config.MyIP = value
	}
	if value, exist := os.LookupEnv("discoveryURL"); exist {
		config.DiscoveryURL = value
	}
	if value, exist := os.LookupEnv("rabbitmq"); exist {
		config.MessageBus = value
	}

	// start the worker to deal with tasks
	var worker = &Worker{}
	ok := worker.Start(&config)
	if ok == false {
		ERROR.Println("failed to start the worker instance")
		return
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	<-c

	worker.Quit()
}
