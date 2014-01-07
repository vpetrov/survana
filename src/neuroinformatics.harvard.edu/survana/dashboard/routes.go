package dashboard

import (
	"log"
	"net/http"
	"neuroinformatics.harvard.edu/survana"
	"time"
)

// registers all route handlers
func (d *Dashboard) RegisterHandlers() {
	app := d.mux

	//must end with slash
	app.Static("/assets/")

	app.Get("/", survana.Protect(d.Index))
	app.Get("/home", survana.Protect(d.Home))
	app.Get("/sidebar", survana.Protect(d.Sidebar))

	//LOGIN
	app.Get("/login", d.Login)
	app.Get("/login/google", d.LoginWithGoogle)
	app.Get("/login/google/response", d.GoogleResponse)
	app.Get("/register", d.Register)

	//LOGOUT
	app.Get("/logout", d.Logout)

	//Form
	app.Get("/forms", survana.Protect(d.FormListPage))
	app.Get("/forms/list", survana.Protect(d.FormList))
	app.Get("/forms/create", survana.Protect(d.CreateFormPage))
	app.Post("/forms/create", survana.Protect(d.CreateForm))
	app.Get("/forms/view", survana.Protect(d.ViewFormPage))
	app.Get("/form", survana.Protect(d.GetForm))
	app.Get("/forms/edit", survana.Protect(d.EditFormPage))
	app.Put("/forms/edit", survana.Protect(d.EditForm))
	app.Delete("/form", survana.Protect(d.DeleteForm))

	//Themes
	app.Get("/theme", d.Theme)

	//Study
	app.Get("/studies", survana.Protect(d.StudyList))
}

// sends the app skeleton to the client
func (d *Dashboard) Index(w http.ResponseWriter, r *survana.Request) {
	d.RenderTemplate(w, "index", nil)
}

// displays the home page
func (d *Dashboard) Home(w http.ResponseWriter, r *survana.Request) {
	d.RenderTemplate(w, "home", nil)
}

func (d *Dashboard) Sidebar(w http.ResponseWriter, r *survana.Request) {
	d.RenderTemplate(w, "sidebar", nil)
}

func (d *Dashboard) StudyList(w http.ResponseWriter, r *survana.Request) {
	d.RenderTemplate(w, "study/list", nil)
}

func (d *Dashboard) FormListPage(w http.ResponseWriter, r *survana.Request) {
	d.RenderTemplate(w, "form/list", nil)
}

func (d *Dashboard) FormList(w http.ResponseWriter, r *survana.Request) {
	forms, err := survana.ListForms(d.Module.Db)

	if err != nil {
		survana.Error(w, err)
		return
	}

	survana.JSONResult(w, true, forms)
}

func (d *Dashboard) CreateFormPage(w http.ResponseWriter, r *survana.Request) {
	d.RenderTemplate(w, "form/create", nil)
}

func (d *Dashboard) CreateForm(w http.ResponseWriter, r *survana.Request) {
	var err error

	form := survana.Form{}

	//parse input data
	err = r.ParseJSON(&form)
	if err != nil {
		survana.Error(w, err)
		return
	}

	form.CreatedOn = time.Now()

	//generate a unique id
	err = form.GenerateId(d.Module.Db)
	if err != nil {
		survana.Error(w, err)
	}

	//save the form
	err = form.Save(d.Module.Db)
	if err != nil {
		survana.Error(w, err)
		return
	}

	//result format is { id: "abcd" }
	result := &struct {
		Id string `json:"id"`
	}{Id: form.Id}

	survana.JSONResult(w, true, result)
}

func (d *Dashboard) ViewFormPage(w http.ResponseWriter, r *survana.Request) {
	d.RenderTemplate(w, "form/view", nil)
}

func (d *Dashboard) GetForm(w http.ResponseWriter, r *survana.Request) {
	query := r.URL.Query()
	form_id := query.Get("id")

	//TODO: Validate alnum
	if len(form_id) == 0 {
		survana.BadRequest(w)
		return
	}

	form, err := survana.FindForm(form_id, d.Module.Db)
	if err != nil {
		survana.Error(w, err)
		return
	}

	//not found?
	if form == nil {
		survana.NotFound(w)
		return
	}

	//return the form as JSON
	survana.JSONResult(w, true, form)
}

func (d *Dashboard) EditFormPage(w http.ResponseWriter, r *survana.Request) {
	d.RenderTemplate(w, "form/edit", nil)
}

func (d *Dashboard) EditForm(w http.ResponseWriter, r *survana.Request) {
	var err error

	//get the form id
	query := r.URL.Query()
	form_id := query.Get("id")

	log.Println("form to update:", form_id)

	//TODO: Validate alnum
	if len(form_id) == 0 {
		survana.BadRequest(w)
		return
	}

	//make sure the form exists
	form, err := survana.FindForm(form_id, d.Module.Db)
	if err != nil {
		survana.Error(w, err)
		return
	}

	//not found?
	if form == nil {
		survana.NotFound(w)
		return
	}

	log.Println("form", form_id, "was found:")
	log.Println(form)

	//parse new form data sent by the client
	user_form := &survana.Form{}
	err = r.ParseJSON(user_form)
	if err != nil {
		survana.Error(w, err)
		return
	}

	//TODO?: validate form fields? validate using a schema?

	log.Printf("%s: %#v\n", "JSON form submitted by the client", user_form)

	//copy properties that should not be changed
	user_form.DBID = form.DBID
	user_form.Id = form.Id
	user_form.CreatedOn = form.CreatedOn

	//update the form
	err = user_form.Save(d.Module.Db)
	if err != nil {
		survana.Error(w, err)
		return
	}

	//success
	survana.NoContent(w)
}

func (d *Dashboard) DeleteForm(w http.ResponseWriter, r *survana.Request) {
	var err error

	//get the form id
	query := r.URL.Query()
	form_id := query.Get("id")

	log.Println("form to delete:", form_id)

	//TODO: Validate alnum
	if len(form_id) == 0 {
		survana.BadRequest(w)
		return
	}

	//make sure the form exists
	form, err := survana.FindForm(form_id, d.Module.Db)
	if err != nil {
		survana.Error(w, err)
		return
	}

	//not found?
	if form == nil {
		survana.NotFound(w)
		return
	}

	err = form.Delete(d.Module.Db)

	if err != nil {
		survana.Error(w, err)
		return
	}

	log.Println("form", form_id, "deleted")
	survana.NoContent(w)
}

func (d *Dashboard) Theme(w http.ResponseWriter, r *survana.Request) {
	template_name := "index"

	//get the form id
	query := r.URL.Query()
	theme_id := query.Get("id")
	theme_version := query.Get("version")
	theme_preview := query.Get("preview")

	log.Println("theme:", theme_id, theme_version)

	//TODO: Validate alnum
	if (len(theme_id) == 0) || (len(theme_version) == 0) {
		survana.BadRequest(w)
		return
	}

	if len(theme_preview) != 0 {
		template_name = "preview"
	}

	d.Module.RenderTemplate(w, "theme/"+theme_id+"/"+theme_version+"/"+template_name, nil)
}
