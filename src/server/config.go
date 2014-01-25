package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"strconv"
)

const (
	DEFAULT_IP       = ""
	DEFAULT_PORT     = 4443
	DEFAULT_SSL_CERT = "ssl/cert.pem"
	DEFAULT_SSL_KEY  = "ssl/key.pem"
	DEFAULT_WWW      = "/www/survana"
	DEFAULT_DB_URL   = "mongodb://localhost/survana"
)

type Config struct {
	IP         string	`json:"ip"`
	Port       string	`json:"-"`
	PortNumber int		`json:"port"`
	Username   string	`json:"username"`
	WWW        string	`json:"www"`
	SSLCert    string	`json:"sslcert"`
	SSLKey     string	`json:"sslkey"`
	DbUrl      string	`json:"db"`
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
