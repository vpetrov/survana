package survana

import (
	"log"
	"net/url"
)

//An object that can be stored in a database.
type DbObject interface {
	DbId() interface{}
	SetDbId(id interface{})

	Collection() string
}

//A database connection
type Database interface {
	Name() string
	URL() *url.URL
	SystemInformation() string
	Version() string

	Connect() error
	Disconnect() error

	HasId(id string, collection string) (bool, error)
	FindId(id string, presult DbObject) error
	Delete(o DbObject) error
	Save(o DbObject) error
	List(collection string, result interface{}) error
	FilteredList(collection string, props []string, result interface{}) error

	UniqueId() string
	IsValidId(id string) bool

	NewLogger(collection string, prefix string) *log.Logger
}

//supported databases
const (
	MONGODB = "mongodb"
)

//factory method to instantiate database drivers based on the ID
func NewDatabase(u *url.URL) Database {
	switch u.Scheme {
	case MONGODB:
		return NewMongoDB(u)
	}

	return nil
}
