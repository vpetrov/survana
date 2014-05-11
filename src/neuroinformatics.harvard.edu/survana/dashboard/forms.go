package dashboard

import (
	"log"
	"net/http"
	"github.com/vpetrov/perfect"
    "neuroinformatics.harvard.edu/survana"
	"time"
)

func (d *Dashboard) FormListPage(w http.ResponseWriter, r *perfect.Request) {
	d.RenderTemplate(w, "form/list", nil)
}

func (d *Dashboard) FormList(w http.ResponseWriter, r *perfect.Request) {

	filter := []string{"id", "name", "title", "version", "created_on", "owner_id"}

	query := r.URL.Query()
	_ = query.Get("ids")

	//decide whether the 'fields' property should be returned
	fields := query.Get("fields")
	if fields == "true" {
		filter = append(filter, "fields")
	}

	forms, err := survana.ListForms(filter, d.Db)

	if err != nil {
		perfect.Error(w, err)
		return
	}

	perfect.JSONResult(w, true, forms)
}

func (d *Dashboard) CreateFormPage(w http.ResponseWriter, r *perfect.Request) {
	d.RenderTemplate(w, "form/create", nil)
}

func (d *Dashboard) CreateForm(w http.ResponseWriter, r *perfect.Request) {
	var err error

	//get the current session
	session, err := r.Session()
	if err != nil {
		perfect.Error(w, err)
	}

	form := survana.NewForm()

	//parse input data
	err = r.ParseJSON(&form)
	if err != nil {
		perfect.Error(w, err)
		return
	}

	form.CreatedOn = time.Now()
	form.OwnerId = session.UserId

	//generate a unique id
	err = form.GenerateId(d.Db)
	if err != nil {
		perfect.Error(w, err)
	}

	//save the form
	err = form.Save(d.Db)
	if err != nil {
		perfect.Error(w, err)
		return
	}

	//result format is { id: "abcd" }
	result := &struct {
		Id string `json:"id"`
	}{Id: form.Id}

	perfect.JSONResult(w, true, result)
}

func (d *Dashboard) ViewFormPage(w http.ResponseWriter, r *perfect.Request) {
	d.RenderTemplate(w, "form/view", nil)
}

func (d *Dashboard) GetForm(w http.ResponseWriter, r *perfect.Request) {
	query := r.URL.Query()
	form_id := query.Get("id")

	//TODO: Validate alnum
	if len(form_id) == 0 {
		perfect.BadRequest(w)
		return
	}

	form, err := survana.FindForm(form_id, d.Db)
	if err != nil {
		perfect.Error(w, err)
		return
	}

	//not found?
	if form == nil {
		perfect.NotFound(w)
		return
	}

	//return the form as JSON
	perfect.JSONResult(w, true, form)
}

func (d *Dashboard) EditFormPage(w http.ResponseWriter, r *perfect.Request) {
	d.RenderTemplate(w, "form/edit", nil)
}

func (d *Dashboard) EditForm(w http.ResponseWriter, r *perfect.Request) {
	var err error

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
	form, err := survana.FindForm(form_id, d.Db)
	if err != nil {
		perfect.Error(w, err)
		return
	}

	//not found?
	if form == nil {
		perfect.NotFound(w)
		return
	}

	log.Println("form", form_id, "was found:")
	log.Println(form)

	//parse new form data sent by the client
	user_form := survana.NewForm()
	err = r.ParseJSON(user_form)
	if err != nil {
		perfect.Error(w, err)
		return
	}

	//TODO?: validate form fields? validate using a schema?

	log.Printf("%s: %#v\n", "JSON form submitted by the client", user_form)

	//copy properties that should not be changed
	user_form.DBID = form.DBID
	user_form.Id = form.Id
	user_form.CreatedOn = form.CreatedOn
	user_form.OwnerId = form.OwnerId

	//update the form
	err = user_form.Save(d.Db)
	if err != nil {
		perfect.Error(w, err)
		return
	}

	//success
	perfect.NoContent(w)
}

func (d *Dashboard) DeleteForm(w http.ResponseWriter, r *perfect.Request) {
	var err error

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
	form, err := survana.FindForm(form_id, d.Db)
	if err != nil {
		perfect.Error(w, err)
		return
	}

	//not found?
	if form == nil {
		perfect.NotFound(w)
		return
	}

    log.Printf("form=%#v", form)

	err = form.Delete(d.Db)

	if err != nil {
		perfect.Error(w, err)
		return
	}

	log.Println("form", form_id, "deleted")
	perfect.NoContent(w)
}
