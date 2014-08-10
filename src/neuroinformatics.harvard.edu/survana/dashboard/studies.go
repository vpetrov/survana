package dashboard

import (
	"github.com/vpetrov/perfect"
	"github.com/vpetrov/perfect/orm"
	"log"
	"net/http"
	"neuroinformatics.harvard.edu/survana"
	"strconv"
	"time"
)

func (d *Dashboard) StudyListPage(w http.ResponseWriter, r *perfect.Request) {
	d.RenderTemplate(w, r, "study/list", nil)
}

func (d *Dashboard) StudyList(w http.ResponseWriter, r *perfect.Request) {
	var (
		err    error
		db     = r.Module.Db
		filter = []string{"id", "name", "title", "description", "version", "created_on", "owner_id", "published", "auth_enabled", "store_url"}
	)

	studies := &[]survana.Study{}
	search := &survana.Study{}
	err = db.Query(search).Select(filter...).All(studies)

	if err != nil {
		perfect.Error(w, r, err)
		return
	}

	perfect.JSONResult(w, r, true, studies)
}

func (d *Dashboard) CreateStudyPage(w http.ResponseWriter, r *perfect.Request) {
	d.RenderTemplate(w, r, "study/create", nil)
}

func (d *Dashboard) CreateStudy(w http.ResponseWriter, r *perfect.Request) {
	var (
		err error
		db  = r.Module.Db
	)

	//get the current session
	session, err := r.Session()
	if err != nil {
		perfect.Error(w, r, err)
		return
	}

	study := &survana.Study{
		CreatedOn: orm.Time(time.Now()),
		OwnerId:   session.ProfileId,
	}

	//parse input data
	err = r.ParseJSON(&study)
	if err != nil {
		perfect.Error(w, r, err)
		return
	}
	//assign default StoreURL
	if study.StoreUrl == nil {
		study.StoreUrl = orm.String(d.Config.StoreUrl)
	}

	//generate a unique id
	err = study.GenerateId(d.Db)
	if err != nil {
		perfect.Error(w, r, err)
		return
	}

	//save the study
	err = db.Save(study)
	if err != nil {
		perfect.Error(w, r, err)
		return
	}

	//result format is { id: "abcd" }
	result := &map[string]string{"id": *study.Id}

	perfect.JSONResult(w, r, true, result)
}

func (d *Dashboard) GetStudy(w http.ResponseWriter, r *perfect.Request) {
	var (
		err      error
		db       = r.Module.Db
		query    = r.URL.Query()
		study_id = query.Get("id")
	)

	//TODO: Validate alnum
	if len(study_id) == 0 {
		perfect.BadRequest(w)
		return
	}

	study := &survana.Study{Id: &study_id}
	err = db.Find(study)
	if err != nil {
		if err == orm.ErrNotFound {
			perfect.NotFound(w)
		} else {
			perfect.Error(w, r, err)
		}
		return
	}

	//return the form as JSON
	perfect.JSONResult(w, r, true, study)
}

func (d *Dashboard) EditStudyPage(w http.ResponseWriter, r *perfect.Request) {
	d.RenderTemplate(w, r, "study/edit", nil)
}

func (d *Dashboard) EditStudy(w http.ResponseWriter, r *perfect.Request) {
	var (
		err      error
		db       = r.Module.Db
		query    = r.URL.Query()
		study_id = query.Get("id")
	)

	//TODO: Validate alnum
	if len(study_id) == 0 {
		perfect.BadRequest(w)
		return
	}

	//make sure the form exists
	study := &survana.Study{Id: &study_id}
	err = db.Peek(study)
	if err != nil {
		if err == orm.ErrNotFound {
			perfect.NotFound(w)
		} else {
			perfect.Error(w, r, err)
		}
		return
	}

	//parse new form data sent by the client into the original 'study' returned
	//by the database.
	err = r.ParseJSON(study)
	if err != nil {
		perfect.Error(w, r, err)
		return
	}

	log.Printf("%s: %#v\n", "JSON study submitted by the client", study)

	//update the study.
	err = db.Save(study)
	if err != nil {
		perfect.Error(w, r, err)
		return
	}

	//success
	perfect.NoContent(w)
}

