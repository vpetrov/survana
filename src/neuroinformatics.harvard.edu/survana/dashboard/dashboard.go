package admin

import (
	"log"
	"net/http"
	"neuroinformatics.harvard.edu/survana"
	"code.google.com/p/goauth2/oauth"
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

	app.Get("/", survana.Protect(a.Index))

	//LOGIN
	app.Get("/login", a.Login)
    app.Get("/login/google", a.LoginWithGoogle)
    app.Get("/login/google/response", a.GoogleResponse)
    app.Get("/register", a.Register)

	//LOGOUT
	app.Get("/logout", a.Logout)
}

// displays the index page
func (a *Admin) Index(w http.ResponseWriter, r *survana.Request) {
	a.RenderTemplate(w, "index", nil)
}

// displays the login page
func (a *Admin) Login(w http.ResponseWriter, r *survana.Request) {
    session, err := r.Session()
    if err != nil {
        survana.Error(w, err)
        return
    }

    //no need to log in if already logged in
    if session != nil && session.Authenticated {
        log.Println("already authenticated as", session.UserId)
        survana.Redirect(w, r, "/")
        return
    }

	a.RenderTemplate(w, "login/index", nil)
}

func (a *Admin) LoginWithGoogle(w http.ResponseWriter, r *survana.Request) {

    session, err := r.Session()
    if err != nil {
        survana.Error(w, err)
        return
    }

    if session != nil && session.Authenticated {
        survana.Redirect(w, r, "/")
        return
    }

	config := &oauth.Config{
		ClientId:     "566666928472-gta9d42i4ac9hf4lkndh6g1tdea3umj0.apps.googleusercontent.com",
		ClientSecret: "NOeyzLMyc9BsjhbvFJieC0sg",
		RedirectURL:  "https://localhost:4443/admin/login/google/response",
		Scope:        "email profile",
		AuthURL:      "https://accounts.google.com/o/oauth2/auth",
		TokenURL:     "https://accounts.google.com/o/oauth2/token",
	}

    survana.FullRedirect(w, r, config.AuthCodeURL(""))
}

func (a *Admin) GoogleResponse(w http.ResponseWriter, r *survana.Request) {

    code := r.FormValue("code")

    if len(code) == 0 {
        //redirect to the login page if Google returns an error
        survana.Redirect(w, r, "/login")
        return
    }
    //session_state := r.FormValue("session_state")

	config := &oauth.Config{
		ClientId:     "566666928472-gta9d42i4ac9hf4lkndh6g1tdea3umj0.apps.googleusercontent.com",
		ClientSecret: "NOeyzLMyc9BsjhbvFJieC0sg",
		RedirectURL:  "https://localhost:4443/admin/login/google/response",
		Scope:        "email profile",
		AuthURL:      "https://accounts.google.com/o/oauth2/auth",
		TokenURL:     "https://accounts.google.com/o/oauth2/token",
	}

    transport := &oauth.Transport{Config: config}
    token, err := transport.Exchange(code)
    if err != nil {
        survana.Error(w, err)
        return
    }

    log.Println("token=", token)

    requestUrl := "https://www.googleapis.com/oauth2/v1/userinfo"

    transport.Token = token

    tr, err := transport.Client().Get(requestUrl)
    if err != nil {
        survana.Error(w, err)
        return
    }

    defer tr.Body.Close()

    //Google's response struct (varies based on user's domain)
    user_data := &struct {
                    Name string         `name`
                    Email string        `email`
                    VerifiedEmail bool  `verified_email,omitempty`
                    GivenName string    `given_name,omitempty`
                    FamilyName string   `family_name,omitempty`
                    ProfileUrl string   `json:"link,omitempty" bson:"profile_url,omitempty"`
                    PictureUrl string   `json:"picture,omitempty" bson:"picture_url,omitempty"`
                    Gender string       `gender,omitempty`
                    Locale string       `locale,omitempty`
                    Domain string       `json:"hd,omitempty" bson:"domain,omitempty"`
                 }{}

    err = r.JSONBody(tr.Body, user_data)
    if err != nil {
        survana.Error(w, err)
        return
    }

    log.Printf("%#v", user_data)

    //see if a user with this email exists
    user, err := survana.FindUser(user_data.Email, a.Module.Db)
    if err != nil {
        survana.Error(w, err)
    }

    //if not found, redirect to the login page
    if user == nil {
        if err != nil {
            survana.Error(w, err)
            return
        }

        message := "We couldn't find an account for " + user_data.Email
        tpl_data := &struct { Message string } { Message: message }

        //display the login page, with a message
        a.RenderTemplate(w, "login/error", tpl_data)
        return
    }

    //get an existing session
    session, err := r.Session()
    if err != nil {
        survana.Error(w, err)
        return
    }

	//mark the session as authenticated
	session.Authenticated = true

    //regenerate the session Id
    session.Id = a.Module.Db.UniqueId()

    //set the current user
    session.UserId = user.Id

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

    //redirect to the index page
    survana.Redirect(w, r, "/")
}

func (a *Admin) Register(w http.ResponseWriter, r *survana.Request) {
    a.RenderTemplate(w, "register", nil)
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
