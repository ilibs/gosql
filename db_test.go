package gosql

import (
	"context"
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

func TestBatchExec(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		{
			// batch insert with structs
			users := []models.Users{
				{Name: "test1", Status: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				{Name: "test2", Status: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				{Name: "test3", Status: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				{Name: "test4", Status: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
			}

			result, err := NamedExec("insert into users(name,status,created_at,updated_at) values(:name,:status,:created_at,:updated_at)", users)

			if err != nil {
				t.Error(err)
			}

			if aff, _ := result.RowsAffected(); aff != 4 {
				t.Error("update set error")
			}
		}

		{
			// batch insert with maps
			users := []map[string]interface{}{
				{"name": "test5", "status": 1, "created_at": "2021-01-25 12:22:22", "updated_at": "2021-01-25 12:22:22"},
				{"name": "test6", "status": 1, "created_at": "2021-01-25 12:22:22", "updated_at": "2021-01-25 12:22:22"},
				{"name": "test7", "status": 1, "created_at": "2021-01-25 12:22:22", "updated_at": "2021-01-25 12:22:22"},
				{"name": "test8", "status": 1, "created_at": "2021-01-25 12:22:22", "updated_at": "2021-01-25 12:22:22"},
			}

			result, err := NamedExec("insert into users(name,status,created_at,updated_at) values(:name,:status,:created_at,:updated_at)", users)

			if err != nil {
				t.Error(err)
			}

			if aff, _ := result.RowsAffected(); aff != 4 {
				t.Error("update set error")
			}
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
			// results := make(map[string]interface{})
			// err = rows.MapScan(results)
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

		post.Url = "http://test.com/2"
		_, err = Use("db2").WithContext(context.Background()).Model(post).Update()
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
		// 1
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

		// 2
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
			err := Tx(func(tx *DB) error {
				for id := 1; id < 10; id++ {
					_, err := tx.Exec("INSERT INTO users(id,name,status,created_at,updated_at) VALUES(?,?,?,?,?)", id, "test"+strconv.Itoa(id), 1, time.Now(), time.Now())
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

			if err != nil {
				t.Fatalf("with transaction failed %s", err)
			}
		}
	})
}

func TestTxx(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {

		ctx, cancel := context.WithCancel(context.Background())

		err := Txx(ctx, func(ctx context.Context, tx *DB) error {
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

		if err == nil {
			t.Fatalf("with transaction must be cancel error")
		}

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
		err := Relation("User", func(b *Builder) {
			b.Where("status = 0")
		}).Get(moment, "select * from moments")

		// b, _ := json.MarshalIndent(moment, "", "	")
		// fmt.Println(string(b), err)

		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestWrapper_Relation2(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		initDatas(t)
		var moments = make([]*MomentList, 0)
		err := Relation("User", func(b *Builder) {
			b.Where("status = 1")
		}).Select(&moments, "select * from moments")

		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestDB_Begin(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		tx, err := Begin()
		if err != nil {
			t.Fatalf("with transaction begin error %s", err)
		}

		var fn = func() error {
			for id := 1; id < 10; id++ {
				_, err := tx.Exec("INSERT INTO users(id,name,status,created_at,updated_at) VALUES(?,?,?,?,?)", id, "test"+strconv.Itoa(id), 1, time.Now(), time.Now())
				if err != nil {
					return err
				}
			}

			var num int
			err = tx.QueryRowx("select count(*) from users").Scan(&num)

			if err != nil {
				return err
			}

			if num != 9 {
				return errors.New("with transaction create failed")
			}
			return nil
		}

		err = fn()

		if err != nil {
			err := tx.Rollback()
			if err != nil {
				t.Fatalf("with transaction rollback error %s", err)
			}
		}

		err = tx.Commit()
		if err != nil {
			t.Fatalf("with transaction commit error %s", err)
		}
	})
}

func TestRelation(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		initDatas(t)

		moment := &MomentList{}
		err := Relation("User", func(b *Builder) {
			// this is builder instance
			b.Where("status = 1")
		}).Get(moment, "select * from moments where id = 1")

		if err != nil {
			t.Fatalf("relation query error %s", err)
		}
	})
}
