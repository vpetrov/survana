package main;

import (
        "errors"
        "log"
        "net"
        "flag"
        "os"
        "os/user"
        "strconv"
        "syscall"
       )

type Config struct {
    IP string
    Port string
    PortNumber int
    Username string
    WWW string
}

const (
        DEFAULT_IP = ""
        DEFAULT_PORT = 80
      )

func main() {

    //cleanup
    defer func () {
        if err := recover(); err != nil {
            log.Println("= ERROR =")
            log.Println(err)
            log.Fatal()
        }
    }()

    //detect current user and effective user id
    cuser, err := user.Current()
    if err != nil {
        panic(err)
    }

    euid := os.Geteuid()

    //parse command-line arguments
    args := parseArguments(cuser)

    //attempt to listen on the chosen port
    if (args.PortNumber < 1024) && (euid != 0) {
        panic(errors.New("You must run this program as an administrator in order to use port " + args.Port));
    }

    //listen on the specified IP and port
    _, err = net.Listen("tcp", args.IP + ":" + args.Port)
    if err != nil {
        panic(err)
    }

    if args.Username != cuser.Username {
        //lookup the new user
        new_euser, err := user.Lookup(args.Username)

        if err != nil {
            panic(errors.New(args.Username + ": Unable to lookup user: " + err.Error()))
        }

        //convert the new effective user id to an int
        new_euid, err := strconv.Atoi(new_euser.Uid)
        if err != nil {
            panic(err)
        }

        //convert the new effective group id to an int
        new_egid, err := strconv.Atoi(new_euser.Gid)

        //set the group id first (if uid were set first, setgid would fail now)
        err = syscall.Setgid(new_egid)
        if err != nil {
            panic(err)
        }

        //reduce the privileges of the process
        err = syscall.Setuid(new_euid)
        if err != nil {
            panic(err)
        }
    }

    log.Println("args=",args)
    log.Println("Listening on port", args.Port)
    log.Println("euid=", os.Geteuid(), ", gid=", os.Getegid())
}

func parseArguments(user *user.User) *Config {

    args :=  &Config{}

    //command line arguments
    flag.StringVar(&args.IP, "ip", DEFAULT_IP, "IP address to listen on")
    flag.IntVar(&args.PortNumber, "port", DEFAULT_PORT, "Port to listen on")
    flag.StringVar(&args.Username, "user", user.Username, "The user to run the server as")
    flag.StringVar(&args.WWW, "www", "../www", "Folder containing all WWW files")
    flag.Parse()

    //convert port string to port number
    args.Port = strconv.Itoa(args.PortNumber)

    return args
}
