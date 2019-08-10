package ngsi

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/ant0ine/go-json-rest/rest"
)

type NotifyContextFunc func(notifyCtxReq *NotifyContextRequest)
type NotifyContextAvailabilityFunc func(notifyCtxAvailReq *NotifyContextAvailabilityRequest)

type NGSIAgent struct {
	Port                        int
	SecurityCfg                 HTTPS
	CtxNotifyHandler            NotifyContextFunc
	CtxAvailbilityNotifyHandler NotifyContextAvailabilityFunc
}

func (agent *NGSIAgent) Start() {
	// start REST API server
	router, err := rest.MakeRouter(
		rest.Post("/notifyContext", agent.handleNotifyContext),
		rest.Post("/notifyContextAvailability", agent.handleNotifyContextAvailability),
	)
	if err != nil {
		ERROR.Fatal(err)
	}

	api := rest.NewApi()
	api.Use(rest.DefaultCommonStack...)
	api.SetApp(router)

	go func() {
		if agent.SecurityCfg.Enabled == true {
			// Create a CA certificate pool and add cert.pem to it
			caCert, err := ioutil.ReadFile(agent.SecurityCfg.CA)
			if err != nil {
				ERROR.Fatal(err)
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
				Addr:      ":" + strconv.Itoa(agent.Port),
				Handler:   api.MakeHandler(),
				TLSConfig: tlsConfig,
			}

			INFO.Printf("Starting IoT agent on port %d for HTTPS requests\n", agent.Port)
			panic(server.ListenAndServeTLS(agent.SecurityCfg.Certificate, agent.SecurityCfg.Key))
		} else {
			INFO.Printf("Starting IoT Discovery on port %d for HTTP requests\n", agent.Port)
			panic(http.ListenAndServe(":"+strconv.Itoa(agent.Port), api.MakeHandler()))
		}
	}()
}

func (agent *NGSIAgent) handleNotifyContext(w rest.ResponseWriter, r *rest.Request) {
	notifyCtxReq := NotifyContextRequest{}
	err := r.DecodeJsonPayload(&notifyCtxReq)
	if err != nil {
		INFO.Println(err)
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if agent.CtxNotifyHandler != nil {
		agent.CtxNotifyHandler(&notifyCtxReq)
	}

	// send out the response
	notifyCtxResp := NotifyContextResponse{}
	w.WriteJson(&notifyCtxResp)
}

func (agent *NGSIAgent) handleNotifyContextAvailability(w rest.ResponseWriter, r *rest.Request) {
	notifyCtxAvailReq := NotifyContextAvailabilityRequest{}
	err := r.DecodeJsonPayload(&notifyCtxAvailReq)
	if err != nil {
		ERROR.Println(err)
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if agent.CtxAvailbilityNotifyHandler != nil {
		agent.CtxAvailbilityNotifyHandler(&notifyCtxAvailReq)
	}

	// send out the response
	notifyCtxAvailResp := NotifyContextAvailabilityResponse{}
	notifyCtxAvailResp.ResponseCode.Code = 200
	notifyCtxAvailResp.ResponseCode.ReasonPhrase = "OK"

	w.WriteJson(&notifyCtxAvailResp)
}

func (agent *NGSIAgent) SetContextNotifyHandler(cb NotifyContextFunc) {
	agent.CtxNotifyHandler = cb
}

func (agent *NGSIAgent) SetContextAvailabilityNotifyHandler(cb NotifyContextAvailabilityFunc) {
	agent.CtxAvailbilityNotifyHandler = cb
}
