package ngsi

import (
	"encoding/json"
	"strconv"
	"strings"
)

type NotifyConditionType int
type UpdateActionType int

const (
	ONTIMEINTERVAL NotifyConditionType = 1
	ONVALUE
	ONCHANGE
)

const (
	UPDATE UpdateActionType = 1
	APPEND
	DELETE
)

type NearBy struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Limit     int     `json:"limit"`
}

type Point struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Circle struct {
	Latitude  float64 `json:"centerLatitude"`
	Longitude float64 `json:"centerLongitude"`
	Radius    float64 `json:"radius"`
}

type Segment struct {
	NW_Corner string `json:"nw_Corner"`
	SE_Corner string `json:"se_Corner"`
}

type MySegment struct {
	NW_Corner Point
	SE_Corner Point
}

func (sg *Segment) Converter() MySegment {
	var mySegment MySegment

	nw := strings.Split(sg.NW_Corner, ",")
	mySegment.NW_Corner.Latitude, _ = strconv.ParseFloat(nw[0], 64)
	mySegment.NW_Corner.Longitude, _ = strconv.ParseFloat(nw[1], 64)

	se := strings.Split(sg.SE_Corner, ",")
	mySegment.SE_Corner.Latitude, _ = strconv.ParseFloat(se[0], 64)
	mySegment.SE_Corner.Longitude, _ = strconv.ParseFloat(se[1], 64)

	return mySegment
}

type Polygon struct {
	Vertices []Point `json:"vertices"`
}

type ContextMetadata struct {
	Name  string      `json:"name"`
	Type  string      `json:"type,omitempty"`
	Value interface{} `json:"value"`
}

func (metadata *ContextMetadata) UnmarshalJSON(b []byte) error {
	type InternalContextMetadata struct {
		Name  string          `json:"name"`
		Type  string          `json:"type,omitempty"`
		Value json.RawMessage `json:"value"`
	}

	m := InternalContextMetadata{}

	err := json.Unmarshal(b, &m)
	if err == nil {
		(*metadata).Name = m.Name
		(*metadata).Type = m.Type

		switch m.Type {
		case "circle":
			var temp Circle
			if err = json.Unmarshal(m.Value, &temp); err == nil {
				(*metadata).Value = temp
			}

		case "point":
			var temp Point
			if err = json.Unmarshal(m.Value, &temp); err == nil {
				(*metadata).Value = temp
			}

		case "polygon":
			var temp Polygon
			if err = json.Unmarshal(m.Value, &temp); err == nil {
				(*metadata).Value = temp
			}

		case "integer":
			var temp int
			if err = json.Unmarshal(m.Value, &temp); err == nil {
				(*metadata).Value = temp
			}

		case "float":
			var temp float64
			if err = json.Unmarshal(m.Value, &temp); err == nil {
				(*metadata).Value = temp
			}

		case "boolean":
			var temp bool
			if err = json.Unmarshal(m.Value, &temp); err == nil {
				(*metadata).Value = temp
			}

		case "string":
			var temp string
			if err = json.Unmarshal(m.Value, &temp); err == nil {
				(*metadata).Value = temp
			}

		case "object":
			var temp map[string]interface{}
			if err = json.Unmarshal(m.Value, &temp); err == nil {
				(*metadata).Value = temp
			}

		default:
			(*metadata).Value = m.Value
		}
	}

	return err
}

type ContextAttribute struct {
	Name     string            `json:"name"`
	Type     string            `json:"type,omitempty"`
	Value    interface{}       `json:"contextValue"`
	Metadata []ContextMetadata `json:"metadata,omitempty"`
}

type OrionContextAttribute struct {
	Name     string            `json:"name"`
	Type     string            `json:"type,omitempty"`
	Value    interface{}       `json:"value"`
	Metadata []ContextMetadata `json:"metadata,omitempty"`
}

func (pAttr *ContextAttribute) UnmarshalJSON(b []byte) error {
	type InternalAttributeObject struct {
		Name     string            `json:"name"`
		Type     string            `json:"type,omitempty"`
		Value    json.RawMessage   `json:"contextValue"`
		Metadata []ContextMetadata `json:"metadata,omitempty"`
	}

	attr := InternalAttributeObject{}

	// handle the attribute value accordingly
	err := json.Unmarshal(b, &attr)
	if err == nil {
		(*pAttr).Name = attr.Name
		(*pAttr).Type = attr.Type

		switch attr.Type {
		case "integer":
			var temp int64
			if err = json.Unmarshal(attr.Value, &temp); err == nil {
				(*pAttr).Value = temp
			}

		case "float":
			var temp float64
			if err = json.Unmarshal(attr.Value, &temp); err == nil {
				(*pAttr).Value = temp
			}

		case "boolean":
			var temp bool
			if err = json.Unmarshal(attr.Value, &temp); err == nil {
				(*pAttr).Value = temp
			}

		case "string":
			var temp string
			if err = json.Unmarshal(attr.Value, &temp); err == nil {
				(*pAttr).Value = temp
			}

		case "object":
			var temp map[string]interface{}
			if err = json.Unmarshal(attr.Value, &temp); err == nil {
				(*pAttr).Value = temp
			}

		default:
			(*pAttr).Value = attr.Value
		}
	}

	// take the metadatas as well
	(*pAttr).Metadata = attr.Metadata

	return err
}

