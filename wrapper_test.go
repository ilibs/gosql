package gosql

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
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

		rows, err = Queryx("select name from users")

		if err != nil {
			t.Error(err)
		}

		for rows.Next() {
			//results := make(map[string]interface{})
			//err = rows.MapScan(results)
			var name string
			err = rows.Scan(&name)
			if err != nil {
				t.Error(err)
			}
			fmt.Println(name)
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

func TestGet(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		insert(1)
		db := Use("default")
		{
			user := &Users{}
			err := db.Get(user, "select * from users where id = ?", 1)

			if err != nil {
				t.Error(err)
			}

			fmt.Println(jsonEncode(user))
		}
	})
}

func TestSelect(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		insert(1)
		insert(2)
		db := Use("default")
		user := make([]*Users, 0)
		err := db.Select(&user, "select * from users")

		if err != nil {
			t.Error(err)
		}

		fmt.Println(jsonEncode(user))
	})
}

func TestTx(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		Tx(func(tx *sqlx.Tx) error {
			for id := 1; id < 10; id++ {
				user := &Users{
					Id:    id,
					Name:  "test" + strconv.Itoa(id),
					Email: "test" + strconv.Itoa(id) + "@test.com",
				}

				Model(user, tx).Create()

				if id == 8 {
					return errors.New("simulation terminated")
				}
			}

			return nil
		})

		num, err := Model(&Users{}).Count()

		if err != nil {
			t.Error(err)
		}

		if num != 0 {
			t.Error("transaction abort failed")
		}

		Tx(func(tx *sqlx.Tx) error {
			for id := 1; id < 10; id++ {
				user := &Users{
					Id:    id,
					Name:  "test" + strconv.Itoa(id),
					Email: "test" + strconv.Itoa(id) + "@test.com",
				}

				Model(user, tx).Create()
			}

			return nil
		})

		num, err = Model(&Users{}).Count()

		if err != nil {
			t.Error(err)
		}

		if num != 9 {
			t.Error("transaction create failed")
		}
	})
}

func TestTxx(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {

		ctx, cancel := context.WithCancel(context.Background())

		Txx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
			for id := 1; id < 10; id++ {
				user := &Users{
					Id:    id,
					Name:  "test" + strconv.Itoa(id),
					Email: "test" + strconv.Itoa(id) + "@test.com",
				}

				Model(user, tx).Create()

				if id == 8 {
					cancel()
					break
				}
			}

			return nil
		})

		num, err := Model(&Users{}).Count()

		if err != nil {
			t.Error(err)
		}

		if num != 0 {
			t.Error("transaction abort failed")
		}
	})
}
