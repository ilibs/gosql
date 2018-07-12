package gosql

import (
	"strconv"
	"testing"
)

func mapInsert(t *testing.T, id int64) int64 {
	id, err := Table("users").Create(map[string]interface{}{
		"id":         id,
		"name":       "test" + strconv.Itoa(int(id)),
		"email":      "test@test.com",
		"created_at": "2018-07-11 11:58:21",
		"updated_at": "2018-07-11 11:58:21",
	})

	if err != nil {
		t.Error(err)
	}

	if id <= 0 {
		t.Error("map insert error")
	}

	return id
}

func TestMapper_Create(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		mapInsert(t, 1)
	})
}

func TestMapper_Update(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		id := mapInsert(t, 1)

		affected, err := Table("users").Where("id = ?", id).Update(map[string]interface{}{
			"name":  "fifsky",
			"email": "fifsky@test.com",
		})

		if err != nil {
			t.Error(err)
		}

		if affected <= 0 {
			t.Error("map update error")
		}
	})
}

func TestMapper_Delete(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		{
			id := mapInsert(t, 1)
			affected, err := Table("users").Where("id = ?", id).Delete()

			if err != nil {
				t.Error(err)
			}

			if affected <= 0 {
				t.Error("map delete error")
			}
		}

		{
			mapInsert(t, 2)
			affected, err := Table("users").Delete()
			if err != nil {
				t.Error(err)
			}

			if affected <= 0 {
				t.Error("map delete error")
			}
		}
	})
}

func TestMapper_Count(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		{
			id := mapInsert(t, 1)
			num, err := Table("users").Where("id = ?", id).Count()

			if err != nil {
				t.Error(err)
			}

			if num != 1 {
				t.Error("map count error")
			}
		}

		{
			mapInsert(t, 2)
			mapInsert(t, 3)
			num, err := Table("users").Count()
			if err != nil {
				t.Error(err)
			}

			if num <= 0 {
				t.Error("map count error")
			}
		}
	})
}
