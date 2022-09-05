package ngsi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type NGSI10Client struct {
	IoTBrokerURL string
	SecurityCfg  *HTTPS
}

func CtxElement2Object(ctxElem *ContextElement) *ContextObject {
	ctxObj := ContextObject{}
	ctxObj.Entity = ctxElem.Entity

	ctxObj.Attributes = make(map[string]ValueObject)
	for _, attr := range ctxElem.Attributes {
		ctxObj.Attributes[attr.Name] = ValueObject{Type: attr.Type, Value: attr.Value}
	}

	ctxObj.Metadata = make(map[string]ValueObject)
	for _, meta := range ctxElem.Metadata {
		ctxObj.Metadata[meta.Name] = ValueObject{Type: meta.Type, Value: meta.Value}
	}

	return &ctxObj
}

func Object2CtxElement(ctxObj *ContextObject) *ContextElement {
	ctxElement := ContextElement{}

	ctxElement.Entity = ctxObj.Entity

	ctxElement.Attributes = make([]ContextAttribute, 0)
	for name, attr := range ctxObj.Attributes {
		ctxAttr := ContextAttribute{Name: name, Type: attr.Type, Value: attr.Value}
		ctxElement.Attributes = append(ctxElement.Attributes, ctxAttr)
	}

	ctxElement.Metadata = make([]ContextMetadata, 0)
	for name, meta := range ctxObj.Metadata {
		ctxMeta := ContextMetadata{Name: name, Type: meta.Type, Value: meta.Value}
		ctxElement.Metadata = append(ctxElement.Metadata, ctxMeta)
	}

	return &ctxElement
}

func (nc *NGSI10Client) UpdateContextObject(ctxObj *ContextObject) error {
	elem := Object2CtxElement(ctxObj)
	return nc.UpdateContext(elem)
}

func (nc *NGSI10Client) UpdateContext(elem *ContextElement) error {
	return nc.sendUpdateContext(elem, false, false)
}

func (nc *NGSI10Client) InternalUpdateContext(elem *ContextElement) error {
	return nc.sendUpdateContext(elem, true, false)
}

func (nc *NGSI10Client) SouthboundUpdateContext(elem *ContextElement, fs string, fsp string) error {
	return nc.sendUpdateContext(elem, false, true, fs, fsp)
}

func (nc *NGSI10Client) sendUpdateContext(elem *ContextElement, internal bool, southbound bool, params ...string) error {
	updateCtxReq := &UpdateContextRequest{
		ContextElements: []ContextElement{*elem},
		UpdateAction:    "UPDATE",
	}

	body, err := json.Marshal(updateCtxReq)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", nc.IoTBrokerURL+"/updateContext", bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	if internal == true {
		req.Header.Add("User-Agent", "lightweight-iot-broker")
	}

	if southbound == true {
		req.Header.Add("fiware-service", params[0])
		req.Header.Add("fiware-servicepath", params[1])
	}

	client := nc.SecurityCfg.GetHTTPClient()
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		ERROR.Println(err)
		return err
	}

	text, _ := ioutil.ReadAll(resp.Body)

	updateCtxResp := UpdateContextResponse{}
	err = json.Unmarshal(text, &updateCtxResp)
	if err != nil {
		return err
	}

	return nil
	// if updateCtxResp.ErrorCode.Code == 200 {
	// 	return nil
	// } else {
	// 	err = errors.New(updateCtxResp.ErrorCode.ReasonPhrase)
	// 	return err
	// }
}

func (nc *NGSI10Client) DeleteContext(eid *EntityId) error {
	return nc.sendDeleteContext(eid, false)
}

func (nc *NGSI10Client) InternalDeleteContext(eid *EntityId) error {
	return nc.sendDeleteContext(eid, true)
}

