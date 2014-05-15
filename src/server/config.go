package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"strconv"
	"neuroinformatics.harvard.edu/survana/dashboard"
	"neuroinformatics.harvard.edu/survana/study"
    "neuroinformatics.harvard.edu/survana/store"
)

const (
	DEFAULT_IP       = ""
	DEFAULT_PORT     = 4443
    DEFAULT_KEY      = "survana.key"
	DEFAULT_SSL_CERT = "ssl/cert.pem"
	DEFAULT_SSL_KEY  = "ssl/key.pem"
	DEFAULT_WWW      = "/www/survana"
	DEFAULT_DB_URL   = "mongodb://localhost/survana"
)

type Config struct {
	IP         string	`json:"ip"`		//web
	Port       string	`json:"-"`		//web
	PortNumber int		`json:"port"`
	Username   string	`json:"username"`	//general
	WWW        string	`json:"www"`		//web
    Key        string   `json:"key"`        //survana private key
	SSLCert    string	`json:"sslcert"`	//web
	SSLKey     string	`json:"sslkey"`		//web
	DbUrl      string	`json:"db"`			//database
    Modules    *ModuleConfig `json:"modules"`//modules
}

type ModuleConfig struct {
    Dashboard   *dashboard.Config   `json:"dashboard,omitempty"`
    Study       *study.Config       `json:"study,omitempty"`
    Store       *store.Config       `json:"store,omitempty"`
}

// Creates a new configuration object and sets empty values to default
func NewConfig(src io.Reader) (config *Config, err error) {

	//read the configuration data
	bytes, err := ioutil.ReadAll(src)
	if err != nil {
		return
	}

	//convert to JSON
	config = &Config{}
	err = json.Unmarshal(bytes, config)
	if err != nil {
		return
	}

	if len(config.IP) == 0 {
		config.IP = DEFAULT_IP
	}

	if len(config.Port) == 0 {
		if config.PortNumber == 0 {
			config.PortNumber = DEFAULT_PORT
		}
		config.Port = strconv.Itoa(DEFAULT_PORT)
	} else {
		config.PortNumber, err = strconv.Atoi(config.Port)
		if err != nil {
			return
		}
	}

	if len(config.WWW) == 0 {
		config.WWW = DEFAULT_WWW
	}

    if len(config.Key) == 0 {
        config.Key = DEFAULT_KEY
    }

	if len(config.SSLCert) == 0 {
		config.SSLCert = DEFAULT_SSL_CERT
	}

	if len(config.SSLKey) == 0 {
		config.SSLKey = DEFAULT_SSL_KEY
	}

	if len(config.DbUrl) == 0 {
		config.DbUrl = DEFAULT_DB_URL
	}

	return
}

//Converts the configuration object into a JSON byte array
func (c *Config) ToJSON() ([]byte, error) {
	return json.MarshalIndent(c, "", "    ")
}
