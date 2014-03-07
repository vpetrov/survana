package dashboard

import (
	"log"
	"net/http"
	"neuroinformatics.harvard.edu/survana"
	"time"
    "strconv"
)

func (d *Dashboard) StudyListPage(w http.ResponseWriter, r *survana.Request) {
	d.RenderTemplate(w, "study/list", nil)
}

func (d *Dashboard) StudyList(w http.ResponseWriter, r *survana.Request) {
	studies, err := survana.ListStudies(d.Db)

	if err != nil {
		survana.Error(w, err)
		return
	}

	survana.JSONResult(w, true, studies)
}

func (d *Dashboard) CreateStudyPage(w http.ResponseWriter, r *survana.Request) {
	d.RenderTemplate(w, "study/create", nil)
}

func (d *Dashboard) CreateStudy(w http.ResponseWriter, r *survana.Request) {
	var err error

	//get the current session
	session, err := r.Session()
	if err != nil {
		survana.Error(w, err)
	}

	study := survana.NewStudy()

	//parse input data
	err = r.ParseJSON(&study)
	if err != nil {
		survana.Error(w, err)
		return
	}

	study.CreatedOn = time.Now()
	study.OwnerId = session.UserId

	//generate a unique id
	err = study.GenerateId(d.Db)
	if err != nil {
		survana.Error(w, err)
	}

	//save the study
	err = study.Save(d.Db)
	if err != nil {
		survana.Error(w, err)
		return
	}

	//result format is { id: "abcd" }
	result := &struct {
		Id string `json:"id"`
	}{Id: study.Id}

	survana.JSONResult(w, true, result)
}

func (d *Dashboard) GetStudy(w http.ResponseWriter, r *survana.Request) {
	query := r.URL.Query()
	study_id := query.Get("id")

	//TODO: Validate alnum
	if len(study_id) == 0 {
		survana.BadRequest(w)
		return
	}

	study, err := survana.FindStudy(study_id, d.Db)
	if err != nil {
		survana.Error(w, err)
		return
	}

	//not found?
	if study == nil {
		survana.NotFound(w)
		return
	}

	//return the form as JSON
	survana.JSONResult(w, true, study)
}

func (d *Dashboard) EditStudyPage(w http.ResponseWriter, r *survana.Request) {
	d.RenderTemplate(w, "study/edit", nil)
}

func (d *Dashboard) EditStudy(w http.ResponseWriter, r *survana.Request) {
	var err error

	//get the study id
	query := r.URL.Query()
	study_id := query.Get("id")

	//TODO: Validate alnum
	if len(study_id) == 0 {
		survana.BadRequest(w)
		return
	}

	//make sure the form exists
	study, err := survana.FindStudy(study_id, d.Db)
	if err != nil {
		survana.Error(w, err)
		return
	}

	//not found?
	if study == nil {
		survana.NotFound(w)
		return
	}

	//parse new form data sent by the client
	user_study := survana.NewStudy()
	err = r.ParseJSON(user_study)
	if err != nil {
		survana.Error(w, err)
		return
	}

	//TODO?: validate form fields? validate using a schema?

	log.Printf("%s: %#v\n", "JSON study submitted by the client", user_study)

	//copy properties that should not be changed
	user_study.DBID = study.DBID
	user_study.Id = study.Id
	user_study.CreatedOn = study.CreatedOn
	user_study.OwnerId = study.OwnerId

	//update the form
	err = user_study.Save(d.Db)
	if err != nil {
		survana.Error(w, err)
		return
	}

	//success
	survana.NoContent(w)
}

func (d *Dashboard) DeleteStudy(w http.ResponseWriter, r *survana.Request) {
	var err error

	//get the form id
	query := r.URL.Query()
	study_id := query.Get("id")

	log.Println("study to delete:", study_id)

	//TODO: Validate alnum
	if len(study_id) == 0 {
		survana.BadRequest(w)
		return
	}

	//make sure the form exists
	study, err := survana.FindStudy(study_id, d.Db)
	if err != nil {
		survana.Error(w, err)
		return
	}

	//not found?
	if study == nil {
		survana.NotFound(w)
		return
	}

	err = study.Delete(d.Db)

	if err != nil {
		survana.Error(w, err)
		return
	}

	log.Println("study", study_id, "deleted")
	survana.NoContent(w)
}

func (d *Dashboard) ViewStudyPage(w http.ResponseWriter, r *survana.Request) {
	d.RenderTemplate(w, "study/view", nil)
}

func (d *Dashboard) PublishStudyPage(w http.ResponseWriter, r *survana.Request) {
	d.RenderTemplate(w, "study/publish", nil)
}

func (d *Dashboard) PublishStudyForm(w http.ResponseWriter, r *survana.Request) {

	query := r.URL.Query()
	study_id := query.Get("id")
    form_index, err := strconv.Atoi(query.Get("f"))

	if (len(study_id) == 0) || err != nil || form_index < 0 {
		survana.BadRequest(w)
		return
	}

	html, err := r.BodyBytes(r.Body)
	if err != nil {
		survana.Error(w, err)
		return
	}

    study, err := survana.FindStudy(study_id, d.Db)
    if err != nil {
        survana.Error(w, err)
        return
    }

    if study == nil || form_index >= len(study.Forms) {
        survana.NotFound(w)
        return
    }

    log.Println("should publish form index", form_index)

    if (study.Html == nil) {
        study.Html = make([][]byte, len(study.Forms))
    }

    //overwrite the HTML
    study.Html[form_index] = html

    //save the study
    err = study.Save(d.Db)
    if err != nil {
        survana.Error(w, err)
        return
    }

    log.Println("done saving published form")

	survana.NoContent(w)
}

func (d *Dashboard) StudySubjectsPage(w http.ResponseWriter, r *survana.Request) {
	d.RenderTemplate(w, "study/subjects", nil)
}

func (d *Dashboard) AddStudySubjects(w http.ResponseWriter, r *survana.Request) {
	var err error

	//get the form id
	query := r.URL.Query()
	study_id := query.Get("id")

	log.Println("study to delete:", study_id)

	//TODO: Validate alnum
	if len(study_id) == 0 {
		survana.BadRequest(w)
		return
	}

	//make sure the form exists
	study, err := survana.FindStudy(study_id, d.Db)
	if err != nil {
		survana.Error(w, err)
		return
	}

	//not found?
	if study == nil {
		survana.NotFound(w)
		return
	}

    //parse json data
    ids := make([]string, 0)
    err = r.JSONBody(r.Body, &ids)
    if err != nil {
        survana.Error(w, err)
        return
    }

    nids := len(ids)
    if nids == 0 {
        survana.BadRequest(w)
        return
    }

    var id string

    //save and enable all IDs
    for i := 0; i < nids; i++ {
        id = ids[i]
        _, exists := study.Subjects[id]

        if !exists {
            study.Subjects[id] = true
        }
    }

    //store the updated Survey
    err = study.Save(d.Db)
    if err != nil {
        survana.Error(w, err)
        return
    }

    survana.JSONResult(w, true, study.Subjects)
}
