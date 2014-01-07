package survana

import (
        "time"
        _ "log"
       )

const (
        nID             = 6
        FORM_COLLECTION = "forms"
      )

type Form struct {
    DBID interface{}        `bson:"_id,omitempty" json:"-"`
    Id   string             `bson:"id,omitempty" json:"id"`
    Name string             `bson:"name,omitempty" json:"name"`
    Title string            `bson:"title,omitempty" json:"title"`
    Version string          `bson:"version,omitempty" json:"version"`
    CreatedOn time.Time     `bson:"created_on,omitempty" json:"created_on"`
    Fields []Field          `bson:"fields,omitempty" json:"fields"`
}

func NewForm() *Form {
    return &Form{
        Fields: make([]Field, 0),
    }
}

//returns a list of forms. if no forms are found, the 'forms' slice will be empty
func ListForms(db Database) (forms []Form, err error) {
    forms = make([]Form, 0)

    filter := []string{"id", "name", "version", "created_on"}

    err = db.FilteredList(FORM_COLLECTION, filter, &forms)
    if err != nil {
        if err == ErrNotFound {
            err = nil
        }
    }

    return
}

func(f *Form) DbId() interface{} {
    return f.DBID
}

func (f *Form) SetDbId(id interface{}) {
    f.DBID = id
}

func (f *Form) Collection() string {
    return FORM_COLLECTION
}

func FindForm(id string, db Database) (form *Form, err error) {
    form = NewForm()

    err = db.FindId(id, form)
    if err != nil {
        if err == ErrNotFound {
            err = nil
        }

        return nil, err
    }

    return
}

func (f *Form) GenerateId(db Database) (err error) {
    var (
            id string
            exists bool = true
        )

    for exists {
        //generate a random id
        id = RandomId(nID)
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

func (f *Form) Delete(db Database) (err error) {
    return db.Delete(f)
}

func (f *Form) Save(db Database) (err error) {
    return db.Save(f)
}
