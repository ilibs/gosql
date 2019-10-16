package gosql

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/ilibs/gosql/v2/internal/example/models"
)

func TestShowSql(t *testing.T) {
	logger.logging = false
	RunWithSchema(t, func(t *testing.T) {
		insert(1)

		ShowSql().Queryx("select * from users")
		user := &models.Users{}
		Model(user).ShowSQL().Where("id = ?", 1).Get()
		Table("users").ShowSQL().Where("id = ?", 1).Update(map[string]interface{}{
			"name": "test2",
		})
	})
}

func TestExec(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		result, err := Exec("insert into users(name,status,created_at,updated_at) value(?,?,?,?)", "test", 1, time.Now(), time.Now())

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

		Exec("update users set status = status + 1 where id = ?", 1)

		result, err = Exec("update users set status = status + 1 where id = ?", 1)
		if err != nil {
			t.Error("update user error", err)
		}

		if aff, _ := result.RowsAffected(); aff == 0 {
			t.Error("update set error")
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
		defer rows.Close()

		for rows.Next() {
			user := &models.Users{}
			err = rows.StructScan(user)
			if err != nil {
				t.Error(err)
			}
		}

		rows, err = Queryx("select name from users")

		if err != nil {
			t.Error(err)
		}
		defer rows.Close()

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
		user := &models.Users{}
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
		db := Use("db2")
		_, err := db.Exec("insert into photos(moment_id,url,created_at,updated_at) value(?,?,?,?)", 1, "http://test.com", time.Now(), time.Now())

		if err != nil {
			t.Error(err)
		}
	})
}

func TestUseTable(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		post := &models.Photos{
			MomentId: 1,
			Url:      "http://test.com",
		}

		_, err := Use("db2").Model(post).Create()

		if err != nil {
			t.Error(err)
		}

		_, err = Use("db2").Table("photos").Where("id = ?", 1).Update(map[string]interface{}{
			"url": "http://test2.com",
		})

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
			user := &models.Users{}
			err := db.Get(user, "select * from users where id = ?", 1)

			if err != nil {
				t.Error(err)
			}

			fmt.Println(jsonEncode(user))
		}
	})
}

func TestGetSingle(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		insert(1)
		db := Use("default")
		{
			var name string
			err := db.Get(&name, "select name from users where id = ?", 1)

			if err != nil {
				t.Error(err)
			}

			fmt.Println(name)
		}
	})
}

func TestSelectSlice(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		insert(1)
		insert(2)
		var users []string
		err := Select(&users, "select name from users")

		if err != nil {
			t.Error(err)
		}

		fmt.Println(jsonEncode(users))
	})
}

func TestSelect(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		insert(1)
		insert(2)
		db := Use("default")
		user := make([]*models.Users, 0)
		err := db.Select(&user, "select * from users")

		if err != nil {
			t.Error(err)
		}

		fmt.Println(jsonEncode(user))
	})
}

func TestQueryxIn(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		insert(1)
		insert(2)
		insert(4)
		insert(5)
		insert(6)

		rows, err := Queryx("select * from users where status = ? and  id in (?)", 0, []int{1, 2, 3})

		if err != nil {
			t.Error(err)
		}
		defer rows.Close()

		for rows.Next() {
			user := &models.Users{}
			err = rows.StructScan(user)
			if err != nil {
				t.Error(err)
			}
		}
	})
}

func TestSelectIn(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		insert(1)
		insert(2)
		insert(4)
		insert(5)
		insert(6)
		db := Use("default")
		user := make([]*models.Users, 0)
		err := db.Select(&user, "select * from users where id in(?)", []int{1, 2, 3})

		if err != nil {
			t.Error(err)
		}

		fmt.Println(jsonEncode(user))
	})
}

func TestTx(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		//1
		{
			Tx(func(tx *DB) error {
				for id := 1; id < 10; id++ {
					user := &models.Users{
						Id:   id,
						Name: "test" + strconv.Itoa(id),
					}

					tx.Model(user).Create()

					if id == 8 {
						return errors.New("simulation terminated")
					}
				}

				return nil
			})

			num, err := Model(&models.Users{}).Count()

			if err != nil {
				t.Error(err)
			}

			if num != 0 {
				t.Error("transaction abort failed")
			}
		}

		//2
		{
			Tx(func(tx *DB) error {
				for id := 1; id < 10; id++ {
					user := &models.Users{
						Id:   id,
						Name: "test" + strconv.Itoa(id),
					}

					tx.Model(user).Create()
				}

				return nil
			})

			num, err := Model(&models.Users{}).Count()

			if err != nil {
				t.Error(err)
			}

			if num != 9 {
				t.Error("transaction create failed")
			}
		}
	})
}

func TestWithTx(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		{
			Tx(func(tx *DB) error {
				for id := 1; id < 10; id++ {
					_, err := tx.Exec("INSERT INTO users(id,name,created_at,updated_at) VALUES(?,?,?,?,?)", id, "test"+strconv.Itoa(id), time.Now(), time.Now())
					if err != nil {
						return err
					}
				}

				var num int
				err := tx.QueryRowx("select count(*) from users").Scan(&num)

				if err != nil {
					return err
				}

				if num != 9 {
					t.Error("with transaction create failed")
				}

				return nil
			})

		}
	})
}

func TestTxx(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {

		ctx, cancel := context.WithCancel(context.Background())

		Txx(ctx, func(ctx context.Context, tx *DB) error {
			for id := 1; id < 10; id++ {
				user := &models.Users{
					Id:   id,
					Name: "test" + strconv.Itoa(id),
				}

				tx.Model(user).Create()

				if id == 8 {
					cancel()
					break
				}
			}

			return nil
		})

		num, err := Model(&models.Users{}).Count()

		if err != nil {
			t.Error(err)
		}

		if num != 0 {
			t.Error("transaction abort failed")
		}
	})
}

func TestWrapper_Relation(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		initDatas(t)
		moment := &MomentList{}
		err := Relation("User", func(b *ModelStruct) {
			b.Where("status = 0")
		}).Get(moment, "select * from moments")

		b, _ := json.MarshalIndent(moment, "", "	")
		fmt.Println(string(b), err)

		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestWrapper_Relation2(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		initDatas(t)
		var moments = make([]*MomentList, 0)
		err := Relation("User", func(b *ModelStruct) {
			b.Where("status = 1")
		}).Select(&moments, "select * from moments")

		if err != nil {
			t.Fatal(err)
		}
	})
}
