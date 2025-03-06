package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"reflect"
	"strings"

	. "fogflow/common/ngsi"
)

func postNotifyContext(ctxElems []ContextElement, subscriptionId string, URL string, ConsumerNGSIVersion string, tenant string, httpsCfg *HTTPS) error {
	//INFO.Println("destination protocol: ", DestinationBrokerType)
	// INFO.Println("ctxElems: ", ctxElems)

	switch ConsumerNGSIVersion {
	case "NGSI-LD":
		return postNGSILDUpsert(ctxElems, URL, tenant)
	case "NGSIv2":
		return postNGSIV2NotifyContext(ctxElems, subscriptionId, URL, tenant)
	default:
		return postNGSIV1NotifyContext(ctxElems, subscriptionId, URL, httpsCfg)
	}
}

// for ngsiv1 consumer
func postNGSIV1NotifyContext(ctxElems []ContextElement, subscriptionId string, URL string, httpsCfg *HTTPS) error {
	// INFO.Println("NGSIv1 NOTIFY: ", URL)

	payload := toNGSIv1Payload(ctxElems)

	notifyCtxReq := &NotifyContextRequest{
		SubscriptionId:   subscriptionId,
		ContextResponses: payload,
	}

	body, err := json.Marshal(notifyCtxReq)
	if err != nil {
		return err
	}

	req, _ := http.NewRequest("POST", URL+"/notifyContext", bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	if strings.HasPrefix(URL, "https") {
		client = httpsCfg.GetHTTPClient()
	}

	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return err
	}

	io.ReadAll(resp.Body)

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

	// INFO.Println(string(body))

	req, _ := http.NewRequest("POST", URL, bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	if strings.HasPrefix(URL, "https") {
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

		// include all attributes
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

		// include all domain metadata
		for _, dmeta := range elem.Metadata {
			attribute := OrionV2Attribute{}
			attribute.Type = dmeta.Type
			attribute.Value = dmeta.Value

			element[dmeta.Name] = attribute
		}

		elementList = append(elementList, element)
	}

	return elementList
}

// for NGSI-LD consumer
func postNGSILDUpsert(ctxElems []ContextElement, URL string, tenant string) error {

	if LoggerIsEnabled(DEBUG) {
		DEBUG.Println("NGSI-LD NOTIFY: ", URL, ctxElems)
	}

	payload := toNGSILDPayload(ctxElems, false)
	if LoggerIsEnabled(DEBUG) {
		DEBUG.Println("NGSI-LD elements: ", payload)
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	if LoggerIsEnabled(DEBUG) {
		DEBUG.Println("NGSI-LD body: ", string(body))
	}

	ngsildConsumerURL := URL + "/ngsi-ld/v1/entityOperations/upsert"
	req, _ := http.NewRequest("POST", ngsildConsumerURL, bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("NGSILD-Tenant", tenant)
	req.Header.Add("Link", "<https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context-v1.3.jsonld>; rel=\"http://www.w3.org/ns/json-ld#context\"; type=\"application/ld+json\"")

	client := &http.Client{}
	if strings.HasPrefix(URL, "https") {
		transCfg := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore expired SSL certificates
		}
		client = &http.Client{Transport: transCfg}
	}

	if LoggerIsEnabled(DEBUG) {
		DEBUG.Println("NGSI-LD req: ", req)
	}

	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		ERROR.Println(err)
		return err
	}
	// Check if status code is 300 or greater
	if resp.StatusCode >= 300 {
		// Read response body
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			INFO.Println("Error reading response body:", err)
		}
		ERROR.Println(
			"Error: Request failed\nStatus Code:", resp.StatusCode,
			"\nDestination URL:", req.URL.String(),
			"\nRequest Body:", string(body),
			"\nResponse Body: ", string(responseBody))
	}

	if LoggerIsEnabled(DEBUG) {
		DEBUG.Println("NGSI-LD resp: ", resp)
	}

	return nil
}

