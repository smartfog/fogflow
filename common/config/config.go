package config

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"

	. "fogflow/common/datamodel"
	. "fogflow/common/ngsi"
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
	CoreSerivceIP string           `json:"coreservice_ip"`
	ExternalIP    string           `json:"external_ip"`
	InternalIP    string           `json:"my_hostip"`
	Location      PhysicalLocation `json:"physical_location"`
	SiteID        string           `json:"site_id"`
	Logging       struct {
		Info     string `json:"info"`
		Protocol string `json:"protocol"`
		Errlog   string `json:"error"`
		Debug    string `json:"debug"`
	} `json:"logging"`
	Discovery struct {
		HostIP    string `json:"host_ip"`
		HTTPPort  int    `json:"http_port"`
		HTTPSPort int    `json:"https_port"`
	} `json:"discovery"`
	Broker struct {
		HostIP            string `json:"host_ip"`
		HTTPPort          int    `json:"http_port"`
		HTTPSPort         int    `json:"https_port"`
		HeartbeatInterval int    `json:"heartbeat_interval"`
	} `json:"broker"`
	Master struct {
		HostIP      string `json:"host_ip"`
		AgentPort   int    `json:"ngsi_agent_port"`
		RESTAPIPort int    `json:"rest_api_port"`
	} `json:"master"`
	Designer struct {
		HostIP     string `json:"host_ip"`
		WebSrvPort int    `json:"webSrvPort"`
		HTTPSPort  int    `json:"https_webSrvPort"`
	} `json:"designer"`
	Worker struct {
		ContainerManagement string                `json:"container_management"`
		AppNameSpace        string                `json:"app_namespace"`
		EdgeControllerPort  int                   `json:"edge_controller_port"`
		Registry            RegistryConfiguration `json:"registry,omitempty"`
		ContainerAutoRemove bool                  `json:"container_autoremove"`
		StartActualTask     bool                  `json:"start_actual_task"`
		Capacity            int                   `json:"capacity"`
		HeartbeatInterval   int                   `json:"heartbeat_interval"`
		DetectionDuration   int                   `json:"detection_duration"`
	} `json:"worker"`
	RabbitMQ struct {
		HostIP   string `json:"host_ip"`
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

func (c *Config) GetDiscoveryURL() string {
	protocol := "http"
	port := c.Discovery.HTTPPort
	if c.HTTPS.Enabled == true {
		protocol = "https"
		port = c.Discovery.HTTPSPort
	}

	hostip := c.InternalIP
	if c.CoreSerivceIP != "" {
		hostip = c.CoreSerivceIP
	}

	if c.Discovery.HostIP != "" {
		hostip = c.Discovery.HostIP
	}

	return fmt.Sprintf("%s://%s:%d/ngsi9", protocol, hostip, port)
}

func (c *Config) GetEdgeControllerURL() string {
	protocol := "http"
	port := c.Worker.EdgeControllerPort

	hostip := c.InternalIP
	if c.ExternalIP != "" {
		hostip = c.ExternalIP
	}

	return fmt.Sprintf("%s://%s:%d", protocol, hostip, port)
}

func (c *Config) GetDesignerURL() string {
	protocol := "http"
	port := c.Designer.WebSrvPort
	if c.HTTPS.Enabled == true {
		protocol = "https"
		port = c.Designer.HTTPSPort
	}

	hostip := c.InternalIP

	if c.Designer.HostIP != "" {
		hostip = c.Designer.HostIP
	}

	return fmt.Sprintf("%s://%s:%d", protocol, hostip, port)
}

func (c *Config) GetBrokerURL4Task() string {
	brokeIP := c.InternalIP

	if c.Broker.HostIP != "" {
		brokeIP = c.Broker.HostIP
	}

	return "http://" + brokeIP + ":" + strconv.Itoa(c.Broker.HTTPPort) + "/ngsi10"
}

func (c *Config) GetBrokerURL() string {
	protocol := "http"
	port := c.Broker.HTTPPort
	if c.HTTPS.Enabled == true {
		protocol = "https"
		port = c.Broker.HTTPSPort
	}

	hostip := c.InternalIP
	if c.ExternalIP != "" {
		hostip = c.ExternalIP
	}

	if c.Broker.HostIP != "" {
		hostip = c.Broker.HostIP
	}

	return fmt.Sprintf("%s://%s:%d/ngsi10", protocol, hostip, port)
}

func (c *Config) GetExternalBrokerURL() string {
	protocol := "http"
	port := c.Broker.HTTPPort
	if c.HTTPS.Enabled == true {
		protocol = "https"
		port = c.Broker.HTTPSPort
	}

	hostip := c.InternalIP
	if c.ExternalIP != "" {
		hostip = c.ExternalIP
	}

	return fmt.Sprintf("%s://%s:%d/ngsi10", protocol, hostip, port)
}

func (c *Config) GetMasterIP() string {
	hostip := c.InternalIP
	if c.ExternalIP != "" {
		hostip = c.ExternalIP
	}

	if c.Master.HostIP != "" {
		hostip = c.Master.HostIP
	}

	return hostip
}

func (c *Config) GetMessageBus() string {
	hostip := c.InternalIP
	if c.CoreSerivceIP != "" {
		hostip = c.CoreSerivceIP
	}

	if c.RabbitMQ.HostIP != "" {
		hostip = c.RabbitMQ.HostIP
	}

	messageBus := fmt.Sprintf("amqp://%s:%s@%s:%d/", c.RabbitMQ.Username, c.RabbitMQ.Password, hostip, c.RabbitMQ.Port)
	return messageBus
}

var logTargets map[string]io.Writer = map[string]io.Writer{
	"stdout":  os.Stdout,
	"stderr":  os.Stderr,
	"discard": ioutil.Discard,
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
