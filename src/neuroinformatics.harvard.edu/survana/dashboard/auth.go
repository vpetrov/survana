package dashboard

import (
	"neuroinformatics.harvard.edu/survana"
    "net/http"
    "log"
	"time"
)

// displays the login page
func (d *Dashboard) LoginPage(w http.ResponseWriter, r *survana.Request) {
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

	d.RenderTemplate(w, "login/index", nil)
}


//Logs out a user.
//returns 204 No Content on success
//returns 500 Internal Server Error on failure
func (d *Dashboard) Logout(w http.ResponseWriter, r *survana.Request) {
	session, err := r.Session()

	if err != nil {
		survana.Error(w, err)
		return
	}

	if !session.Authenticated {
		survana.NoContent(w)
		return
	}

	err = session.Delete(d.Module.Db)
	if err != nil {
		survana.Error(w, err)
		return
	}

	//To delete the cookie, we set its value to some bogus string,
	//and the expiration to one second past the beginning of unix time.
	http.SetCookie(w, &http.Cookie{
		Name:     survana.SESSION_ID,
		Value:    "Homer",
		Path:     d.Module.MountPoint,
		Expires:  time.Unix(1, 0),
		Secure:   true,
		HttpOnly: true,
	})

	//return 204 No Content on success
	survana.NoContent(w)

	//note that the user has logged out
	go d.Module.Log.Printf("logout")
}
