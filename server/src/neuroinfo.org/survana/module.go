package survana

import (
    "errors"
    "runtime"
	"log"
	"net/http"
    "html/template"
    "time"
    "labix.org/v2/mgo"
)

type Module struct {
	Id           string
	Prefix       string
	Dir          string
	StaticPrefix string
	StaticDir    string
    EnableSessions bool
    DbSession    *mgo.Session
    Db           *mgo.Database
}

const (
    ADMIN       = "admin"
    DASHBOARD   = "dashboard"
    STUDY       = "study"
    STORE       = "store"

    STATIC_DIR  = "assets"
    TEMPLATE_DIR= "templates"
    TEMPLATE_EXT= ".tpl.html"

    )

type RequestHandler interface {
	Mount()
	Index(w http.ResponseWriter, req *http.Request)
}

//wrapper around http.HandleFunc which prefixes all paths with the module Prefix
func (m *Module) HandleFunc(path string, handler http.HandlerFunc) {
    http.HandleFunc(m.Prefix + path, handler)
}

func (m *Module) StaticHandler() {
    log.Printf("[%s:static]\t%s\t-> %s\n", m.Id, m.StaticPrefix, m.StaticDir)
	//static file handler
	http.Handle(m.StaticPrefix, http.StripPrefix(m.StaticPrefix, http.FileServer(http.Dir(m.StaticDir))))
}

func (m *Module) Protect(handler http.HandlerFunc) http.HandlerFunc {
    return func (w http.ResponseWriter, req *http.Request) {
        log.Println("Checking cookie for URL", req.URL.Path)

        sessionCookie, err := req.Cookie(SESSION_ID)

        if err != nil {
            //if no cookie found, unauthorized!
            Unauthorized(w, errors.New("cookie " + SESSION_ID + ": " + err.Error()))
            return
        }

        session_id := sessionCookie.Value

        //if no session id was found, unauthorized!
        if !IsValidSessionId(session_id) {
             Unauthorized(w, errors.New("Empty or invalid session ID"))
             return
        }

        //if no session information was found, unauthorized!
        session, err := NewSession(m.Db, session_id)

        if err != nil {
            Unauthorized(w, errors.New("Session " + session_id + ": " + err.Error()))
            return
        }

        //if the session hasn't been authorized, unauthorized!
        if !session.Authenticated {
            Unauthorized(w, errors.New("Session " + session_id + ": not logged in"))
            return
        }

        //must be authenticated at this point
        handler(w, req)
    }
}

func (m *Module) Render(name string) http.HandlerFunc {

    log.Println("Registering new function to render template '",name,"'")

    return func (w http.ResponseWriter, req *http.Request) {

        log.Println("[" + m.Id+"] ==> serving auto request", req.URL.Path)

        err := m.RenderTemplate(w, name, nil)

        if err != nil {
            m.Error(w, err)
            return
        }
    }
}

func (m *Module) RenderTemplate(w http.ResponseWriter, name string, data interface{}) (err error) {

    tpl, err := template.ParseFiles(m.Dir + name + TEMPLATE_EXT)
    if err != nil {
        return
    }

    err = tpl.Execute(w, nil)

    return
}

func (m *Module) UserSession(req *http.Request) (session *Session, err error) {

    //get the session cookie
    session_cookie, err := req.Cookie(SESSION_ID)

    // return if an error has occurred (it's ok if the cookie wasn't found)`
    if err != nil && err != http.ErrNoCookie {
        return
    }

    var session_id string

    //either use the session id from the cookie (if it's valid) or generate a new one
    if session_cookie != nil && IsValidSessionId(session_cookie.Value) {
        session_id = session_cookie.Value
    } else {
        session_id = UniqueId()
    }

    // load session by id
    session, err = NewSession(m.Db, session_id)

    if err != nil {
        return
    }

    return
}

func (m *Module) SessionCookie(session *Session) *http.Cookie {

    return &http.Cookie{
        Name:  SESSION_ID,
        Value: session.Id,
        Path: m.Prefix[:len(m.Prefix)-1],
        Expires: time.Now().AddDate(1,0,0),
        Secure: true,
        HttpOnly: true,
    }
}

func Unauthorized(w http.ResponseWriter, err error) {
    log.Println("ERROR:", err)
    http.Error(w, "Unauthorized: " + err.Error(), http.StatusUnauthorized)
}

func (m *Module) Error(w http.ResponseWriter, err error) {
    _, file, line, _ := runtime.Caller(1)
    log.Printf("ERROR:%s:%d: %s\n", file, line, err)
    http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

func (m *Module) Redirect(w http.ResponseWriter, req *http.Request, redirectPath string) {
    http.Redirect(w, req, m.Prefix + redirectPath, http.StatusSeeOther)
}

func (m *Module) BadRequest(w http.ResponseWriter) {
    http.Error(w, "Bad Request", http.StatusBadRequest)
}
