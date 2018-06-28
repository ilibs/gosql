package gosql

import (
	"testing"
	"time"
)

func TestExec(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		result, err := Exec("insert into users(name,email,created_at,updated_at) value(?,?,?,?)", "test", "test@gmail.com", time.Now(), time.Now())

		if err != nil {
			t.Error(err)
		}

		id, err := result.LastInsertId()

		if err != nil {
			t.Error(err)
		}

		if id != 1 {
			t.Error("lastInsertId error")
		}
	})
}

func TestQueryx(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		insert(1)
		insert(2)

		rows, err := Queryx("select * from users")

		if err != nil {
			t.Error(err)
		}

		for rows.Next() {
			user := &Users{}
			err = rows.StructScan(user)
			if err != nil {
				t.Error(err)
			}
		}
	})
}

func TestQueryRowx(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		insert(1)
		user := &Users{}
		err := QueryRowx("select * from users where id = 1").StructScan(user)

		if err != nil {
			t.Error(err)
		}

		if user.Id != 1 {
			t.Error("wraper QueryRowx error")
		}
	})
}

func TestUse(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		db := Use("default")
		_, err := db.Exec("insert into users(name,email,created_at,updated_at) value(?,?,?,?)", "test", "test@gmail.com", time.Now(), time.Now())

		if err != nil {
			t.Error(err)
		}
	})
}