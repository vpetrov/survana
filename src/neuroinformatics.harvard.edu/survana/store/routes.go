package store

import (
	"log"
	"net/http"
	"neuroinformatics.harvard.edu/survana"
    "github.com/vpetrov/perfect"
)

// registers all route handlers
func (s *Store) RegisterHandlers() {
	app := s.Mux

    log.Println("REGISTERING STORE HANDLERS")

	app.Get("/response", s.NewResponse)
}

func (s *Store) NewResponse(w http.ResponseWriter, r *perfect.Request) {
    log.Println("[store] handling", r.URL)
    var (
            err         error
            study_id    string
            //subject_id  string
            result      []string
        )

    query := r.URL.Query()


    study_id = query.Get("s")

    log.Println("new store response request", study_id)
    //subject_id = query.Get("id")

    if len(study_id) == 0 {
        perfect.BadRequest(w)
        return
    }

    //otherwise, fetch the study
    study, err := survana.FindStudy(study_id, s.Db)
    if err != nil {
        perfect.Error(w, err)
        return
    }

    if study == nil {
        perfect.NotFound(w)
        return
    }

    response_queue := make(map[string]*survana.Response, 0)

    err = r.ParseJSON(response_queue)
    if err != nil {
        perfect.Error(w, err)
        return
    }

    result = make([]string, 0)
    for r_id, v := range response_queue {
        log.Println("Saving " + r_id, v)
    }

    perfect.JSONResult(w, true, result)
}

