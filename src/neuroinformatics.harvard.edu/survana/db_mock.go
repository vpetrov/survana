package survana

import (
	"log"
	"net/url"
	"os"
)

type MockDatabase struct {
	Calls    map[string]int
    OnList   func(string, interface{})
    OnFilteredList func(string, []string, interface{})
    OnHasId func(string) (bool, error)
	OnFindId func(string, DbObject)
	OnDelete func(DbObject)
	OnSave   func(DbObject)
}

func NewMockDatabase() *MockDatabase {
	return &MockDatabase{
		Calls: make(map[string]int, 0),
	}
}

func (db *MockDatabase) Name() string {
	return "mock"
}

func (db *MockDatabase) URL() *url.URL {
	u, _ := url.Parse("mock://localhost:1/mock")
	return u
}

func (db *MockDatabase) SystemInformation() string {
	return "Mock DB on SomeOS"
}

func (db *MockDatabase) Version() string {
	return "Mock DB v.1.0"
}

func (db *MockDatabase) Connect() error {
	return nil
}

func (db *MockDatabase) Disconnect() error {
	return nil
}

func (db *MockDatabase) List(collection string, result interface{}) error {
    db.Calls["List"]++

    if db.OnList != nil {
        db.OnList(collection, result)
    }

    return nil
}

func (db *MockDatabase) FilteredList(collection string, props []string, result interface{}) error {
    db.Calls["FilteredList"]++
    if db.OnFilteredList != nil {
        db.OnFilteredList(collection, props, result)
    }

    return nil
}

func (db *MockDatabase) FindId(id string, presult DbObject) error {
	db.Calls["FindId"]++
	if db.OnFindId != nil {
		db.OnFindId(id, presult)
	}
	return nil
}

func (db *MockDatabase) Delete(o DbObject) error {
	db.Calls["Delete"]++
	if db.OnDelete != nil {
		db.OnDelete(o)
	}

	return nil
}

func (db *MockDatabase) Save(o DbObject) error {
	db.Calls["Save"]++
	if db.OnSave != nil {
		db.OnSave(o)
	}

	return nil
}

func (db *MockDatabase) HasId(id string, collection string) (result bool, err error) {
    db.Calls["HasId"]++
    if db.OnHasId != nil {
        result, err = db.OnHasId(id)
    }

    return
}

func (db *MockDatabase) UniqueId() string {
	return "ABCD"
}

func (db *MockDatabase) IsValidId(id string) bool {
	return true
}

func (db *MockDatabase) NewLogger(collection, prefix string) *log.Logger {
	return log.New(os.Stdout, "mock", 0)
}
