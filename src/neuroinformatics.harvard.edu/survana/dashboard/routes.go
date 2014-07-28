package dashboard

import (
	"errors"
	"github.com/vpetrov/perfect"
	"github.com/vpetrov/perfect/auth"
	"net/http"
)

// registers all route handlers
func (dashboard *Dashboard) RegisterHandlers() {

	//must end with slash
	dashboard.Static("/assets/")

	dashboard.Get("/", auth.Protect(dashboard.Index))
	dashboard.Get("/home", auth.Protect(dashboard.Home))
	dashboard.Get("/sidebar", auth.Protect(dashboard.Sidebar))

	/*dashboard.Get("/login/google", d.LoginWithGoogle)
	dashboard.Get("/login/google/response", d.GoogleResponse)
	dashboard.Get("/register", d.Register)
	*/

	//LOGOUT
	dashboard.Get("/logout", dashboard.Auth.Logout)

	//Form
	dashboard.Get("/forms", auth.Protect(dashboard.FormListPage))
	dashboard.Get("/forms/list", auth.Protect(dashboard.FormList))
	dashboard.Get("/forms/create", auth.Protect(dashboard.CreateFormPage))
	dashboard.Post("/forms/create", auth.Protect(dashboard.CreateForm))
	dashboard.Get("/forms/view", auth.Protect(dashboard.ViewFormPage))
	dashboard.Get("/form", auth.Protect(dashboard.GetForm))
	dashboard.Get("/forms/edit", auth.Protect(dashboard.EditFormPage))
	dashboard.Put("/forms/edit", auth.Protect(dashboard.EditForm))
	dashboard.Delete("/form", auth.Protect(dashboard.DeleteForm))

	//Themes
	dashboard.Get("/theme", dashboard.Theme)

	//Study
	dashboard.Get("/studies", auth.Protect(dashboard.StudyListPage))
	dashboard.Get("/studies/list", auth.Protect(dashboard.StudyList))
	dashboard.Get("/studies/create", auth.Protect(dashboard.CreateStudyPage))
	dashboard.Post("/studies/create", auth.Protect(dashboard.CreateStudy))
	dashboard.Get("/studies/view", auth.Protect(dashboard.ViewStudyPage))
	dashboard.Get("/study", auth.Protect(dashboard.GetStudy))
	dashboard.Get("/studies/edit", auth.Protect(dashboard.EditStudyPage))
	dashboard.Put("/studies/edit", auth.Protect(dashboard.EditStudy))
	dashboard.Delete("/study", auth.Protect(dashboard.DeleteStudy))
	dashboard.Get("/studies/publish", auth.Protect(dashboard.PublishStudyPage))
	dashboard.Post("/studies/publish", auth.Protect(dashboard.PublishStudyForm))
	dashboard.Get("/studies/subjects", auth.Protect(dashboard.StudySubjectsPage))
	dashboard.Put("/studies/subjects", auth.Protect(dashboard.AddStudySubjects))
}

// sends the app skeleton to the client
func (d *Dashboard) Index(w http.ResponseWriter, r *perfect.Request) {
	profile, err := r.Profile()
	if err != nil {
		perfect.Error(w, r, err)
		return
	}

	//profile not found?
	if profile == nil {
		perfect.Error(w, r, errors.New("User profile not found"))
		return
	}

	data := &struct {
		Module *perfect.Module
		User   *perfect.Profile
	}{
		Module: d.Module,
		User:   profile,
	}

	d.RenderTemplate(w, r, "index", data)
}

// displays the home page
func (d *Dashboard) Home(w http.ResponseWriter, r *perfect.Request) {
	d.RenderTemplate(w, r, "home", nil)
}

func (d *Dashboard) Sidebar(w http.ResponseWriter, r *perfect.Request) {
	d.RenderTemplate(w, r, "sidebar", nil)
}
