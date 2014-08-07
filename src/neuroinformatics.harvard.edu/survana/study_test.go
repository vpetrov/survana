package survana

import (
	"github.com/vpetrov/perfect/json"
	"github.com/vpetrov/perfect/orm"
	"reflect"
	"testing"
	"time"
)

type testCase struct {
	json string
}

var (
	actual = &Study{
		Object:    orm.Object{Id: "123"},
		Id:        orm.String("xyz"),
		CreatedOn: orm.Time(time.Date(2014, 7, 1, 12, 12, 12, 12, time.Local)),
		OwnerId:   orm.String("root@localhost"),
		StoreUrl:  orm.String("https://localhost/store"),
	}

	expected = &Study{
		Object:    orm.Object{Id: "123"},
		Id:        orm.String("xyz"),
		CreatedOn: orm.Time(time.Date(2014, 7, 1, 12, 12, 12, 12, time.Local)),
		OwnerId:   orm.String("root@localhost"),
		StoreUrl:  orm.String("https://localhost/store"),
	}
)

var test_cases = []testCase{
	{json: `{}`},
	{json: `{"id":"abc"}`},
	{json: `{"created_on":"2014-07-01 13:13:13.000000013 -0400 EDT"}`},
	{json: `{"id":"abc", "created_on":"2014-07-01 13:13:13.000000013 -0400 EDT"}`},
	{json: `{"owner_id":"admin@example.com"}`},
	{json: `{"id":"abc123", "owner_id":"admin@example.com"}`},
	{json: `{"id":"abc123", "owner_id":"admin@example.com", "created_on":"2014-07-01 13:13:13.000000013 -0400 EDT"}`},
	{json: `{"store_url":"http://example.com/store"}`},
	{json: `{"id":"abc1234", "store_url":"http://example.com/store"}`},
	{json: `{"id":"abc1234", "store_url":"http://example.com/store","owner_id":"admin1@example.com"}`},
	{json: `{"id":"abc1234", "store_url":"http://example.com/store","owner_id":"admin1@example.com","created_on":"2014-07-01 13:13:13.000000013 -0400 EDT"}`},
}

func TestReadOnlyStudyFields(t *testing.T) {
	var err error

	for i, tc := range test_cases {
		err = json.Unmarshal([]byte(tc.json), actual)
		if err != nil {
			t.Fatalf("test case %v: json unmarshal error: %v", i, err)
		}

		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("test case %v: Study objects are not equal.\n actual: %v\n expected: %v\n", i, actual, expected)
		}
	}
}
