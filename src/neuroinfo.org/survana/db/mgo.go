package db

import (
        "labix.org/v2/mgo"
        "labix.org/v2/mgo/bson"
        "net/url"
        "log"
       )

type MongoDB struct {
    Url *url.URL
    Database  *mgo.Database
}

var (
        session *mgo.Session
        sessionInfo mgo.BuildInfo
    )

func NewMongoDB(u *url.URL) *MongoDB {
    return &MongoDB{
        Url: u,
    }
}

func (db *MongoDB) Connect() (err error) {

    //create a single instance of the session object
    //which will get reused by subsequent calls to Connect()
    if session == nil {
        session, err = mgo.Dial(db.Url.String())
        if err != nil {
            return
        }
    }

    //create an mgo.Database object
    db.Database = session.DB(db.Url.Path[1:])

    //fetch the session information
    sessionInfo, err = session.BuildInfo()

    return
}

func (db *MongoDB) Disconnect() error {
    if session != nil {
        session.Close()
        session = nil
    }

    return nil
}

func (db *MongoDB) Name() string {
    return db.Url.Path[1:]
}

func (db *MongoDB) URL() *url.URL {
    return db.Url
}

func (db *MongoDB) SystemInformation() string {
    return sessionInfo.SysInfo
}

func (db *MongoDB) Version() string {
    return "MongoDB " + sessionInfo.Version
}

func (db *MongoDB) FindId(id string, result Object) (err error) {
    collection := result.Collection()
    if len(collection) == 0 {
        return ErrInvalidCollection
    }

	err = db.Database.C(collection).Find(bson.M{"id": id}).One(result)

    if err == mgo.ErrNotFound {
        err = ErrNotFound
    } else {
        log.Printf("Found id=%v, _id=%#v\n", id, result.DbId())
    }

    return
}

// Stores objects in the database
// new objects won't have valid IDs. Providing an empty/invalid ID to
// to mgo.UpsertId will cause an error to be returned. Since MongoDB will
// do this exact same operation and the IDs are unique, we can safely
// generate an ID here and use mgo.UpsertId(), instead of Insert/Update.
func (db *MongoDB) Save(obj Object) (err error) {
    dbid := obj.DbId()
    var mgoid bson.ObjectId
    var ok bool

    collection := obj.Collection()

    if len(collection) == 0 {
        return ErrInvalidCollection
    }

	if dbid != nil {
        mgoid, ok = dbid.(bson.ObjectId)
        if !ok || !mgoid.Valid() {
            return ErrInvalidId
        }
        //reuse 'ok' as signal to not update the ID
        ok = false
    } else {
		mgoid = bson.NewObjectId()
        //resuse 'ok' as signal to update the ID on successful save
        ok = true
	}

	// generating our own ID allows us to use UpsertId here
	_, err = db.Database.C(collection).UpsertId(mgoid, bson.M{"$set": obj})
	if err != nil {
		return
	}

    //update the dbid if necessary
    if ok {
        obj.SetDbId(mgoid)
    }

    return
}

//TODO: stub
func (db *MongoDB) Delete(obj Object) (err error) {
    //get the interface{} value
    dbid := obj.DbId()
    if dbid == nil {
        return
    }

    //get the object's collection. ignore if it's invalid
    collection := obj.Collection()
    if len(collection) == 0 {
        return
    }

    //remove the object by id
    err = db.Database.C(collection).RemoveId(dbid)
    return
}

//generates a unique id; safe to use across multiple machines
func (db *MongoDB) UniqueId() string {
	return bson.NewObjectId().Hex()
}

func (db *MongoDB) IsValidId(id string) bool {
    return bson.IsObjectIdHex(id)
}
