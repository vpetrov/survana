package dashboard

import (
        _ "log"
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
    app.Get("/forms", survana.Protect(d.FormList))
    app.Get("/forms/create", survana.Protect(d.CreateForm))

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

func (d *Dashboard) FormList(w http.ResponseWriter, r *survana.Request) {
    d.RenderTemplate(w, "form/list", nil)
}

func (d *Dashboard) CreateForm(w http.ResponseWriter, r *survana.Request) {
    time.Sleep(time.Second * 4)
    d.RenderTemplate(w, "form/create", nil)
}
