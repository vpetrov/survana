package study

import (
	"log"
    "time"
    "strconv"
    "strings"
	"net/http"
	"neuroinformatics.harvard.edu/survana"
)

// registers all route handlers
func (s *Study) RegisterHandlers() {
	app := s.Mux

	//must end with slash
	app.Static("/assets/")

	app.Get("/", s.NewIndex)
    app.Get("/form", s.Form)
}

func (s *Study) NewIndex(w http.ResponseWriter, r *survana.Request) {
    var (
            err error
            study_id string
            subject_id string
        )

    //check to see if we have ?study or ?study/index
    if strings.Contains(r.URL.RawQuery, "/") {
        params := strings.SplitN(r.URL.RawQuery, "/", 2)
        study_id = params[0]
        subject_id = strings.ToUpper(params[1])
    } else {
        study_id = r.URL.RawQuery
    }

    //render the home page if no study was mentioned
    if (len(study_id) == 0) {
        s.RenderTemplate(w, "index", nil)
        return
    }

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

    log.Printf("auth_enabled=%v subjects=%#v", study.AuthEnabled, study.Subjects)

    //no auth? just render the study index page
    if !study.AuthEnabled {
        s.RenderTemplate(w, "study/index", study)
        return
    }

    /* auth is required */
    //get the session
    session, err := r.Session()
    if err != nil {
        survana.Error(w, err)
        return
    }

    //if the session has been authorized, render the login screen
    if session.Authenticated {
        s.RenderTemplate(w, "study/index", study)
        return
    }

    if len(subject_id) == 0 {
        s.Auth.LoginPage(w, r)
        return
    }

    enabled, ok := study.Subjects[subject_id]

    if !ok {
        s.Auth.LoginPage(w, r)
        return
    }

    if !enabled {
        log.Println("subject id", subject_id, "is not enabled")
        s.Auth.LoginPage(w, r)
        return
    }

    session.Authenticated = true
    session.Values["study_id"] = study_id;
    session.Values["subject_id"] = subject_id;

    err = session.Save(s.Db)
    if err != nil {
        survana.Error(w, err)
        return
    }

	//set the cookie and make it valid for a month
	http.SetCookie(w, &http.Cookie{
		Name:     survana.SESSION_ID,
		Value:    session.Id,
		Path:     r.Module.MountPoint,
		Expires:  time.Now().Add(time.Hour * 24 * 30),
		Secure:   true,
		HttpOnly: true,
	})

    //finally render the template
    s.RenderTemplate(w, "study/index", study)
}

// sends the app skeleton to the client
func (s *Study) Index(w http.ResponseWriter, r *survana.Request) {
    var (
            err error
            study_id string
            subject_id string
        )

    if strings.Contains(r.URL.RawQuery, "/") {
        params := strings.SplitN(r.URL.RawQuery, "/", 2)
        study_id = params[0]
        subject_id = strings.ToUpper(params[1])
    } else {
        study_id = r.URL.RawQuery
    }

    //render the home page if no study was mentioned
    if (len(study_id) == 0) {
        s.RenderTemplate(w, "index", nil)
        return
    }

    log.Println("study id", study_id, "subject id", subject_id)

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
