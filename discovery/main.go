package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/ant0ine/go-json-rest/rest"

	. "fogflow/common/config"
)

func main() {
	cfgFile := flag.String("f", "config.json", "A configuration file")
	flag.Parse()
	config, err := LoadConfig(*cfgFile)
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("%s\n", err.Error()))
		fmt.Println("please specify the configuration file, for example, \r\n\t./discovery -f config.json")
		os.Exit(-1)
	}

	// overwrite the configuration with environment variables
	if hostip, exist := os.LookupEnv("postgresql_host"); exist {
		config.Discovery.DBCfg.Host = hostip
	}
	if port, exist := os.LookupEnv("postgresql_port"); exist {
		config.Discovery.DBCfg.Port, _ = strconv.Atoi(port)
	}

	// initialize IoT Discovery
	iotDiscovery := FastDiscovery{}
	iotDiscovery.Init(&config.Discovery.DBCfg)

	// start REST API server
	router, err := rest.MakeRouter(
		// standard ngsi9 API
		rest.Post("/ngsi9/registerContext", iotDiscovery.RegisterContext),
		rest.Post("/ngsi9/discoverContextAvailability", iotDiscovery.DiscoverContextAvailability),
		rest.Post("/ngsi9/subscribeContextAvailability", iotDiscovery.SubscribeContextAvailability),
		rest.Post("/ngsi9/unsubscribeContextAvailability", iotDiscovery.UnsubscribeContextAvailability),

		// convenient ngsi9 API
		rest.Get("/ngsi9/registration/#eid", iotDiscovery.getRegisteredEntity),
		rest.Delete("/ngsi9/registration/#eid", iotDiscovery.deleteRegisteredEntity),
		rest.Get("/ngsi9/subscription/#sid", iotDiscovery.getSubscription),
		rest.Get("/ngsi9/subscription", iotDiscovery.getSubscriptions),

		rest.Get("/ngsi9/status", iotDiscovery.getStatus),
	)
	if err != nil {
		log.Fatal(err)
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

	go func() {
		INFO.Printf("Starting IoT Discovery on port %d\n", config.Discovery.Port)
		panic(http.ListenAndServe(":"+strconv.Itoa(config.Discovery.Port), api.MakeHandler()))
	}()

	// wait for Control +C to quit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	<-c

	iotDiscovery.Stop()
}
