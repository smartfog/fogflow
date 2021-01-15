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
		rest.Post("/ngsi10/notifyLDContextAvailability", broker.NotifyLDContextAvailability),
		// ngsiv2 API
		rest.Post("/v2/subscriptions", broker.Subscriptionv2Context),
		// api for iot-agent
		// Fiware Entity Update API
		rest.Post("/v1/updateContext", broker.UpdateContext),

		//Southbound feature addition- Device Registration API
		rest.Post("/NGSI9/registerContext", broker.RegisterContext),
		rest.Delete("/NGSI9/registration/#rid", apisrv.deleteRegistration),
		rest.Get("/NGSI9/registration/#rid", apisrv.getRegistration),

		// convenient ngsi10 API
		rest.Get("/ngsi10/entity", apisrv.getEntities),
		rest.Get("/v2/entities", apisrv.getEntities),
		rest.Get("/ngsi10/entity/#eid", apisrv.getEntity),
		rest.Get("/ngsi10/entity/#eid/#attr", apisrv.getAttribute),
		rest.Delete("/ngsi10/entity/#eid", apisrv.deleteEntity),

		rest.Get("/ngsi10/subscription", apisrv.getSubscriptions),
		rest.Get("/ngsi10/subscription/#sid", apisrv.getSubscription),
		rest.Delete("/ngsi10/subscription/#sid", apisrv.deleteSubscription),

		//NGSIV2 APIs
		rest.Get("/v2/subscriptions", apisrv.getv2Subscriptions),
		rest.Get("/v2/subscription/#sid", apisrv.getv2Subscription),
		rest.Delete("/v2/subscription/#sid", apisrv.deletev2Subscription),

		//NGSI-LD APIs
		rest.Post("/ngsi-ld/v1/notifyContext/", broker.NotifyLdContext),
		rest.Post("/ngsi-ld/v1/entities/", broker.LDCreateEntity),
		rest.Get("/ngsi-ld/v1/entities/#eid", apisrv.LDGetEntity),
		rest.Get("/ngsi-ld/v1/entities", apisrv.GetQueryParamsEntities),
		rest.Post("/ngsi-ld/v1/entities/#eid/attrs", broker.LDAppendEntityAttributes),
		rest.Patch("/ngsi-ld/v1/entities/#eid/attrs", broker.LDUpdateEntityAttributes),
		rest.Patch("/ngsi-ld/v1/entities/#eid/attrs/#attr", broker.LDUpdateEntityByAttribute),
		rest.Delete("/ngsi-ld/v1/entities/#eid", apisrv.DeleteLDEntity),
		rest.Delete("/ngsi-ld/v1/entities/#eid/attrs/#attr", broker.LDDeleteEntityAttribute),

		rest.Post("/ngsi-ld/v1/subscriptions/", broker.LDCreateSubscription),
		rest.Get("/ngsi-ld/v1/subscriptions/", broker.GetLDSubscriptions),
		rest.Get("/ngsi-ld/v1/subscriptions/#sid", apisrv.GetLDSubscription),
		rest.Patch("/ngsi-ld/v1/subscriptions/#sid", broker.UpdateLDSubscription),
		rest.Delete("/ngsi-ld/v1/subscriptions/#sid", apisrv.DeleteLDSubscription),
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

//Southbound feature addition
func (apisrv *RestApiSrv) getRegistration(w rest.ResponseWriter, r *rest.Request) {
	var rid = r.PathParam("rid")

	if r.Header.Get("fiware-service") != "" && r.Header.Get("fiware-servicepath") != "" {
		rid = apisrv.broker.createIdWithFiwareHeaders(rid, r.Header.Get("fiware-service"), r.Header.Get("fiware-servicepath"))
	} /*else {
	          rest.Error(w, "Bad Request! Fiware-Service and/or Fiware-ServicePath Headers are Missing!", 400)
	          return
	  }
	  Commented because other registrations also exist, which do not have Fiware headers, like Worker, Broker, etc.*/

	registration := apisrv.broker.getRegistration(rid)
	if registration == nil {
		w.WriteHeader(404)
		w.WriteJson(nil)
	} else {
		w.WriteHeader(200)
		w.WriteJson(registration)
	}
}

func (apisrv *RestApiSrv) deleteRegistration(w rest.ResponseWriter, r *rest.Request) {
	var rid = r.PathParam("rid")

	if r.Header.Get("fiware-service") != "" && r.Header.Get("fiware-servicepath") != "" {
		rid = apisrv.broker.createIdWithFiwareHeaders(rid, r.Header.Get("fiware-service"), r.Header.Get("fiware-servicepath"))
	} else {
		rest.Error(w, "Bad Request! Fiware-Service and/or Fiware-ServicePath Headers are Missing!", 400)
		return
	}

	err := apisrv.broker.deleteRegistration(rid)
	if err == nil {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(400)
	}
}

// NGSI-LD starts from here.