func toNGSILDPayload(ctxElems []ContextElement, addAtContext bool) []map[string]interface{} {
	elementListLD := make([]map[string]interface{}, 0)
	for _, elem := range ctxElems {
		// convert it to NGSI-LD
		elementLD := make(map[string]interface{})

		if strings.HasPrefix(elem.Entity.ID, "urn:") || strings.HasPrefix(elem.Entity.ID, "URN:") {
			elementLD["id"] = elem.Entity.ID
		} else {
			elementLD["id"] = "urn:" + elem.Entity.ID
		}

		elementLD["type"] = elem.Entity.Type

		if addAtContext {
			elementLD["@context"] = NGSILD_CORE_CONTEXT
		}

		// include all attributes from the ngsi v1 entity
		for _, attr := range elem.Attributes {

			// fmt.Println("  attr: ", attr)
			if LoggerIsEnabled(DEBUG) {
				DEBUG.Println("attribute to put into ngsi-ld entity: ", attr)
			}
			if (attr.Name == "") || (attr.Type == "") || (attr.Value == nil) {
				continue
			}

			propertyValue := make(map[string]interface{})

			switch strings.ToLower(attr.Type) {
			case "relationship":
				propertyValue["type"] = "Relationship"
				propertyValue["object"] = attr.Value

			case "datetime":
				datatimeValue := make(map[string]interface{})
				datatimeValue["@type"] = "datetime"
				datatimeValue["@value"] = attr.Value

				propertyValue["type"] = "Property"
				propertyValue["value"] = datatimeValue
			default:
				propertyValue["type"] = "Property"
				propertyValue["value"] = attr.Value
				for _, metadata := range attr.Metadata {
					propertyValue[metadata.Name] = metadata.Value
				}
			}

			if len(attr.Metadata) > 0 {
				for _, metadata := range attr.Metadata {
					propertyValue[metadata.Name] = metadata.Value
				}

			}

			if existingValue, ok := elementLD[attr.Name]; ok {

				if existingSlice, isSlice := existingValue.([]interface{}); isSlice {
					// Append the new value to the existing slice
					elementLD[attr.Name] = append(existingSlice, propertyValue)
					if LoggerIsEnabled(DEBUG) {
						DEBUG.Println("isSlice attribute to put into ngsi-ld entity: ", attr.Name, elementLD[attr.Name], propertyValue)
					}
				} else {
					// tempValue := element[attr.Name]
					// If it's not a slice, create a new slice and append both the old and new values
					// element[attr.Name] = []interface{}{existingValue, propertyValue}
					// element[attr.Name] = append(existingSlice, tempValue)
					// element[attr.Name] = append(existingSlice, propertyValue)
					newSlice := []interface{}{existingValue, propertyValue}
					elementLD[attr.Name] = newSlice
					if LoggerIsEnabled(DEBUG) {
						DEBUG.Println("add new propertyValue: ", attr.Name, propertyValue)
						DEBUG.Println("so now is: ", attr.Name, elementLD[attr.Name])
					}
				}
			} else {
				elementLD[attr.Name] = propertyValue
				if LoggerIsEnabled(DEBUG) {
					DEBUG.Println("attribute to put into ngsi-ld entity: ", attr.Name, propertyValue)
				}
			}
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

			elementLD[meta.Name] = propertyValue
		}

		elementListLD = append(elementListLD, elementLD)
	}

	return elementListLD
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

	req, _ := http.NewRequest("POST", ProviderURL+"/subscribeContext", bytes.NewBuffer(body))
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

	text, _ := io.ReadAll(resp.Body)

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

func subscribeContextProviderNGSILD(sub *SubscribeContextRequest, ProviderURL string, subID string, httpsCfg *HTTPS) (string, error) {

	// u1, err := uuid.NewUUID()
	// if err != nil {
	// 	return "", err
	// }
	// subID := "urn:subscription:" + u1.String()

	subscribeContextProviderNGSILD := sub.ToNGSILD(subID)

	body, err := json.Marshal(subscribeContextProviderNGSILD)
	if LoggerIsEnabled(DEBUG) {
		DEBUG.Println("body ", string(body), "| err ", err)
	}
	if err != nil {
		return "", err
	}

	req, _ := http.NewRequest("POST", ProviderURL+"/ngsi-ld/v1/subscriptions", bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/ld+json")
	req.Header.Add("User-Agent", "lightweight-iot-broker")

	if LoggerIsEnabled(DEBUG) {
		DEBUG.Println("req", req)
	}

	client := httpsCfg.GetHTTPClient()
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return "", err
	}

	return subID, nil

}

func unsubscribeContextProvider(sid string, ProviderURL string, httpsCfg *HTTPS) error {
	unsubscription := &UnsubscribeContextRequest{
		SubscriptionId: sid,
	}

	body, err := json.Marshal(unsubscription)
	if err != nil {
		return err
	}

	req, _ := http.NewRequest("POST", ProviderURL+"/unsubscribeContext", bytes.NewBuffer(body))
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

	text, _ := io.ReadAll(resp.Body)

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

func hasMetadataChange(curMetadata *ContextMetadata, ctxElement *ContextElement) bool {
	for _, meta := range (*ctxElement).Metadata {
		sameName := (meta.Name == curMetadata.Name)
		sameValue := reflect.DeepEqual(meta.Value, curMetadata.Value)

		if LoggerIsEnabled(DEBUG) {
			DEBUG.Println(curMetadata.Name, "Domain metadata name is the same: ", sameName)
			DEBUG.Println(curMetadata.Value, "Domain metadata value is the same: ", sameValue)
		}

		if sameName && sameValue {
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

	// check if there is any new attribute name, no check for the change of attribute values
	for _, attr := range recvElement.Attributes {
		if isNewAttribute(attr.Name, curElement) {
			return true
		}
	}

	// check if there is any update on the domain metadata, including the change of metadata values
	for _, metadata := range recvElement.Metadata {
		if hasMetadataChange(&metadata, curElement) {
			return true
		}
	}

	return false
}

// func updateAttributes(currentAttributes []ContextAttribute, newAttributes []ContextAttribute) {
// 	attributeNames := createSet()
// 	for _, newAttribute := range newAttributes {
// 		addToSet(attributeNames, newAttribute.Name)
// 	}

// 	for i := range currentAttributes {
// 		if setContains(attributeNames, currentAttributes[i].Name) {
// 			// Remove element at index 'i'
// 			currentAttributes = append(currentAttributes[:i], currentAttributes[i+1:]...)
// 			i-- // Adjust index to prevent skipping the next element
// 		}
// 	}

// 	currentAttributes = append(currentAttributes, newAttributes)

// }

func updateAttributes(currentAttributes *[]ContextAttribute, newAttributes []ContextAttribute) {
	attributeNames := createSet()
	for _, newAttribute := range newAttributes {
		addToSet(attributeNames, newAttribute.Name)
	}

	// Create a new slice for filtered attributes
	filteredAttributes := (*currentAttributes)[:0] // Keeps same capacity, avoids new allocation

	for _, attr := range *currentAttributes {
		if !setContains(attributeNames, attr.Name) {
			filteredAttributes = append(filteredAttributes, attr) // Keep only non-matching attributes
		}
	}

	// Append new attributes
	*currentAttributes = append(filteredAttributes, newAttributes...)
}

// func updateAttribute(attr *ContextAttribute, curElement *ContextElement) {
// 	attrNameSeen := false
// 	for i := range (*curElement).Attributes { //we need to keep this for to cleanup all the other instances of the same attribute.Name
// 		pCurAttr := &(*curElement).Attributes[i]
// 		if pCurAttr.Name == attr.Name {
// 			if attrNameSeen {
// 				// Remove element at index 'i'
// 				(*curElement).Attributes = append((*curElement).Attributes[:i], (*curElement).Attributes[i+1:]...)
// 				i-- // Adjust index to prevent skipping the next element
// 			} else {
// 				if LoggerIsEnabled(DEBUG) {
// 					DEBUG.Println("currentAttribute: ", pCurAttr, " to be update with attribute: ", attr)
// 				}
// 				//update the value of existing attribute
// 				pCurAttr.Value = attr.Value

// 				//update the metadata list for the existing attribute
// 				for _, metadata := range attr.Metadata {
// 					updateAttributeMetadata(&metadata, pCurAttr)
// 				}
// 				attrNameSeen = true
// 			}
// 		}
// 	}

// 	// if same attibute name not found: add it as new attribute
// 	if !attrNameSeen {
// 		(*curElement).Attributes = append((*curElement).Attributes, *attr)
// 	}
// }

// func updateAttributeMetadata(metadata *ContextMetadata, attr *ContextAttribute) {
// 	for i := range (*attr).Metadata {
// 		pCurMetadata := &(*attr).Metadata[i]
// 		if pCurMetadata.Name == metadata.Name {
// 			// update the value of existing metadata
// 			pCurMetadata.Value = metadata.Value
// 			return
// 		}
// 	}

// 	// add it as new metadata
// 	(*attr).Metadata = append((*attr).Metadata, *metadata)
// }

// TOFIX implement this function similar to updateAttributes function above
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
