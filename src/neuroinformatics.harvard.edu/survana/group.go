package survana

import (
	"github.com/vpetrov/perfect/orm"
)

type Group struct {
	orm.Object `bson:",inline,omitempty" json:"-"`
	Id         *string `bson:"id,omitempty" json:"id,omitempty"`
	Name       *string `bson:"name,omitempty" json:"name,omitempty"`
}
