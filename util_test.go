package gosql

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/ilibs/gosql/v2/internal/example/models"
)

type MyString string

type MyStruct struct {
	num  int
	text MyString
}

var (
	zeroPtr    *string
	zeroSlice  []int
	zeroFunc   func() string
	zeroMap    map[string]string
	emptyIface interface{}
	zeroIface  fmt.Stringer
	zeroValues = []interface{}{
		nil,

		// bool
		false,

		// int
		0,
		int8(0),
		int16(0),
		int32(0),
		int64(0),
		uint(0),
		uint8(0),
		uint16(0),
		uint32(0),
		uint64(0),

		// float
		0.0,
		float32(0.0),
		float64(0.0),

		// string
		"",

		// alias
		MyString(""),

		// func
		zeroFunc,

		// array / slice
		[0]int{},
		zeroSlice,

		// map
		zeroMap,

		// interface
		emptyIface,
		zeroIface,

		// pointer
		zeroPtr,

		// struct
		MyStruct{},
		time.Time{},
		MyStruct{num: 0},
		MyStruct{text: MyString("")},
		sql.NullString{String: "", Valid: false},
		sql.NullInt64{Int64: 0, Valid: false},
	}
	nonZeroIface  fmt.Stringer = time.Now()
	nonZeroValues              = []interface{}{
		// bool
		true,

		// int
		1,
		int8(1),
		int16(1),
		int32(1),
		int64(1),
		uint8(1),
		uint16(1),
		uint32(1),
		uint64(1),

		// float
		1.0,
		float32(1.0),
		float64(1.0),

		// string
		"test",

		// alias
		MyString("test"),

		// func
		time.Now,

		// array / slice
		[]int{},
		[]int{42},
		[1]int{42},

		// map
		make(map[string]string, 1),

		// interface
		nonZeroIface,

		// pointer
		&nonZeroIface,

		// struct
		MyStruct{num: 1},
		time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
		sql.NullString{String: "", Valid: true},
		sql.NullInt64{Int64: 0, Valid: true},
	}
)

func TestIsZero(t *testing.T) {
	for _, value := range zeroValues {
		if !IsZero(reflect.ValueOf(value)) {
			t.Errorf("expected '%v' (%T) to be recognized as zero value", value, value)
		}
	}

	for _, value := range nonZeroValues {
		if IsZero(reflect.ValueOf(value)) {
			t.Errorf("did not expect '%v' (%T) to be recognized as zero value", value, value)
		}
	}
}

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
	user := &models.Users{
		Id:   1,
		Name: "test",
	}

	rv := reflect.Indirect(reflect.ValueOf(user))
	fields := mapper.FieldMap(rv)

	m := zeroValueFilter(fields, nil)

	if _, ok := m["status"]; ok {
		t.Error("status value not filter")
	}

	if _, ok := m["creatd_at"]; ok {
		t.Error("creatd_at zero value not filter")
	}

	if _, ok := m["updated_at"]; ok {
		t.Error("updated_at zero value not filter")
	}

	m2 := zeroValueFilter(fields, []string{"status"})
	if _, ok := m2["status"]; !ok {
		t.Error("status shouldn't be filtered")
	}
}

func TestUtil_structAutoTime(t *testing.T) {
	user := &models.Users{
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
		"created_at": "2018-07-11 11:58:21",
		"updated_at": "2018-07-11 11:58:21",
	}

	keySort := []string{"created_at", "id", "name", "updated_at"}

	s := sortedParamKeys(m)

	for i, k := range s {
		if k != keySort[i] {
			t.Error("sort error", k)
		}
	}
}
