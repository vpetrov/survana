package survana

import (
	_ "log"
)

const (
	USER_COLLECTION = "users"
)

type User struct {
    DBO             `bson:",inline,omitempty" json:"-"`
	Id     string   `bson:"id,omitempty" json:"id,omitempty"`
	Name   string   `bson:"name,omitempty" json:"name,omitempty"`
	Groups []string `bson:"groups,omitempty" json:"groups,omitempty"`
    AuthType string `bson:"auth_type,omitempty" json:"auth_type,omitempty"`
}

func NewUser(email, name string) *User {
	return &User{
        DBO: DBO{
                DBID:nil,
                Collection: USER_COLLECTION,
        },
		Id:   email,
		Name: name,
	}
}

func EmptyUser() *User {
    return &User{
        DBO: DBO { Collection: USER_COLLECTION },
    }
}

func FindUser(email string, db Database) (user *User, err error) {
	user = EmptyUser()
	err = db.FindId(email, user)

	if err != nil {
		if err == ErrNotFound {
			err = nil
		}

		return nil, err
	}

	return
}
