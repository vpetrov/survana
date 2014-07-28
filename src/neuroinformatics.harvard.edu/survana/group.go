package survana

import (
	"github.com/vpetrov/perfect/orm"
)

type Group struct {
	orm.Object `bson:",inline,omitempty" json:"-"`
	Id         *string `bson:"id,omitempty" json:"id,omitempty"`
	Name       *string `bson:"name,omitempty" json:"name,omitempty"`
}

//creates a new group
func NewGroup(name string) *Group {
	return &Group{
		Name: &name,
	}
}

func FindGroup(id string, db orm.Database) (group *Group, err error) {
	group = &Group{Id: &id}
	err = db.Find(group)

	if err != nil {
		if err == orm.ErrNotFound {
			err = nil
		}

		return nil, err
	}

	return
}
