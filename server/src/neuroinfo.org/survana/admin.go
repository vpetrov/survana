package survana

import (
        "log"
        "net/http"
        "neuroinfo.org/survana"
       )

type Dashboard struct {
    Module survana.Module
}

func (d *Dashboard) Mount() {
    //::index
    http.HandleFunc(d.Module.Prefix, d.Index)

    d.Module.StaticHandler()
}

//implements RequestHandler
func (d *Module) Index(w http.ResponseWriter, req *http.Request) {
    log.Println(d.Module.Name + " ==> serving request", req.URL.Path)
    w.Write([]byte("hello from admin! my prefix is " + d.Module.Prefix))
}
