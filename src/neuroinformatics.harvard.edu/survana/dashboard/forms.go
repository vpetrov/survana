package dashboard

import (
	"github.com/vpetrov/perfect"
	"github.com/vpetrov/perfect/orm"
	"log"
	"net/http"
	"neuroinformatics.harvard.edu/survana"
	"time"
)

func (d *Dashboard) FormListPage(w http.ResponseWriter, r *perfect.Request) {
	d.RenderTemplate(w, r, "form/list", nil)
}

func (d *Dashboard) FormList(w http.ResponseWriter, r *perfect.Request) {
	var db = r.Module.Db

	//filter := []string{"id", "name", "title", "version", "created_on", "owner_id"}

	query := r.URL.Query()
	_ = query.Get("ids")

	//decide whether the 'fields' property should be returned
	/*fields := query.Get("fields")
	if fields == "true" {
		filter = append(filter, "fields")
	}*/

	forms := &[]survana.Form{}
	search := &survana.Form{}
	err := db.Query(search).All(forms)

	if err != nil && err != orm.ErrNotFound {
		perfect.Error(w, r, err)
		return
	}

	perfect.JSONResult(w, r, true, forms)
}

func (d *Dashboard) CreateFormPage(w http.ResponseWriter, r *perfect.Request) {
	d.RenderTemplate(w, r, "form/create", nil)
}

func (d *Dashboard) CreateForm(w http.ResponseWriter, r *perfect.Request) {
	var (
		err error
		db  = r.Module.Db
	)

	//get the current session
	session, err := r.Session()
	if err != nil {
		perfect.Error(w, r, err)
	}

	form := &survana.Form{}

	//parse input data
	err = r.ParseJSON(form)
	if err != nil {
		perfect.Error(w, r, err)
		return
	}

	now := time.Now()
	form.CreatedOn = &now
	form.OwnerId = session.ProfileId

	//generate a unique id
	err = form.GenerateId(db)
	if err != nil {
		perfect.Error(w, r, err)
	}

	//save the form
	err = db.Save(form)
	if err != nil {
		perfect.Error(w, r, err)
		return
	}

	//result format is { id: "abcd" }
	result := &struct {
		Id *string `json:"id"`
	}{Id: form.Id}

	perfect.JSONResult(w, r, true, result)
}

func (d *Dashboard) ViewFormPage(w http.ResponseWriter, r *perfect.Request) {
	d.RenderTemplate(w, r, "form/view", nil)
}

func (d *Dashboard) GetForm(w http.ResponseWriter, r *perfect.Request) {
	var (
		err error
		db  = r.Module.Db
	)
	query := r.URL.Query()
	form_id := query.Get("id")

	//TODO: Validate alnum
	if len(form_id) == 0 {
		perfect.BadRequest(w)
		return
	}

	form := &survana.Form{Id: &form_id}
	err = db.Find(form)
	if err != nil {
		if err == orm.ErrNotFound {
			perfect.NotFound(w)
		} else {
			perfect.Error(w, r, err)
		}
		return
	}

	//return the form as JSON
	perfect.JSONResult(w, r, true, form)
}

func (d *Dashboard) EditFormPage(w http.ResponseWriter, r *perfect.Request) {
	d.RenderTemplate(w, r, "form/edit", nil)
}

func (d *Dashboard) EditForm(w http.ResponseWriter, r *perfect.Request) {
	var (
		err error
		db  = r.Module.Db
	)

	//get the form id
	query := r.URL.Query()
	form_id := query.Get("id")

	log.Println("form to update:", form_id)

	//TODO: Validate alnum
	if len(form_id) == 0 {
		perfect.BadRequest(w)
		return
	}

	//make sure the form exists
	form := &survana.Form{Id: &form_id}
	err = db.Find(form)
	if err != nil {
		if err == orm.ErrNotFound {
			perfect.NotFound(w)
		} else {
			perfect.Error(w, r, err)
		}
		return
	}

	log.Println("form", form_id, "was found:")
	log.Println(form)

	//parse new form data sent by the client
	user_form := &survana.Form{}
	err = r.ParseJSON(user_form)
	if err != nil {
		perfect.Error(w, r, err)
		return
	}

	//TODO?: validate form fields? validate using a schema?

	log.Printf("%s: %#v\n", "JSON form submitted by the client", user_form)

	//unset properties that should not be changed
	user_form.Object = form.Object
	user_form.Id = nil
	user_form.CreatedOn = nil
	user_form.OwnerId = nil

	//update the form
	err = db.Save(user_form)
	if err != nil {
		perfect.Error(w, r, err)
		return
	}

	//success
	perfect.NoContent(w)
}

func (d *Dashboard) DeleteForm(w http.ResponseWriter, r *perfect.Request) {
	var (
		err error
		db  = r.Module.Db
	)

	//get the form id
	query := r.URL.Query()
	form_id := query.Get("id")

	log.Println("form to delete:", form_id)

	//TODO: Validate alnum
	if len(form_id) == 0 {
		perfect.BadRequest(w)
		return
	}

	//make sure the form exists
	form := &survana.Form{Id: &form_id}
	err = db.Find(form)
	if err != nil {
		if err == orm.ErrNotFound {
			perfect.NotFound(w)
		} else {
			perfect.Error(w, r, err)
		}
		return
	}

	log.Printf("form=%#v", form)

	err = db.Remove(form)

	if err != nil {
		perfect.Error(w, r, err)
		return
	}

	log.Println("form", form_id, "deleted")
	perfect.NoContent(w)
}
