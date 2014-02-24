package study

import (
	"log"
    "strconv"
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
    study, err := survana.FindStudy(study_id, s.Db)
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
    index := query.Get("f")

    form_index, err := strconv.Atoi(index)
    if err != nil || len(study_id) == 0 || form_index < 0 {
        survana.BadRequest(w)
        return
    }

    study, err := survana.FindStudy(study_id, s.Db)
    if err != nil {
        survana.Error(w, err)
        return
    }

    if study == nil {
        survana.NotFound(w)
        return
    }

    //make sure the study has been published
    if !study.Published || form_index >= len(study.Html) {
        survana.NotFound(w)
        return
    }

    //fetch the HTML code
    html := study.Html[form_index]

    if len(html) == 0 {
        survana.NotFound(w)
        return
    }

    //write the HTML
    w.Write(html)
}
