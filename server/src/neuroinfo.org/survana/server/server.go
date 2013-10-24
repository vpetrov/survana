package main

import (
        "errors"
        "log"
        "net"
        "crypto/tls"
        "net/http"
        "flag"
        "os"
        "os/user"
        "strconv"
        "syscall"
        "neuroinfo.org/survana"
        "neuroinfo.org/survana/admin"
       )

type Config struct {
    ConfigFile string
    IP string
    Port string
    PortNumber int
    Username string
    WWW string
    SSLCert string
    SSLKey string
}

const (
        DEFAULT_CONFIG = "./config.json"
        DEFAULT_IP = ""
        DEFAULT_PORT = 443
        MAX_PRIVILEGED_PORT = 1024
        ROOT_UID = 0
        DEFAULT_SSL_CERT = "ssl/cert.pem"
        DEFAULT_SSL_KEY  = "ssl/key.pem"
        DEFAULT_WWW = "../../../../../www"
      )

func main() {

    //cleanup and error reporting
    defer func () {
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

    err = checkArguments(args)
    if err != nil {
        panic(err)
    }

    //bind the address:port
    listener, err := listen(args.IP + ":" + args.Port, args.SSLCert, args.SSLKey)
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

    log.Println("Listening on ", args.IP + ":" + args.Port, "as", args.Username)

    //modules
    modules := make(map[string]survana.RequestHandler, 3)
    enableModules(modules, args.WWW)

    err = http.Serve(listener, nil)
    if err != nil {
        panic(err)
    }
}

func enableModules(modules map[string]survana.RequestHandler, www_root string) {
    //_ = NewModule("/admin/", www_root)
    _ = newModule("admin", www_root, "");
}

func newModule(name string, www_root, prefix string) survana.RequestHandler {

    if len(prefix) == 0 {
        prefix = "/" + name + "/"
    }

    module_dir := www_root + name + "/"
    static_dir := module_dir + survana.STATIC_DIR + "/"
    static_prefix := prefix + survana.STATIC_DIR + "/"

    var mod interface{}

    switch name {
        case "admin":   mod = &admin.Module{
                            Module: survana.Module {
                                    Name: name,
                                    Prefix: prefix,
                                    Dir: module_dir,
                                    StaticPrefix: static_prefix,
                                    StaticDir: static_dir,
                              },
                        }
    }

    handler := mod.(survana.RequestHandler)

    handler.Mount()

    return handler
}

func parseArguments(current_user *user.User) *Config {

    args :=  &Config{}

    //command line arguments
    flag.StringVar(&args.ConfigFile, "config", "", "Configuration file")
    flag.StringVar(&args.IP, "ip", DEFAULT_IP, "IP address to listen on")
    flag.IntVar(&args.PortNumber, "port", DEFAULT_PORT, "SSL port to listen on")
    flag.StringVar(&args.SSLCert, "sslcert", DEFAULT_SSL_CERT, "SSL PEM certificate file")
    flag.StringVar(&args.SSLKey, "sslkey", DEFAULT_SSL_KEY, "SSL PEM key file")
    flag.StringVar(&args.Username, "user", current_user.Username, "The user to run the server as")
    flag.StringVar(&args.WWW, "www", DEFAULT_WWW, "Folder containing all WWW files")
    flag.Parse()

    //convert port string to port number
    args.Port = strconv.Itoa(args.PortNumber)

    return args
}

func checkArguments(args *Config) error {

    euid := os.Geteuid()

    //verify port number
    if (args.PortNumber < MAX_PRIVILEGED_PORT) && (euid != ROOT_UID) {
        return errors.New("You must run this program as an administrator in order to start the server on port " + args.Port)
    }

    //if a configuration file was specified, verify its existence
    if len(args.ConfigFile) > 0 {
        _, err := os.Stat(args.ConfigFile)

        if os.IsNotExist(err) {
            panic(err)
        }
    }

    //verify the existence of the SSL certificate and key
    _, err := os.Stat(args.SSLCert)
    if os.IsNotExist(err) {
        panic(err)
    }

    _, err = os.Stat(args.SSLKey)
    if os.IsNotExist(err) {
        panic(err)
    }

    //verify that the WWW folder exists and is a folder
    www, err := os.Stat(args.WWW)
    if os.IsNotExist(err) {
        panic(err)
    }

    if !www.IsDir() {
        panic(errors.New(args.WWW + ": not a directory"))
    }

    //ensure trailing slash
    nwww := len(args.WWW)
    if args.WWW[nwww-1 : nwww] != "/" {
        args.WWW += "/"
    }

    return nil
}

/* Starts a TLS listener on the specified address, using the specified SSL certificate and key */
func listen(address, sslcert, sslkey string) (tlsListener net.Listener, err error) {

    tlsConfig := &tls.Config{}
    //attach SSL certificates to the TLS configuration
    tlsConfig.Certificates = make([]tls.Certificate, 1)
    tlsConfig.Certificates[0], err = tls.LoadX509KeyPair(sslcert, sslkey)
    if err != nil {
        panic(err)
    }

    //listen on the specified IP and port
    socket, err := net.Listen("tcp", address)
    if err != nil {
        panic(err)
    }

    //wrap the listener with a TLS listener
    tlsListener = tls.NewListener(socket, tlsConfig)

    return
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
