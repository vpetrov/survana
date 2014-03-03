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

    if config == nil {
        config = &Config{}
    }

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

    //by default, use the subject_id auth strategy
    m.Auth = auth.NewSubjectIdStrategy(nil)

	m.ParseTemplates()

	m.RegisterHandlers()

	return m
}
