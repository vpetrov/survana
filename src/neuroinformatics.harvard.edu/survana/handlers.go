package survana

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime"
    "runtime/debug"
)

type JSONResponse struct {
	Success bool        `json:"success"`
	Message interface{} `json:"message,omitempty"`
}

func NotLoggedIn(handler RequestHandler) RequestHandler {
	return func(w http.ResponseWriter, r *Request) {
		//get the session
		session, err := r.Session()
		if err != nil {
			Error(w, err)
			return
		}

		//if the session has already been authorized, redirect
		if session.Authenticated {
			Redirect(w, r, "/")
			return
		}

		//must not be authenticated at this point
		handler(w, r)
    }
}

func NotFound(w http.ResponseWriter) {
	http.Error(w, "Not Found", http.StatusNotFound)
}

//returns 204 No Content
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

//returns 500 Internal Server Error, and prints the error to the server log
func Error(w http.ResponseWriter, err error) {
    debug.PrintStack()
	_, file, line, _ := runtime.Caller(1)
	log.Printf("ERROR:%s:%d: %s\n", file, line, err)
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

//redirects a request to a path relative to the module
//if the request is an
func Redirect(w http.ResponseWriter, r *Request, redirectPath string) {
	FullRedirect(w, r, r.Module.MountPoint+redirectPath)
}

//redirects an external resource
func FullRedirect(w http.ResponseWriter, r *Request, url string) {
	if r.Request.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		XHRRedirect(w, r, url)
	} else {
		http.Redirect(w, r.Request, url, http.StatusSeeOther)
	}
}

//redirects to a full URL
func XHRRedirect(w http.ResponseWriter, r *Request, url string) {
	data := &struct {
		Redirect string `json:"redirect,omitempty"`
	}{
		Redirect: url,
	}

    w.Header().Set("X-Survana-Redirect", url)

	JSONResult(w, false, data)
}

//sends a { 'success': <bool>, 'message': <custom data> } response to the client
func JSONResult(w http.ResponseWriter, success bool, data interface{}) {

	result := &JSONResponse{
		Success: success,
		Message: data,
	}

	jsondata, err := json.Marshal(result)
	if err != nil {
		Error(w, err)
		return
	}

	w.Write(jsondata)
}

//returns 401 Bad Request
func BadRequest(w http.ResponseWriter) {
	http.Error(w, "Bad Request", http.StatusBadRequest)
}

//returns 401 Unauthorized
func Unauthorized(w http.ResponseWriter, err error) {
	http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
}
