package store

import (
	"github.com/vpetrov/perfect"
)

const (
	NAME = "store"
)

//The Store component
type Store struct {
	*perfect.Module
    Config *Config
}

// creates a new Admin module
func NewModule(path string, db perfect.Database, config *Config, key *perfect.PrivateKey) *Store {
	mux := perfect.NewRESTMux()

    if config == nil {
        config = &Config{}
    }

	m := &Store{
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

//	m.ParseTemplates()
	m.RegisterHandlers()

	return m
}
