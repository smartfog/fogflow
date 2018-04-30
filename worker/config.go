package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	. "fogflow/common/datamodel"
)

//loggers
var (
	INFO     *log.Logger
	PROTOCOL *log.Logger
	ERROR    *log.Logger
	DEBUG    *log.Logger
)

//The default output for all the loggers is set to ioutil.Discard
func init() {
	INFO = log.New(ioutil.Discard, "", 0)
	PROTOCOL = log.New(ioutil.Discard, "", 0)
	ERROR = log.New(ioutil.Discard, "", 0)
	DEBUG = log.New(ioutil.Discard, "", 0)
}

type Logging struct {
	Info     string `json:"info"`
	Protocol string `json:"protocol"`
	Errlog   string `json:"error"`
	Debug    string `json:"debug"`
}

type RegistryConfiguration struct {
	Username      string `json:"username,omitempty"`
	Password      string `json:"password,omitempty"`
	Email         string `json:"email,omitempty"`
	ServerAddress string `json:"serveraddress,omitempty"`
}

func (r *RegistryConfiguration) IsConfigured() bool {
	if r.Username != "" && r.Password != "" && r.Email != "" && r.ServerAddress != "" {
		return true
	}

	return false
}

//Current configuration struct, maxQueueDepth sets the maximum number of unacknowledged mesages
//for a client. Listeners is a slice of ListenerConfigs
type Config struct {
	MyIP                string                `json:"my_ip"`
	MessageBus          string                `json:"message_bus"`
	DiscoveryURL        string                `json:"iot_discovery_url"`
	MyRole              string                `json:"my_role"`
	ContainerAutoRemove bool                  `json:"container_autoremove"`
	Registry            RegistryConfiguration `json:"registry,omitempty"`
	BrokerURL           string                `json:"broker_url"` // known from IoT Discovery
	PLocation           PhysicalLocation      `json:"physical_location"`
	LLocation           LogicalLocation       `json:"logical_location"`
	Log                 Logging               `json:"logging"`
}

var logTargets map[string]io.Writer = map[string]io.Writer{
	"stdout":  os.Stdout,
	"stderr":  os.Stderr,
	"discard": ioutil.Discard,
}

func (c *Config) SetLogTargets() {
	target, ok := logTargets[c.Log.Info]
	if !ok {
		target = os.Stdout
	}
	INFO = log.New(target, "INFO: ", log.Ldate|log.Ltime)

	target, ok = logTargets[c.Log.Protocol]
	if !ok {
		target = ioutil.Discard
	}
	PROTOCOL = log.New(target, "PROTOCOL: ", log.Ldate|log.Ltime)

	target, ok = logTargets[c.Log.Errlog]
	if !ok {
		target = os.Stderr
	}
	ERROR = log.New(target, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	target, ok = logTargets[c.Log.Debug]
	if !ok {
		target = ioutil.Discard
	}
	DEBUG = log.New(target, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func ParseConfig(confFile string, confVar *Config) error {
	file, err := os.Open(confFile)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(file)

	err = decoder.Decode(confVar)
	if err != nil {
		return err
	}

	return nil
}

func createConfig(configFile string) Config {
	var config Config

	if configFile == "" {
		config.MessageBus = "amqp://guest:guest@localhost:5672/"
	} else {
		err := ParseConfig(configFile, &config)
		if err != nil {
			os.Stderr.WriteString(fmt.Sprintf("%s\n", err.Error()))
		}
	}

	config.SetLogTargets()
	return config
}

func LoadConfig() Config {
	configurationFile := flag.String("f", "config.json", "A configuration file")

	myNID := flag.Int("nodeid", 0, "the current node ID")
	mySID := flag.Int("siteid", 0, "the current site ID")
	myLID := flag.Int("layer", 0, "the current layer ID")
	myPID := flag.Int("parent", 0, "the ID of the parent site")

	flag.Parse()

	config := createConfig(*configurationFile)

	if *myNID != 0 && *mySID != 0 && *myLID != 0 && *myPID != 0 {
		config.LLocation.NodeNo = *myNID
		config.LLocation.SiteNo = *mySID
		config.LLocation.LayerNo = *myLID
		config.LLocation.ParentSiteNo = *myPID
	}

	return config
}