func (d *Dashboard) DeleteStudy(w http.ResponseWriter, r *perfect.Request) {
	var (
		err      error
		db       = r.Module.Db
		query    = r.URL.Query()
		study_id = query.Get("id")
	)

	log.Println("study to delete:", study_id)

	//TODO: Validate alnum
	if len(study_id) == 0 {
		perfect.BadRequest(w)
		return
	}

	//make sure the form exists
	study := &survana.Study{Id: &study_id}
	err = db.Peek(study)
	if err != nil {
		if err == orm.ErrNotFound {
			perfect.NotFound(w)
		} else {
			perfect.Error(w, r, err)
		}
		return
	}

	err = db.Remove(study)

	if err != nil {
		perfect.Error(w, r, err)
		return
	}

	log.Println("study", study_id, "deleted")
	perfect.NoContent(w)
}

func (d *Dashboard) ViewStudyPage(w http.ResponseWriter, r *perfect.Request) {
	d.RenderTemplate(w, r, "study/view", nil)
}

func (d *Dashboard) PublishStudyPage(w http.ResponseWriter, r *perfect.Request) {
	d.RenderTemplate(w, r, "study/publish", nil)
}

func (d *Dashboard) PublishStudyForm(w http.ResponseWriter, r *perfect.Request) {

	var (
		err      error
		db       = r.Module.Db
		query    = r.URL.Query()
		study_id = query.Get("id")
	)

	form_index, err := strconv.Atoi(query.Get("f"))

	if (len(study_id) == 0) || err != nil || form_index < 0 {
		perfect.BadRequest(w)
		return
	}

	html, err := r.BodyBytes(r.Body)
	if err != nil {
		perfect.Error(w, r, err)
		return
	}

	study := &survana.Study{Id: &study_id}
	err = db.Query(study).Select("html", "forms").One(study)
	if err != nil {
		if err == orm.ErrNotFound {
			perfect.NotFound(w)
		} else {
			perfect.Error(w, r, err)
		}
		return
	}

	//count the total number of forms in the study
	var nforms int = 0
	if study.Forms != nil {
		nforms = len(*study.Forms)
	}

	log.Println("study=", study, "form_index", form_index, "study.Forms.length=", nforms)

	if nforms == 0 || form_index >= nforms {
		perfect.BadRequest(w)
		return
	}

	//make the Html array have the same number of elements as study.Forms
	if study.Html == nil || len(*study.Html) != nforms {
		html := make([][]byte, nforms, nforms)
		//preserve any existing elements
		if study.Html != nil {
			copy(html, *study.Html)
		}
		//switch the pointer to the new array
		study.Html = &html
	}

	//assign the html
	(*study.Html)[form_index] = html

	//avoid re-sending the Forms array to the DB
	study.Forms = nil

	//save the study
	err = db.Save(study)
	if err != nil {
		perfect.Error(w, r, err)
		return
	}

	log.Println("done saving published form")

	perfect.NoContent(w)
}

func (d *Dashboard) StudySubjectsPage(w http.ResponseWriter, r *perfect.Request) {
	d.RenderTemplate(w, r, "study/subjects", nil)
}

func (d *Dashboard) AddStudySubjects(w http.ResponseWriter, r *perfect.Request) {
	var (
		err      error
		db       = r.Module.Db
		query    = r.URL.Query()
		study_id = query.Get("id")
	)

	//TODO: Validate alnum
	if len(study_id) == 0 {
		perfect.BadRequest(w)
		return
	}

	//make sure the form exists
	study := &survana.Study{Id: &study_id}
	err = db.Query(study).Select("subjects", "auth_enabled").One(study)
	if err != nil {
		if err == orm.ErrNotFound {
			perfect.NotFound(w)
		} else {
			perfect.Error(w, r, err)
		}
		return
	}

	//parse json data
	ids := make([]string, 0)
	err = r.JSONBody(r.Body, &ids)
	if err != nil {
		perfect.Error(w, r, err)
		return
	}

	nids := len(ids)
	if nids == 0 {
		perfect.BadRequest(w)
		return
	}

	var id string

	if study.Subjects == nil {
		study.Subjects = &map[string]bool{}
	}

	//save and enable all IDs
	for i := 0; i < nids; i++ {
		id = ids[i]
		_, exists := (*study.Subjects)[id]

		if !exists {
			(*study.Subjects)[id] = true
		}
	}

	//auto-enable authentication for this study
	study.AuthEnabled = orm.Bool(true)

	//store the updated Survey
	err = db.Save(study)
	if err != nil {
		perfect.Error(w, r, err)
		return
	}

	perfect.JSONResult(w, r, true, study.Subjects)
}
