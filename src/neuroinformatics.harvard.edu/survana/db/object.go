package db

type Object interface {
	DbId() interface{}
	SetDbId(id interface{})

	Collection() string
}
