package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/piprate/json-gold/ld"
	. "github.com/smartfog/fogflow/common/constants"
	. "github.com/smartfog/fogflow/common/ngsi"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func postNotifyContext(ctxElems []ContextElement, subscriptionId string, URL string, IsOrionBroker bool, httpsCfg *HTTPS) error {
	//INFO.Println("NOTIFY: ", URL)
	elementRespList := make([]ContextElementResponse, 0)

	if IsOrionBroker == true {
		return postOrionV2NotifyContext(ctxElems, URL, subscriptionId)
	}

	for _, elem := range ctxElems {
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

// send an notify to orion broker based on v2, for the compatability reason
func postOrionV2NotifyContext(ctxElems []ContextElement, URL string, subscriptionId string) error {
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

	notifyCtxReq := &OrionV2NotifyContextRequest{
		SubscriptionId: subscriptionId,
		Entities:       elementList,
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
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return err
	}

	return nil
}

func subscriptionLDContextProvider(sub *LDSubscriptionRequest, ProviderURL string, httpsCfg *HTTPS) (string, error) {
	body, err := json.Marshal(*sub)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", ProviderURL+"/ngsi-ld/v1/subscriptions/", bytes.NewBuffer(body)) // add NGSILD url
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", "lightweight-iot-broker")
	req.Header.Add("Require-Reliability", "true")
	req.Header.Add("Link", DEFAULT_CONTEXT)
	// add link header
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

//NGSIV2
func subscriptionProvider(sub *SubscriptionRequest, ProviderURL string, httpsCfg *HTTPS) (string, error) {
	body, err := json.Marshal(*sub)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", ProviderURL+"/v2/subscriptions", bytes.NewBuffer(body))
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

	subscribeCtxResp := Subscribev2Response{}
	err = json.Unmarshal(text, &subscribeCtxResp)
	if err != nil {
		return "", err
	}

	if subscribeCtxResp.SubscriptionResponse.SubscriptionId != "" {
		return subscribeCtxResp.SubscriptionResponse.SubscriptionId, nil
	} else {
		err = errors.New(subscribeCtxResp.SubscriptionError.ErrorCode.ReasonPhrase)
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

//NGSIV2
func unsubscribev2ContextProvider(sid string, ProviderURL string, httpsCfg *HTTPS) error {
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

// NGSI-LD starts here.

func compactData(entity map[string]interface{}, context interface{}) (interface{}, error) {
	proc := ld.NewJsonLdProcessor()
	options := ld.NewJsonLdOptions("")
	compacted, err := proc.Compact(entity, context, options)
	return compacted, err
}

func ldPostNotifyContext(ldCtxElems []map[string]interface{}, subscriptionId string, URL string, httpsCfg *HTTPS) error {
	INFO.Println("NOTIFY: ", URL)
	ldCompactedElems := make([]map[string]interface{}, 0)
	for k, _ := range ldCtxElems {
		resolved, _ := compactData(ldCtxElems[k], "https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld")
		ldCompactedElems = append(ldCompactedElems, resolved.(map[string]interface{}))
	}
	LdElementList := make([]interface{}, 0)
	for _, ldEle := range ldCompactedElems {
		element := make(map[string]interface{})
		element["id"] = ldEle["id"]
		element["type"] = ldEle["type"]
		for k, _ := range ldEle {
			if k != "id" && k != "type" && k != "modifiedAt" && k != "createdAt" && k != "observationSpace" && k != "operationSpace" && k != "location" && k != "@context" {
				element[k] = ldEle[k]
			}
		}
		LdElementList = append(LdElementList, element)
	}

	notifyCtxReq := &LDNotifyContextRequest{
		SubscriptionId: subscriptionId,
		Data:           LdElementList,
		Type:           "Notification",
		Id:             "fogflow:notification",
		NotifyAt:       time.Now().String(),
	}
	body, err := json.Marshal(notifyCtxReq)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Link", "https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld")

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

func ldCloneWithSelectedAttributes(ldElem map[string]interface{}, selectedAttributes []string) map[string]interface{} {
	if len(selectedAttributes) == 0 {
		return ldElem
	} else {
		preparedCopy := make(map[string]interface{})
		for key, val := range ldElem {
			if key == "id" {
				preparedCopy["id"] = val
			} else if key == "type" {
				preparedCopy["type"] = val
			} else if key == "createdAt" {
				preparedCopy["createdAt"] = val
			} else if key == "modifiedAt" {
				//preparedCopy["modifiedAt"] = val
			} else if key == "location" {
				preparedCopy["location"] = val
			} else if key == "@context" {
				preparedCopy["@context"] = val
			} else {
				// add attribute only if present in selectedAttributes
				for _, requiredAttrName := range selectedAttributes {
					if key == requiredAttrName {
						preparedCopy[key] = val
						break
					}
				}
			}
		}
		return preparedCopy
	}
}

func isNewLdAttribute(name string, currEle map[string]interface{}) bool {
	for attr, _ := range currEle {
		if attr == name {
			return false
		}
	}

	return true
}

func getId(updateCtxEle map[string]interface{}) string {
	var eid string
	if _, ok := updateCtxEle["id"]; ok == true {
		eid = updateCtxEle["id"].(string)
	} else if _, ok := updateCtxEle["@id"]; ok == true {
		eid = updateCtxEle["@id"].(string)
	}
	return eid
}

func hasLdUpdatedMetadata(recCtxEle interface{}, currCtxEle interface{}) bool {
	if recCtxEle == nil && currCtxEle == nil {
		return false
	}
	if currCtxEle == nil {
		return true
	}
	recCtxEleMap := recCtxEle.(map[string]interface{})
	currCtxEleMap := currCtxEle.(map[string]interface{})
	for attr, _ := range recCtxEleMap {
		if attr != "@id" && attr != "id" && attr != "type" && attr != "modifiedAt" && attr != "createdAt" && attr != "observationSpace" && attr != "operationSpace" && attr != "location" && attr != "@context" {
			if isNewLdAttribute(attr, currCtxEleMap) == true {
				return true
			}
		}
	}
	return false
}

/*
   Header validation for upsert api
*/

func contentTypeValidator(cType string) error {
	if cType == "application/x-www-form-urlencoded" {
		err := errors.New("No content type header provided")
		return err
	}
	cTypeInLower := strings.ToLower(cType)
	if cTypeInLower != "application/json" && cTypeInLower != "application/ld + json" {
		err := errors.New("Unsupported content type. Allowed are application/json and application/ld+json.")
		return err
	}
	return nil
}
