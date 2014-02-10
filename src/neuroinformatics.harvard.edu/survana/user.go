package survana

import (
	_ "log"
)

const (
	USER_COLLECTION = "users"
)

type User struct {
	Id     string   `bson:"id,omitempty" json:"id,omitempty"`
	Name   string   `bson:"name,omitempty" json:"name,omitempty"`
	Groups []string `bson:"groups,omitempty" json:"groups,omitempty"`
    AuthType string `bson:"auth_type,omitempty" json:"auth_type,omitempty"`

	//DbObject
	DBID interface{} `bson:"_id,omitempty" json:"-"`
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
