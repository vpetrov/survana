package survana

import (
        "log"
        "net/http"
       )

type Module struct {
    Name string
    Prefix string
    Dir string
    StaticPrefix string
    StaticDir string
}

const STATIC_DIR = "static"

type RequestHandler interface {
    Mount()
    Index(w http.ResponseWriter, req *http.Request)
}

func (m *Module) StaticHandler() {
    //static file handler
    http.Handle(m.StaticPrefix, http.StripPrefix(m.StaticPrefix, http.FileServer(http.Dir(m.StaticDir))))
}

func (m *Module) Info() {
    log.Println("WOAH!")
}
