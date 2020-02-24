package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	. "github.com/smartfog/fogflow/common/config"
	. "github.com/smartfog/fogflow/common/ngsi"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type RestApiSrv struct {
	broker *ThinBroker
}

func (apisrv *RestApiSrv) Start(cfg *Config, broker *ThinBroker) {
	apisrv.broker = broker

	// start REST API server
	router, err := rest.MakeRouter(
		// standard ngsi10 API
		rest.Post("/ngsi10/updateContext", broker.UpdateContext),
		rest.Post("/ngsi10/queryContext", broker.QueryContext),
		rest.Post("/ngsi10/notifyContext", broker.NotifyContext),
		rest.Post("/ngsi10/subscribeContext", broker.SubscribeContext),
		rest.Post("/ngsi10/unsubscribeContext", broker.UnsubscribeContext),
		rest.Post("/ngsi10/notifyContextAvailability", broker.NotifyContextAvailability),
		rest.Post("/ngsi10/notifyContextAvailabilityv2", broker.Notifyv2ContextAvailability),
		// ngsiv2 API
		rest.Post("/v2/subscriptions", broker.Subscriptionv2Context),
		// api for iot-agent
		rest.Post("/v1/updateContext", broker.UpdateContext),

		// convenient ngsi10 API
		rest.Get("/ngsi10/entity", apisrv.getEntities),
		rest.Get("/v2/entities", apisrv.getEntities),
		rest.Get("/ngsi10/entity/#eid", apisrv.getEntity),
		rest.Get("/ngsi10/entity/#eid/#attr", apisrv.getAttribute),
		rest.Delete("/ngsi10/entity/#eid", apisrv.deleteEntity),

		rest.Get("/ngsi10/subscription", apisrv.getSubscriptions),
		rest.Get("/ngsi10/subscription/#sid", apisrv.getSubscription),
		rest.Delete("/ngsi10/subscription/#sid", apisrv.deleteSubscription),

		//NGSIV2
		rest.Get("/v2/subscriptions", apisrv.getv2Subscriptions),
		rest.Get("/v2/subscription/#sid", apisrv.getv2Subscription),
		rest.Delete("/v2/subscription/#sid", apisrv.deletev2Subscription),
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
		AllowedMethods:                []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:                []string{"Accept", "Content-Type", "X-Custom-Header", "Origin", "Destination"},
		AccessControlAllowCredentials: true,
		AccessControlMaxAge:           3600,
	})

	api.SetApp(router)

	// for internal HTTP-based communication
	go func() {
		INFO.Printf("Starting IoT Broker on port %d for HTTP requests\n", cfg.Broker.HTTPPort)
		panic(http.ListenAndServe(":"+strconv.Itoa(cfg.Broker.HTTPPort), api.MakeHandler()))
	}()

	// for external HTTPS-based communication
	go func() {
		if cfg.HTTPS.Enabled == true {
			// Create a CA certificate pool and add cert.pem to it
			caCert, err := ioutil.ReadFile(cfg.HTTPS.CA)
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
				Addr:      ":" + strconv.Itoa(cfg.Broker.HTTPSPort),
				Handler:   api.MakeHandler(),
				TLSConfig: tlsConfig,
			}

			fmt.Printf("Starting IoT Broker on port %d for HTTPS requests\n", cfg.Broker.HTTPSPort)
			panic(server.ListenAndServeTLS(cfg.HTTPS.Certificate, cfg.HTTPS.Key))
		}
	}()
}

func (apisrv *RestApiSrv) Stop() {

}

func (apisrv *RestApiSrv) getEntities(w rest.ResponseWriter, r *rest.Request) {
	entities := apisrv.broker.getEntities()
	w.WriteJson(entities)
}

func (apisrv *RestApiSrv) getEntity(w rest.ResponseWriter, r *rest.Request) {
	var eid = r.PathParam("eid")

	entity := apisrv.broker.getEntity(eid)
	if entity == nil {
		w.WriteHeader(404)
	} else {
		w.WriteHeader(200)
		w.WriteJson(entity)
	}
}

func (apisrv *RestApiSrv) getAttribute(w rest.ResponseWriter, r *rest.Request) {
	var eid = r.PathParam("eid")
	var attrname = r.PathParam("attr")

	attribute := apisrv.broker.getAttribute(eid, attrname)
	if attribute == nil {
		w.WriteHeader(404)
	} else {
		w.WriteHeader(200)
		w.WriteJson(attribute)
	}
}

func (apisrv *RestApiSrv) deleteEntity(w rest.ResponseWriter, r *rest.Request) {
	var eid = r.PathParam("eid")

	err := apisrv.broker.deleteEntity(eid)
	if err == nil {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(404)
	}
}

func (apisrv *RestApiSrv) getSubscriptions(w rest.ResponseWriter, r *rest.Request) {
	subscriptions := apisrv.broker.getSubscriptions()
	w.WriteHeader(200)
	w.WriteJson(subscriptions)
}

/*
	Handler to delete NGSIV2 subscription by Id
*/

func (apisrv *RestApiSrv) getv2Subscriptions(w rest.ResponseWriter, r *rest.Request) {
	v2subscriptions := apisrv.broker.getv2Subscriptions()
	w.WriteHeader(200)
	w.WriteJson(v2subscriptions)
}

func (apisrv *RestApiSrv) getSubscription(w rest.ResponseWriter, r *rest.Request) {
	var sid = r.PathParam("sid")

	subscription := apisrv.broker.getSubscription(sid)
	if subscription == nil {
		w.WriteHeader(404)
	} else {
		w.WriteHeader(200)
		w.WriteJson(subscription)
	}
}

/*
	Handler to get NGSIV2 subscription by SubscriptionId
*/

func (apisrv *RestApiSrv) getv2Subscription(w rest.ResponseWriter, r *rest.Request) {
	var sid = r.PathParam("sid")

	v2subscription := apisrv.broker.getv2Subscription(sid)

	if v2subscription == nil {
		w.WriteHeader(404)
	} else {
		w.WriteHeader(200)
		w.WriteJson(v2subscription)
	}
}

func (apisrv *RestApiSrv) deleteSubscription(w rest.ResponseWriter, r *rest.Request) {
	var sid = r.PathParam("sid")

	err := apisrv.broker.deleteSubscription(sid)
	if err == nil {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(404)
	}
}

/*
	Handler to delete NGSIV2 subscription by SubscriptionId
*/

func (apisrv *RestApiSrv) deletev2Subscription(w rest.ResponseWriter, r *rest.Request) {
	var sid = r.PathParam("sid")

	err := apisrv.broker.deletev2Subscription(sid)
	if err == nil {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(404)
	}
}
