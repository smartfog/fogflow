package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	. "fogflow/common/ngsi"
)

func postNotifyContext(ctxElems []ContextElement, subscriptionId string, URL string, DestinationBrokerType string, tenant string, httpsCfg *HTTPS) error {
	INFO.Println("destionation protocol: ", DestinationBrokerType)

	switch DestinationBrokerType {
	case "NGSI-LD":
		return postNGSILDUpsert(ctxElems, subscriptionId, URL, tenant)
	case "NGSIv2":
		return postNGSIV2NotifyContext(ctxElems, subscriptionId, URL, tenant)
	default:
		return postNGSIV1NotifyContext(ctxElems, subscriptionId, URL, httpsCfg)
	}
}

// for ngsiv1 consumer
func postNGSIV1NotifyContext(ctxElems []ContextElement, subscriptionId string, URL string, httpsCfg *HTTPS) error {
	INFO.Println("NGSIv1 NOTIFY: ", URL)

	payload := toNGSIv1Payload(ctxElems)

	notifyCtxReq := &NotifyContextRequest{
		SubscriptionId:   subscriptionId,
		ContextResponses: payload,
	}

	body, err := json.Marshal(notifyCtxReq)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", URL+"/notifyContext", bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	if strings.HasPrefix(URL, "https") == true {
		client = httpsCfg.GetHTTPClient()
	}

	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return err
	}

	ioutil.ReadAll(resp.Body)

	return nil
}

func toNGSIv1Payload(ctxElems []ContextElement) []ContextElementResponse {
	elementRespList := make([]ContextElementResponse, 0)

	for _, elem := range ctxElems {
		elementResponse := ContextElementResponse{}
		elementResponse.ContextElement = elem
		elementResponse.StatusCode.Code = 200
		elementResponse.StatusCode.ReasonPhrase = "OK"

		elementRespList = append(elementRespList, elementResponse)
	}

	return elementRespList
}

// for NGSIv2 consumer
func postNGSIV2NotifyContext(ctxElems []ContextElement, subscriptionId string, URL string, tenant string) error {
	INFO.Println("NGSIv2 NOTIFY: ", URL)

	payload := toNGSIv2Payload(ctxElems)

	notifyCtxReq := &OrionV2NotifyContextRequest{
		SubscriptionId: subscriptionId,
		Entities:       payload,
	}

	body, err := json.Marshal(notifyCtxReq)
	if err != nil {
		return err
	}

	INFO.Println(string(body))

	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	if strings.HasPrefix(URL, "https") == true {
		transCfg := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore expired SSL certificates
		}
		client = &http.Client{Transport: transCfg}
	}

	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return err
	}

	return nil
}

func toNGSIv2Payload(ctxElems []ContextElement) []map[string]interface{} {
	elementList := make([]map[string]interface{}, 0)
	for _, elem := range ctxElems {
		// convert it to NGSI v2
		element := make(map[string]interface{})

		element["id"] = elem.Entity.ID
		element["type"] = elem.Entity.Type

		for _, attr := range elem.Attributes {
			attribute := OrionV2Attribute{}
			attribute.Type = attr.Type
			attribute.Value = attr.Value

			attribute.Metadata = make(map[string]interface{})
			for _, meta := range attr.Metadata {
				m := OrionV2Metadata{}
				m.Type = meta.Type
				m.Value = meta.Value

				attribute.Metadata[meta.Name] = m
			}

			element[attr.Name] = attribute
		}

		elementList = append(elementList, element)
	}

	return elementList
}

