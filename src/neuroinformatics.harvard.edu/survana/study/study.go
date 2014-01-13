package study

import (
	"neuroinformatics.harvard.edu/survana"
)

const (
	NAME = "study"
)

//The Admin component
type Study struct {
	*survana.Module
	mux *survana.RESTMux
}

// creates a new Admin module
func NewModule(path string, db survana.Database) *Study {

	mux := survana.NewRESTMux()

	m := &Study{
		Module: &survana.Module{
			Name:   NAME,
			Path:   path,
			Db:     db,
			Router: mux,
			Log:    db.NewLogger("logs", NAME),
		},
		mux: mux,
	}

	m.ParseTemplates()

	m.RegisterHandlers()

	return m
}
