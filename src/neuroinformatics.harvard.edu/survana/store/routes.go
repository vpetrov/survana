package store

import (
	"encoding/json"
	"github.com/vpetrov/perfect"
	"github.com/vpetrov/perfect/orm"
	"log"
	"net/http"
	"neuroinformatics.harvard.edu/survana"
)

// registers all route handlers
func (store *Store) RegisterHandlers() {
	log.Println("REGISTERING STORE HANDLERS")

	store.Post("/response", store.NewResponse)
	store.Get("/download", store.Download)
}

func (store *Store) NewResponse(w http.ResponseWriter, r *perfect.Request) {
	var (
		err      error
		study_id string
		db       orm.Database = r.Module.Db
	)

	query := r.URL.Query()

	study_id = query.Get("s")

	log.Println("new store response request", study_id)

	if len(study_id) == 0 {
		perfect.BadRequest(w)
		return
	}

	//otherwise, fetch the study
	study := &survana.Study{Id: &study_id}
	err = db.Find(study)
	if err != nil {
		if err == orm.ErrNotFound {
			perfect.NotFound(w)
		} else {
			perfect.Error(w, r, err)
		}
		return
	}

	response_queue := &map[string]json.RawMessage{}

	err = r.ParseJSON(response_queue)
	if err != nil {
		//TODO: need a perfect.ServerError
		log.Println(err)
		perfect.BadRequest(w)
		return
	}

	var (
		response *survana.Response
		result   = map[string]bool{}
	)

	for r_id, v := range *response_queue {
		//unmarshal each response into a Response object
		response = &survana.Response{}

		err = json.Unmarshal(v, response)
		if err != nil {
			log.Println(err)
			result[r_id] = false
			continue
		}

		log.Println("Saving "+r_id, response)

		//update study_id
		response.StudyId = &study_id

		//save response
		//todo: parallelize this loop
		err = db.Save(response)
		if err != nil {
			result[r_id] = false
			log.Println(err)
			continue
		}

		result[r_id] = true
	}

	perfect.JSONResult(w, r, true, result)
}

func (s *Store) Download(w http.ResponseWriter, r *perfect.Request) {
	var (
		err      error
		study_id string
		db       = r.Module.Db
	)

	query := r.URL.Query()
	study_id = query.Get("s")

	log.Println("new store download request", study_id)

	if len(study_id) == 0 {
		perfect.BadRequest(w)
		return
	}

	//otherwise, fetch the study
	study := &survana.Study{Id: &study_id}
	err = db.Find(study)
	if err != nil {
		if err == orm.ErrNotFound {
			perfect.NotFound(w)
		} else {
			perfect.Error(w, r, err)
		}
		return
	}

	//fetch all responses for this study
	var (
		responses = &[]survana.Response{}
		search    = &survana.Response{StudyId: &study_id}
	)

	err = db.Query(search).All(responses)
	if err != nil {
		perfect.Error(w, r, err)
		return
	}

	perfect.JSONResult(w, r, true, responses)
}