type EntityId struct {
	ID        string `json:"id"`
	Type      string `json:"type,omitempty"`
	IsPattern bool   `json:"isPattern,omitempty"`
}

type ValueObject struct {
	Type  string      `json:"type,omitempty"`
	Value interface{} `json:"value"`
}

type ContextObject struct {
	Entity              EntityId               `json:"entityId"`
	Attributes          map[string]ValueObject `json:"attributes,omitempty"`
	Metadata            map[string]ValueObject `json:"metadata,omitempty"`
	AttributeDomainName string                 `json:"attributeDomainName,omitempty"`
}

func (ctxObj *ContextObject) IsEmpty() bool {
	if len(ctxObj.Attributes) == 0 && len(ctxObj.Metadata) == 0 {
		return true
	} else {
		return false
	}
}

type ContextElement struct {
	Entity              EntityId           `json:"entityId"`
	ID                  string             `json:"id"`
	Type                string             `json:"type,omitempty"`
	IsPattern           string             `json:"isPattern"`
	AttributeDomainName string             `json:"attributeDomainName,omitempty"`
	Attributes          []ContextAttribute `json:"attributes,omitempty"`
	Metadata            []ContextMetadata  `json:"domainMetadata,omitempty"`
}

func (ce *ContextElement) GetAttribute(name string) *ContextAttribute {
	for _, attr := range ce.Attributes {
		if attr.Name == name {
			return &attr
		}
	}

	return nil
}

func (ce *ContextElement) GetMetadata(name string) *ContextMetadata {
	for _, meta := range ce.Metadata {
		if meta.Name == name {
			return &meta
		}
	}

	return nil
}

func (ce *ContextElement) IsEmpty() bool {
	if len(ce.Attributes) == 0 && len(ce.Metadata) == 0 {
		return true
	} else {
		return false
	}
}

func (ce *ContextElement) Clone(orig *ContextElement) {
	ce.Entity.ID = orig.Entity.ID
	ce.Entity.Type = orig.Entity.Type
	ce.Entity.IsPattern = orig.Entity.IsPattern

	ce.AttributeDomainName = orig.AttributeDomainName
}

type ContextElementOrion struct {
	ID                  string                  `json:"id"`
	Type                string                  `json:"type"`
	IsPattern           string                  `json:"isPattern"`
	AttributeDomainName string                  `json:"attributeDomainName,omitempty"`
	Attributes          []OrionContextAttribute `json:"attributes,omitempty"`
	Metadatas           []ContextMetadata       `json:"metadatas,omitempty"`
}

func (element *ContextElement) MarshalJSON() ([]byte, error) {
	if element.ID != "" || element.Type != "" {
		convertedElement := ContextElementOrion{}

		convertedElement.ID = element.ID
		convertedElement.Type = element.Type
		convertedElement.IsPattern = element.IsPattern

		convertedElement.AttributeDomainName = element.AttributeDomainName

		convertedElement.Attributes = make([]OrionContextAttribute, 0)

		for _, attr := range element.Attributes {
			orionAttr := OrionContextAttribute{}
			orionAttr.Name = attr.Name
			orionAttr.Type = attr.Type
			orionAttr.Value = attr.Value

			orionAttr.Metadata = make([]ContextMetadata, len(attr.Metadata))
			copy(orionAttr.Metadata, attr.Metadata)

			convertedElement.Attributes = append(convertedElement.Attributes, orionAttr)
		}

		/* Orion is not using domain context metadata
		convertedElement.Metadatas = make([]ContextMetadata, len(element.Metadata))
		copy(convertedElement.Metadatas, element.Metadata)

			convertedElement.Metadatas = make([]ContextMetadata, 0)
			for _, meta := range element.Metadata {
				orionMeta := ContextMetadata{}
				orionMeta.Name = meta.Name
				orionMeta.Type = "string"

				bytes, err := json.Marshal(&meta.Value)
				if err == nil {
					orionMeta.Value = string(bytes)
				} else {
					orionMeta.Value = ""
				}

				convertedElement.Metadatas = append(convertedElement.Metadatas, orionMeta)
			} */

		return json.Marshal(&convertedElement)
	} else {
		return json.Marshal(&struct {
			Entity              EntityId           `json:"entityId"`
			AttributeDomainName string             `json:"attributeDomainName,omitempty"`
			Attributes          []ContextAttribute `json:"attributes,omitempty"`
			Metadata            []ContextMetadata  `json:"domainMetadata,omitempty"`
		}{
			Entity:              element.Entity,
			AttributeDomainName: element.AttributeDomainName,
			Attributes:          element.Attributes,
			Metadata:            element.Metadata,
		})
	}
}

