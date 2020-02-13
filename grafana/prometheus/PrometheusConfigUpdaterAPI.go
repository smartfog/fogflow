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
	"sync"
)

var (
	IncomingPortNumber int = 4545
	//Store the targets in this file
	TargetsFileStorage="/etc/prometheus/tgroups/target_groups.json"
	FileMutex sync.Mutex
)

//Each target for prometheus has a hostname and port.
type Target struct {
	Address string
	Port string
}


//A simple API to read the config, and write to config
func main() {
	MainAPI := rest.NewApi()
	MainAPI.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/config", GetConfig),
		rest.Post("/config", PostConfig),
	)
	if err != nil {
		log.Fatal(err)
	}
	MainAPI.SetApp(router)

	panic(http.ListenAndServe(":"+strconv.Itoa(IncomingPortNumber), MainAPI.MakeHandler()))

}


func GetConfig(w rest.ResponseWriter, r *rest.Request) {
	var content string
	content = ReadFromFile(TargetsFileStorage)
	w.WriteJson(content)
}
func PostConfig(w rest.ResponseWriter, r *rest.Request) {
	var NewConfig string

	var BodyBytes []byte
	var error interface{}
	if r.Body != nil {
		BodyBytes, error = ioutil.ReadAll(r.Body)
	}

	CheckError(error)
	NewConfig=string(BodyBytes)

	SaveToFile(TargetsFileStorage,NewConfig)
}



func ReadFromFile(path string) string{
	filerc, error := os.Open(path)
	CheckError(error)

	defer filerc.Close()

	buffer := new(bytes.Buffer)
	buffer.ReadFrom(filerc)
	Content := buffer.String()

	return Content
}

func SaveToFile(path string, content string){
	FileMutex.Lock()
	defer FileMutex.Unlock()
	//The save should be atomic, hence, we create a new file and do the rename as the atomic operation
	//Create a new file
	TemporaryFile, error := os.Create(path+".temp")
	CheckError(error)

	//Write to the new file
	WrittenData, err := TemporaryFile.WriteString(content)
	CheckError(error)
	if err != nil {
		TemporaryFile.Close()
	}

	error=os.Rename(path+".temp", path)
	CheckError(error)

	fmt.Println(WrittenData, "bytes written successfully")

	CheckError(TemporaryFile.Close())

}

func CheckError(err interface{}) {
	if err != nil {
		fmt.Println(err)
		return
	}
}
