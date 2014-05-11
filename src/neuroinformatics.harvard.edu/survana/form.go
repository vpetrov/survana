package survana

import (
	_ "log"
	"time"
    "github.com/vpetrov/perfect"
)

const (
	nID             = 6 //length of a form ID
	FORM_COLLECTION = "forms"
)

type Form struct {
    perfect.DBO             `bson:",inline,omitempty" json:"-"`
	Id          string      `bson:"id,omitempty" json:"id"`
	Name        string      `bson:"name,omitempty" json:"name"`
	Title       string      `bson:"title,omitempty" json:"title"`
	Description string      `bson:"description,omitempty" json:"description"`
	Version     string      `bson:"version,omitempty" json:"version"`
	CreatedOn   time.Time   `bson:"created_on,omitempty" json:"created_on"`
	Fields      []Field     `bson:"fields,omitempty" json:"fields"`

	//ACL
	OwnerId string `bson:"owner_id,omitempty" json:"owner_id,omitempty"`
}

func NewForm() *Form {
	return &Form{
        DBO: perfect.DBO { Collection: FORM_COLLECTION },
		Fields: make([]Field, 0),
	}
}

//returns a list of forms. if no forms are found, the 'forms' slice will be empty
func ListForms(filter []string, db perfect.Database) (forms []Form, err error) {
	forms = make([]Form, 0)

	err = db.FilteredList(FORM_COLLECTION, filter, &forms)
	if err != nil {
		if err == perfect.ErrNotFound {
			err = nil
		}
	}

	return
}

func FindForm(id string, db perfect.Database) (form *Form, err error) {
	form = NewForm()

	err = db.FindId(id, form)
	if err != nil {
		if err == perfect.ErrNotFound {
			err = nil
		}

		return nil, err
	}

	return
}

func (f *Form) GenerateId(db perfect.Database) (err error) {
	var (
		id     string
		exists bool = true
	)

	for exists {
		//generate a random id
		id = perfect.RandomId(nID)
		//check if it exists
		exists, err = db.HasId(id, FORM_COLLECTION)
		if err != nil {
			return
		}
	}

	//if a unique id was found, assign it to this object's Id
	f.Id = id

	return
}

func (f *Form) Delete(db perfect.Database) (err error) {
	return db.Delete(f)
}

func (f *Form) Save(db perfect.Database) (err error) {
	return db.Save(f)
}
