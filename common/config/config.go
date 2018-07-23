package config

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	. "fogflow/common/datamodel"
)

var (
	INFO     *log.Logger
	PROTOCOL *log.Logger
	ERROR    *log.Logger
	DEBUG    *log.Logger
)

type DatabaseCfg struct {
	DBReset  bool   `json:"dbreset"`
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	DBname   string `json:"dbname"`
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

type Config struct {
	ExternalIP string           `json:"external_ip"`
	InternalIP string           `json:"internal_ip"`
	PLocation  PhysicalLocation `json:"physical_location"`
	LLocation  LogicalLocation  `json:"logical_location"`
	Logging    struct {
		Info     string `json:"info"`
		Protocol string `json:"protocol"`
		Errlog   string `json:"error"`
		Debug    string `json:"debug"`
	} `json:"logging"`
	Discovery struct {
		Port  int         `json:"port"`
		DBCfg DatabaseCfg `json:"postgresql"`
	} `json:"discovery"`
	Broker struct {
		Port          int `json:"port"`
		WebSocketPort int `json:"websocket"`
	} `json:"broker"`
	Master struct {
		AgentPort int `json:"ngsi_agent_port"`
	} `json:"master"`
	Worker struct {
		Registry            RegistryConfiguration `json:"registry,omitempty"`
		ContainerAutoRemove bool                  `json:"container_autoremove"`
	} `json:"worker"`
	RabbitMQ struct {
		Port     int    `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"rabbitmq"`
}

var logTargets map[string]io.Writer = map[string]io.Writer{
	"stdout":  os.Stdout,
	"stderr":  os.Stderr,
	"discard": ioutil.Discard,
}

func (c *Config) GetDiscoveryURL() string {
	discoveryURL := fmt.Sprintf("http://%:%s/ngsi9", c.InternalIP, c.Discovery.Port)
	return discoveryURL
}

func (c *Config) GetMessageBus() string {
	messageBus := fmt.Sprintf("amqp://%:%s@%s:%s/", c.RabbitMQ.Username, c.RabbitMQ.Password, c.InternalIP, c.RabbitMQ.Port)
	return messageBus
}

func (c *Config) SetLogTargets() {
	target, ok := logTargets[c.Logging.Info]
	if !ok {
		target = os.Stdout
	}
	INFO = log.New(target, "INFO: ", log.Ldate|log.Ltime)
	target, ok = logTargets[c.Logging.Protocol]
	if !ok {
		target = ioutil.Discard
	}
	PROTOCOL = log.New(target, "PROTOCOL: ", log.Ldate|log.Ltime)
	target, ok = logTargets[c.Logging.Errlog]
	if !ok {
		target = os.Stderr
	}
	ERROR = log.New(target, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	target, ok = logTargets[c.Logging.Debug]
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

func LoadConfig(configFile string) (Config, error) {
	var config Config

	abspath, _ := filepath.Abs(configFile)
	err := ParseConfig(abspath, &config)
	if err != nil {
		return config, err
	}

	config.SetLogTargets()
	return config, nil
}
