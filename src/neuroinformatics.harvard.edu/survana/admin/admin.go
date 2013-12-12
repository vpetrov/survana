package admin

import (
	_ "log"
	"net/http"
	"neuroinformatics.harvard.edu/survana"
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
func NewModule(path string, db survana.Database) *Admin {

	mux := survana.NewRESTMux()

	m := &Admin{
		Module: &survana.Module{
			Name:   NAME,
			Path:   path,
			Db:     db,
			Router: mux,
			Log:    db.NewLogger("logs", NAME),
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

	//LOGIN
	app.Get("/login", a.LoginPage)
	app.Post("/login", a.Login)

	//LOGOUT
	app.Get("/logout", a.Logout)
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
	} else {
		//    user := survana.NewUser("victor.petrov@survana.org", "Victor Petrov")
		//    user.Login()
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

	//mark the session as authenticated
	session.Authenticated = true

	// update the session
	err = session.Save(a.Module.Db)
	if err != nil {
		survana.Error(w, err)
		return
	}

	//set the cookie
	http.SetCookie(w, &http.Cookie{
		Name:     survana.SESSION_ID,
		Value:    session.Id,
		Path:     a.Module.MountPoint,
		Expires:  time.Now().Add(survana.SESSION_TIMEOUT),
		Secure:   true,
		HttpOnly: true,
	})

	//return 204 No Content to indicate success
	survana.NoContent(w)
}

//Logs out a user.
//returns 204 No Content on success
//returns 500 Internal Server Error on failure
func (a *Admin) Logout(w http.ResponseWriter, r *survana.Request) {
	session, err := r.Session()

	if err != nil {
		survana.Error(w, err)
		return
	}

	if !session.Authenticated {
		survana.NoContent(w)
		return
	}

	err = session.Delete(a.Module.Db)
	if err != nil {
		survana.Error(w, err)
		return
	}

	//To delete the cookie, we set its value to some bogus string,
	//and the expiration to one second past the beginning of unix time.
	http.SetCookie(w, &http.Cookie{
		Name:     survana.SESSION_ID,
		Value:    "Homer",
		Path:     a.Module.MountPoint,
		Expires:  time.Unix(1, 0),
		Secure:   true,
		HttpOnly: true,
	})

	//return 204 No Content on success
	survana.NoContent(w)

	//note that the user has logged out
	go a.Module.Log.Printf("logout")
}