type StatusCode struct {
	Code         int    `json:"code"`
	ReasonPhrase string `json:"reasonPhrase,omitempty"`
	Details      string `json:"details,omitempty"`
}

type SubscribeError struct {
	SubscriptionId string     `json:"subscriptionId,omitempty"`
	ErrorCode      StatusCode `json:"errorCode"`
}

type NotifyCondition struct {
	Type        string   `json:"type"`
	CondValues  []string `json:"condValueList,omitempty"`
	Restriction string   `json:"restriction,omitempty"`
}

type OperationScope struct {
	Type  string      `json:"scopeType"`
	Value interface{} `json:"scopeValue"`
}

func (scope *OperationScope) UnmarshalJSON(b []byte) error {
	type InternalOperationScope struct {
		Type  string          `json:"scopeType"`
		Value json.RawMessage `json:"scopeValue"`
	}

	s := InternalOperationScope{}
	err := json.Unmarshal(b, &s)
	if err == nil {
		(*scope).Type = s.Type

		switch s.Type {
		case "simplegeolocation":
			var temp Segment
			if err = json.Unmarshal(s.Value, &temp); err == nil {
				(*scope).Value = temp
			}
		case "circle":
			var temp Circle
			if err = json.Unmarshal(s.Value, &temp); err == nil {
				(*scope).Value = temp
			}
		case "point":
			var temp Point
			if err = json.Unmarshal(s.Value, &temp); err == nil {
				(*scope).Value = temp
			}
		case "polygon":
			var temp Polygon
			if err = json.Unmarshal(s.Value, &temp); err == nil {
				(*scope).Value = temp
			}
		case "nearby":
			var temp NearBy
			if err = json.Unmarshal(s.Value, &temp); err == nil {
				(*scope).Value = temp
			}
		case "stringQuery":
			var temp string
			if err = json.Unmarshal(s.Value, &temp); err == nil {
				(*scope).Value = temp
			}
		default:
			(*scope).Value = s.Value
		}
	}

	return err
}

type Restriction struct {
	AttributeExpression string           `json:"attributeExpression, omitempty"`
	Scopes              []OperationScope `json:"scopes,omitempty"`
}

type SubscribeResponse struct {
	SubscriptionId string `json:"subscriptionId"`
	Duration       string `json:"duration,omitempty"`
	Throttling     string `json:"throttling,omitempty"`
}

type ContextRegistrationAttribute struct {
	Name     string            `json:"name"`
	Type     string            `json:"type,omitempty"`
	IsDomain bool              `json:"isDomain"`
	Metadata []ContextMetadata `json:"metadata,omitempty"`
}

type ContextRegistration struct {
	EntityIdList                  []EntityId                     `json:"entities,omitempty"`
	ContextRegistrationAttributes []ContextRegistrationAttribute `json:"attributes,omitempty"`
	Metadata                      []ContextMetadata              `json:"contextMetadata,omitempty"`
	ProvidingApplication          string                         `json:"providingApplication"`
}

type ContextRegistrationResponse struct {
	ContextRegistration ContextRegistration `json:"contextRegistration,omitempty"`
	ErrorCode           StatusCode          `json:"errorCode,omitempty"`
}

type ContextElementResponse struct {
	ContextElement ContextElement `json:"contextElement"`
	StatusCode     StatusCode     `json:"statusCode"`
}

// NGSI10
type QueryContextRequest struct {
	Entities    []EntityId  `json:"entities"`
	Attributes  []string    `json:"attributes,omitempty"`
	Restriction Restriction `json:"restriction,omitempty"`
}

type QueryContextResponse struct {
	ContextResponses []ContextElementResponse `json:"contextResponses,omitempty"`
	ErrorCode        StatusCode               `json:"errorCode,omitempty"`
}

type Subscriber struct {
	IsOrion            bool
	IsInternal         bool
	RequireReliability bool
	BrokerURL          string
	NotifyCache        []*ContextElement
}

type SubscribeContextRequest struct {
	Entities         []EntityId        `json:"entities"`
	Attributes       []string          `json:"attributes,omitempty"`
	Reference        string            `json:"reference"`
	Duration         string            `json:"duration,omitempty"`
	Restriction      Restriction       `json:"restriction,omitempty"`
	NotifyConditions []NotifyCondition `json:"notifyConditions,omitempty"`
	Throttling       string            `json:"throttling,omitempty"`
	Subscriber       Subscriber
}