func (nc *NGSI10Client) sendDeleteContext(eid *EntityId, internal bool) error {
	element := ContextElement{}

	entity := EntityId{}
	entity.ID = eid.ID
	entity.Type = eid.Type
	entity.IsPattern = eid.IsPattern

	element.Entity = entity

	updateCtxReq := &UpdateContextRequest{
		ContextElements: []ContextElement{element},
		UpdateAction:    "DELETE",
	}

	body, err := json.Marshal(updateCtxReq)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", nc.IoTBrokerURL+"/updateContext", bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	if internal == true {
		req.Header.Add("User-Agent", "lightweight-iot-broker")
	}

	client := nc.SecurityCfg.GetHTTPClient()
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		ERROR.Println(err)
		return err
	}

	text, _ := ioutil.ReadAll(resp.Body)

	updateCtxResp := UpdateContextResponse{}
	err = json.Unmarshal(text, &updateCtxResp)
	if err != nil {
		return err
	}

	return nil
	// if updateCtxResp.ErrorCode.Code == 200 {
	// 	return nil
	// } else {
	// 	err = errors.New(updateCtxResp.ErrorCode.ReasonPhrase)
	// 	return err
	// }
}

func (nc *NGSI10Client) NotifyContext(elem *ContextElement) error {
	elementResponse := ContextElementResponse{}
	elementResponse.ContextElement = *elem
	elementResponse.StatusCode.Code = 200
	elementResponse.StatusCode.ReasonPhrase = "OK"

	notifyCtxReq := &NotifyContextRequest{
		ContextResponses: []ContextElementResponse{elementResponse},
	}

	body, err := json.Marshal(notifyCtxReq)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", nc.IoTBrokerURL+"/notifyContext", bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	client := nc.SecurityCfg.GetHTTPClient()
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		ERROR.Println(err)
		return err
	}

	text, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(text))

	notifyCtxResp := NotifyContextResponse{}
	err = json.Unmarshal(text, &notifyCtxResp)
	if err != nil {
		return err
	}

	if notifyCtxResp.ResponseCode.Code == 200 {
		return nil
	} else {
		err = errors.New(notifyCtxResp.ResponseCode.ReasonPhrase)
		return err
	}
}

func (nc *NGSI10Client) GetEntity(id string) (*ContextObject, error) {
	req, err := http.NewRequest("GET", nc.IoTBrokerURL+"/entity/"+id, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	client := nc.SecurityCfg.GetHTTPClient()
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		ERROR.Println(err)
		return nil, err
	}

	text, _ := ioutil.ReadAll(resp.Body)

	ctxElement := ContextElement{}
	err = json.Unmarshal(text, &ctxElement)
	if err != nil {
		return nil, err
	}

	ctxObj := CtxElement2Object(&ctxElement)

	return ctxObj, nil
}

func (nc *NGSI10Client) QueryContext(query *QueryContextRequest) ([]*ContextObject, error) {
	body, err := json.Marshal(*query)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", nc.IoTBrokerURL+"/queryContext", bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	client := nc.SecurityCfg.GetHTTPClient()
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		ERROR.Println(err)
		return nil, err
	}

	text, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(text))

	queryCtxResp := QueryContextResponse{}
	err = json.Unmarshal(text, &queryCtxResp)
	if err != nil {
		return nil, err
	}

	ctxObjectList := make([]*ContextObject, 0)
	for _, contextElementResponse := range queryCtxResp.ContextResponses {
		ctxObj := CtxElement2Object(&contextElementResponse.ContextElement)
		ctxObjectList = append(ctxObjectList, ctxObj)
	}

	return ctxObjectList, nil
}

func (nc *NGSI10Client) InternalQueryContext(query *QueryContextRequest) ([]ContextElement, error) {
	body, err := json.Marshal(*query)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", nc.IoTBrokerURL+"/queryContext", bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", "lightweight-iot-broker")

	client := nc.SecurityCfg.GetHTTPClient()
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		ERROR.Println(err)
		return nil, err
	}

	text, _ := ioutil.ReadAll(resp.Body)

	queryCtxResp := QueryContextResponse{}
	err = json.Unmarshal(text, &queryCtxResp)
	if err != nil {
		return nil, err
	}

	ctxElements := make([]ContextElement, 0)
	for _, contextElementResponse := range queryCtxResp.ContextResponses {
		ctxElements = append(ctxElements, contextElementResponse.ContextElement)
	}

	return ctxElements, nil
}

