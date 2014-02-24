package study

import (
	"neuroinformatics.harvard.edu/survana"
    "neuroinformatics.harvard.edu/survana/auth"
)

const (
	NAME = "study"
)

//The Admin component
type Study struct {
	*survana.Module
    Config *Config
    Auth auth.Strategy
}

// creates a new Admin module
func NewModule(path string, db survana.Database, config *Config, key *survana.PrivateKey) *Study {
	mux := survana.NewRESTMux()

	m := &Study{
		Module: &survana.Module{
			Name:   NAME,
			Path:   path,
			Db:     db,
			Router: mux,
            Mux: mux,
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