type SubscribeContextResponse struct {
	SubscribeResponse SubscribeResponse `json:"subscribeResponse,omitempty"`
	SubscribeError    SubscribeError    `json:"subscribeError,omitempty"`
}

type UpdateContextSubscriptionRequest struct {
	SubscriptionId   string            `json:"subscriptionId"`
	Duration         string            `json:"duration,omitempty"`
	Restriction      Restriction       `json:"restriction,omitempty"`
	NotifyConditions []NotifyCondition `json:"notifyConditions,omitempty"`
	Throttling       string            `json:"throttling,omitempty"`
}

type UpdateContextSubscriptionResponse struct {
	SubscribeResponse `json:"subscribeResponse,omitempty"`
	SubscribeError    `json:"subscribeError,omitempty"`
}

type UnsubscribeContextRequest struct {
	SubscriptionId string `json:"subscriptionId"`
}

type UnsubscribeContextResponse struct {
	SubscriptionId string     `json:"subscriptionId"`
	StatusCode     StatusCode `json:"statusCode"`
}

type NotifyContextRequest struct {
	SubscriptionId   string                   `json:"subscriptionId"`
	Originator       string                   `json:"originator"`
	ContextResponses []ContextElementResponse `json:"contextResponses,omitempty"`
}

type NotifyContextResponse struct {
	ResponseCode StatusCode `json:"responseCode"`
}

type UpdateContextRequest struct {
	ContextElements []ContextElement `json:"contextElements"`
	UpdateAction    string           `json:"updateAction"`
}

type UpdateContextResponse struct {
	ContextResponses []ContextElementResponse `json:"contextResponses"`
	ErrorCode        StatusCode               `json:"errorCode,omitempty"`
}

// NGSI9
type RegisterContextRequest struct {
	ContextRegistrations []ContextRegistration `json:"contextRegistrations,omitempty"`
	Duration             string                `json:"duration,omitempty"`
	RegistrationId       string                `json:"registrationId,omitempty"`
}

type RegisterContextResponse struct {
	Duration       string     `json:"duration,omitempty"`
	RegistrationId string     `json:"registrationId"`
	ErrorCode      StatusCode `json:"errorCode,omitempty"`
}

type DiscoverContextAvailabilityRequest struct {
	Entities    []EntityId  `json:"entities"`
	Attributes  []string    `json:"attributes,omitempty"`
	Restriction Restriction `json:"restriction,omitempty"`
}

type DiscoverContextAvailabilityResponse struct {
	ContextRegistrationResponses []ContextRegistrationResponse `json:"contextRegistrationResponses,omitempty"`
	ErrorCode                    StatusCode                    `json:"errorCode,omitempty"`
}

type SubscribeContextAvailabilityRequest struct {
	Entities       []EntityId  `json:"entities"`
	Attributes     []string    `json:"attributes,omitempty"`
	Reference      string      `json:"reference"`
	Duration       string      `json:"duration,omitempty"`
	Restriction    Restriction `json:"restriction,omitempty"`
	SubscriptionId string      `json:"subscriptionId,omitempty"`
}

type SubscribeContextAvailabilityResponse struct {
	SubscriptionId string     `json:"subscribeId"`
	Duration       string     `json:"duration,omitempty"`
	ErrorCode      StatusCode `json:"errorCode,omitempty"`
}

type UpdateContextAvailabilitySubscriptionRequest struct {
	Entities       []EntityId         `json:"entities"`
	Attributes     []ContextAttribute `json:"attributes,omitempty"`
	Duration       string             `json:"duration,omitempty"`
	Restriction    Restriction        `json:"restriction,omitempty"`
	SubscriptionId string             `json:"subscriptionId,omitempty"`
}

type UpdateContextAvailabilitySubscriptionResponse struct {
	SubscriptionId string     `json:"subscriptionId"`
	Duration       string     `json:"duration,omitempty"`
	ErrorCode      StatusCode `json:"errorCode,omitempty"`
}

type UnsubscribeContextAvailabilityRequest struct {
	SubscriptionId string `json:"subscriptionId"`
}

type UnsubscribeContextAvailabilityResponse struct {
	SubscriptionId string     `json:"subscriptionId"`
	StatusCode     StatusCode `json:"statusCode"`
}

type NotifyContextAvailabilityRequest struct {
	SubscriptionId                  string                        `json:"subscribeId"`
	ContextRegistrationResponseList []ContextRegistrationResponse `json:"contextRegistrationResponses,omitempty"`
	ErrorCode                       StatusCode                    `json:"errorCode,omitempty"`
}

type NotifyContextAvailabilityResponse struct {
	ResponseCode StatusCode `json:"responseCode"`
}
