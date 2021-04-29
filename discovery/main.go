package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/ant0ine/go-json-rest/rest"

	. "github.com/smartfog/fogflow/common/config"
	. "github.com/smartfog/fogflow/common/ngsi"
)

func main() {
	cfgFile := flag.String("f", "config.json", "A configuration file")
	flag.Parse()
	config, err := LoadConfig(*cfgFile)
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("%s\n", err.Error()))
		INFO.Println("please specify the configuration file, for example, \r\n\t./discovery -f config.json")
		os.Exit(-1)
	}

	// load the certificate
	config.HTTPS.LoadConfig()

	// initialize IoT Discovery
	iotDiscovery := FastDiscovery{}
	iotDiscovery.Init(&config.HTTPS)

	// start REST API server
	router, err := rest.MakeRouter(
		// standard ngsi9 API
		rest.Post("/ngsi9/registerContext", iotDiscovery.RegisterContext),
		rest.Post("/ngsi9/discoverContextAvailability", iotDiscovery.DiscoverContextAvailability),
		rest.Post("/ngsi9/subscribeContextAvailability", iotDiscovery.SubscribeContextAvailability),
		rest.Post("/ngsi9/unsubscribeContextAvailability", iotDiscovery.UnsubscribeContextAvailability),
		rest.Delete("/ngsi9/registration/#eid", iotDiscovery.deleteRegisteredEntity),
		rest.Post("/ngsi9/UpdateLDContextAvailability/#sid", iotDiscovery.UpdateLDContextAvailability),

		// convenient ngsi9 API
		rest.Get("/ngsi9/registration/#eid", iotDiscovery.getRegisteredEntity),
		rest.Get("/ngsi9/ngsi-ld/registration/#eid", iotDiscovery.getRegisteredLDEntity),
		rest.Get("/ngsi9/subscription/#sid", iotDiscovery.getSubscription),
		rest.Get("/ngsi9/subscription", iotDiscovery.getSubscriptions),

		// for health check
		rest.Get("/ngsi9/status", iotDiscovery.getStatus),
		rest.Get("/ngsi9/broker", iotDiscovery.getBrokerList),

		// hearbeat from active brokers
		rest.Post("/ngsi9/broker", iotDiscovery.onBrokerHeartbeat),
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

	// for internal HTTP-based communication
	go func() {
		INFO.Printf("Starting IoT Discovery on port %d for internal HTTP requests\n", config.Discovery.HTTPPort)
		panic(http.ListenAndServe(":"+strconv.Itoa(config.Discovery.HTTPPort), api.MakeHandler()))
	}()

	// for external HTTPS-based communication
	go func() {
		if config.HTTPS.Enabled == true {
			// Create a CA certificate pool and add cert.pem to it
			caCert, err := ioutil.ReadFile(config.HTTPS.CA)
			if err != nil {
				log.Fatal(err)
			}
			caCertPool := x509.NewCertPool()
			caCertPool.AppendCertsFromPEM(caCert)

			// Create the TLS Config with the CA pool and enable Client certificate validation
			tlsConfig := &tls.Config{
				ClientCAs:  caCertPool,
				ClientAuth: tls.RequireAndVerifyClientCert,
			}
			tlsConfig.BuildNameToCertificate()

			// Create a Server instance to listen on the port with the TLS config
			server := &http.Server{
				Addr:      ":" + strconv.Itoa(config.Discovery.HTTPSPort),
				Handler:   api.MakeHandler(),
				TLSConfig: tlsConfig,
			}

			INFO.Printf("Starting IoT Discovery on port %d for HTTPS requests\n", config.Discovery.HTTPSPort)
			panic(server.ListenAndServeTLS(config.HTTPS.Certificate, config.HTTPS.Key))
		}
	}()

	// start a timer to do something periodically
	ticker := time.NewTicker(2 * time.Second)
	go func() {
		for _ = range ticker.C {
			iotDiscovery.OnTimer()
		}
	}()

	// wait for Control +C to quit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	<-c

	fmt.Println("Stoping IoT Discovery")

	iotDiscovery.Stop()
}
