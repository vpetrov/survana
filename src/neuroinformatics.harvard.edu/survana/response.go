package survana

import (
        "github.com/vpetrov/perfect"
        "encoding/json"
        "log"
       )

const (
        RESPONSE_COLLECTION = "responses"
      )

type Response struct {
    perfect.DBO             `bson:",inline,omitempty" json:"-"`
	Id          string      `bson:"id,omitempty" json:"id"`
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
    log.Println("UNMARSHALLING:",string(data))
    err := json.Unmarshal(data, &r.Data)

    log.Println("RESULT:", r.Data)

    return err
}
