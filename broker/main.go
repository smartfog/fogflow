package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfgFile := flag.String("f", "config.json", "A configuration file")
	flag.Parse()
	config := CreateConfig(*cfgFile)

	// overwrite the configuration with environment variables
	if value, exist := os.LookupEnv("host"); exist {
		config.Host = value
	}
	if value, exist := os.LookupEnv("discoveryURL"); exist {
		config.IoTDiscoveryURL = value
	}
	if value, exist := os.LookupEnv("mylocation"); exist {
		json.Unmarshal([]byte(value), &config.PLocation)
	}

	// check if IoT Discovery is ready
	for {
		resp, err := http.Get(config.IoTDiscoveryURL + "/status")
		if err != nil {
			ERROR.Println(err)
		} else {
			INFO.Println(resp.StatusCode)
		}

		if (err == nil) && (resp.StatusCode == 200) {
			break
		} else {
			time.Sleep(2 * time.Second)
		}
	}

	// initialize broker
	broker := ThinBroker{}
	broker.Start(&config)

	// start the REST API server
	restapi := &RestApiSrv{}
	restapi.Start(&config, &broker)

	// start a timer to do something periodically
	ticker := time.NewTicker(2 * time.Second)
	go func() {
		for _ = range ticker.C {
			broker.OnTimer()
		}
	}()

	// wait for Control+C to quit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	<-c

	// stop the timer
	ticker.Stop()

	// stop the REST API server
	restapi.Stop()

	// stop the broker
	broker.Stop()
}
