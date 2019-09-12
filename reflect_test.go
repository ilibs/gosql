package gosql

import (
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/ilibs/gosql/v2/internal/example/models"
)

func TestReflectMapper_FieldMap(t *testing.T) {
	mapper := NewReflectMapper("db")

	{
		user := &models.Users{
			Id:   1,
			Name: "test",
			SuccessTime: sql.NullString{
				String: "2018-09-03 00:00:00",
				Valid:  false,
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		fields := mapper.FieldMap(reflect.ValueOf(user))
		if len(fields) != 6 {
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
		photos := &models.Photos{}
		fields := mapper.FieldMap(reflect.ValueOf(photos))

		if len(fields) != 5 {
			t.Error("FieldMap length error")
		}

		if v := fields["url"].Interface().(string); v != photos.Url {
			t.Errorf("Expecting %s, got %s", photos.Url, v)
		}

		if _, ok := fields["created_at"].Interface().(time.Time); !ok {
			t.Error("Expecting true, got false")
		}

		//fmt.Println(fields)
	}
}
