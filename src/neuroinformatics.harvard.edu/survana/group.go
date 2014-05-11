package survana

import (
	_ "log"
    "github.com/vpetrov/perfect"
)

const (
	GROUP_COLLECTION = "groups"
)

type Group struct {
    perfect.DBO         `bson:",inline,omitempty" json:"-"`
	Id   string `bson:"id,omitempty" json:"id,omitempty"`
	Name string `bson:"name,omitempty" json:"name,omitempty"`
}

//creates a new group
func NewGroup(name string) *Group {
	return &Group{
        DBO: perfect.DBO { Collection: GROUP_COLLECTION },
		Name: name,
	}
}

func EmptyGroup() *Group {
    return &Group{
        DBO: perfect.DBO { Collection: GROUP_COLLECTION },
    }
}

func FindGroup(id string, db perfect.Database) (group *Group, err error) {
	group = EmptyGroup()
	err = db.FindId(id, group)

	if err != nil {
		if err == perfect.ErrNotFound {
			err = nil
		}

		return nil, err
	}

	return
}
