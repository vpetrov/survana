package study

import (
	"log"
	"net/http"
	"neuroinformatics.harvard.edu/survana"
)

// registers all route handlers
func (s *Study) RegisterHandlers() {
	app := s.mux

	//must end with slash
	app.Static("/assets/")

	app.Get("/", s.Index)
    app.Get("/form", s.Form)
}

// sends the app skeleton to the client
func (s *Study) Index(w http.ResponseWriter, r *survana.Request) {
    var err error

    //render the home page if no study was mentioned
    if (len(r.URL.RawQuery) == 0) {
        s.RenderTemplate(w, "index", nil)
        return
    }

    //set the study id
    study_id := r.URL.RawQuery

    log.Println("study id", study_id)

    //otherwise, fetch the study
    study, err := survana.FindStudy(study_id, s.Module.Db)
    if err != nil {
        survana.Error(w, err)
        return
    }

    if study == nil {
        survana.NotFound(w)
        return
    }

    s.RenderTemplate(w, "study/index", study)
}

func (s *Study) Form(w http.ResponseWriter, r *survana.Request) {
    var err error

    query := r.URL.Query()

    study_id := query.Get("s")
    form_id := query.Get("f")

    if len(study_id) == 0 || len(form_id) == 0 {
        survana.BadRequest(w)
        return
    }

    study, err := survana.FindStudy(study_id, s.Module.Db)
    if err != nil {
        survana.Error(w, err)
        return
    }

    if study == nil {
        survana.NotFound(w)
        return
    }

    //make sure the study has been published
    if !study.Published {
        survana.NotFound(w)
        return
    }

    //fetch the HTML code
    html, ok := study.Html[form_id]

    //if no such form has been published, it hasn't been found
    if !ok {
        survana.NotFound(w)
        return
    }

    if len(html) == 0 {
        survana.NotFound(w)
        return
    }

    //write the HTML
    w.Write(html)
}