func (nc *NGSI10Client) SubscribeContext(sub *SubscribeContextRequest, correlatorID string, requireReliability bool) (string, error) {
	body, err := json.Marshal(*sub)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", nc.IoTBrokerURL+"/subscribeContext", bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Fiware-Correlator", correlatorID)

	if requireReliability == true {
		req.Header.Add("Require-Reliability", "true")
	}

	client := nc.SecurityCfg.GetHTTPClient()
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		ERROR.Println(err)
		return "", err
	}

	text, _ := ioutil.ReadAll(resp.Body)

	subscribeCtxResp := SubscribeContextResponse{}
	err = json.Unmarshal(text, &subscribeCtxResp)
	if err != nil {
		return "", err
	}

	if subscribeCtxResp.SubscribeResponse.SubscriptionId != "" {
		return subscribeCtxResp.SubscribeResponse.SubscriptionId, nil
	} else {
		err = errors.New(subscribeCtxResp.SubscribeError.ErrorCode.ReasonPhrase)
		return "", err
	}
}

func (nc *NGSI10Client) UnsubscribeContext(sid string) error {
	unsubscription := &UnsubscribeContextRequest{
		SubscriptionId: sid,
	}

	body, err := json.Marshal(unsubscription)
	if err != nil {
		return err
	}

	//fmt.Println(string(body))

	req, err := http.NewRequest("POST", nc.IoTBrokerURL+"/unsubscribeContext", bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	client := nc.SecurityCfg.GetHTTPClient()
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		ERROR.Println(err)
		return err
	}

	text, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(text))

	unsubscribeCtxResp := UnsubscribeContextResponse{}
	err = json.Unmarshal(text, &unsubscribeCtxResp)
	if err != nil {
		return err
	}

	if unsubscribeCtxResp.StatusCode.Code == 200 {
		return nil
	} else {
		err = errors.New(unsubscribeCtxResp.StatusCode.ReasonPhrase)
		return err
	}
}

type NGSI9Client struct {
	IoTDiscoveryURL string
	SecurityCfg     *HTTPS
}

type NGSIV2Client struct {
	IoTDiscoveryURL string
	SecurityCfg     *HTTPS
}

func (nc *NGSI9Client) RegisterContext(registerCtxReq *RegisterContextRequest) (string, error) {
	body, err := json.Marshal(registerCtxReq)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", nc.IoTDiscoveryURL+"/registerContext", bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	client := nc.SecurityCfg.GetHTTPClient()
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		ERROR.Println(err)
		return "", err
	}

	text, _ := ioutil.ReadAll(resp.Body)
	registerCtxResp := RegisterContextResponse{}
	err = json.Unmarshal(text, &registerCtxResp)
	if err != nil {
		return "", err
	}

	if registerCtxResp.ErrorCode.Code == 200 {
		return registerCtxResp.RegistrationId, nil
	} else {
		err = errors.New(registerCtxResp.ErrorCode.ReasonPhrase)
		return "", err
	}
}

func (nc *NGSI9Client) UnregisterEntity(eid string) error {
	req, err := http.NewRequest("DELETE", nc.IoTDiscoveryURL+"/registration/"+eid, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	client := nc.SecurityCfg.GetHTTPClient()
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		ERROR.Println(err)
		return err
	}

	return nil
}

