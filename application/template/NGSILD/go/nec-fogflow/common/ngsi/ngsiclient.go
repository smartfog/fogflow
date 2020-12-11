package ngsi

import (
        "encoding/json"
        "io/ioutil"
        "net/http"
        "fmt"
)

type NGSILdClient struct {
        IoTBrokerURL string
}

func (ld *NGSILdClient)QueryContext(id string) (map[string]interface{},error){
        req, _ := http.NewRequest("GET", ld.IoTBrokerURL+"/ngsi-ld/v1/entities/" + id, nil)
        req.Header.Add("Content-Type", "application/ld+json")
        req.Header.Add("Accept", "application/ld+json")
        res, err := http.DefaultClient.Do( req )
        fmt.Println(res)
        if err != nil {
                return nil,err
        }
        data, err := ioutil.ReadAll(res.Body)
        if err != nil {
                return nil, err
        }
        var reqData interface{}
        err = json.Unmarshal(data, &reqData)
        if err != nil {
                return nil ,err
        }
        itemsMap := reqData.(map[string]interface{})
        fmt.Println(itemsMap)
        res.Body.Close()
        return itemsMap,err
}

