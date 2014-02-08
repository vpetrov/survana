package dashboard

import (
    "log"
	"net/http"
	"neuroinformatics.harvard.edu/survana"
    "crypto/sha512"
    "crypto/rand"
    "encoding/base64"
    "bytes"
    )

const (
	BUILTIN_USER_COLLECTION = "builtin_users"
    SALT_ENTROPY = 3
)


type builtinUser struct {
    Id          string  `bson:"id,omitempty"`
    Password    []byte  `bson:"password,omitempty"`
    Salt        string  `bson:"salt,omitempty"`
    UserId      string  `bson:"user_id,omitempty"`

    /* DbObject */
    DBID        interface{} `bson:"_id,omitempty"`
}

func (d *Dashboard) BuiltinAuthPage(w http.ResponseWriter, r *survana.Request) {
    d.RenderTemplate(w, "auth/builtin/login", d)
}

func (d *Dashboard) BuiltinRegistrationPage(w http.ResponseWriter, r *survana.Request) {
    d.RenderTemplate(w, "auth/builtin/register", d)

}

func (d *Dashboard) BuiltinAuth(w http.ResponseWriter, r *survana.Request) {

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

    if !ok1 || !ok2 || len(username) == 0 || len(password) == 0 {
        survana.BadRequest(w)
        return
    }

    //find this user in the built-in user database
    bUser, err := findBuiltinUser(username, d.Module.Db)
    if err != nil {
        survana.Error(w, err)
        return
    }

    //not found?
    if bUser == nil {
        log.Println("No such builtin user: %v", username)
        survana.JSONResult(w, false, "Invalid username or password")
        return
    }

    log.Println("password=%v salt=%v", password, bUser.Salt)

    sha512_password := hash(password, bUser.Salt)

    log.Printf("sha512_password=%v", sha512_password)
    log.Printf("  user_password=%v", bUser.Password)


    //wrong password?
    if !bytes.Equal(sha512_password, bUser.Password) {
        log.Println("hashes are not equal :/")
        survana.JSONResult(w, false, "Invalid username or password")
        return
    }

    //success
    survana.JSONResult(w, true, d.Module.MountPoint + "/")
}

func (d *Dashboard) BuiltinRegister(w http.ResponseWriter, r *survana.Request) {

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

    if !ok1 || !ok2 || len(username) == 0 || len(password) == 0 {
        survana.JSONResult(w, false, "Please complete all fields")
        return
    }

    //hash the password 
    password_salt := generateRandomSalt(SALT_ENTROPY)
    password_hash := hash(password, password_salt)

    bUser := newBuiltinUser(username, password_hash, password_salt)
    err = d.Module.Db.Save(bUser)
    if err != nil {
        survana.Error(w, err)
        return
    }

    survana.JSONResult(w, true, d.Module.MountPoint + "/")
}

func newBuiltinUser(username string, password []byte, password_salt string) *builtinUser {
	return &builtinUser{
		Id:   username,
		Password: password,
        Salt: password_salt,
	}
}

func (u *builtinUser) DbId() interface{} {
	return u.DBID
}

func (u *builtinUser) SetDbId(v interface{}) {
	u.DBID = v
}

func (u *builtinUser) Collection() string {
	return BUILTIN_USER_COLLECTION
}


func findBuiltinUser(username string, db survana.Database) (user *builtinUser, err error) {
	user = &builtinUser{}
	err = db.FindId(username, user)

	if err != nil {
		if err == survana.ErrNotFound {
			err = nil
		}

		return nil, err
	}

	return
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
