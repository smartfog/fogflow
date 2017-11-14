package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	. "fogflow/common/ngsi"
)

func postNotifyContext(ctxElems []ContextElement, subscriptionId string, URL string, IsOrionBroker bool) error {
	//INFO.Println("NOTIFY: ", URL)
	elementRespList := make([]ContextElementResponse, 0)

	if IsOrionBroker == true {
		// reset the subscriptionId due to the limited length in Orion Broker
		subscriptionId = ""
	}

	for _, elem := range ctxElems {
		if IsOrionBroker == true {
			// convert it to orion-compatible format
			elem.ID = elem.Entity.ID
			elem.Type = elem.Entity.Type
			if elem.Entity.IsPattern == true {
				elem.IsPattern = "true"
			} else {
				elem.IsPattern = "false"
			}
		}

		elementResponse := ContextElementResponse{}
		elementResponse.ContextElement = elem
		elementResponse.StatusCode.Code = 200
		elementResponse.StatusCode.ReasonPhrase = "OK"

		elementRespList = append(elementRespList, elementResponse)
	}

	notifyCtxReq := &NotifyContextRequest{
		SubscriptionId:   subscriptionId,
		ContextResponses: elementRespList,
	}

	body, err := json.Marshal(notifyCtxReq)
	if err != nil {
		return err
	}

	DEBUG.Println(string(body))

	req, err := http.NewRequest("POST", URL+"/notifyContext", bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func subscribeContextProvider(sub *SubscribeContextRequest, ProviderURL string) (string, error) {
	body, err := json.Marshal(*sub)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", ProviderURL+"/subscribeContext", bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", "lightweight-iot-broker")
	req.Header.Add("Require-Reliability", "true")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

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

func unsubscribeContextProvider(sid string, ProviderURL string) error {
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

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

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
