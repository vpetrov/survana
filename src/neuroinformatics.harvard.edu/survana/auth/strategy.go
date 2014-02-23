package auth

import (
        "neuroinformatics.harvard.edu/survana"
        "net/http"
       )

type Strategy interface {
    LoginPage(w http.ResponseWriter, r *survana.Request)
    RegistrationPage(w http.ResponseWriter, r *survana.Request)

    Login(w http.ResponseWriter, r *survana.Request)
    Register(w http.ResponseWriter, r *survana.Request)
    Logout(w http.ResponseWriter, r *survana.Request)
}