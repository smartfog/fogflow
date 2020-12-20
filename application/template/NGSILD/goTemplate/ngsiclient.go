package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func queryContext(id string, IoTBrokerURL string) (map[string]interface{}, error) {
	req, _ := http.NewRequest("GET", IoTBrokerURL+"/ngsi-ld/v1/entities/"+id, nil)
	req.Header.Add("Content-Type", "application/ld+json")
	req.Header.Add("Accept", "application/ld+json")
	res, err := http.DefaultClient.Do(req)
	fmt.Println(res)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var reqData interface{}
	err = json.Unmarshal(data, &reqData)
	if err != nil {
		return nil, err
	}
	itemsMap := reqData.(map[string]interface{})
	fmt.Println(itemsMap)
	res.Body.Close()
	return itemsMap, err
}

func UpdateLdContext(updateCtx map[string]interface{}, IoTBrokerURL string) error {
	body, err := json.Marshal(updateCtx)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", IoTBrokerURL+"/ngsi-ld/v1/entities/", bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/ld+json")
	req.Header.Add("Link", "<{{link}}>; rel=\"https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld\"; type=\"application/ld+json\"")
	res, err := http.DefaultClient.Do(req)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
