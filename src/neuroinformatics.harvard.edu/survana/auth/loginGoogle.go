package auth
/*
import (
	"code.google.com/p/goauth2/oauth"
	"log"
	"net/http"
	"neuroinformatics.harvard.edu/survana"
    "time"
    )
func (d *Dashboard) LoginWithGoogle(w http.ResponseWriter, r *survana.Request) {

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
		RedirectURL:  "https://localhost:4443/dashboard/login/google/response",
		Scope:        "email profile",
		AuthURL:      "https://accounts.google.com/o/oauth2/auth",
		TokenURL:     "https://accounts.google.com/o/oauth2/token",
	}

	survana.FullRedirect(w, r, config.AuthCodeURL(""))
}

func (d *Dashboard) GoogleResponse(w http.ResponseWriter, r *survana.Request) {

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
		RedirectURL:  "https://localhost:4443/dashboard/login/google/response",
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
		Name          string `name`
		Email         string `email`
		VerifiedEmail bool   `verified_email,omitempty`
		GivenName     string `given_name,omitempty`
		FamilyName    string `family_name,omitempty`
		ProfileUrl    string `json:"link,omitempty" bson:"profile_url,omitempty"`
		PictureUrl    string `json:"picture,omitempty" bson:"picture_url,omitempty"`
		Gender        string `gender,omitempty`
		Locale        string `locale,omitempty`
		Domain        string `json:"hd,omitempty" bson:"domain,omitempty"`
	}{}

	err = r.JSONBody(tr.Body, user_data)
	if err != nil {
		survana.Error(w, err)
		return
	}

	log.Printf("%#v", user_data)

	//see if a user with this email exists
	user, err := survana.FindUser(user_data.Email, d.Module.Db)
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
		tpl_data := &struct{ Message string }{Message: message}

		//display the login page, with a message
		d.RenderTemplate(w, "login/error", tpl_data)
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
	session.Id = d.Module.Db.UniqueId()

	//set the current user
	session.UserId = user.Id

	// update the session
	err = session.Save(d.Module.Db)
	if err != nil {
		survana.Error(w, err)
		return
	}

	//set the cookie
	http.SetCookie(w, &http.Cookie{
		Name:     survana.SESSION_ID,
		Value:    session.Id,
		Path:     d.Module.MountPoint,
		Expires:  time.Now().Add(survana.SESSION_TIMEOUT),
		Secure:   true,
		HttpOnly: true,
	})

	//redirect to the index page
	survana.Redirect(w, r, "/")
}

func (d *Dashboard) Register(w http.ResponseWriter, r *survana.Request) {
	d.RenderTemplate(w, "register", nil)
}
*/
