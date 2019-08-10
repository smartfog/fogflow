package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	. "github.com/smartfog/fogflow/common/config"
	. "github.com/smartfog/fogflow/common/ngsi"
)

func main() {
	cfgFile := flag.String("f", "config.json", "A configuration file")
	id := flag.String("i", "0", "its ID in the current site")

	flag.Parse()
	config, err := LoadConfig(*cfgFile)
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("%s\n", err.Error()))
		ERROR.Println("please specify the configuration file, for example, \r\n\t./broker -f config.json")
		os.Exit(-1)
	}

	// load the certificates
	config.HTTPS.LoadConfig()

	myID := "Broker." + config.SiteID
	if (*id) != "0" {
		myID = myID + "." + (*id)
	}

	// check if IoT Discovery is ready
	for {
		httpClient := config.HTTPS.GetHTTPClient()
		resp, err := httpClient.Get(config.GetDiscoveryURL(true) + "/status")
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
	broker := ThinBroker{id: myID}
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
