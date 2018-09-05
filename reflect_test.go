package gosql

import (
	"database/sql"
	"reflect"
	"testing"
	"time"
)

func TestReflectMapper_FieldMap(t *testing.T) {
	mapper := NewReflectMapper("db")

	{
		user := &Users{
			Id:    1,
			Name:  "test",
			Email: "test@test.com",
			SuccessTime: sql.NullString{
				String: "2018-09-03 00:00:00",
				Valid:  false,
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}


		fields := mapper.FieldMap(reflect.ValueOf(user))
		if len(fields) != 7 {
			t.Error("FieldMap length error")
		}

		if v := fields["name"].Interface().(string); v != user.Name {
			t.Errorf("Expecting %s, got %s", user.Name, v)
		}

		if v := fields["success_time"].Interface().(sql.NullString).String; v != user.SuccessTime.String {
			t.Errorf("Expecting %s, got %s", user.Name, v)
		}
	}

	{
		user := &UserCombs{}
		fields := mapper.FieldMap(reflect.ValueOf(user))

		if len(fields) != 7 {
			t.Error("FieldMap length error")
		}

		if v := fields["name"].Interface().(string); v != user.Name {
			t.Errorf("Expecting %s, got %s", user.Name, v)
		}

		if v := fields["success_time"].Interface().(sql.NullString).String; v != user.SuccessTime.String {
			t.Errorf("Expecting %s, got %s", user.Name, v)
		}

		//fmt.Println(fields)
	}
}
