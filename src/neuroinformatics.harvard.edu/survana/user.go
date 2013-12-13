package survana

import (
	"code.google.com/p/goauth2/oauth"
	"log"
)

const (
	USER_COLLECTION = "users"
)

type User struct {
	Id   string
	Name string

	//DbObject
	DBID interface{} `bson:"_id,omitempty"`
}

func NewUser(email, name string) *User {
	return &User{
		Id:   email,
		Name: name,
	}
}

func (u *User) DbId() interface{} {
	return u.DBID
}

func (u *User) SetDbId(v interface{}) {
	u.DBID = v
}

func (u *User) Collection() string {
	return USER_COLLECTION
}

func FindUser(email string, db Database) (user *User, err error) {
	user = &User{}
	err = db.FindId(email, user)

	if err != nil {
		if err == ErrNotFound {
			err = nil
		}

		return nil, err
	}

	return
}

func (u *User) Login() {
	// Set up a configuration.
	config := &oauth.Config{
		ClientId:     "566666928472-gta9d42i4ac9hf4lkndh6g1tdea3umj0.apps.googleusercontent.com",
		ClientSecret: "NOeyzLMyc9BsjhbvFJieC0sg",
		RedirectURL:  "https://localhost:4443/admin/oauth-response",
		Scope:        "openid email",
		AuthURL:      "https://accounts.google.com/o/oauth2/auth",
		TokenURL:     "https://accounts.google.com/o/oauth2/token",
		TokenCache:   oauth.CacheFile("cache.json"),
	}

	_ = &oauth.Transport{Config: config}

	_, err := config.TokenCache.Token()

	if err != nil {
		log.Println("Visit this URL:", config.AuthCodeURL(""))
	}
}
