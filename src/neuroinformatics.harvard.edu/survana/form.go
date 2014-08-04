package survana

import (
	"github.com/vpetrov/perfect"
	"github.com/vpetrov/perfect/orm"
	"time"
)

const (
	nID = 6 //length of a form ID
)

type Form struct {
	orm.Object  `bson:",inline,omitempty" json:"-"`
	Id          *string    `bson:"id,omitempty" json:"id,omitempty"`
	Name        *string    `bson:"name,omitempty" json:"name,omitempty"`
	Title       *string    `bson:"title,omitempty" json:"title,omitempty"`
	Description *string    `bson:"description,omitempty" json:"description,omitempty"`
	Version     *string    `bson:"version,omitempty" json:"version,omitempty"`
	CreatedOn   *time.Time `bson:"created_on,omitempty" json:"created_on,omitempty"`
	Fields      *[]Field   `bson:"fields,omitempty" json:"fields,omitempty"`

	//ACL
	OwnerId *string `bson:"owner_id,omitempty" json:"owner_id,omitempty"`
}

//TODO: this is vulnerable to a race condition. It would be better to
//set a unique constraint on the Id property and attempt to write to the db
func (f *Form) GenerateId(db orm.Database) (err error) {
	var (
		exists bool  = true
		search *Form = &Form{}
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
