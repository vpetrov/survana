package auth

import (
        "neuroinformatics.harvard.edu/survana"
        "net/http"
       )

type Strategy interface {
    Attach(module *survana.Module)

    LoginPage(w http.ResponseWriter, r *survana.Request)
    RegistrationPage(w http.ResponseWriter, r *survana.Request)

    Login(w http.ResponseWriter, r *survana.Request) (profile_id string, err error)
    Register(w http.ResponseWriter, r *survana.Request)
    Logout(w http.ResponseWriter, r *survana.Request)
}
