package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net"
	"net/http"
	"net/url"
	"neuroinformatics.harvard.edu/survana"
	"neuroinformatics.harvard.edu/survana/dashboard"
	"neuroinformatics.harvard.edu/survana/study"
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

var (
	configFile string
	DB         survana.Database
)

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

    private_key, err := GetPrivateKey(config.Key)
    if err != nil {
        panic(err)
    }

    log.Println("Survana ID:", private_key.Id)

	//Mount all modules
	err = EnableModules(private_key, config)
	if err != nil {
		panic(err)
	}

	log.Println("Listening on ", config.IP+":"+config.Port, "as", config.Username)

	//Go!
	err = http.Serve(listener, survana.Modules)
	if err != nil {
		panic(err)
	}

	DB.Disconnect()
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
func EnableModules(private_key *survana.PrivateKey, config *Config) (err error) {

    log.Println("%#v", config);

	//dashboard
	dashboard_module := dashboard.NewModule(config.WWW+"/dashboard",
                                            GetDB(config.DbUrl, "dashboard"),
                                            config.Modules.Dashboard,
                                            private_key)
	survana.Modules.Mount(dashboard_module.Module, "/dashboard")

    //study
    //TODO: figure out how the dashboard should share published studies with the study module
    //for now, let them use the same database
    study_module := study.NewModule(config.WWW + "/study", GetDB(config.DbUrl, "dashboard"))
    survana.Modules.Mount(study_module.Module, "/study")

	return nil
}

func GetDB(u string, dbname string) survana.Database {
	dburl, err := url.Parse(u)
	if err != nil {
		panic(err)
	}

	if len(dbname) == 0 {
		panic("Invalid database name")
	}

	DB = survana.NewDatabase(dburl, dbname)
	if err != nil {
		panic(err)
	}

	err = DB.Connect()
	if err != nil {
		panic(err)
	}

	log.Println("Connected to", DB.SystemInformation())
	log.Println("Database version:", DB.Version())

	return DB
}