func (nc *NGSI9Client) GetProviderURL(id string) (string, *EntityRegistration) {
	//resp, err := http.Get(nc.IoTDiscoveryURL + "/registration/" + id)
	//defer resp.Body.Close()

	req, err := http.NewRequest("GET", nc.IoTDiscoveryURL+"/registration/"+id, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	client := nc.SecurityCfg.GetHTTPClient()
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	registeredEntity := EntityRegistration{}
	if err != nil {
		ERROR.Println(err)
		return "", nil
	}

	text, _ := ioutil.ReadAll(resp.Body)

	if text == nil {
		return "", nil
	}
	err = json.Unmarshal(text, &registeredEntity)
	registration := ContextRegistration{}
	err = json.Unmarshal(text, &registration)
	if err != nil {
		return "", nil
	} else {
		return registration.ProvidingApplication, &registeredEntity
	}
}

func (nc *NGSI9Client) QuerySiteList(geoscope OperationScope) ([]SiteInfo, error) {
	body, err := json.Marshal(geoscope)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", nc.IoTDiscoveryURL+"/querysite", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	client := nc.SecurityCfg.GetHTTPClient()
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	text, _ := ioutil.ReadAll(resp.Body)

	siteList := make([]SiteInfo, 0)
	err = json.Unmarshal(text, &siteList)
	if err != nil {
		return nil, err
	} else {
		return siteList, nil
	}
}

func (nc *NGSI9Client) DiscoverContextAvailability(discoverCtxAvailabilityReq *DiscoverContextAvailabilityRequest) ([]ContextRegistration, error) {
	body, err := json.Marshal(discoverCtxAvailabilityReq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", nc.IoTDiscoveryURL+"/discoverContextAvailability", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	client := nc.SecurityCfg.GetHTTPClient()
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	text, _ := ioutil.ReadAll(resp.Body)

	discoverCtxAvailResp := DiscoverContextAvailabilityResponse{}
	err = json.Unmarshal(text, &discoverCtxAvailResp)
	if err != nil {
		return nil, err
	}

	registrationList := make([]ContextRegistration, 0)
	for _, registration := range discoverCtxAvailResp.ContextRegistrationResponses {
		registrationList = append(registrationList, registration.ContextRegistration)
	}

	return registrationList, nil
}

func (nc *NGSI9Client) SubscribeContextAvailability(sub *SubscribeContextAvailabilityRequest) (string, error) {
	body, err := json.Marshal(*sub)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", nc.IoTDiscoveryURL+"/subscribeContextAvailability", bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	client := nc.SecurityCfg.GetHTTPClient()
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		ERROR.Println(err)
		return "", err
	}

	text, _ := ioutil.ReadAll(resp.Body)

	subscribeCtxAvailResp := SubscribeContextAvailabilityResponse{}
	err = json.Unmarshal(text, &subscribeCtxAvailResp)
	if err != nil {
		return "", err
	}

	if subscribeCtxAvailResp.SubscriptionId != "" {
		return subscribeCtxAvailResp.SubscriptionId, nil
	} else {
		err = errors.New(subscribeCtxAvailResp.ErrorCode.ReasonPhrase)
		return "", err
	}
}

func (nc *NGSIV2Client) Subscribev2ContextAvailability(sub *SubscribeContextAvailabilityRequest) (string, error) {
	body, err := json.Marshal(*sub)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", nc.IoTDiscoveryURL+"/subscribeContextAvailability", bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	client := nc.SecurityCfg.GetHTTPClient()
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		ERROR.Println(err)
		return "", err
	}

	text, _ := ioutil.ReadAll(resp.Body)

	subscribeCtxAvailResp := Subscribev2ContextAvailabilityResponse{}
	err = json.Unmarshal(text, &subscribeCtxAvailResp)
	if err != nil {
		return "", err
	}

	if subscribeCtxAvailResp.SubscriptionId != "" {
		return subscribeCtxAvailResp.SubscriptionId, nil
	} else {
		err = errors.New(subscribeCtxAvailResp.ErrorCode.ReasonPhrase)
		return "", err
	}
}
func (nc *NGSI9Client) UnsubscribeContextAvailability(sid string) error {
	unsubscription := &UnsubscribeContextAvailabilityRequest{
		SubscriptionId: sid,
	}

	body, err := json.Marshal(unsubscription)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", nc.IoTDiscoveryURL+"/unsubscribeContextAvailability", bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	client := nc.SecurityCfg.GetHTTPClient()
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		ERROR.Println(err)
		return err
	}

	text, _ := ioutil.ReadAll(resp.Body)

	unsubscribeCtxAvailResp := UnsubscribeContextAvailabilityResponse{}
	err = json.Unmarshal(text, &unsubscribeCtxAvailResp)
	if err != nil {
		return err
	}

	if unsubscribeCtxAvailResp.StatusCode.Code == 200 {
		return nil
	} else {
		err = errors.New(unsubscribeCtxAvailResp.StatusCode.ReasonPhrase)
		return err
	}
}

