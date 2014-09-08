package dashboard

import (
	"github.com/vpetrov/perfect"
	"github.com/vpetrov/perfect/auth"
	"github.com/vpetrov/perfect/orm"
	"log"
)

const (
	NAME = "dashboard"
)

//The Admin component
type Dashboard struct {
	*perfect.Module
	Config *Config //dashboard.Config
	Auth   auth.Strategy
}

// creates a new Admin module
func NewModule(path string, db orm.Database, config *Config, key *perfect.PrivateKey) *Dashboard {

	m := &Dashboard{
		Module: &perfect.Module{
			Mux:  perfect.NewHTTPMux(),
			Name: NAME,
			Path: path,
			Db:   db,
			Log:  log.New(db.NewLogger("logs", ""), NAME, log.LstdFlags),
		},
		Config: config,
	}

	if config.Authentication != nil {
		m.Auth = auth.New(config.Authentication)
		m.Auth.Attach(m.Module)
	}

	//Parse all templates
	err := m.ParseTemplates()
	if err != nil {
		log.Fatalln(err)
	}

	m.RegisterHandlers()

	return m
}
