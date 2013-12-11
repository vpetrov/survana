package mock

import (
	"log"
	"net/url"
	"neuroinfo.org/survana/db"
	"os"
)

type Database struct {
	Calls    map[string]int
	OnFindId func(db.Object)
	OnDelete func(db.Object)
	OnSave   func(db.Object)
}

func NewDatabase() *Database {
	return &Database{
		Calls: make(map[string]int, 0),
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

func (db *Database) FindId(id string, presult db.Object) error {
	db.Calls["FindId"]++
	if db.OnFindId != nil {
		db.OnFindId(presult)
	}
	return nil
}

func (db *Database) Delete(o db.Object) error {
	db.Calls["Delete"]++
	if db.OnDelete != nil {
		db.OnDelete(o)
	}

	return nil
}

func (db *Database) Save(o db.Object) error {
	db.Calls["Save"]++
	if db.OnSave != nil {
		db.OnSave(o)
	}

	return nil
}

func (db *Database) UniqueId() string {
	return "ABCD"
}

func (db *Database) IsValidId(id string) bool {
	return true
}

func (db *Database) NewLogger(collection, prefix string) *log.Logger {
	return log.New(os.Stdout, "mock", 0)
}
