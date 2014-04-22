package dashboard

import (
	"net/http"
    "errors"
	"neuroinformatics.harvard.edu/survana"
    "neuroinformatics.harvard.edu/survana/auth"
)

// registers all route handlers
func (d *Dashboard) RegisterHandlers() {
	app := d.Mux

	//must end with slash
	app.Static("/assets/")

	app.Get("/", auth.Protect(d.Index))
	app.Get("/home", auth.Protect(d.Home))
	app.Get("/sidebar", auth.Protect(d.Sidebar))

	/*app.Get("/login/google", d.LoginWithGoogle)
	app.Get("/login/google/response", d.GoogleResponse)
	app.Get("/register", d.Register)
    */

	//LOGOUT
	app.Get("/logout", d.Auth.Logout)

	//Form
	app.Get("/forms", auth.Protect(d.FormListPage))
	app.Get("/forms/list", auth.Protect(d.FormList))
	app.Get("/forms/create", auth.Protect(d.CreateFormPage))
	app.Post("/forms/create", auth.Protect(d.CreateForm))
	app.Get("/forms/view", auth.Protect(d.ViewFormPage))
	app.Get("/form", auth.Protect(d.GetForm))
	app.Get("/forms/edit", auth.Protect(d.EditFormPage))
	app.Put("/forms/edit", auth.Protect(d.EditForm))
	app.Delete("/form", auth.Protect(d.DeleteForm))

	//Themes
	app.Get("/theme", d.Theme)

	//Study
	app.Get("/studies", auth.Protect(d.StudyListPage))
	app.Get("/studies/list", auth.Protect(d.StudyList))
	app.Get("/studies/create", auth.Protect(d.CreateStudyPage))
	app.Post("/studies/create", auth.Protect(d.CreateStudy))
	app.Get("/studies/view", auth.Protect(d.ViewStudyPage))
	app.Get("/study", auth.Protect(d.GetStudy))
	app.Get("/studies/edit", auth.Protect(d.EditStudyPage))
	app.Put("/studies/edit", auth.Protect(d.EditStudy))
	app.Delete("/study", auth.Protect(d.DeleteStudy))
	app.Get("/studies/publish", auth.Protect(d.PublishStudyPage))
	app.Post("/studies/publish", auth.Protect(d.PublishStudyForm))
    app.Get("/studies/subjects", auth.Protect(d.StudySubjectsPage))
    app.Put("/studies/subjects", auth.Protect(d.AddStudySubjects))
}

// sends the app skeleton to the client
func (d *Dashboard) Index(w http.ResponseWriter, r *survana.Request) {
    user, err := r.User()
    if err != nil {
        survana.Error(w, err)
        return
    }

    //profile not found?
    if user == nil {
        survana.Error(w, errors.New("User profile not found"))
        return
    }

    data := &struct {
                Module *survana.Module
                User *survana.User
            }{
                Module: d.Module,
                User:user,
            }

	d.RenderTemplate(w, "index", data)
}

// displays the home page
func (d *Dashboard) Home(w http.ResponseWriter, r *survana.Request) {
	d.RenderTemplate(w, "home", nil)
}

func (d *Dashboard) Sidebar(w http.ResponseWriter, r *survana.Request) {
	d.RenderTemplate(w, "sidebar", nil)
}
