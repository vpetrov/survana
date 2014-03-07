package survana

import (
	_ "log"
)

const (
	GROUP_COLLECTION = "groups"
)

type Group struct {
    DBO         `bson:",inline,omitempty" json:"-"`
	Id   string `bson:"id,omitempty" json:"id,omitempty"`
	Name string `bson:"name,omitempty" json:"name,omitempty"`
}

//creates a new group
func NewGroup(name string) *Group {
	return &Group{
        DBO: DBO { Collection: GROUP_COLLECTION },
		Name: name,
	}
}

func EmptyGroup() *Group {
    return &Group{
        DBO: DBO { Collection: GROUP_COLLECTION },
    }
}

func FindGroup(id string, db Database) (group *Group, err error) {
	group = EmptyGroup()
	err = db.FindId(id, group)

	if err != nil {
		if err == ErrNotFound {
			err = nil
		}

		return nil, err
	}

	return
}
