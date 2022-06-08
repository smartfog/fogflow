package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	. "fogflow/common/config"
	. "fogflow/common/ngsi"

	"github.com/ant0ine/go-json-rest/rest"
)

type RestApiSrv struct {
	broker *ThinBroker
}

func (apisrv *RestApiSrv) Start(cfg *Config, broker *ThinBroker) {
	apisrv.broker = broker

	// start REST API server
	router, err := rest.MakeRouter(
		// convenient ngsi10 API
		rest.Get("/version", apisrv.getVersion),
		rest.Get("/ngsi10/entity", apisrv.getEntities),
		rest.Get("/ngsi10/entity/#eid", apisrv.getEntity),
		rest.Get("/ngsi10/entity/#eid/#attr", apisrv.getAttribute),
		rest.Delete("/ngsi10/entity/#eid", apisrv.deleteEntity),

		rest.Get("/ngsi10/subscription", apisrv.getSubscriptions),
		rest.Get("/ngsi10/subscription/#sid", apisrv.getSubscription),
		rest.Delete("/ngsi10/subscription/#sid", apisrv.deleteSubscription),

		//============= standard ngsi10 API ===========================

		rest.Post("/ngsi10/updateContext", broker.NGSIV1_UpdateContext),

		rest.Post("/ngsi10/queryContext", broker.NGSIV1_QueryContext),

		rest.Post("/ngsi10/notifyContext", broker.NGSIV1_NotifyContext),
		rest.Post("/ngsi10/subscribeContext", broker.NGSIV1_SubscribeContext),
		rest.Post("/ngsi10/unsubscribeContext", broker.NGSIV1_UnsubscribeContext),
		rest.Post("/ngsi10/notifyContextAvailability", broker.NGSIV1_NotifyContextAvailability),

		//============= NGSIV2 APIs ===========================

		rest.Post("/v2/entities", broker.NGSIV2_createEntities),
		rest.Get("/v2/entities", broker.NGSIV2_getEntities),
		rest.Get("/v2/entities/", broker.NGSIV2_queryEntities),
		rest.Get("/v2/entities/#eid", broker.NGSIV2_getEntity),
		rest.Post("/v2/notify", broker.NGSIV2_notify),

		//============= NGSI-LD APIs ===========================

		// update
		rest.Post("/ngsi-ld/v1/entityOperations/upsert", broker.NGSILD_UpdateContext),
		rest.Post("/ngsi-ld/v1/entities/", broker.NGSILD_CreateEntity),
		rest.Post("/ngsi-ld/v1/entities", broker.NGSILD_CreateEntity),
		rest.Delete("/ngsi-ld/v1/entities/#eid", broker.NGSILD_DeleteEntity),
		rest.Delete("/ngsi-ld/v1/entities/#eid/attrs/#attr", broker.NGSILD_DeleteAttribute),

		// query
		rest.Post("/ngsi-ld/v1/entityOperations/query", broker.NGSILD_QueryByPostedFilters),
		rest.Get("/ngsi-ld/v1/entities", broker.NGSILD_QueryByParameters),
		rest.Get("/ngsi-ld/v1/entities/#eid", broker.NGSILD_QueryById),

		// subscrie and notify
		rest.Post("/ngsi-ld/v1/subscriptions/", broker.NGSILD_SubcribeContext),
		rest.Post("/ngsi-ld/v1/notifyContext/", broker.NGSILD_NotifyContext),
		rest.Delete("/ngsi-ld/v1/subscriptions/#sid", broker.NGSILD_UnsubscribeLDContext),
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

func (apisrv *RestApiSrv) getVersion(w rest.ResponseWriter, r *rest.Request) {
	version := make(map[string]string)

	version["version"] = "3.0"
	version["date"] = "2021-01-31"

	w.WriteJson(version)
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

func (apisrv *RestApiSrv) deleteSubscription(w rest.ResponseWriter, r *rest.Request) {
	var sid = r.PathParam("sid")

	err := apisrv.broker.deleteSubscription(sid)
	if err == nil {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(404)
	}
}