func (nc *NGSIV2Client) Unsubscribev2ContextAvailability(sid string) error {
	unsubscription := &Unsubscribev2ContextAvailabilityRequest{
		SubscriptionId: sid,
	}

	body, err := json.Marshal(unsubscription)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", nc.IoTDiscoveryURL+"/unsubscribeContextAvailability", bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	client := nc.SecurityCfg.GetHTTPClient()
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		ERROR.Println(err)
		return err
	}

	text, _ := ioutil.ReadAll(resp.Body)

	unsubscribeCtxAvailResp := UnsubscribeContextAvailabilityResponse{}
	err = json.Unmarshal(text, &unsubscribeCtxAvailResp)
	if err != nil {
		return err
	}

	if unsubscribeCtxAvailResp.StatusCode.Code == 200 {
		return nil
	} else {
		err = errors.New(unsubscribeCtxAvailResp.StatusCode.ReasonPhrase)
		return err
	}
}

func (nc *NGSI9Client) DiscoveryNearbyIoTBroker(nearby NearBy) (string, error) {
	discoverReq := DiscoverContextAvailabilityRequest{}

	entity := EntityId{}
	entity.Type = "IoTBroker"
	entity.IsPattern = true
	discoverReq.Entities = make([]EntityId, 0)

	discoverReq.Entities = append(discoverReq.Entities, entity)

	scope := OperationScope{}
	scope.Type = "nearby"
	scope.Value = nearby

	discoverReq.Restriction.Scopes = make([]OperationScope, 0)
	discoverReq.Restriction.Scopes = append(discoverReq.Restriction.Scopes, scope)

	registerationList, err := nc.DiscoverContextAvailability(&discoverReq)

	if err != nil {
		return "", err
	}

	if registerationList == nil {
		return "", nil
	} else {
		for _, reg := range registerationList {
			return reg.ProvidingApplication, nil
		}
	}
	return "", nil
}

func (nc *NGSI9Client) SendHeartBeat(brokerProfile *BrokerProfile) error {
	body, err := json.Marshal(brokerProfile)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", nc.IoTDiscoveryURL+"/broker/heartbeat", bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	client := nc.SecurityCfg.GetHTTPClient()
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}

	return err
}

// NGSI-LD feature addition

func (nc *NGSI10Client) CreateLDEntityOnRemote(elem map[string]interface{}, link string) error {
	/*body, err := json.Marshal(elem)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", nc.IoTBrokerURL+"/ngsi-ld/v1/entities/", bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/ld+json")

	if link != "" {
		req.Header.Add("Link", link)
	}

	client := nc.SecurityCfg.GetHTTPClient()
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		ERROR.Println(err)
		return err
	}
	return nil*/

	id := elem["id"].(string)
	idSplit := strings.Split(id, "@")
	elem["id"] = idSplit[0]
	fiwareService := idSplit[1]
	fiwareServicePath := elem["fiwareServicePath"].(string)
	delete(elem, "fiwareServicePath")
	var upsertBody [1]interface{}
	fmt.Println("Ele", elem)
	upsertBody[0] = elem
	fmt.Println("upsertBody", upsertBody)
	body, err := json.Marshal(upsertBody)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", nc.IoTBrokerURL+"/ngsi-ld/v1/entityOperations/upsert", bytes.NewBuffer(body))
	req.Header.Add("Accept", "application/ld+json")

	if link != "" {
		req.Header.Add("Link", link)
		req.Header.Add("Content-Type", "application/json")
	} else {
		req.Header.Add("Content-Type", "application/ld+json")
	}
	if fiwareService != "default" {
		req.Header.Add("fiware-service", fiwareService)
	}
	if fiwareServicePath != "default" {
		req.Header.Add("fiware-servicepath", fiwareServicePath)
	}
	client := nc.SecurityCfg.GetHTTPClient()
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		ERROR.Println(err)
		return err
	}
	return nil

}

