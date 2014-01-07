package survana

import (
	"testing"
    "time"
)

var mock_form *Form = &Form {
    DBID:   2,
    Id:     "FORM_ABCD",
    Name:   "MockForm",
    Title:  "Mock Form ABCD",
    Version:"1.0",
    CreatedOn: time.Now(),
    Fields: make([]Field, 0),
}

func TestNewForm(t *testing.T) {
    form := NewForm()

    if len(form.Id) != 0 {
		t.Errorf("len(form.Id) = %v ('%v'), want %v", len(form.Id), form.Id, 0)
    }
}

func TestListForms(t *testing.T) {
    db := NewMockDatabase()
    db.OnFilteredList = func (collection string, props []string, result interface{}) {

        if collection != FORM_COLLECTION {
            t.Errorf("db.List() collection supplied was %v, expected %v", collection, FORM_COLLECTION)
        }

        if len(props) == 0 {
            t.Errorf("filter list length is %v, expected non-zero", len(props))
        }

        if result == nil {
            t.Errorf("filter list result is nil, expected pointer")
        }

        list, ok := result.(*[]Form)

        if (!ok) {
            t.Errorf("result variable is %#v, expected pointer")
        }

        if len(*list) > 0 {
            t.Errorf("expected an empty result list, got a list with %v items", len(*list))
        }

        *list = append(*list, *mock_form);
        (*list)[0].Fields = nil;
    }

    forms, err := ListForms(db)

    if db.Calls["FilteredList"] != 1 {
		t.Errorf("db.FilteredList() was called %v time(s), expected %v call(s)", db.Calls["FilteredList"], 1)
    }

    if err != nil {
        t.Errorf("err = %v", err)
    }

    if len(forms) != 1 {
        t.Errorf("number of forms returned is %v, expected %v", len(forms), 1)
    }

    if forms[0].DBID != mock_form.DBID {
		t.Errorf("database id is %v, expected %v", forms[0].DBID, mock_form.DBID)
    }

    if forms[0].Fields != nil {
        t.Errorf("Fields property is %v, expected %v", forms[0].Fields, nil)
    }
}

