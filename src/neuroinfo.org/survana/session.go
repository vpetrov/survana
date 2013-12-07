package survana

import (
	"errors"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
)

const (
	SESSION_ID         = "SSESSIONID"
	SESSION_COLLECTION = "sessions"
)

var (
	ErrSessionNotFound = errors.New("Session not found")
)

//Represents a user's session. Id and _id are kept separate so that in
//the future, Id's can be regenerated on every request.
//Id and Authenticated are aliases for Values['id'] and Values['authenticated']
type Session struct {
	db            *mgo.Database     //pointer to the datastore
	_id           bson.ObjectId     //the ID in the database
	Id            string            //the publicly visible session id
	Authenticated bool              //whether the user has logged in or not
	Values        map[string]string //all other values go here
}

//creates a new Session object with no Id.
func NewSession(db *mgo.Database) *Session {
	return &Session{
		db:            db,
		Id:            "",
		Authenticated: false,
		Values:        make(map[string]string, 0),
	}
}

//Creates a new session or resumes a previous session.
func CreateSession(db *mgo.Database, id string) (session *Session, err error) {
	//create an empty session object
	session = NewSession(db)
	validId := IsValidSessionId(id)

	if validId {
		//attempt to load an existing session from the database
		err = session.Load(id)

		//if the session was found
		if err == nil {
			return
		}
	}

	//if the session was not found, create a new one
	if err == ErrSessionNotFound || !validId {
		err = nil
		session.Id = UniqueId()
		return
	}

	//otherwise, just return the error
	return
}

// Loads session info from the database
func (s *Session) Load(id string) (err error) {

	err = s.db.C(SESSION_COLLECTION).Find(bson.M{"id": id}).One(&s.Values)

	//if the session doesn't exist, return error
	if err != nil {
		if err == mgo.ErrNotFound {
			err = ErrSessionNotFound
		}

		return
	}

	//auth status
	if s.Values["authenticated"] == "1" {
		s.Authenticated = true
	} else {
		s.Authenticated = false
	}

	s.Values["id"] = id
	s.Id = id

	return
}

// Saves the session. Generates a value for _id if it doesn't exist.
func (s *Session) Save() (err error) {

	if !IsValidSessionId(s.Id) {
		err = errors.New("Invalid session id")
		return
	}

	//sync auth state
	if s.Authenticated == true {
		s.Values["authenticated"] = "1"
	} else {
		s.Values["authenticated"] = "0"
	}

	//sync Id
	s.Values["id"] = s.Id

	log.Println("saving session", s.Id, " to collection", SESSION_COLLECTION, " in database", s.db.Name)

	// new sessions won't have valid IDs. Providing an empty/invalid ID to
	// to mgo.UpsertId will cause an error to be returned. Since MongoDB will
	// do this exact same operation and the IDs are unique, we can safely
	// generate an ID here and use mgo.UpsertId(), instead of Insert/Update.
	if !s._id.Valid() {
		s._id = bson.NewObjectId()
	}

	// generating our own ID allows us to use UpsertId here
	_, err = s.db.C(SESSION_COLLECTION).UpsertId(s._id, bson.M{"$set": s.Values})
	if err != nil {
		return
	}

	return
}

// checks to see whether the session id is valid
func IsValidSessionId(id string) bool {
	return bson.IsObjectIdHex(id)
}

//generates a unique id; safe to use across multiple machines
func UniqueId() string {
	return bson.NewObjectId().Hex()
}