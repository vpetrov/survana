package survana

import (
        "neuroinfo.org/survana/db"
        "net/url"
        "log"
       )

type Database interface {
    Name() string
    URL() *url.URL
    SystemInformation() string
    Version() string

    Connect() error
    Disconnect() error

    FindId(id string, presult db.Object) error
    Delete(o db.Object) error
    Save(o db.Object) error

    UniqueId() string
    IsValidId(id string) bool

    NewLogger(collection string, prefix string) *log.Logger
}

const (
        MONGODB = "mongodb"
      )

//factory method to instantiate database drivers based on the ID
func NewDatabase(u *url.URL) Database {
    switch (u.Scheme) {
        case MONGODB: return db.NewMongoDB(u)
    }

    return nil
}