func (nc *NGSI10Client) AppendLDEntityOnRemote(elem map[string]interface{}, eid string) error {
	body, err := json.Marshal(elem)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", nc.IoTBrokerURL+"/ngsi-ld/v1/entities/"+eid+"/attrs", bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")

	client := nc.SecurityCfg.GetHTTPClient()
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		ERROR.Println(err)
		return err
	}
	return nil
}

func (nc *NGSI10Client) UpdateLDEntityAttributeOnRemote(elem map[string]interface{}, eid string) (error, int) {
	body, err := json.Marshal(elem)
	if err != nil {
		return err, 0
	}
	req, err := http.NewRequest("PATCH", nc.IoTBrokerURL+"/ngsi-ld/v1/entities/"+eid+"/attrs", bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")

	client := nc.SecurityCfg.GetHTTPClient()
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		ERROR.Println(err)
		return err, 0
	}
	if resp.StatusCode == 207 {
		return err, resp.StatusCode
	}
	return nil, 0
}

func (nc *NGSI10Client) UpdateLDEntityspecificAttributeOnRemote(elem map[string]interface{}, eid string, attribute string) (error, int) {
	body, err := json.Marshal(elem)
	if err != nil {
		return err, 0
	}
	req, err := http.NewRequest("PATCH", nc.IoTBrokerURL+"/ngsi-ld/v1/entities/"+eid+"/attrs/"+attribute, bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")

	client := nc.SecurityCfg.GetHTTPClient()
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		ERROR.Println(err)
		return err, 0
	}
	if resp.StatusCode == 404 {
		return err, resp.StatusCode
	}
	return nil, 0
}

func (nc *NGSI10Client) SubscribeLdContext(sub *LDSubscriptionRequest, requireReliability bool, fs string, fsp string) (string, error) {
	body, err := json.Marshal(*sub)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", nc.IoTBrokerURL+"/ngsi-ld/v1/subscriptions/", bytes.NewBuffer(body))
	if fs != "" {
		req.Header.Add("fiware-service", fs)
	}
	if fsp != "" {
		req.Header.Add("fiware-servicepath", fsp)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/ld+json")
	req.Header.Add("Link", "<https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context-v1.3.jsonld>; rel=\"http://www.w3.org/ns/json-ld#context\"; type=\"application/ld+json\"")
	if requireReliability == true {
		req.Header.Add("Require-Reliability", "true")
	}

	client := nc.SecurityCfg.GetHTTPClient()
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		ERROR.Println(err)
		return "", err
	}

	text, _ := ioutil.ReadAll(resp.Body)

	subscribeCtxResp := SubscribeContextResponse{}
	err = json.Unmarshal(text, &subscribeCtxResp)
	if err != nil {
		return "", err
	}

	if subscribeCtxResp.SubscribeResponse.SubscriptionId != "" {
		return subscribeCtxResp.SubscribeResponse.SubscriptionId, nil
	} else {
		err = errors.New(subscribeCtxResp.SubscribeError.ErrorCode.ReasonPhrase)
		return "", err
	}
}

//Query for NGSILD entity with entityId

func (nc *NGSI10Client) QueryForNGSILDEntity(eid string) int {
	req, _ := http.NewRequest("GET", nc.IoTBrokerURL+"/ngsi-ld/v1/entities/"+eid, nil)
	req.Header.Add("Content-Type", "application/ld+ijson")
	req.Header.Add("Accept", "application/ld+json")
	client := nc.SecurityCfg.GetHTTPClient()
	resp, _ := client.Do(req)
	return resp.StatusCode

}

