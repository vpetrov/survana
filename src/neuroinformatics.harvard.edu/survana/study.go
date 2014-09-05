package survana

import (
	"github.com/vpetrov/perfect"
	"github.com/vpetrov/perfect/orm"
	"time"
)

type Study struct {
	orm.Object  `bson:",inline,omitempty" json:"-"`
	Id          *string          `bson:"id,omitempty" json:"id,omitempty,readonly"`
	Name        *string          `bson:"name,omitempty" json:"name,omitempty"`
	Title       *string          `bson:"title,omitempty" json:"title,omitempty"`
	Description *string          `bson:"description,omitempty" json:"description,omitempty"`
	Version     *string          `bson:"version,omitempty" json:"version,omitempty"`
	CreatedOn   *time.Time       `bson:"created_on,omitempty" json:"created_on,omitempty,readonly"`
	FormIds     *[]string        `bson:"form_ids,omitempty" json:"form_ids,omitempty"`
	Html        *[][]byte        `bson:"html,omitempty" json:"-"`
	Published   *bool            `bson:"published,omitempty" json:"published,omitempty,readonly"`
	Subjects    *map[string]bool `bson:"subjects,omitempty" json:"subjects,omitempty"`
	AuthEnabled *bool            `bson:"auth_enabled,omitempty" json:"auth_enabled,omitempty"`
	StoreUrl    *string          `bson:"store_url,omitempty" json:"store_url,omitempty,readonly"`
	Revision    int              `bson:"revision,omitempty" json:"revision,omitempty,readonly"`
	//note: revision is not a pointer because 0 is not a valid revision number, which starts at 1

	//ACL
	OwnerId *string `bson:"owner_id,omitempty" json:"owner_id,omitempty,readonly"`
}

//TODO: this method has a race condition on Id. In addition, it's exactly the same as Form.GenerateId()
func (f *Study) GenerateId(db orm.Database) (err error) {
	var (
		exists bool   = true
		search *Study = &Study{}
	)

	for exists {
		//generate a random id
		search.Id = orm.String(perfect.RandomId(nID))
		//check if it exists
		err = db.Find(search)
		if err != nil {
			if err != orm.ErrNotFound {
				return
			}
			err = nil

			break
		}
	}

	//if a unique id was found, assign it to this object's Id
	f.Id = search.Id

	return
}
