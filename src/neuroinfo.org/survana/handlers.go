package survana

import (
        "net/http"
        "errors"
        "log"
        "runtime"
       )

//A handler that filters all requests that have not been authenticated
//returns 401 Unauthorized if the user's session hasn't been marked as authenticated
func Protect(handler RequestHandler) RequestHandler {
    return func(w http.ResponseWriter, r *Request) {
        //get the session
        session, err := r.Session()

        if err != nil {
            Error(w, err)
        }

        //if the session hasn't been authorized, unauthorized!
        if !session.Authenticated {
            Unauthorized(w, errors.New("not logged in"))
            return
        }

        //must be authenticated at this point
        handler(w, r)
    }
}

//returns 204 No Content
func NoContent(w http.ResponseWriter) {
    w.WriteHeader(http.StatusNoContent)
}

//returns 500 Internal Server Error, and prints the error to the server log
func Error(w http.ResponseWriter, err error) {
    _, file, line, _ := runtime.Caller(1)
    log.Printf("ERROR:%s:%d: %s\n", file, line, err)
    http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

//redirects a request to a path
func Redirect(w http.ResponseWriter, r *Request, redirectPath string) {
    http.Redirect(w, r.Request, r.Module.MountPoint + redirectPath, http.StatusSeeOther)
}

//returns 401 Bad Request
func BadRequest(w http.ResponseWriter) {
    http.Error(w, "Bad Request", http.StatusBadRequest)
}

//returns 401 Unauthorized
func Unauthorized(w http.ResponseWriter, err error) {
    http.Error(w, "Unauthorized: " + err.Error(), http.StatusUnauthorized)
}
