package survana

import (
	"encoding/json"
	"github.com/vpetrov/perfect/orm"
)

type Response struct {
	orm.Object `bson:",inline,omitempty" json:"-"`
	Id         *string     `bson:"id,omitempty" json:"id,omitempty"`
	StudyId    *string     `bson:"study_id,omitempty" json:"study_id,omitempty"`
	Data       interface{} `bson:"data,omitempty" json:"data,omitempty"`
}

func (r *Response) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &r.Data)
}
