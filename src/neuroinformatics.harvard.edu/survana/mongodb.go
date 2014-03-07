package survana

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net/url"
)

type MongoDB struct {
	Url      *url.URL
	Database *mgo.Database
	name	 string
}

var (
	session     *mgo.Session
	sessionInfo mgo.BuildInfo
)

func NewMongoDB(u *url.URL, name string) *MongoDB {
	if len(name) != 0 {
		u.Path = "/" + name
	}

	return &MongoDB{
		Url: u,
		name: name,
	}
}

func (db *MongoDB) Connect() (err error) {

	//create a single instance of the session object
	//which will get reused by subsequent calls to Connect()
	if session == nil {
		session, err = mgo.Dial(db.Url.String())
		if err != nil {
			return
		}
	}

	//create an mgo.Database object
	db.Database = session.DB(db.Url.Path[1:])

	//fetch the session information
	sessionInfo, err = session.BuildInfo()

	return
}

func (db *MongoDB) Disconnect() error {
	if session != nil {
		session.Close()
		session = nil
	}

	return nil
}

func (db *MongoDB) Name() string {
	return db.name
}

func (db *MongoDB) URL() *url.URL {
	return db.Url
}

func (db *MongoDB) SystemInformation() string {
	return sessionInfo.SysInfo
}

func (db *MongoDB) Version() string {
	return "MongoDB " + sessionInfo.Version
}

func (db *MongoDB) HasId(id string, collection string) (bool, error) {
	if len(collection) == 0 {
		return false, ErrInvalidCollection
	}

	count, err := db.Database.C(collection).Find(bson.M{"id": id}).Count()

	if err != nil {
		return false, err
	}

	return (count > 0), err
}

func (db *MongoDB) List(collection string, result interface{}) (err error) {
	if len(collection) == 0 {
		return ErrInvalidCollection
	}

	err = db.Database.C(collection).Find(nil).All(result)
	if err == mgo.ErrNotFound {
		err = ErrNotFound
	}

	return
}

func (db *MongoDB) FilteredList(collection string, props []string, result interface{}) (err error) {
	if len(collection) == 0 {
		return ErrInvalidCollection
	}

	filter := bson.M{}
	for p := range props {
		filter[props[p]] = 1
	}

	err = db.Database.C(collection).Find(nil).Select(filter).All(result)
	if err == mgo.ErrNotFound {
		err = ErrNotFound
	}

	return
}

func (db *MongoDB) FindId(id string, result DBI) (err error) {
	collection := result.DbCollection()
	if len(collection) == 0 {
		return ErrInvalidCollection
	}

	err = db.Database.C(collection).Find(bson.M{"id": id}).One(result)

    //restore the collection name, since bson erases it during unmarshalling
    result.SetDbCollection(collection)

	if err == mgo.ErrNotFound {
		err = ErrNotFound
	} else {
		log.Printf("Found id=%v, _id=%#v\n", id, result.DbId())
	}

	return
}

// Stores objects in the database. If the objects don't return a DbId, a new
// _id will be generated and assigned to the object. If a valid DbId exists,
// an Update operation will be performed, otherwise - an Insert()
// On success, the DBI will have a valid DbId. On error, the DbId will
// be the same it used to be, or nil if there was no DbId (and an error will
// be returned)
func (db *MongoDB) Save(obj DBI) (err error) {
	dbid := obj.DbId()
	var mgoid bson.ObjectId
	var ok bool

	collection := obj.DbCollection()

	if len(collection) == 0 {
		return ErrInvalidCollection
	}

	if dbid != nil {
		/* UPDATE */
		mgoid, ok = dbid.(bson.ObjectId)
		if !ok || !mgoid.Valid() {
			return ErrInvalidId
		}

		// Remove the _id property while we update the object.
		// This is necessary because MongoDB complains if the updates contain an _id
		// even if updated _id is the same as the original.
		obj.SetDbId(nil)

		// Restore _id when this function exits
		defer obj.SetDbId(mgoid)

		log.Printf("%s %#v\n", "UPDATING object", obj)
		// perform the update
		err = db.Database.C(collection).UpdateId(mgoid, bson.M{"$set": obj})
		if err != nil {
			return
		}
	} else {
		/* INSERT */

		//generate a new _id
		mgoid = bson.NewObjectId()

		//set the new _id
		obj.SetDbId(mgoid)

		log.Printf("%s %#v\n", "INSERTING new object", obj)

		//insert the object
		err = db.Database.C(collection).Insert(obj)
		if err != nil {
			//remove the _id if there was an error
			obj.SetDbId(nil)
			return
		}
	}

	return
}

func (db *MongoDB) Delete(obj DBI) error {
	//get the interface{} value
	dbid := obj.DbId()
	if dbid == nil {
		return ErrInvalidId
	}

	//get the object's collection. ignore if it's invalid
	collection := obj.DbCollection()
	if len(collection) == 0 {
		return ErrInvalidCollection
	}

	//remove the object by id
	return db.Database.C(collection).RemoveId(dbid)
}

//generates a unique id; safe to use across multiple machines
func (db *MongoDB) UniqueId() string {
	return bson.NewObjectId().Hex()
}

func (db *MongoDB) IsValidId(id string) bool {
	return bson.IsObjectIdHex(id)
}
