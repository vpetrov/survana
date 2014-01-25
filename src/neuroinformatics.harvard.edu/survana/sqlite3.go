package survana

import (
	"log"
	"net/url"
    "database/sql"
    "github.com/coopernurse/gorp"
    _ "github.com/mattn/go-sqlite3"
)

type SQLite3 struct {
	Url      *url.URL
	Database *gorp.DbMap
	name	 string
}

var (
        sqlSession *sql.DB
    )

func NewSQLite3(u *url.URL, name string) *SQLite3 {
	return &SQLite3{
		Url: u,
		name: name,
	}
}

func (s *SQLite3) Name() string {
	return s.name;
}

func (s *SQLite3) URL() *url.URL {
	return s.Url
}

func (s *SQLite3) SystemInformation() string {
	return "SYSINFO: stub"
}

func (s *SQLite3) Version() string {
	return "SQLITE v3 STUB"
}

func (s *SQLite3) Connect() (err error) {
	log.Println("Connecting to ", s.Url.Path + "_" + s.name);
    sqlSession, err := sql.Open("sqlite3", s.Url.Path + "_" + s.name)
    if err != nil {
		log.Println("SQLLITE OPEN ERROR :(")
        return
    }

    s.Database = &gorp.DbMap{Db: sqlSession, Dialect: gorp.SqliteDialect{}}

    s.Database.AddTable(User{}).SetKeys(false, "Id")
    s.Database.AddTable(Group{}).SetKeys(false, "Id")
    s.Database.AddTable(Form{}).SetKeys(false, "Id")
    s.Database.AddTable(Study{}).SetKeys(false, "Id")
    //s.Database.AddTable(Field{}).SetKeys(true, "Id")

	log.Println("Creating tables in SQLITE")

    err = s.Database.CreateTablesIfNotExists()

	return
}


func (s *SQLite3) Disconnect() error {
	log.Println("Disconnect: stub")
	return nil
}

func (s *SQLite3) HasId(id string, collection string) (bool, error) {
	log.Println("HasId: stub", id, collection)
	return true, nil
}

func (s *SQLite3) FindId(id string, result DbObject) (err error) {
	collection := result.Collection()
	if len(collection) == 0 {
		return ErrInvalidCollection
	}

	err = s.Database.SelectOne(result, "SELECT * FROM ? WHERE id=?", collection, id );
	if err != nil {
		log.Println("SQLITE ERROR: ---------", err);
	}

	if err == sql.ErrNoRows {
		err = ErrNotFound
	} else {
		log.Printf("Found id=%v, _id=%#v\n", id, result.DbId())
	}

	return
}

func (s *SQLite3) Delete(o DbObject) error {
	log.Println("Delete: stub")
	return nil
}

func (s *SQLite3) Save(o DbObject) error {
	log.Println("Save: stub")
	return nil
}

func (s *SQLite3) List(collection string, result interface{}) error {
	log.Println("List: stub", collection, result)
	return nil
}

func (s *SQLite3) FilteredList(collection string, props []string, result interface{}) error {
	log.Println("FilteredList: stub", props, result)
	return nil
}

func (s *SQLite3) UniqueId() string {
	log.Println("Unique ID: stub");

	return "1"
}

func (s *SQLite3) IsValidId(id string) bool {
	log.Println("IsValidId: stub");
	return true
}

func (s *SQLite3) NewLogger(collection string, prefix string) *log.Logger {
	log.Println("NewLogger: stub", collection, prefix)
	return nil
}
