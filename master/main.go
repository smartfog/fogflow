package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/ant0ine/go-json-rest/rest"

	. "fogflow/common/config"
	. "fogflow/common/ngsi"
)

func main() {
	configurationFile := flag.String("f", "config.json", "A configuration file")
	flag.Parse()
	config, err := LoadConfig(*configurationFile)
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("%s\n", err.Error()))
		INFO.Println("please specify the configuration file, for example, \r\n\t./master -f config.json")
		os.Exit(-1)
	}

	config.HTTPS.Enabled = false

	myID := "Master." + config.SiteID

	master := Master{id: myID}
	master.Start(&config)

	// start REST API server
	router, err := rest.MakeRouter(
		rest.Get("/workers", master.GetWorkerList),
		rest.Get("/tasks", master.GetTaskList),
		rest.Get("/status", master.GetStatus),
	)
	if err != nil {
		ERROR.Fatal(err)
		os.Exit(-1)
	}

	api := rest.NewApi()
	api.Use(rest.DefaultCommonStack...)

	api.Use(&rest.CorsMiddleware{
		RejectNonCorsRequests: false,
		OriginValidator: func(origin string, request *rest.Request) bool {
			return true
		},
		AllowedMethods:                []string{"GET", "POST", "PUT"},
		AllowedHeaders:                []string{"Accept", "Content-Type", "X-Custom-Header", "Origin"},
		AccessControlAllowCredentials: true,
		AccessControlMaxAge:           3600,
	})

	api.SetApp(router)

	// for internal HTTP-based communication
	go func() {
		INFO.Printf("Starting REST API server on port %d\n", config.Master.RESTAPIPort)
		panic(http.ListenAndServe(":"+strconv.Itoa(config.Master.RESTAPIPort), api.MakeHandler()))
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	<-c

	master.Quit()
}
