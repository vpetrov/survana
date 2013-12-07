package survana

import (
        "labix.org/v2/mgo"
        "labix.org/v2/mgo/bson"
        "log"
       )

const (
        SESSION_ID = "SSESSIONID"
        SESSION_COLLECTION = "sessions"
      )

type Session struct {
    db      *mgo.Database
    _id     bson.ObjectId
    Id      string
    Authenticated bool
    Values  map[string]string
}

//creates a new Session object and populates it with values from the dataabse
//some sessions may not be present in the database
func NewSession(db *mgo.Database, id string) (session *Session, err error) {
    session = &Session{
        db: db,
        Id: id,
        Values: make(map[string]string, 0),
    }

    err = db.C(SESSION_COLLECTION).Find(bson.M{"id": session.Id}).One(&session.Values)

    //if the session doesn't exist, create one
    if err == mgo.ErrNotFound {
        log.Println("session",id,"not found")
        err = nil
    } else {
        log.Println("session",id,"found")
    }

    //update convenience properties
    if err == nil {
        //auth status
        if session.Values["authenticated"] == "1" {
            session.Authenticated = true
        } else {
            session.Authenticated = false
        }

        session.Values["id"] = session.Id
    }

    return
}

// Saves the session. Generates a value for _id if it doesn't exist.
func (s *Session) Save() (err error) {
    //sync auth state
    if s.Authenticated == true {
        s.Values["authenticated"] = "1"
    } else {
        s.Values["authenticated"] = "0"
    }

    //sync Id
    s.Values["id"] = s.Id

    log.Println("saving session", s.Id," to collection", SESSION_COLLECTION, " in database", s.db.Name)

    // new sessions won't have valid IDs. Providing an empty/invalid ID to
    // to mgo.UpsertId will cause an error to be returned. Since MongoDB will
    // do this exact same operation and the IDs are unique, we can safely
    // generate an ID here and use mgo.UpsertId(), instead of Insert/Update.
    if !s._id.Valid() {
        s._id = bson.NewObjectId()
    }

    // generating our own ID allows us to use UpsertId here
    _, err = s.db.C(SESSION_COLLECTION).UpsertId(s._id, bson.M{"$set":s.Values})
    if err != nil {
        return
    }

    return
}

func IsValidSessionId(id string) bool {
    return bson.IsObjectIdHex(id)
}

func UniqueId() string {
    return bson.NewObjectId().Hex()
}
