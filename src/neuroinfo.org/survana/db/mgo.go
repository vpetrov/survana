package db

import (
        "labix.org/v2/mgo"
        "labix.org/v2/mgo/bson"
        "net/url"
        "errors"
       )

type MongoDB struct {
    Url url.URL
    Database  *mgo.Database
}

var (
        session *mgo.Session
        sessionInfo mgo.BuildInfo
        ErrNotFound = errors.New("Not found")
    )

func NewMongoDB(u url.URL) *MongoDB {
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

func (db *MongoDB) URL() url.URL {
    return db.Url
}

func (db *MongoDB) SystemInformation() string {
    return sessionInfo.SysInfo
}

func (db *MongoDB) Version() string {
    return sessionInfo.Version
}

func (db *MongoDB) FindId(id string, collection string, result interface{}) (err error) {
	err = db.Database.C(collection).Find(bson.M{"id": id}).One(result)

    if err == mgo.ErrNotFound {
        err = ErrNotFound
    }

    return
}
