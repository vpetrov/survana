package admin

import (
    "errors"
	"log"
    "encoding/json"
    "io/ioutil"
	"net/http"
	"neuroinfo.org/survana"
)

type Module struct {
	Module survana.Module
	Config
}

func (d *Module) Mount() {

    d.Module.HandleFunc("",         d.Module.Render("index"))
    d.Module.HandleFunc("login",    d.Module.Render("login/index"))
    d.Module.HandleFunc("home",     d.Module.Protect(d.Module.Render("home")))
    d.Module.HandleFunc("session",  d.Session)

	d.Module.StaticHandler()
}

func (d *Module) Session(w http.ResponseWriter, req *http.Request) {
    //expected POST fields
    type POST struct {
        Key string
    }

    // fetch the session
    session, err := d.Module.UserSession(req)

    if err != nil {
        d.Module.Error(w, err)
    }

    // if the session is authenticated already, return success
    if session.Authenticated {
        w.WriteHeader(http.StatusNoContent)
        return
    }

    //otherwise, decode the response body and see what the user provided
    body, err := ioutil.ReadAll(req.Body)
    if err != nil {
        d.Module.Error(w, err)
        return
    }

    //verify that the body contains data
    if len(body) == 0 {
        d.Module.BadRequest(w)
        return
    }

    //space for user provided values
    var form *POST = &POST{}

    err = json.Unmarshal(body, form)
    if err != nil {
        d.Module.Error(w, err)
        return
    }

    // if the key is invalid, return HTTP Unauthorized
    if form.Key != d.Config.Key {
        survana.Unauthorized(w, errors.New("Invalid Key"))
        return
    } else {
        // mark the session as authenticated
        session.Authenticated = true
    }

    // save the session
    err = session.Save()
    if err != nil {
        d.Module.Error(w, err)
        return
    }

    //set cookie header
    http.SetCookie(w, d.Module.SessionCookie(session))
    w.WriteHeader(http.StatusNoContent)
}
