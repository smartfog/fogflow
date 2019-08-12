package config

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	. "github.com/smartfog/fogflow/common/datamodel"
	. "github.com/smartfog/fogflow/common/ngsi"
)

type DatabaseCfg struct {
	UseOnlyCache bool   `json:"use_only_cache"`
	DBReset      bool   `json:"dbreset"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Host         string `json:"host"`
	Port         int    `json:"port"`
	DBname       string `json:"dbname"`
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
	WebPortalIP   string           `json:"webportal_ip"`
	CoreSerivceIP string           `json:"coreservice_ip"`
	ExternalIP    string           `json:"external_hostip"`
	InternalIP    string           `json:"internal_hostip"`
	Location      PhysicalLocation `json:"physical_location"`
	SiteID        string           `json:"site_id"`
	Logging       struct {
		Info     string `json:"info"`
		Protocol string `json:"protocol"`
		Errlog   string `json:"error"`
		Debug    string `json:"debug"`
	} `json:"logging"`
	Discovery struct {
		HTTPPort  int `json:"http_port"`
		HTTPSPort int `json:"https_port"`
	} `json:"discovery"`
	Broker struct {
		HTTPPort  int `json:"http_port"`
		HTTPSPort int `json:"https_port"`
	} `json:"broker"`
	Master struct {
		AgentPort int `json:"ngsi_agent_port"`
	} `json:"master"`
	Worker struct {
		Registry            RegistryConfiguration `json:"registry,omitempty"`
		ContainerAutoRemove bool                  `json:"container_autoremove"`
		StartActualTask     bool                  `json:"start_actual_task"`
		Capacity            int                   `json:"capacity"`
		//EdgeAddress         string                `json:"edge_address"`
		CAdvisorPort int `json:"cadvisor_port"`
	} `json:"worker"`
	RabbitMQ struct {
		Port     int    `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"rabbitmq"`
	HTTPS      HTTPS
	Prometheus struct {
		Address   string `json:"address"`
		DataPort  int    `json:"data_port"`
		AdminPort int    `json:"admin_port"`
	} `json:"prometheus"`
}

var logTargets map[string]io.Writer = map[string]io.Writer{
	"stdout":  os.Stdout,
	"stderr":  os.Stderr,
	"discard": ioutil.Discard,
}

func (c *Config) GetDiscoveryURL(flag4HTTPS bool) string {
	if flag4HTTPS == false {
		discoveryURL := fmt.Sprintf("http://%s:%d/ngsi9", c.CoreSerivceIP, c.Discovery.HTTPPort)
		return discoveryURL
	}

	if c.HTTPS.Enabled == true {
		discoveryURL := fmt.Sprintf("https://%s:%d/ngsi9", c.CoreSerivceIP, c.Discovery.HTTPSPort)
		return discoveryURL
	} else {
		discoveryURL := fmt.Sprintf("http://%s:%d/ngsi9", c.CoreSerivceIP, c.Discovery.HTTPPort)
		return discoveryURL
	}
}

func (c *Config) GetMessageBus() string {
	messageBus := fmt.Sprintf("amqp://%s:%s@%s:%d/", c.RabbitMQ.Username, c.RabbitMQ.Password, c.CoreSerivceIP, c.RabbitMQ.Port)
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