// Query for NGSIV1 entity with entityId

func (nc *NGSI10Client) QueryForNGSIV1Entity(eid string) int {
	req, _ := http.NewRequest("GET", nc.IoTBrokerURL+"/entity/"+eid, nil)
	client := nc.SecurityCfg.GetHTTPClient()
	resp, _ := client.Do(req)
	return resp.StatusCode

}

// client to update subscribe Context availbility on discovery
func (nc *NGSI9Client) UpdateLDContextAvailability(sub *SubscribeContextAvailabilityRequest, sid string) (string, error) {
	body, err := json.Marshal(*sub)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", nc.IoTDiscoveryURL+"/UpdateLDContextAvailability"+"/"+sid, bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	client := nc.SecurityCfg.GetHTTPClient()
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		ERROR.Println(err)
		return "", err
	}

	text, _ := ioutil.ReadAll(resp.Body)

	subscribeCtxAvailResp := SubscribeContextAvailabilityResponse{}
	err = json.Unmarshal(text, &subscribeCtxAvailResp)
	if err != nil {
		return "", err
	}

	if subscribeCtxAvailResp.SubscriptionId != "" {
		return subscribeCtxAvailResp.SubscriptionId, nil
	} else {
		err = errors.New(subscribeCtxAvailResp.ErrorCode.ReasonPhrase)
		return "", err
	}
}

// LD QueryContext

func (nc *NGSI10Client) InternalLDQueryContext(query *LDQueryContextRequest, fs string, fsp string) ([]interface{}, error) {
	/*for index, entity := *query.Entities {
		id := entity.ID
		idSplit := strings.Split(id, "@")
		entity.ID = idSplit[0]
		query.Entities[index] = entity
	}
	fmt.Println(query)*/
	body, err := json.Marshal(*query)
	if err != nil {
		return nil, err
	}
	BrokerURL := strings.TrimSuffix(nc.IoTBrokerURL, "/ngsi10")
	req, err := http.NewRequest("POST", BrokerURL+"/ngsi-ld/v1/entityOperations/query", bytes.NewBuffer(body))
	if fs != "default" {
		req.Header.Add("fiware-service", fs)
	}
	if fsp != "" {
		req.Header.Add("fiware-servicepath", fsp)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", "lightweight-iot-broker")

	client := nc.SecurityCfg.GetHTTPClient()
	resp, err := client.Do(req)
	fmt.Println("resp", resp)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		ERROR.Println(err)
		return nil, err
	}

	text, _ := ioutil.ReadAll(resp.Body)
	var queryCtxResp []interface{}
	err = json.Unmarshal(text, &queryCtxResp)
	if err != nil {
		return nil, err
	}
	return queryCtxResp, nil
}

func (nc *NGSI10Client) QueryLdContext(query *LDQueryContextRequest, fs, fsp string) ([]interface{}, error) {
	body, err := json.Marshal(*query)
	if err != nil {
		return nil, err
	}
	BrokerURL := strings.TrimSuffix(nc.IoTBrokerURL, "/ngsi10")
	req, err := http.NewRequest("POST", BrokerURL+"/ngsi-ld/v1/entityOperations/query", bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	if fs != "default" {
		req.Header.Add("fiware-service", fs)
	}

	if fsp != "" {
		req.Header.Add("fiware-servicepath", fsp)
	}

	client := nc.SecurityCfg.GetHTTPClient()
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		ERROR.Println(err)
		return nil, err
	}

	text, _ := ioutil.ReadAll(resp.Body)
	var queryCtxResp []interface{}
	err = json.Unmarshal(text, &queryCtxResp)
	if err != nil {
		return nil, err
	}
	fmt.Println("queryCtxResp", queryCtxResp)
	return queryCtxResp, nil

}
