package survana

import (
        "github.com/vpetrov/perfect"
        "encoding/json"
       )

const (
        RESPONSE_COLLECTION = "responses"
      )

type Response struct {
    perfect.DBO             `bson:",inline,omitempty" json:"-"`
	Id          string      `bson:"id,omitempty" json:"id"`
    StudyId     string      `bson:"study_id,omitempty" json:"study_id"`
    Data        interface{} `bson:"data,omitempty" json:"data"`
}

func NewResponse() *Response {
    return &Response{
        DBO: perfect.DBO{ Collection: RESPONSE_COLLECTION },
    }
}

func (r *Response) Save(db perfect.Database) error {
    return db.Save(r)
}

func (r *Response) UnmarshalJSON(data []byte) error {
    return json.Unmarshal(data, &r.Data)
}

func FindResponsesByStudy(study_id string, db perfect.Database) (result []*Response, err error) {
    result = make([]*Response, 0)
    err = db.SearchByField(RESPONSE_COLLECTION, "study_id", study_id, nil, &result)
    if err != nil {
        return
    }

    return
}
