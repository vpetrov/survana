package admin

import (
	"labix.org/v2/mgo"
	"net/http"
	"neuroinfo.org/survana"
	"time"
)

const (
	NAME = "admin"
)

//The Admin component
type Admin struct {
	*survana.Module
	mux *survana.RESTMux
}

// creates a new Admin module
func NewModule(path string, dbsession *mgo.Session) *Admin {

	mux := survana.NewRESTMux()

	m := &Admin{
		Module: &survana.Module{
			Name:      NAME,
			Path:      path,
			DbSession: dbsession,
			Db:        dbsession.DB(NAME),
			Router:    mux,
		},
		mux: mux,
	}

	m.ParseTemplates()

	m.RegisterHandlers()

	return m
}

// registers all route handlers
func (a *Admin) RegisterHandlers() {
	app := a.mux

	//must end with slash
	app.Static("/assets/")

	app.Get("/", a.Index)
	app.Get("/home", survana.Protect(a.Home))

	app.Get("/login", a.LoginPage)
	app.Post("/login", a.Login)
}

// displays the index page
func (a *Admin) Index(w http.ResponseWriter, r *survana.Request) {
	a.RenderTemplate(w, "index", nil)
}

// displays the home page
func (a *Admin) Home(w http.ResponseWriter, r *survana.Request) {
	a.RenderTemplate(w, "home", nil)
}

// displays the login page
func (a *Admin) LoginPage(w http.ResponseWriter, r *survana.Request) {
	a.RenderTemplate(w, "login/index", nil)
}

// checks the login details and creates a user session
// returns 204 if the user is already logged in or if the request succeeded
// returns 401 if the key was incorrect
// returns 500 on all other errors
func (a *Admin) Login(w http.ResponseWriter, r *survana.Request) {

	session, err := r.Session()

	if err != nil {
		survana.Error(w, err)
		return
	}

	// if the user was already authenticated, return early
	if session.Authenticated {
		survana.NoContent(w)
		return
	}

	// attempt to read the login details
	form := &struct {
		Key string
	}{}

	//read and decode the form
	err = r.ParseForm(form)

	if err != nil {
		survana.Error(w, err)
		return
	}

	//if the key is incorrect, return bad request
	if form.Key != "secret" {
		survana.BadRequest(w)
		return
	}

	// update the session
	session.Authenticated = true
	err = session.Save()
	if err != nil {
		survana.Error(w, err)
	}

	//set the cookie
	http.SetCookie(w, &http.Cookie{
		Name:     survana.SESSION_ID,
		Value:    session.Id,
		Path:     a.Module.MountPoint,
		Expires:  time.Now().Add(survana.SESSION_TIMEOUT),
		HttpOnly: true,
	})

	//return 204 No Content to indicate success
	survana.NoContent(w)
}
