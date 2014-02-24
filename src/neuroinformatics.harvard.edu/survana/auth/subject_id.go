package auth

import (
    _ "log"
	"net/http"
	"neuroinformatics.harvard.edu/survana"
    )

const (
	SUBJECT_ID_COLLECTION = "auth_subject_id"
)

type studySubjects struct {
    Id          string  `bson:"id,omitempty"`
    StudyId     string  `bson:"study_id,omitempty"`
    Subjects    map[string]bool `bson:"subjects,omitempty"`

    /* DbObject */
    DBID        interface{} `bson:"_id,omitempty"`
}

type SubjectIdStrategy struct {
}

func NewSubjectIdStrategy(config *Config) SubjectIdStrategy {
    return SubjectIdStrategy{}
}

func (sid SubjectIdStrategy) Attach(module *survana.Module) {
    app := module.Mux

    app.Get("/login", survana.NotLoggedIn(sid.LoginPage))
    app.Post("/login", survana.NotLoggedIn(sid.Login))
}


func (sid SubjectIdStrategy) LoginPage(w http.ResponseWriter, r *survana.Request) {
    r.Module.RenderTemplate(w, "auth/subject_id/login", nil)
}

func (sid SubjectIdStrategy) RegistrationPage(w http.ResponseWriter, r *survana.Request) {
    survana.NotFound(w)
}

func (sid SubjectIdStrategy) Login(w http.ResponseWriter, r *survana.Request) {
    survana.NotFound(w)
}

func (sid SubjectIdStrategy) Register(w http.ResponseWriter, r *survana.Request) {
    survana.NotFound(w)
}

func (sid SubjectIdStrategy) Logout(w http.ResponseWriter, r *survana.Request) {
    survana.NotFound(w)
}
