package main

import (
	"crypto/tls"
	"flag"
	"labix.org/v2/mgo"
	"log"
	"net"
	"net/http"
	"neuroinfo.org/survana"
	"neuroinfo.org/survana/admin"
	"os"
	"os/user"
	"strconv"
	"syscall"
)

const (
	DEFAULT_CONFIG       = "survana.json"
	MAX_PRIVILEGED_PORT  = 1024
	ROOT_UID             = 0
	DEFAULT_CONFIG_PERMS = 0600
)

var configFile string

func main() {
	log.Println("Starting Survana")
	//detect current user and effective user id
	cuser, err := user.Current()
	if err != nil {
		panic(err)
	}

	//parse command-line arguments
	ParseArguments()

	//read configuration file
	config, err := ReadConfiguration(configFile)
	if err != nil {
		panic(err)
	}

	if len(config.Username) == 0 {
		config.Username = cuser.Username
	}

	//open the port as the current user
	listener, err := Listen(config)
	if err != nil {
		panic(err)
	}

	//switch to unprivileged mode
	if config.Username != cuser.Username {
		err = DecreasePrivilegesTo(config.Username)
		if err != nil {
			panic(err)
		}
	}

	//Mount all modules
	err = EnableModules(config)
	if err != nil {
		panic(err)
	}

	log.Println("Listening on ", config.IP+":"+config.Port, "as", config.Username)

	//Go!
	err = http.Serve(listener, nil)
	if err != nil {
		panic(err)
	}
}

func ParseArguments() {

	//command line arguments
	flag.StringVar(&configFile, "config", DEFAULT_CONFIG, "Configuration file")
	flag.Parse()
}

//reads a configuration file and returns a Config object
func ReadConfiguration(path string) (config *Config, err error) {

	//open the configuration file
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	config, err = NewConfig(file)

	return
}

func DecreasePrivilegesTo(username string) (err error) {

	//lookup the new user
	cuser, err := user.Lookup(username)

	if err != nil {
		return
	}

	//convert the new effective user id to an int
	new_euid, err := strconv.Atoi(cuser.Uid)
	if err != nil {
		return
	}

	//convert the new effective group id to an int
	new_egid, err := strconv.Atoi(cuser.Gid)

	//set the group id first (if uid were set first, setgid would fail now)
	err = syscall.Setgid(new_egid)
	if err != nil {
		return
	}

	//reduce the privileges of the process
	err = syscall.Setuid(new_euid)
	if err != nil {
		return
	}

	return
}

// Starts a net.Listener on the specified address, using the specified SSL certificate and key
func Listen(config *Config) (tlsListener net.Listener, err error) {
	log.Printf("Reading SSL certificate (%s) and SSL key (%s)", config.SSLCert, config.SSLKey)

	tlsConfig := &tls.Config{}
	//attach SSL certificates to the TLS configuration
	tlsConfig.Certificates = make([]tls.Certificate, 1)
	tlsConfig.Certificates[0], err = tls.LoadX509KeyPair(config.SSLCert, config.SSLKey)
	if err != nil {
		return
	}

	//listen on the specified IP and port
	socket, err := net.Listen("tcp", config.IP+":"+config.Port)
	if err != nil {
		return
	}

	//wrap the listener with a net.Listener
	tlsListener = tls.NewListener(socket, tlsConfig)

	return
}

//Create and mount all known modules
func EnableModules(config *Config) error {

	dbSession, err := mgo.Dial(config.DbUrl)
	if err != nil {
		return err
	}

	//close db connection
	defer dbSession.Close()

	//ADMIN
	admin_module := admin.NewModule(config.WWW+"/admin", dbSession)
	survana.Modules.Mount(admin_module.Module, "/admin")

	return nil
}
