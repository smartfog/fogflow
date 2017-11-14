package ngsi

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/ant0ine/go-json-rest/rest"
)

type NotifyContextFunc func(notifyCtxReq *NotifyContextRequest)
type NotifyContextAvailabilityFunc func(notifyCtxAvailReq *NotifyContextAvailabilityRequest)

type NGSIAgent struct {
	Port                        int
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
		log.Fatal(err)
	}

	api := rest.NewApi()
	api.Use(rest.DefaultCommonStack...)
	api.SetApp(router)

	go func() {
		fmt.Printf("Starting IoT agent, listening on %d\n", agent.Port)
		log.Fatal(http.ListenAndServe(":"+strconv.Itoa(agent.Port), api.MakeHandler()))
	}()
}

func (agent *NGSIAgent) handleNotifyContext(w rest.ResponseWriter, r *rest.Request) {
	notifyCtxReq := NotifyContextRequest{}
	err := r.DecodeJsonPayload(&notifyCtxReq)
	if err != nil {
		fmt.Println(err)
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
		fmt.Println(err)
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
