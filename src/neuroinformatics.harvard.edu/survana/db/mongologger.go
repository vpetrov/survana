package db

import (
	"labix.org/v2/mgo/bson"
	"log"
	"time"
)

type MongoLogger struct {
	mongodb    *MongoDB
	collection string
}

func (ml MongoLogger) Write(p []byte) (n int, err error) {
	//will always return that all bytes were written
	n = len(p)
	//create the log document
	doc := bson.M{
		"timestamp": time.Now(),
		"message":   string(p),
	}

	err = ml.mongodb.Database.C(ml.collection).Insert(doc)
	return
}

func (db *MongoDB) NewLogger(collection string, prefix string) *log.Logger {
	return log.New(MongoLogger{
		collection: collection,
		mongodb:    db,
	},
		prefix+":",
		0)
}
