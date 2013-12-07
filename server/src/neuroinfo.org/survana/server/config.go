package server

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"neuroinfo.org/survana/admin"
	"strconv"
)

const (
	DEFAULT_IP       = ""
	DEFAULT_PORT     = 4443
	DEFAULT_SSL_CERT = "ssl/cert.pem"
	DEFAULT_SSL_KEY  = "ssl/key.pem"
	DEFAULT_WWW      = "../../../../../www"
    DEFAULT_DB_URL   = "mongodb://localhost/survana"
)

type Config struct {
	Ip         string
	Port       string
	PortNumber int
	Username   string
	Www        string
	Sslcert    string
	Sslkey     string
    Dburl      string
	Admin      admin.Config
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

	if len(config.Ip) == 0 {
		config.Ip = DEFAULT_IP
	}

	if len(config.Port) == 0 {
		config.PortNumber = DEFAULT_PORT
		config.Port = strconv.Itoa(DEFAULT_PORT)
	} else {
		config.PortNumber, err = strconv.Atoi(config.Port)
		if err != nil {
			return
		}
	}

	if len(config.Www) == 0 {
		config.Www = DEFAULT_WWW
	}

	if len(config.Sslcert) == 0 {
		config.Sslcert = DEFAULT_SSL_CERT
	}

	if len(config.Sslkey) == 0 {
		config.Sslkey = DEFAULT_SSL_KEY
	}

    if len(config.Dburl) == 0 {
        config.Dburl = DEFAULT_DB_URL
    }

	return
}

//Converts the configuration object into a JSON byte array
func (c *Config) ToJSON() ([]byte, error) {
	return json.MarshalIndent(c, "", "    ")
}