// for NGSI-LD consumer
func postNGSILDUpsert(ctxElems []ContextElement, subscriptionId string, URL string, tenant string) error {
	INFO.Println("NGSI-LD NOTIFY: ", URL)

	payload := toNGSILDPayload(ctxElems)

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	INFO.Println(string(body))

	brokerURL := URL + "/ngsi-ld/v1/entityOperations/upsert"
	req, err := http.NewRequest("POST", brokerURL, bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("NGSILD-Tenant", tenant)
	req.Header.Add("Link", "<https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context-v1.3.jsonld>; rel=\"http://www.w3.org/ns/json-ld#context\"; type=\"application/ld+json\"")

	client := &http.Client{}
	if strings.HasPrefix(URL, "https") == true {
		transCfg := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore expired SSL certificates
		}
		client = &http.Client{Transport: transCfg}
	}

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

func toNGSILDPayload(ctxElems []ContextElement) []map[string]interface{} {
	elementList := make([]map[string]interface{}, 0)
	for _, elem := range ctxElems {
		// convert it to NGSI-LD
		element := make(map[string]interface{})

		element["id"] = "urn:" + elem.Entity.ID
		element["type"] = elem.Entity.Type

		// include all attributes from the ngsi v1 entity
		for _, attr := range elem.Attributes {
			propertyValue := make(map[string]interface{})

			propertyValue["type"] = "Property"

			switch attr.Type {
			case "datetime":
				datatimeValue := make(map[string]interface{})
				datatimeValue["@type"] = "datetime"
				datatimeValue["@value"] = attr.Value

				propertyValue["value"] = datatimeValue
			default:
				propertyValue["value"] = attr.Value
			}

			element[attr.Name] = propertyValue
		}

		// include all domain metadata from the ngsi v1 entity as extra properities
		for _, meta := range elem.Metadata {
			propertyValue := make(map[string]interface{})

			INFO.Println(meta.Type)

			switch meta.Type {
			case "point":
				propertyValue["type"] = "GeoProperty"

				location := meta.Value.(Point)

				pointLocation := make(map[string]interface{})
				pointLocation["type"] = "Point"
				coordinates := [2]interface{}{
					location.Longitude,
					location.Latitude,
				}
				pointLocation["coordinates"] = coordinates
				propertyValue["value"] = pointLocation

			default:
				propertyValue["type"] = "Property"
				propertyValue["value"] = meta.Value
			}

			element[meta.Name] = propertyValue
		}

		elementList = append(elementList, element)
	}

	return elementList
}

type OrionV2NotifyContextRequest struct {
	SubscriptionId string                   `json:"subscriptionId"`
	Entities       []map[string]interface{} `json:"data"`
}

type OrionV2Attribute struct {
	Type     string                 `json:"type"`
	Value    interface{}            `json:"value"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

type OrionV2Metadata struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

func subscribeContextProvider(sub *SubscribeContextRequest, ProviderURL string, httpsCfg *HTTPS) (string, error) {
	body, err := json.Marshal(*sub)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", ProviderURL+"/subscribeContext", bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", "lightweight-iot-broker")
	req.Header.Add("Require-Reliability", "true")

	client := httpsCfg.GetHTTPClient()
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
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

func unsubscribeContextProvider(sid string, ProviderURL string, httpsCfg *HTTPS) error {
	unsubscription := &UnsubscribeContextRequest{
		SubscriptionId: sid,
	}

	body, err := json.Marshal(unsubscription)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", ProviderURL+"/unsubscribeContext", bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", "lightweight-iot-broker")

	client := httpsCfg.GetHTTPClient()
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return err
	}

	text, _ := ioutil.ReadAll(resp.Body)

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

func isNewAttribute(name string, ctxElement *ContextElement) bool {
	for _, attr := range (*ctxElement).Attributes {
		if attr.Name == name {
			return false
		}
	}

	return true
}

func isNewMetadata(name string, ctxElement *ContextElement) bool {
	for _, meta := range (*ctxElement).Metadata {
		if meta.Name == name {
			return false
		}
	}

	return true
}

func hasUpdatedMetadata(recvElement *ContextElement, curElement *ContextElement) bool {
	if recvElement == nil && curElement == nil {
		return false
	}

	if curElement == nil {
		return true
	}

	for _, attr := range recvElement.Attributes {
		if isNewAttribute(attr.Name, curElement) == true {
			return true
		}
	}

	for _, metadata := range recvElement.Metadata {
		if isNewMetadata(metadata.Name, curElement) == true {
			return true
		}
	}

	return false
}

func updateAttribute(attr *ContextAttribute, ctxElement *ContextElement) {
	for i := range (*ctxElement).Attributes {
		pCurAttr := &(*ctxElement).Attributes[i]
		if pCurAttr.Name == attr.Name {
			//update the value of existing attribute
			pCurAttr.Value = attr.Value

			//update the metadata list for the existing attribute
			for _, metadata := range attr.Metadata {
				updateAttributeMetadata(&metadata, pCurAttr)
			}

			return
		}
	}

	// add it as new attribute
	(*ctxElement).Attributes = append((*ctxElement).Attributes, *attr)
}

func updateAttributeMetadata(metadata *ContextMetadata, attr *ContextAttribute) {
	for i := range (*attr).Metadata {
		pCurMetadata := &(*attr).Metadata[i]
		if pCurMetadata.Name == metadata.Name {
			// update the value of existing metadata
			pCurMetadata.Value = metadata.Value
			return
		}
	}

	// add it as new metadata
	(*attr).Metadata = append((*attr).Metadata, *metadata)
}

func updateDomainMetadata(metadata *ContextMetadata, ctxElement *ContextElement) {
	for i := range (*ctxElement).Metadata {
		pCurMetadata := &(*ctxElement).Metadata[i]
		if pCurMetadata.Name == metadata.Name {
			// update the value of existing metadata
			pCurMetadata.Value = metadata.Value
			return
		}
	}

	// add it as new metadata
	(*ctxElement).Metadata = append((*ctxElement).Metadata, *metadata)
}
