package mock

import (
        "net/url"
        "neuroinfo.org/survana/db"
       )

type Database struct {
    Calls map[string]int
    OnFindId func(db.Object)
}

func NewDatabase() *Database {
    return &Database{
        Calls: make(map[string]int,0),
    }
}

func (db *Database) Name() string {
    return "mock"
}

func (db *Database) URL() *url.URL {
    u, _ := url.Parse("mock://localhost:1/mock")
    return u
}

func (db *Database) SystemInformation() string {
    return "Mock DB on SomeOS"
}

func (db *Database) Version() string {
    return "Mock DB v.1.0"
}

func (db *Database) Connect() error {
    return nil
}

func (db *Database) Disconnect() error {
    return nil
}

func (db *Database) FindId(id string, collection string, presult db.Object) error {
    db.Calls["FindId"]++
    if db.OnFindId != nil {
        db.OnFindId(presult)
    }
    return nil
}

func (db *Database) UniqueId() string {
    return "ABCD"
}

func (db *Database) IsValidId(id string) bool {
    return true
}
