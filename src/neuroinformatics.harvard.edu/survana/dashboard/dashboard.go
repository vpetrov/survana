package dashboard

import (
	"github.com/vpetrov/perfect"
    "github.com/vpetrov/perfect/auth"
)

const (
	NAME = "dashboard"
)

//The Admin component
type Dashboard struct {
    *perfect.Module
    Config  *Config //dashboard.Config
    Auth    auth.Strategy
}


// creates a new Admin module
func NewModule(path string, db perfect.Database, config *Config, key *perfect.PrivateKey) *Dashboard {

	mux := perfect.NewRESTMux()

	m := &Dashboard{
		Module: &perfect.Module{
			Name:   NAME,
			Path:   path,
			Db:     db,
			Router: mux,
            Mux:    mux,
			Log:    db.NewLogger("logs", NAME),
		},
        Config: config,
	}

    if config.Authentication != nil {
        m.Auth = auth.New(config.Authentication)
        m.Auth.Attach(m.Module)
    }

	m.ParseTemplates()

	m.RegisterHandlers()

	return m
}
