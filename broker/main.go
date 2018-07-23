package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	. "fogflow/common/config"
)

func main() {
	cfgFile := flag.String("f", "config.json", "A configuration file")
	flag.Parse()
	config, err := LoadConfig(*cfgFile)
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("%s\n", err.Error()))
		fmt.Println("please specify the configuration file, for example, \r\n\t./broker -f config.json")
		os.Exit(-1)
	}

	myID := "Broker." + strconv.Itoa(config.LLocation.LayerNo) + "." + strconv.Itoa(config.LLocation.SiteNo)

	// check if IoT Discovery is ready
	for {
		resp, err := http.Get(config.GetDiscoveryURL() + "/status")
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
