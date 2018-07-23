package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	. "fogflow/common/config"
)

func main() {
	configurationFile := flag.String("f", "config.json", "A configuration file")
	flag.Parse()
	config, err := LoadConfig(*configurationFile)
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("%s\n", err.Error()))
		fmt.Println("please specify the configuration file, for example, \r\n\t./master -f config.json")
		os.Exit(-1)
	}

	myID := "Master." + strconv.Itoa(config.LLocation.LayerNo) + "." + strconv.Itoa(config.LLocation.SiteNo)

	master := Master{id: myID}
	master.Start(&config)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	<-c

	master.Quit()
}
