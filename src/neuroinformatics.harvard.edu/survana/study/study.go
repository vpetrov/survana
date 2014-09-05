package study

import (
	"github.com/vpetrov/perfect"
	"github.com/vpetrov/perfect/auth"
	"github.com/vpetrov/perfect/orm"
	"log"
)

const (
	NAME = "study"
)

//The Admin component
type Study struct {
	*perfect.Module
	Config *Config
	Auth   auth.Strategy
}

// creates a new Admin module
func NewModule(path string, db orm.Database, config *Config, key *perfect.PrivateKey) *Study {
	if config == nil {
		config = &Config{}
	}

	m := &Study{
		Module: &perfect.Module{
			Mux:  perfect.PrettyMux(),
			Name: NAME,
			Path: path,
			Db:   db,
			Log:  log.New(db.NewLogger("logs", ""), NAME, log.LstdFlags),
		},
		Config: config,
	}

	//by default, use the subject_id auth strategy
	m.Auth = NewSubjectIdStrategy(nil)

	err := m.ParseTemplates()
	if err != nil {
		log.Fatalln(err)
	}

	m.RegisterHandlers()

	return m
}
