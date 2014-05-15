package study

import (
	"github.com/vpetrov/perfect"
    "github.com/vpetrov/perfect/auth"
    "log"
)

const (
	NAME = "study"
)

//The Admin component
type Study struct {
	*perfect.Module
    Config *Config
    Auth auth.Strategy
}

// creates a new Admin module
func NewModule(path string, db perfect.Database, config *Config, key *perfect.PrivateKey) *Study {
	mux := perfect.NewRESTMux()

    if config == nil {
        config = &Config{}
    }

	m := &Study{
		Module: &perfect.Module{
			Name:   NAME,
			Path:   path,
			Db:     db,
			Router: mux,
            Mux: mux,
			Log:    db.NewLogger("logs", NAME),
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
