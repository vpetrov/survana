package survana

import (
	"testing"
)

var mock_session *Session = &Session{
	DBID:          1,
	Id:            "ABCD",
	Authenticated: true,
	Values:        map[string]string{"id": "ABCD", "authenticated": "1"},
}

func TestNewSession(t *testing.T) {
	session := NewSession()

	if len(session.Id) != 0 {
		t.Errorf("len(session.Id) = %v ('%v'), want %v", len(session.Id), session.Id, 0)
	}

	if session.Authenticated {
		t.Errorf("session.Authenticated is %v, want %v", session.Authenticated, false)
	}

	if session.Values == nil {
		t.Errorf("session.Values is not allocated")
	}
}

func TestFindSession(t *testing.T) {
	db := NewMockDatabase()
	db.OnFindId = func(id string, v DbObject) {
		s, ok := v.(*Session)
		if !ok {
			t.Fatalf("MockDatabase did not return a *Session")
		}

		*s = *mock_session
	}

	session, err := FindSession("ABCD", db)

	if db.Calls["FindId"] != 1 {
		t.Errorf("db.FindId() was called %v time(s), expected %v call(s)", db.Calls["FindId"], 1)
	}

	if err != nil {
		t.Errorf("err = %v", err)
	}

	if session.DBID != mock_session.DBID {
		t.Errorf("database id is %v, expected %v", session.DBID, mock_session.DBID)
	}

	if session.Id != mock_session.Id {
		t.Errorf("len(session.Id) is %v, want %v", session.Id, mock_session.Id)
	}

	if !session.Authenticated {
		t.Errorf("session.Authenticated is %v, want %v", session.Authenticated, true)
	}
}

func TestDelete(t *testing.T) {
	//create new database mock
	db := NewMockDatabase()

	//create a new session object
	session := NewSession()

	//attempt to delete the session
	err := session.Delete(db)

	if db.Calls["Delete"] != 1 {
		t.Errorf("db.Delete() was called %v time(s), expected %v calls(s)", db.Calls["Delete"], 1)
	}

	if err != nil {
		t.Errorf("err = %v", err)
	}
}

func TestSave(t *testing.T) {
	//create new database mock
	db := NewMockDatabase()

	//create a new session object
	session := NewSession()

	//attempt to delete the session
	err := session.Save(db)

	if db.Calls["Save"] != 1 {
		t.Errorf("db.Save() was called %v time(s), expected %v calls(s)", db.Calls["Save"], 1)
	}

	if err != nil {
		t.Errorf("err = %v", err)
	}
}

