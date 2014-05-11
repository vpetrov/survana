package dashboard

import (
	"log"
	"net/http"
	"neuroinformatics.harvard.edu/survana"
    "github.com/vpetrov/perfect"
	"time"
    "strconv"
)

func (d *Dashboard) StudyListPage(w http.ResponseWriter, r *perfect.Request) {
	d.RenderTemplate(w, "study/list", nil)
}

func (d *Dashboard) StudyList(w http.ResponseWriter, r *perfect.Request) {
	studies, err := survana.ListStudies(d.Db)

	if err != nil {
		perfect.Error(w, err)
		return
	}

	perfect.JSONResult(w, true, studies)
}

func (d *Dashboard) CreateStudyPage(w http.ResponseWriter, r *perfect.Request) {
	d.RenderTemplate(w, "study/create", nil)
}

func (d *Dashboard) CreateStudy(w http.ResponseWriter, r *perfect.Request) {
	var err error

	//get the current session
	session, err := r.Session()
	if err != nil {
		perfect.Error(w, err)
	}

	study := survana.NewStudy()

	//parse input data
	err = r.ParseJSON(&study)
	if err != nil {
		perfect.Error(w, err)
		return
	}
    now := time.Now()
	study.CreatedOn = &now
	study.OwnerId = session.UserId

	//generate a unique id
	err = study.GenerateId(d.Db)
	if err != nil {
		perfect.Error(w, err)
	}

	//save the study
	err = study.Save(d.Db)
	if err != nil {
		perfect.Error(w, err)
		return
	}

	//result format is { id: "abcd" }
	result := &struct {
		Id string `json:"id"`
	}{Id: study.Id}

	perfect.JSONResult(w, true, result)
}

func (d *Dashboard) GetStudy(w http.ResponseWriter, r *perfect.Request) {
	query := r.URL.Query()
	study_id := query.Get("id")

	//TODO: Validate alnum
	if len(study_id) == 0 {
		perfect.BadRequest(w)
		return
	}

	study, err := survana.FindStudy(study_id, d.Db)
	if err != nil {
		perfect.Error(w, err)
		return
	}

	//not found?
	if study == nil {
		perfect.NotFound(w)
		return
	}

	//return the form as JSON
	perfect.JSONResult(w, true, study)
}

func (d *Dashboard) EditStudyPage(w http.ResponseWriter, r *perfect.Request) {
	d.RenderTemplate(w, "study/edit", nil)
}

func (d *Dashboard) EditStudy(w http.ResponseWriter, r *perfect.Request) {
	var err error

	//get the study id
	query := r.URL.Query()
	study_id := query.Get("id")

	//TODO: Validate alnum
	if len(study_id) == 0 {
		perfect.BadRequest(w)
		return
	}

	//make sure the form exists
	study, err := survana.FindStudy(study_id, d.Db)
	if err != nil {
		perfect.Error(w, err)
		return
	}

	//not found?
	if study == nil {
		perfect.NotFound(w)
		return
	}

	//parse new form data sent by the client into the original 'study' returned
    //by the database. This has the potential of letting the user overwrite
    //any fields. We're relying on a future call to study.RemoveInternalAttributes()
    //to set all internal attributes to their zero value, so that the database
    //will ignore them (assuming the fields are declared as ',omitempty'
	err = r.ParseJSON(study)
	if err != nil {
		perfect.Error(w, err)
		return
	}

    //restore read-only attributes
    study.RemoveInternalAttributes()

	log.Printf("%s: %#v\n", "JSON study submitted by the client", study)

	//make sure 'Html' stays in sync with 'published'
    if !study.Published {
        study.Html = nil
    }

	//update the study. This needs to be refactored, because it's now sending back
    //the ENTIRE study, when, really, we just need to somehow send only the updates.
    //This would be possible by using map[string]interface{}, but it's not trivial
    //to validate that only the fields that are allowed to be changed, are going to
    //be changed, especially since the names of the fields involved in serialization
    //are specified as struct tags (so we either use 'reflect', or come up with some
    //other system). I'm leaving this issue for another refactoring session.
	err = study.Save(d.Db)
	if err != nil {
		perfect.Error(w, err)
		return
	}

	//success
	perfect.NoContent(w)
}

