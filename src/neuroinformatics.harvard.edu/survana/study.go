package survana

import (
	_ "log"
	"time"
)

const (
	STUDY_COLLECTION = "studies"
)

type Study struct {
	Id          string    `bson:"id,omitempty" json:"id"`
	Name        string    `bson:"name,omitempty" json:"name"`
	Title       string    `bson:"title,omitempty" json:"title"`
	Description string    `bson:"description,omitempty" json:"description"`
	Version     string    `bson:"version,omitempty" json:"version"`
	CreatedOn   time.Time `bson:"created_on,omitempty" json:"created_on"`
	Forms       []Form    `bson:"forms,omitempty" json:"forms"`
    Html        [][]byte  `bson:"html,omitempty" json:"-"`
	Published   bool      `bson:"published" json:"published"`

	//DbObject
	DBID interface{} `bson:"_id,omitempty" json:"-"`

	//ACL
	OwnerId string `bson:"owner_id,omitempty" json:"owner_id,omitempty"`
}

func NewStudy() *Study {
	return &Study{
        Html: make([][]byte, 0),
    }
}

func (s *Study) DbId() interface{} {
	return s.DBID
}

func (s *Study) SetDbId(id interface{}) {
	s.DBID = id
}

func (s *Study) Collection() string {
	return STUDY_COLLECTION
}

func FindStudy(id string, db Database) (study *Study, err error) {
	study = NewStudy()

	err = db.FindId(id, study)
	if err != nil {
		if err == ErrNotFound {
			err = nil
		}

		return nil, err
	}

	return
}

//returns a list of studies.
func ListStudies(db Database) (studies []Study, err error) {
	studies = make([]Study, 0)

	filter := []string{"id", "name", "title", "version", "created_on", "owner_id", "forms", "published"}

	err = db.FilteredList(STUDY_COLLECTION, filter, &studies)
	if err != nil {
		if err == ErrNotFound {
			err = nil
		}
	}

	return
}

func (s *Study) Delete(db Database) (err error) {
	return db.Delete(s)
}

func (s *Study) Save(db Database) (err error) {
	return db.Save(s)
}

func (f *Study) GenerateId(db Database) (err error) {
	var (
		id     string
		exists bool = true
	)

	for exists {
		//generate a random id
		id = RandomId(nID)
		//check if it exists
		exists, err = db.HasId(id, STUDY_COLLECTION)
		if err != nil {
			return
		}
	}

	//if a unique id was found, assign it to this object's Id
	f.Id = id

	return
}
