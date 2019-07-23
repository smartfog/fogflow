package ngsi

import (
	"encoding/json"
	"testing"
)

func TestContextMetadata(t *testing.T) {
	registerCtxReq := RegisterContextRequest{}

	registerCtxReq.Duration = "10M"
	registerCtxReq.RegistrationId = "0"
	registerCtxReq.ContextRegistrations = make([]ContextRegistration, 0)

	registeration := ContextRegistration{}

	registeration.ProvidingApplication = "http://127.0.0.1:8080/ngsi10"
	registeration.EntityIdList = make([]EntityId, 0)

	eid := EntityId{}
	eid.ID = "001"
	eid.IsPattern = false
	eid.Type = "Test"

	registeration.EntityIdList = append(registeration.EntityIdList, eid)

	registeration.Metadata = make([]ContextMetadata, 0)

	point := Point{}
	point.Latitude = 86.0
	point.Longitude = 30.0

	meta := ContextMetadata{}
	meta.Name = "location"
	meta.Type = "point"
	meta.Value = point

	registeration.Metadata = append(registeration.Metadata, meta)

	meta2 := ContextMetadata{}
	meta2.Name = "layer"
	meta2.Type = "integer"
	meta2.Value = 3

	registeration.Metadata = append(registeration.Metadata, meta2)

	registerCtxReq.ContextRegistrations = append(registerCtxReq.ContextRegistrations, registeration)

	jsonText, _ := json.Marshal(registerCtxReq)
	t.Logf("%+v\n", registerCtxReq)
	t.Log(string(jsonText))

	testObj := RegisterContextRequest{}

	err := json.Unmarshal(jsonText, &testObj)
	if err == nil {
		t.Logf("%+v\n", testObj)
		t.Log(testObj)
	} else {
		t.Fatal(err)
	}
}