func (apisrv *RestApiSrv) DeleteLDEntity(w rest.ResponseWriter, r *rest.Request) {
	var eid = r.PathParam("eid")
	if ctype, accept := r.Header.Get("Content-Type"), r.Header.Get("Accept"); (ctype == "application/json" || ctype == "application/ld+json") && accept == "application/ld+json" {
		err := apisrv.broker.ldDeleteEntity(eid)
		if err == nil {
			w.WriteHeader(204)
		} else {
			rest.Error(w, err.Error(), 404)
		}
	} else {
		rest.Error(w, "Missing Headers or Incorrect Header values!", http.StatusBadRequest)
		return
	}
}

func (apisrv *RestApiSrv) LDGetEntity(w rest.ResponseWriter, r *rest.Request) {
	var eid = r.PathParam("eid")
	if ctype, accept := r.Header.Get("Content-Type"), r.Header.Get("Accept"); ctype == "application/ld+json" || accept == "application/ld+json" || accept == "application/ld+json" || accept == "application/*" || accept == "application/json" || accept == "*/*" {
		entity := apisrv.broker.ldGetEntity(eid)
		if entity != nil {
			if accept == "application/json" || accept == " " {
				w.Header().Set("Content-Type", "application/json")
			} else {
				w.Header().Set("Content-Type", "application/ld+json")
			}
			w.WriteHeader(200)
			w.WriteJson(entity)
		} else {
			w.WriteHeader(404)
		}
	} else {
		rest.Error(w, "Missing Headers or Incorrect Header values!", http.StatusBadRequest)
		return
	}
}

// To get query parameters from NGSI-LD Entity Query Requests
func (apisrv *RestApiSrv) GetQueryParamsEntities(w rest.ResponseWriter, r *rest.Request) {
	queryValues := r.URL.Query()
	if ctype, accept := r.Header.Get("Content-Type"), r.Header.Get("Accept"); ctype == "application/ld+json" && accept == "application/ld+json" {
		if attrs, ok := queryValues["attrs"]; ok == true {
			if len(queryValues) > 1 {
				rest.Error(w, "Query parameters other than attrs are not supported in this request!", 400)
				return
			}
			if attrs[0] == "" {
				rest.Error(w, "Incorrect or missing query parameters!", http.StatusBadRequest)
				return
			}
			entities := apisrv.broker.ldEntityGetByAttribute(attrs)
			if len(entities) > 0 {
				w.WriteHeader(200)
				w.WriteJson(entities)
			} else {
				w.WriteHeader(404)
			}
			return
		} else if typ, ok := queryValues["type"]; ok == true {
			if typ[0] == "" {
				rest.Error(w, "Incorrect or missing query parameters!", http.StatusBadRequest)
				return
			}
			if eid, ok := queryValues["id"]; ok == true {
				if len(queryValues) > 2 {
					rest.Error(w, "Query parameters other than id and type are not supported in this request!", 400)
					return
				}
				if eid[0] == "" {
					rest.Error(w, "Incorrect or missing query parameters!", http.StatusBadRequest)
					return
				}
				entities := apisrv.broker.ldEntityGetById(eid, typ)
				if len(entities) > 0 {
					w.WriteHeader(200)
					w.WriteJson(entities)
				} else {
					w.WriteHeader(404)
				}
				return
			}
			if idPattern, ok := queryValues["idPattern"]; ok == true {
				if len(queryValues) > 2 {
					rest.Error(w, "Query parameters other than idPattern and type are not supported in this request!", 400)
					return
				}
				if idPattern[0] == "" {
					rest.Error(w, "Incorrect or missing query parameters!", http.StatusBadRequest)
					return
				}
				entities := apisrv.broker.ldEntityGetByIdPattern(idPattern, typ)
				if len(entities) > 0 {
					w.WriteHeader(200)
					w.WriteJson(entities)
				} else {
					w.WriteHeader(404)
				}
				return
			}

			if len(queryValues) > 1 {
				rest.Error(w, "Query parameters other than type are not supported in this request!", 400)
				return
			}
			link := r.Header.Get("Link")
			entities, err := apisrv.broker.ldEntityGetByType(typ, link)
			if err != nil {
				w.WriteHeader(500)
			} else {
				if len(entities) > 0 {
					w.WriteHeader(200)
					w.WriteJson(entities)
				} else {
					w.WriteHeader(404)
				}
				return
			}
		} else {
			rest.Error(w, "Incorrect or missing query parameters!", http.StatusBadRequest)
			return
		}
	} else {
		rest.Error(w, "Missing Headers or Incorrect Header values!", http.StatusBadRequest)
		return
	}
}

func (apisrv *RestApiSrv) GetLDSubscription(w rest.ResponseWriter, r *rest.Request) {
	if accept := r.Header.Get("Accept"); accept == "application/ld+json" {
		sid := r.PathParam("sid")
		subscription := apisrv.broker.getLDSubscription(sid)
		if subscription != nil {
			w.WriteHeader(200)
			w.WriteJson(subscription)
		} else {
			w.WriteHeader(404)
		}
	} else {
		rest.Error(w, "Missing Headers or Incorrect Header values!", http.StatusBadRequest)
	}
}

func (apisrv *RestApiSrv) DeleteLDSubscription(w rest.ResponseWriter, r *rest.Request) {
	sid := r.PathParam("sid")
	err := apisrv.broker.deleteLDSubscription(sid)
	if err != nil {
		if err.Error() == "NotFound" {
			w.WriteHeader(404)
		} else {
			rest.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(204)
	}
}
