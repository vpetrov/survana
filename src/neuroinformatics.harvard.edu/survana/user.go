package survana

import (
	_ "log"
)

const (
	USER_COLLECTION = "users"
)

type User struct {
	Id   string
	Name string

	//DbObject
	DBID interface{} `bson:"_id,omitempty"`
}

func NewUser(email, name string) *User {
	return &User{
		Id:   email,
		Name: name,
	}
}

func (u *User) DbId() interface{} {
	return u.DBID
}

func (u *User) SetDbId(v interface{}) {
	u.DBID = v
}

func (u *User) Collection() string {
	return USER_COLLECTION
}

func FindUser(email string, db Database) (user *User, err error) {
	user = &User{}
	err = db.FindId(email, user)

	if err != nil {
		if err == ErrNotFound {
			err = nil
		}

		return nil, err
	}

	return
}