func (d *Dashboard) DeleteStudy(w http.ResponseWriter, r *perfect.Request) {
	var err error

	//get the form id
	query := r.URL.Query()
	study_id := query.Get("id")

	log.Println("study to delete:", study_id)

	//TODO: Validate alnum
	if len(study_id) == 0 {
		perfect.BadRequest(w)
		return
	}

	//make sure the form exists
	study, err := survana.FindStudy(study_id, d.Db)
	if err != nil {
		perfect.Error(w, err)
		return
	}

	//not found?
	if study == nil {
		perfect.NotFound(w)
		return
	}

	err = study.Delete(d.Db)

	if err != nil {
		perfect.Error(w, err)
		return
	}

	log.Println("study", study_id, "deleted")
	perfect.NoContent(w)
}

func (d *Dashboard) ViewStudyPage(w http.ResponseWriter, r *perfect.Request) {
	d.RenderTemplate(w, "study/view", nil)
}

func (d *Dashboard) PublishStudyPage(w http.ResponseWriter, r *perfect.Request) {
	d.RenderTemplate(w, "study/publish", nil)
}

func (d *Dashboard) PublishStudyForm(w http.ResponseWriter, r *perfect.Request) {

	query := r.URL.Query()
	study_id := query.Get("id")
    form_index, err := strconv.Atoi(query.Get("f"))

	if (len(study_id) == 0) || err != nil || form_index < 0 {
		perfect.BadRequest(w)
		return
	}

	html, err := r.BodyBytes(r.Body)
	if err != nil {
		perfect.Error(w, err)
		return
	}

    study, err := survana.FindStudy(study_id, d.Db)
    if err != nil {
        perfect.Error(w, err)
        return
    }

    //count the total number of forms in the study
    nforms := len(study.Forms)

    log.Println("study=", study, "form_index", form_index, "study.Forms.length=", nforms)

    if study == nil || form_index >= nforms {
        perfect.NotFound(w)
        return
    }

    //make the Html array have the same number of elements as study.Forms
    if (len(study.Html) != nforms) {
        html := make([][]byte, nforms, nforms)
        //preserve any existing elements
        copy(html, study.Html)
        //switch the pointer to the new array
        study.Html = html
    }

    //assign the html
    study.Html[form_index] = html

    //save the study
    err = study.Save(d.Db)
    if err != nil {
        perfect.Error(w, err)
        return
    }

    log.Println("done saving published form")

	perfect.NoContent(w)
}

func (d *Dashboard) StudySubjectsPage(w http.ResponseWriter, r *perfect.Request) {
	d.RenderTemplate(w, "study/subjects", nil)
}

func (d *Dashboard) AddStudySubjects(w http.ResponseWriter, r *perfect.Request) {
	var err error

	//get the form id
	query := r.URL.Query()
	study_id := query.Get("id")

	//TODO: Validate alnum
	if len(study_id) == 0 {
		perfect.BadRequest(w)
		return
	}

	//make sure the form exists
	study, err := survana.FindStudy(study_id, d.Db)
	if err != nil {
		perfect.Error(w, err)
		return
	}

	//not found?
	if study == nil {
		perfect.NotFound(w)
		return
	}

    //parse json data
    ids := make([]string, 0)
    err = r.JSONBody(r.Body, &ids)
    if err != nil {
        perfect.Error(w, err)
        return
    }

    nids := len(ids)
    if nids == 0 {
        perfect.BadRequest(w)
        return
    }

    var id string

    //save and enable all IDs
    for i := 0; i < nids; i++ {
        id = ids[i]
        _, exists := study.Subjects[id]

        if !exists {
            study.AddSubject(id, true)
        }
    }

    //auto-enable authentication for this study
    study.AuthEnabled = true;

    //store the updated Survey
    err = study.Save(d.Db)
    if err != nil {
        perfect.Error(w, err)
        return
    }

    perfect.JSONResult(w, true, study.Subjects)
}
