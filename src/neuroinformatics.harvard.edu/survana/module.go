package survana

import (
	"errors"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	TEMPLATE_DIR = "templates"
	TEMPLATE_EXT = ".html"
)

var (
	SESSION_TIMEOUT time.Duration = time.Hour
)

//A Survana Module is a standalone component that can be mounted on
//URL paths and can decide how requests are routed.
type Module struct {
	Name           string
	MountPoint     string
	Path           string
	SessionTimeout time.Duration

	Db  Database
	Log *log.Logger

	Router    Router
	Templates *template.Template
}

//parses all template files from the 'templates' folder of the module
func (m *Module) ParseTemplates() error {

	log.Println("Parsing templates from", m.Path)

	m.Templates = template.New(m.Name)
    //set start/end tags (delimiters)
    _ = m.Templates.Delims("<%", "%>")

	pathlen := len(m.Path)
	tpldirlen := len(TEMPLATE_DIR)
	tplextlen := len(TEMPLATE_EXT)

	tplParser := func(currentPath string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(currentPath) == TEMPLATE_EXT {
			//the template name is anything after 'path/templates/'
			// '2' is for the 2 slashes in the path above
			from := pathlen + tpldirlen + 2
			to := len(currentPath) - tplextlen
			requestPath := currentPath[from:to]

			//read the template file
			data, err := ioutil.ReadFile(currentPath)
			if err != nil {
				return err
			}

			_, err = m.Templates.New(requestPath).Parse(string(data))
			if err != nil {
				return err
			}
		}

		return nil
	}

	err := filepath.Walk(m.Path, tplParser)

	return err
}

// renders a template file
func (m *Module) RenderTemplate(w http.ResponseWriter, path string, data interface{}) {
	tpl := m.Templates.Lookup(path)
	if tpl == nil {
		Error(w, errors.New("Template not found: "+path))
		return
	}

	err := tpl.Execute(w, data)

	if err != nil {
		Error(w, err)
		return
	}
}
