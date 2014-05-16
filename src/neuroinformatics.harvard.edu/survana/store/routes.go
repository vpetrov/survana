package store

import (
	"log"
	"net/http"
    "encoding/json"
	"neuroinformatics.harvard.edu/survana"
    "github.com/vpetrov/perfect"
)

// registers all route handlers
func (s *Store) RegisterHandlers() {
	app := s.Mux

    log.Println("REGISTERING STORE HANDLERS")

	app.Post("/response", s.NewResponse)
    app.Get("/download", s.Download)
}

func (s *Store) NewResponse(w http.ResponseWriter, r *perfect.Request) {
    var (
            err         error
            study_id    string
            result      map[string]bool
            response    *survana.Response
        )

    query := r.URL.Query()

    study_id = query.Get("s")

    log.Println("new store response request", study_id)

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

    response_queue := make(map[string]json.RawMessage, 0)

    err = r.ParseJSON(&response_queue)
    if err != nil {
        //TODO: need a perfect.ServerError
        log.Println(err)
        perfect.BadRequest(w)
        return
    }

    result = map[string]bool{}

    for r_id, v := range response_queue {
        //unmarshal each response into a Response object
        response = survana.NewResponse()

        err = json.Unmarshal(v, response)
        if err != nil {
            log.Println(err)
            result[r_id] = false
            continue
        }

        log.Println("Saving " + r_id, response)

        //update study_id
        response.StudyId = study_id

        //save response
        //todo: parallelize this loop
        err = response.Save(r.Module.Db)
        if err != nil {
            result[r_id] = false
            log.Println(err)
            continue
        }

        result[r_id] = true
    }

    perfect.JSONResult(w, true, result)
}


func (s *Store) Download(w http.ResponseWriter, r *perfect.Request) {
    var (
            err         error
            study_id    string
        )

    query := r.URL.Query()
    study_id = query.Get("s")

    log.Println("new store download request", study_id)

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

    //fetch all responses for this study
    responses, err := survana.FindResponsesByStudy(study_id, r.Module.Db)
    if err != nil {
        perfect.Error(w, err)
        return
    }

    perfect.JSONResult(w, true, responses)
}
