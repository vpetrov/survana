package dashboard

import (
	"log"
	"net/http"
	"neuroinformatics.harvard.edu/survana"
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
	app.Get("/login", d.BuiltinAuthPage)
    //TODO:app.Post("/login", d.Login)
    app.Post("/login", d.BuiltinAuth)

    //Registration is optional
    if (d.Config.AllowRegistration) {
        app.Get("/register", d.BuiltinRegistrationPage)
        app.Post("/register", d.BuiltinRegister)
    }

	/*app.Get("/login/google", d.LoginWithGoogle)
	app.Get("/login/google/response", d.GoogleResponse)
	app.Get("/register", d.Register)
    */


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
	app.Get("/studies", survana.Protect(d.StudyListPage))
	app.Get("/studies/list", survana.Protect(d.StudyList))
	app.Get("/studies/create", survana.Protect(d.CreateStudyPage))
	app.Post("/studies/create", survana.Protect(d.CreateStudy))
	app.Get("/studies/view", survana.Protect(d.ViewStudyPage))
	app.Get("/study", survana.Protect(d.GetStudy))
	app.Get("/studies/edit", survana.Protect(d.EditStudyPage))
	app.Put("/studies/edit", survana.Protect(d.EditStudy))
	app.Delete("/study", survana.Protect(d.DeleteStudy))
	app.Get("/studies/publish", survana.Protect(d.PublishStudyPage))
	app.Post("/studies/publish", survana.Protect(d.PublishStudyForm))

}

// sends the app skeleton to the client
func (d *Dashboard) Index(w http.ResponseWriter, r *survana.Request) {
    log.Println("DASHBOARD INDEX")
	d.RenderTemplate(w, "index", nil)
}

// displays the home page
func (d *Dashboard) Home(w http.ResponseWriter, r *survana.Request) {
	d.RenderTemplate(w, "home", nil)
}

func (d *Dashboard) Sidebar(w http.ResponseWriter, r *survana.Request) {
	d.RenderTemplate(w, "sidebar", nil)
}
