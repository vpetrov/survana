package main

import (
	"errors"
	"flag"
	"log"
	"neuroinfo.org/survana/server"
	"os"
	"os/user"
	"reflect"
	"strconv"
	"strings"
	"syscall"
)

const (
	DEFAULT_CONFIG       = "survana.json"
	MAX_PRIVILEGED_PORT  = 1024
	ROOT_UID             = 0
	DEFAULT_CONFIG_PERMS = 0600
)

type Arguments struct {
	ConfigFile string
	Username   string
	Config     server.Config
}

func main() {

	//cleanup and error reporting
	defer func() {
		/*        if err := recover(); err != nil {
		          log.Println("ERROR:")
		          log.Fatal(err)
		      } */
	}()

	//detect current user and effective user id
	cuser, err := user.Current()
	if err != nil {
		panic(err)
	}

	//parse command-line arguments
	args := parseArguments(cuser)

	//open the configuration file
	file, err := os.Open(args.ConfigFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	//parse the configuration file
	config, err := server.NewConfig(file)
	if err != nil {
		panic(err)
	}

	err = checkArguments(args)
	if err != nil {
		panic(err)
	}

	//apply all arguments to this config
	overrideConfig(config)

    err = checkConfig(config)
    if err != nil {
        panic(err)
    }

	//create a new server instance based on the configuration
	srv, err := server.NewServer(config)

    if err != nil {
        panic(err)
    }

	//bind the address:port
	err = srv.Listen()
	if err != nil {
		panic(err)
	}

	//switch to unprivileged mode
	if args.Username != cuser.Username {
		err = decreasePrivilegesTo(args.Username)
		if err != nil {
			panic(err)
		}
	}

	log.Println("Listening on ", args.Config.Ip+":"+args.Config.Port, "as", args.Username)

	//modules
	srv.EnableModules()

	err = srv.Serve()
	if err != nil {
		panic(err)
	}

}

func parseArguments(current_user *user.User) *Arguments {

	args := &Arguments{}

	//command line arguments
	flag.StringVar(&args.ConfigFile, "config", DEFAULT_CONFIG, "Configuration file")
	flag.StringVar(&args.Username, "user", current_user.Username, "The user to run the server as")

    flag.StringVar(&args.Config.Dburl, "dburl", server.DEFAULT_DB_URL, "Database URL")
	flag.StringVar(&args.Config.Ip, "ip", server.DEFAULT_IP, "IP address to listen on")
	flag.IntVar(&args.Config.PortNumber, "port", server.DEFAULT_PORT, "SSL port to listen on")
	flag.StringVar(&args.Config.Sslcert, "sslcert", server.DEFAULT_SSL_CERT, "SSL PEM certificate file")
	flag.StringVar(&args.Config.Sslkey, "sslkey", server.DEFAULT_SSL_KEY, "SSL PEM key file")
	flag.StringVar(&args.Config.Www, "www", server.DEFAULT_WWW, "Folder containing all WWW files")
	flag.Parse()

	//convert port string to port number
	args.Config.Port = strconv.Itoa(args.Config.PortNumber)

	return args
}

func checkArguments(args *Arguments) error {

	euid := os.Geteuid()

	//verify port number
	if (args.Config.PortNumber < MAX_PRIVILEGED_PORT) && (euid != ROOT_UID) {
		return errors.New("You must run this program as an administrator in order to start the server on port " + args.Config.Port)
	}

	//if a configuration file was specified, verify its existence
	if len(args.ConfigFile) > 0 {
		_, err := os.Stat(args.ConfigFile)

		if os.IsNotExist(err) {
			panic(err)
		}
	}

    return nil
}

func checkConfig(config *server.Config) error {

	//verify the existence of the SSL certificate and key
	_, err := os.Stat(config.Sslcert)
	if os.IsNotExist(err) {
		panic(err)
	}

	_, err = os.Stat(config.Sslkey)
	if os.IsNotExist(err) {
		panic(err)
	}

	//verify that the WWW folder exists and is a folder
	www, err := os.Stat(config.Www)
	if os.IsNotExist(err) {
		panic(err)
	}

	if !www.IsDir() {
		panic(errors.New(config.Www + ": not a directory"))
	}

	//ensure trailing slash
	nwww := len(config.Www)
	if config.Www[nwww-1:nwww] != "/" {
		config.Www += "/"
	}

	return nil
}

// Iterate over each flag that was set and override the config value
// by using reflection to look up the name of the argument in the config fields
func overrideConfig(config *server.Config) {

	//only visit arguments that have been set
	flag.Visit(func(f *flag.Flag) {

		//use reflection to get to the definition of the config struct
		ps := reflect.ValueOf(config)
		//dereference pointer
		s := ps.Elem()

		//our struct must have type reflect.Struct
		if s.Kind() != reflect.Struct {
			panic(errors.New("Configuration object must be a struct"))
		}

		//ucfirst() field name
		fieldName := strings.ToUpper(f.Name[:1]) + f.Name[1:]

		//find the struct field by name
		field := s.FieldByName(fieldName)

		//skip invalid or read-only fields
		if !field.IsValid() || !field.CanSet() {
			return
		}

		//TODO: Figure out how to set other types
		switch field.Kind() {
		case reflect.String:
			field.SetString(f.Value.String())
		}
	})
}

func decreasePrivilegesTo(username string) (err error) {

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
