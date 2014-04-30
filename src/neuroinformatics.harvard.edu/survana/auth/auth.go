package auth

import (
        "log"
        "time"
        "crypto/sha512"
        "crypto/rand"
        "encoding/base64"
        "net/http"
        "neuroinformatics.harvard.edu/survana"
       )

const (
        BUILTIN = "built-in"
        SUBJECT_ID = "subject_id"
        OAUTH2  = "oauth2"
        LDAP    = "ldap"
        NIS     = "nis"
      )

const (
        LOGIN_PATH = "/login"
      )

func New(config *Config) Strategy {

    switch config.Type {
        case BUILTIN: return NewBuiltinStrategy(config)
        case SUBJECT_ID: return NewSubjectIdStrategy(config)
        default: log.Printf("WARNING: Authentication type '%v' is not yet supported.", config.Type)
    }

    return nil
}

func Login(w http.ResponseWriter, r *survana.Request) {
    //get the session
	session, err := r.Session()
	if err != nil {
		survana.Error(w, err)
		return
	}

    //if the user is already authenticated, redirect to home
	if session != nil && session.Authenticated {
		survana.Redirect(w, r, "/")
		return
	}

    //TODO: make this work! (and delete the hard-coded built-in strategy!)
    //user, err := r.Module.Auth.Login(w, r);
    bauth := NewBuiltinStrategy(&Config{Type:BUILTIN})
    profile_id, err := bauth.Login(w, r);

    if err != nil {
        log.Println("login error:", err)
        survana.JSONResult(w, false, err.Error())
        return
    }

	//mark the session as authenticated
	session.Authenticated = true

	//regenerate the session Id
	session.Id = r.Module.Db.UniqueId()

	//set the current user
	session.UserId = profile_id

	// update the session
	err = session.Save(r.Module.Db)
	if err != nil {
		survana.Error(w, err)
		return
	}

	//set the cookie
	http.SetCookie(w, &http.Cookie{
		Name:     survana.SESSION_ID,
		Value:    session.Id,
		Path:     r.Module.MountPoint,
		Expires:  time.Now().Add(survana.SESSION_TIMEOUT),
		Secure:   true,
		HttpOnly: true,
	})

    //success
    survana.JSONResult(w, true, r.Module.MountPoint + "/")
}

//Logs out a user.
//returns 204 No Content on success
//returns 500 Internal Server Error on failure
func logout(w http.ResponseWriter, r *survana.Request) {
	session, err := r.Session()

	if err != nil {
		survana.Error(w, err)
		return
	}

	if !session.Authenticated {
		survana.NoContent(w)
		return
	}

	err = session.Delete(r.Module.Db)
	if err != nil {
		survana.Error(w, err)
		return
	}

	//To delete the cookie, we set its value to some bogus string,
	//and the expiration to one second past the beginning of unix time.
	http.SetCookie(w, &http.Cookie{
		Name:     survana.SESSION_ID,
		Value:    "Homer",
		Path:     r.Module.MountPoint,
		Expires:  time.Unix(1, 0),
		Secure:   true,
		HttpOnly: true,
	})

	//return 204 No Content on success
	survana.Redirect(w, r, "/")

	//note that the user has logged out
	go r.Module.Log.Printf("logout")
}

func generateRandomSalt(nbytes int) string {
    //allocate a slice of nbytes bytes
    random_bytes := make([]byte, nbytes)
    //read random data from a crypto rng
    rand.Read(random_bytes)

    return base64.StdEncoding.EncodeToString(random_bytes)
}

//salt prevents rainbow attacks. We insert some characters to block trivial attempts
//at guessing the salting function. That said, the attack could just look at the source
//code to find this function, but we still try to make it as hard as possible to guess it.
func salt(password, salt string) string {
    return "$" + salt + "::" + password + "::" + salt + "$"
}

func hash(password, password_salt string) []byte {
    //salt the password string
    salted_password := salt(password, password_salt)

    hash := sha512.New()

    return hash.Sum([]byte(salted_password))
}

//A handler that filters all requests that have not been authenticated
//returns 401 Unauthorized if the user's session hasn't been marked as authenticated
func Protect(handler survana.RequestHandler) survana.RequestHandler {
	return func(w http.ResponseWriter, r *survana.Request) {
		//get the session
		session, err := r.Session()

		if err != nil {
			survana.Error(w, err)
			return
		}

		//if the session hasn't been authorized, redirect
		if !session.Authenticated {
            redirect_path := LOGIN_PATH

            w.Header().Set("X-Session-Expired", "1")

            //forward query params
            if len(r.URL.RawQuery) > 0 {
                redirect_path += "?" + r.URL.RawQuery
            }

			survana.Redirect(w, r, redirect_path)
			return
		}

		//must be authenticated at this point
		handler(w, r)
	}
}

