package dashboard

import (
	"github.com/vpetrov/perfect"
	"net/http"
)

func (d *Dashboard) Theme(w http.ResponseWriter, r *perfect.Request) {
	template_name := "index"

	//get the form id
	query := r.URL.Query()
	theme_id := query.Get("id")
	theme_version := query.Get("version")
	theme_publish := query.Get("publish")
	theme_study := query.Get("study")

	//TODO: Validate alnum
	if (len(theme_id) == 0) || (len(theme_version) == 0) {
		perfect.BadRequest(w)
		return
	}

	if len(theme_publish) != 0 {
		template_name = "publish"
	} else {
		template_name = "preview"
	}

	if len(theme_study) != 0 {
		template_name += "-study"
	}

	d.Module.RenderTemplate(w, "theme/"+theme_id+"/"+theme_version+"/"+template_name, nil)
}
