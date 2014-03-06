package survana

import (
	"log"
	"net/url"
)

//An object that can be stored in a database.
type DBO struct {
	DBID interface{}    `bson:"_id,omitempty" json:"-"`
	Collection string   `bson:"-" json:"-"`
}

type DBI interface {
	DbId() interface{}
	SetDbId(id interface{})
	DbCollection() string
}

func (dbo *DBO) DbId() interface{} {
    return dbo.DBID
}

func (dbo *DBO) SetDbId(id interface{}) {
    dbo.DBID = id
}

func (dbo *DBO) DbCollection() string {
    return dbo.Collection
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
	FindId(id string, presult DBI) error
	Delete(o DBI) error
	Save(o DBI) error
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
func NewDatabase(u *url.URL, name string) Database {
	switch u.Scheme {
        case MONGODB:
            return NewMongoDB(u, name)
	}

	return nil
}
