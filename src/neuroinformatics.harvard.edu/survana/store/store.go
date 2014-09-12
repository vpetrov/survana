package store

import (
	"github.com/vpetrov/perfect"
	"github.com/vpetrov/perfect/orm"
	"log"
)

const (
	NAME = "store"
)

//The Store component
type Store struct {
	*perfect.Module
	Config *Config
	Key    *perfect.PrivateKey
}

// creates a new Admin module
func NewModule(path string, db orm.Database, config *Config, key *perfect.PrivateKey) *Store {
	if config == nil {
		config = &Config{}
	}

	m := &Store{
		Module: &perfect.Module{
			Mux:  perfect.NewPrettyMux(),
			Name: NAME,
			Path: path,
			Db:   db,
			Log:  log.New(db.NewLogger("logs", ""), NAME, log.LstdFlags|log.Lmicroseconds),
		},
		Config: config,
		Key:    key,
	}

	//	err := m.ParseTemplates()
	//  if err != nil {
	//      log.Fatalln(err)
	//  }

	m.RegisterHandlers()

	return m
}
