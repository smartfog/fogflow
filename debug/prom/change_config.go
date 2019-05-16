package main

import (
	"bytes"
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	LISTEN_PORT:=4545
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/config", GetConfig),
		rest.Post("/config", PostConfig),
		rest.Post("/add-target", AddTarget),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)

	panic(http.ListenAndServe(":"+strconv.Itoa(LISTEN_PORT), api.MakeHandler()))

}

type Target struct {
	Address string
	Port string
}

type PromConfig struct {
	content string
}
var oldConfig = PromConfig{}
var configPath="/etc/prometheus/tgroups/target_groups.json"

func GetConfig(w rest.ResponseWriter, r *rest.Request) {
	var content string
	content = ReadFromFile(configPath)
	w.WriteJson(content)
}
func PostConfig(w rest.ResponseWriter, r *rest.Request) {
	//newConfig := PromConfig{}
	var newConfig string
	//r.DecodeJsonPayload(newConfig)


	var bodyBytes []byte
	if r.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(r.Body)
	}
	newConfig=string(bodyBytes)
	fmt.Printf("saving this:\n%v",newConfig)

	SaveToFile(configPath,newConfig)
}
func AddTarget(w rest.ResponseWriter, r *rest.Request) {
	oldContent := ReadFromFile(configPath)
	newContent := oldContent
	//TODO add data to oldcontent
	SaveToFile(configPath,newContent)
}

func ReadFromFile(path string) string{
	filerc, err := os.Open(path)
	if err != nil{
		log.Fatal(err)
	}
	defer filerc.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(filerc)
	contents := buf.String()

	fmt.Printf("\nContents of file %v is: \n%v\n",path,contents)
	return contents
}

func SaveToFile(path string, content string){
	//The save should be atomic, hence we create a new file and do the rename as the atomic operation
	f, err := os.Create(path+".temp")
	if err != nil {
		fmt.Println(err)
		return
	}

	l, err := f.WriteString(content)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	rerr:=os.Rename(path+".temp", path)
	if rerr != nil {
		log.Fatal(rerr)
	}
	fmt.Println(l, "bytes written successfully")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}


//
//var store = map[string]*Country{}
//
//var lock = sync.RWMutex{}
//
//func GetCountry(w rest.ResponseWriter, r *rest.Request) {
//	code := r.PathParam("code")
//
//	lock.RLock()
//	var country *Country
//	if store[code] != nil {
//		country = &Country{}
//		*country = *store[code]
//	}
//	lock.RUnlock()
//
//	if country == nil {
//		rest.NotFound(w, r)
//		return
//	}
//	w.WriteJson(country)
//}
//
//func GetAllCountries(w rest.ResponseWriter, r *rest.Request) {
//	lock.RLock()
//	countries := make([]Country, len(store))
//	i := 0
//	for _, country := range store {
//		countries[i] = *country
//		i++
//	}
//	lock.RUnlock()
//	w.WriteJson(&countries)
//}
//
//func PostCountry(w rest.ResponseWriter, r *rest.Request) {
//	country := Country{}
//	err := r.DecodeJsonPayload(&country)
//	if err != nil {
//		rest.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	if country.Code == "" {
//		rest.Error(w, "country code required", 400)
//		return
//	}
//	if country.Name == "" {
//		rest.Error(w, "country name required", 400)
//		return
//	}
//	lock.Lock()
//	store[country.Code] = &country
//	lock.Unlock()
//	w.WriteJson(&country)
//}
//
//func DeleteCountry(w rest.ResponseWriter, r *rest.Request) {
//	code := r.PathParam("code")
//	lock.Lock()
//	delete(store, code)
//	lock.Unlock()
//	w.WriteHeader(http.StatusOK)
//}
