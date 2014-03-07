package auth

import (
    "log"
	"net/http"
	"neuroinformatics.harvard.edu/survana"
    "bytes"
    "time"
    )

const (
	BUILTIN_USER_COLLECTION = "builtin_users"
    SALT_ENTROPY = 3
)

type builtinUser struct {
    survana.DBO         `bson:",inline,omitempty"`
    Id          string  `bson:"id,omitempty"`
    Password    []byte  `bson:"password,omitempty"`
    Salt        string  `bson:"salt,omitempty"`
    UserId      string  `bson:"user_id,omitempty"`
}

type BuiltinStrategy struct {
    Config *Config
}

func NewBuiltinStrategy(config *Config) BuiltinStrategy {
    return BuiltinStrategy{
        Config: config,
    }
}

func (b BuiltinStrategy) Attach(module *survana.Module) {
    app := module.Mux

    app.Get("/login", survana.NotLoggedIn(b.LoginPage))
    app.Post("/login", survana.NotLoggedIn(b.Login))

    //Registration is optional
    if b.Config.AllowRegistration {
        app.Get("/register", survana.NotLoggedIn(b.RegistrationPage))
        app.Post("/register", survana.NotLoggedIn(b.Register))
    }
}

func (b BuiltinStrategy) LoginPage(w http.ResponseWriter, r *survana.Request) {
    r.Module.RenderTemplate(w, "auth/builtin/login", b.Config)
}

func (b BuiltinStrategy) RegistrationPage(w http.ResponseWriter, r *survana.Request) {
    r.Module.RenderTemplate(w, "auth/builtin/register", nil)
}

func (b BuiltinStrategy) Login(w http.ResponseWriter, r *survana.Request) {

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

    //this is why each strategy needs to be able to render its
    //login screens, so that it can ask for custom fields.
    //here we have a simple username/password combo, but the
    //other strategies could show various options based on the
    //auth configuration
    data := make(map[string]string)

    err = r.ParseJSON(&data)
    if err != nil {
        survana.Error(w, err)
        return
    }

    username, ok1 := data["username"]
    password, ok2 := data["password"]

    if !ok1 || !ok2 || len(username) == 0 || len(password) == 0 {
        survana.BadRequest(w)
        return
    }

    //find this user in the built-in user database
    bUser, err := findBuiltinUser(username, r.Module.Db)
    if err != nil {
        survana.Error(w, err)
        return
    }

    //not found?
    if bUser == nil {
        log.Printf("No such builtin user: %v", username)
        survana.JSONResult(w, false, "Invalid username or password")
        return
    }

    log.Printf("password=%v salt=%v", password, bUser.Salt)

    sha512_password := hash(password, bUser.Salt)

    log.Printf("sha512_password=%v", sha512_password)
    log.Printf("  user_password=%v", bUser.Password)


    //wrong password?
    if !bytes.Equal(sha512_password, bUser.Password) {
        log.Println("hashes are not equal :/")
        survana.JSONResult(w, false, "Invalid username or password")
        return
    }

    /* TODO: the code below sets the authentication cookie.
       This should be moved to auth.Login(), which calls strategy.Login()
       and then sets this cookie.
   */

	//mark the session as authenticated
	session.Authenticated = true

	//regenerate the session Id
	session.Id = r.Module.Db.UniqueId()

	//set the current user
	session.UserId = bUser.UserId

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

func (b BuiltinStrategy) Register(w http.ResponseWriter, r *survana.Request) {

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

    data := make(map[string]string)

    err = r.ParseJSON(&data)
    if err != nil {
        survana.Error(w, err)
        return
    }

    username, ok1 := data["username"]
    password, ok2 := data["password"]
    name, ok3 := data["name"]
    email, ok4 := data["email"]

    //TODO: this needs to be refactored into something better
    if !ok1 || !ok2 || !ok3 || !ok4 || len(username) == 0 || len(password) == 0 || len(name) == 0 || len(email) == 0 {
        survana.JSONResult(w, false, "Please complete all fields")
        return
    }

    //query the database to check if the username exists
    username_exists, err := r.Module.Db.HasId(username, BUILTIN_USER_COLLECTION)
    if err != nil {
        survana.Error(w, err)
        return
    }

    //make sure users can't register duplicate usernames
    if username_exists {
        survana.JSONResult(w, false, "This username already exists")
        return
    }

    //hash the password 
    password_salt := generateRandomSalt(SALT_ENTROPY)
    password_hash := hash(password, password_salt)

    //create a Survana user (profile)
    user := survana.NewUser(email, name)
    user.AuthType = BUILTIN;
    err = r.Module.Db.Save(user)
    if err != nil {
        survana.Error(w, err)
        return
    }

    //create an entry to store auth details
    bUser := newBuiltinUser(username, password_hash, password_salt)
    bUser.UserId = user.Id
    err = r.Module.Db.Save(bUser)
    if err != nil {
        survana.Error(w, err)
        return
    }

    survana.JSONResult(w, true, r.Module.MountPoint + "/")
}

//default logout
func (b BuiltinStrategy) Logout(w http.ResponseWriter, r *survana.Request) {
    logout(w, r)
}

func newBuiltinUser(username string, password []byte, password_salt string) *builtinUser {
	return &builtinUser{
        DBO: survana.DBO { Collection: BUILTIN_USER_COLLECTION },
		Id:   username,
		Password: password,
        Salt: password_salt,
	}
}

func emptyBuiltinUser() *builtinUser {
    return &builtinUser {
        DBO: survana.DBO { Collection: BUILTIN_USER_COLLECTION },
    }
}

func findBuiltinUser(username string, db survana.Database) (user *builtinUser, err error) {
	user = emptyBuiltinUser()
	err = db.FindId(username, user)

	if err != nil {
		if err == survana.ErrNotFound {
			err = nil
		}

		return nil, err
	}

	return
}

