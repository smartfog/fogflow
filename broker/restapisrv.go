package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/ant0ine/go-json-rest/rest"

	. "github.com/smartfog/fogflow/common/config"
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

		// convenient ngsi10 API
		rest.Get("/ngsi10/entity", apisrv.getEntities),
		rest.Get("/ngsi10/entity/#eid", apisrv.getEntity),
		rest.Get("/ngsi10/entity/#eid/#attr", apisrv.getAttribute),
		rest.Delete("/ngsi10/entity/#eid", apisrv.deleteEntity),

		rest.Get("/ngsi10/subscription", apisrv.getSubscriptions),
		rest.Get("/ngsi10/subscription/#sid", apisrv.getSubscription),
		rest.Delete("/ngsi10/subscription/#sid", apisrv.deleteSubscription),
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

	go func() {
		INFO.Printf("Starting IoT Broker on %d\n", cfg.Broker.Port)
		panic(http.ListenAndServe(":"+strconv.Itoa(cfg.Broker.Port), api.MakeHandler()))
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
	w.WriteJson(subscriptions)
}

func (apisrv *RestApiSrv) getSubscription(w rest.ResponseWriter, r *rest.Request) {
	var sid = r.PathParam("sid")

	subscription := apisrv.broker.getSubscription(sid)
	if subscription == nil {
		w.WriteHeader(404)
	} else {
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
