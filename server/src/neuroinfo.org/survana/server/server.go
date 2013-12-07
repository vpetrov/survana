package server

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"neuroinfo.org/survana"
	"neuroinfo.org/survana/admin"
    "labix.org/v2/mgo"
    _ "labix.org/v2/mgo/bson"
)

type Server struct {
	Config      *Config
	tlsListener net.Listener
    dbSession   *mgo.Session
    Modules     map[string]survana.RequestHandler
}

func NewServer(config *Config) (server *Server, err error) {
    dbSession, err := mgo.Dial(config.Dburl)
    if err != nil {
        return
    }

	 server = &Server{
                Config: config,
                dbSession: dbSession,
              }
     return
}

func (srv *Server) EnableModules() {
	_ = srv.newModule(survana.ADMIN, srv.Config.Admin.Prefix)
}

func (srv *Server) newModule(id string, prefix string) survana.RequestHandler {

	if len(prefix) == 0 {
		prefix = "/" + id + "/"
	}

	module_dir := srv.Config.Www + id + "/"
	static_dir := module_dir + survana.STATIC_DIR + "/"
	static_prefix := prefix + survana.STATIC_DIR + "/"

	var mod interface{}

	switch id {
	case survana.ADMIN:
		mod = &admin.Module{
			Module: survana.Module{
				Id:           id,
				Prefix:       prefix,
				Dir:          module_dir,
				StaticPrefix: static_prefix,
				StaticDir:    static_dir,
                EnableSessions: true,
                DbSession:    srv.dbSession,
                Db:           srv.dbSession.DB(id),
			},
			Config: srv.Config.Admin,
		}
	}

	handler := mod.(survana.RequestHandler)

	handler.Mount()

	return handler
}

/* Starts a TLS listener on the specified address, using the specified SSL certificate and key */
func (srv *Server) Listen() (err error) {
	log.Printf("Reading SSL certificate (%s) and SSL key (%s)", srv.Config.Sslcert, srv.Config.Sslkey)

	tlsConfig := &tls.Config{}
	//attach SSL certificates to the TLS configuration
	tlsConfig.Certificates = make([]tls.Certificate, 1)
	tlsConfig.Certificates[0], err = tls.LoadX509KeyPair(srv.Config.Sslcert, srv.Config.Sslkey)
	if err != nil {
		return
	}

	//listen on the specified IP and port
	socket, err := net.Listen("tcp", srv.Config.Ip+":"+srv.Config.Port)
	if err != nil {
		return
	}

	//wrap the listener with a TLS listener
	srv.tlsListener = tls.NewListener(socket, tlsConfig)

	return
}

func (srv *Server) Serve() error {
    //close db connection
    defer srv.dbSession.Close()

	return http.Serve(srv.tlsListener, nil)
}
