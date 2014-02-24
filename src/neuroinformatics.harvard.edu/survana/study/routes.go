package study

import (
	"log"
    "strconv"
	"net/http"
	"neuroinformatics.harvard.edu/survana"
    "neuroinformatics.harvard.edu/survana/auth"
)

// registers all route handlers
func (s *Study) RegisterHandlers() {
	app := s.Mux

	//must end with slash
	app.Static("/assets/")

	app.Get("/", auth.Protect(s.Index))
    app.Get("/form", auth.Protect(s.Form))
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

    //if this study is using Subject IDs, render the login screen
    if len(study.Subjects) > 0 {
        s.RenderTemplate(w, "study/login", study)
        return
    }

    s.RenderTemplate(w, "study/index", study)
}

// sends the app skeleton to the client
func (s *Study) Login(w http.ResponseWriter, r *survana.Request) {
    var err error

    //render the home page if no study was mentioned
    if (len(r.URL.RawQuery) == 0) {
        survana.BadRequest(w)
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

    //read form data
    form := make(map[string]string)

    err = r.ParseJSON(&form)
    if err != nil {
        survana.Error(w, err)
        return
    }

    //read the subject id
    subject_id, ok := form["subject_id"]
    if !ok || len(subject_id) == 0 {
        survana.JSONResult(w, false, "Please complete all the fields.")
        return
    }

    //check that the subject id exists in the study.Subjects and it's enabled
    enabled, ok := study.Subjects[subject_id]
    if !ok {
        survana.JSONResult(w, false, "We were unable to find this ID.")
        return
    }

    if !enabled {
        survana.JSONResult(w, false, "This ID has already been used.")
        return
    }

    survana.JSONResult(w, true, s.MountPoint + "/go?" + study_id)
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
