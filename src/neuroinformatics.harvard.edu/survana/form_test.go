package survana

import (
	"github.com/vpetrov/perfect/orm"
	"testing"
	"time"
)

var now = time.Now()

var mock_form *Form = &Form{
	Object:    orm.Object{Id: 2},
	Id:        orm.String("FORM_ABCD"),
	Name:      orm.String("MockForm"),
	Title:     orm.String("Mock Form ABCD"),
	Version:   orm.String("1.0"),
	CreatedOn: &now,
	Fields:    &[]Field{},
}
