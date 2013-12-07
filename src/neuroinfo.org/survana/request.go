package survana

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

//An interface for any type that can route Survana requests
type Router interface {
	Route(w http.ResponseWriter, r *Request)
}

//An HTTP request to a Survana component.
type RequestHandler func(http.ResponseWriter, *Request)

//A struct that wraps http.Request and provides additional fields for all
//RequestHandlers to use. All methods from http.Request should be promoted
//for use with survana.Request.
type Request struct {
	*http.Request
	URL     *url.URL
	Module  *Module // the module that's handling the request
	session *Session
}

// returns a new Request object
func NewRequest(r *http.Request, path string, module *Module) *Request {
	rurl := &url.URL{
		Scheme:   r.URL.Scheme,
		Opaque:   r.URL.Opaque,
		User:     r.URL.User,
		Host:     r.URL.Host,
		Path:     path,
		RawQuery: r.URL.RawQuery,
		Fragment: r.URL.Fragment,
	}

	return &Request{
		Request: r,
		URL:     rurl,
		Module:  module,
	}
}

// returns either an existing session, or a new session
func (r *Request) Session() (*Session, error) {

	var err error

	//if the session exists already, return it
	if r.session != nil {
		return r.session, nil
	}

	//get the session id cookie, if it exists
	session_id, _ := r.Cookie(SESSION_ID)

	//create a new session. if the session_id is invalid,
	//the function will generate a new id
	r.session, err = CreateSession(r.Module.Db, session_id)

	return r.session, err
}

// returns the value of the cookie by name
func (r *Request) Cookie(name string) (value string, ok bool) {
	ok = false
	cookie, err := r.Request.Cookie(name)

	if cookie != nil && err == nil {
		ok = true
		value = cookie.Value
	}

	return
}

// Parses the request body as a JSON-encoded form
func (r *Request) ParseForm(v interface{}) (err error) {
	// read the body
	body, err := ioutil.ReadAll(r.Request.Body)

	if err != nil {
		return
	}

	if len(body) == 0 {
		err = errors.New("empty request body")
		return
	}

	//parse the JSON body
	err = json.Unmarshal(body, v)

	return
}
