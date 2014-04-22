package dashboard

import (
	"neuroinformatics.harvard.edu/survana"
    "neuroinformatics.harvard.edu/survana/auth"
)

const (
	NAME = "dashboard"
)

//The Admin component
type Dashboard struct {
	*survana.Module
    Config *Config //dashboard.Config
    Auth           auth.Strategy
}


// creates a new Admin module
func NewModule(path string, db survana.Database, config *Config, key *survana.PrivateKey) *Dashboard {

	mux := survana.NewRESTMux()

	m := &Dashboard{
		Module: &survana.Module{
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
