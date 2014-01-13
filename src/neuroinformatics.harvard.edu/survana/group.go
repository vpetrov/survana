package survana

import (
	_ "log"
)

const (
	GROUP_COLLECTION = "groups"
)

type Group struct {
	Id   string `bson:"id,omitempty" json:"id,omitempty"`
	Name string `bson:"name,omitempty" json:"name,omitempty"`

	//DbObject
	DBID interface{} `bson:"_id,omitempty" json:"-"`
}

//creates a new group
func NewGroup(name string) *Group {
	return &Group{
		Name: name,
	}
}

//DbObject
func (g *Group) DbId() interface{} {
	return g.DBID
}

func (g *Group) SetDbId(v interface{}) {
	g.DBID = v
}

func (g *Group) Collection() string {
	return GROUP_COLLECTION
}

func FindGroup(id string, db Database) (group *Group, err error) {
	group = &Group{}
	err = db.FindId(id, group)

	if err != nil {
		if err == ErrNotFound {
			err = nil
		}

		return nil, err
	}

	return
}
