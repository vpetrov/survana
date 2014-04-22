package auth

import (
    "log"
	"net/http"
	"neuroinformatics.harvard.edu/survana"
    "bytes"
    "errors"
    )

const (
	BUILTIN_USER_COLLECTION = "builtin_users"
    SALT_ENTROPY = 3
    BERR_INVALID_CREDENTIALS = "Invalid username or password"

    ADMIN_USERNAME = "admin"
    ADMIN_PASSWORD = "admin"
    ADMIN_NAME = "Administrator"
    ADMIN_EMAIL = "root@localhost"
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
    app.Post("/login", survana.NotLoggedIn(Login))

    //Registration is optional
    if b.Config.AllowRegistration {
        app.Get("/register", survana.NotLoggedIn(b.RegistrationPage))
        app.Post("/register", survana.NotLoggedIn(b.Register))
    }

    //setup the admin account
    setupAdminAccount(module)
}

func setupAdminAccount(module *survana.Module) {
    var (
            admin_user      *builtinUser
            admin_profile   *survana.User
            err             error
            db              survana.Database = module.Db
        )

    admin_user, err = findBuiltinUser(ADMIN_USERNAME, db)
    if err != nil {
        log.Fatal(err)
        return
    }

    //create a new user if it doesn't exist, then create Survana profile 
    if admin_user == nil {
        //create builtin user
        password_salt := generateRandomSalt(SALT_ENTROPY)
        password_hash := hash(ADMIN_PASSWORD, password_salt)
        admin_user = newBuiltinUser(ADMIN_USERNAME, password_hash, password_salt)

        //create profile
        admin_profile = survana.NewUser(ADMIN_EMAIL, ADMIN_NAME)
        admin_profile.AuthType = BUILTIN;

        //link the builtin user and Survana profile
        admin_user.UserId = admin_profile.Id

        //save the admin user
        err = db.Save(admin_user)
        if err != nil {
            log.Fatal(err)
            return
        }

        //save the admin profile
        err = db.Save(admin_profile)
        if err != nil {
            log.Fatal(err)
            return
        }
    } 
}

func (b BuiltinStrategy) LoginPage(w http.ResponseWriter, r *survana.Request) {
    r.Module.RenderTemplate(w, "auth/builtin/login", b.Config)
}

func (b BuiltinStrategy) RegistrationPage(w http.ResponseWriter, r *survana.Request) {
    r.Module.RenderTemplate(w, "auth/builtin/register", nil)
}

func (b BuiltinStrategy) Login(w http.ResponseWriter, r *survana.Request) (profile_id string, err error) {

    //this is why each strategy needs to be able to render its
    //login screens, so that it can ask for custom fields.
    //here we have a simple username/password combo, but the
    //other strategies could show various options based on the
    //auth configuration
    data := make(map[string]string)

    err = r.ParseJSON(&data)
    if err != nil {
        log.Println(err)
        return
    }

    username, ok1 := data["username"]
    password, ok2 := data["password"]

    if !ok1 || !ok2 || len(username) == 0 || len(password) == 0 {
        err = errors.New("Invalid request")
        return
    }

    //find this user in the built-in user database
    user, err := findBuiltinUser(username, r.Module.Db)
    if err != nil {
        return
    }

    //not found?
    if user == nil {
        log.Printf("No such builtin user: %v", username)
        err = errors.New("Invalid username or password")
        return
    }

    sha512_password := hash(password, user.Salt)

    //wrong password?
    if !bytes.Equal(sha512_password, user.Password) {
        err = errors.New(BERR_INVALID_CREDENTIALS)
        return
    }

    return user.UserId, nil
}

func createBuiltinProfile(username, password, name, email string, db survana.Database) (user *builtinUser, profile *survana.User, err error) {
    //query the database to check if the username exists
    username_exists, err := db.HasId(username, BUILTIN_USER_COLLECTION)
    if err != nil {
        log.Println(err)
        return
    }

    //make sure users can't register duplicate usernames
    if username_exists {
        err = errors.New("This username already exists")
        return
    }

    //create a Survana user (profile)
    profile = survana.NewUser(email, name)
    profile.AuthType = BUILTIN;
    err = db.Save(profile)
    if err != nil {
        log.Println(err)
        return
    }

    //hash the password 
    password_salt := generateRandomSalt(SALT_ENTROPY)
    password_hash := hash(password, password_salt)

    //create an entry to store auth details
    user = newBuiltinUser(username, password_hash, password_salt)
    user.UserId = user.Id //TODO: change to ProfileId
    err = db.Save(user)
    if err != nil {
        log.Println(err)
        return
    }

    return user, profile, nil
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

    _, _, err = createBuiltinProfile(username, password, email, name, r.Module.Db)
    if err != nil {
        survana.JSONResult(w, true, r.Module.MountPoint + "/")
        return
    }

    survana.JSONResult(w, false, err)
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

