package survana

import (
        "neuroinfo.org/survana/db"
        "net/url"
       )

type Database interface {
    Name() string
    URL() url.URL
    SystemInformation() string
    Version() string

    Connect() error
    Disconnect() error

    FindSession(string, *map[string]string) error
}

type DbObject interface {
    DbId() interface{}
}

const (
        MONGODB = "mongodb"
      )

//factory method to instantiate database drivers based on the ID
func NewDatabase(u url.URL) Database {
    switch (u.Scheme) {
        case MONGODB: return db.NewMongoDB(u)
    }

    return nil
}
