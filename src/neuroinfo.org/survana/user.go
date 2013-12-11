package survana

import (
	"code.google.com/p/goauth2/oauth"
	"log"
)

type User struct {
	Email string
	Name  string
}

func NewUser(email, name string) *User {
	return &User{
		Email: email,
		Name:  name,
	}
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
