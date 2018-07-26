package gosql

import (
	"reflect"
	"testing"
)

func TestUtil_inSlice(t *testing.T) {
	s := []string{"a", "b", "c"}

	if !inSlice("a", s) {
		t.Error("in slice find error")
	}

	if inSlice("d", s) {
		t.Error("in slice exist error")
	}
}

func TestUtil_zeroValueFilter(t *testing.T) {
	user := &Users{
		Id:   1,
		Name: "test",
	}

	rv := reflect.Indirect(reflect.ValueOf(user))
	fields := mapper.FieldMap(rv)

	m := zeroValueFilter(fields, nil)

	if _, ok := m["email"]; ok {
		t.Error("email zero value not filter")
	}

	if _, ok := m["creatd_at"]; ok {
		t.Error("creatd_at zero value not filter")
	}

	if _, ok := m["updated_at"]; ok {
		t.Error("updated_at zero value not filter")
	}

	m2 := zeroValueFilter(fields, []string{"email"})
	if _, ok := m2["email"]; !ok {
		t.Error("email shouldn't be filtered")
	}
}

func TestUtil_structAutoTime(t *testing.T) {
	user := &Users{
		Id:   1,
		Name: "test",
	}
	rv := reflect.Indirect(reflect.ValueOf(user))
	fields := mapper.FieldMap(rv)

	structAutoTime(fields, AUTO_CREATE_TIME_FIELDS)

	if user.CreatedAt.IsZero() {
		t.Error("auto time fail")
	}
}

func TestUtil_sortedParamKeys(t *testing.T) {
	m := map[string]interface{}{
		"id":         1,
		"name":       "test",
		"email":      "test@test.com",
		"created_at": "2018-07-11 11:58:21",
		"updated_at": "2018-07-11 11:58:21",
	}

	keySort := []string{"created_at", "email", "id", "name", "updated_at"}

	s := sortedParamKeys(m)

	for i, k := range s {
		if k != keySort[i] {
			t.Error("sort error", k)

		}
	}
}
