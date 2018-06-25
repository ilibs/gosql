package gosql

import (
	"os"
	"testing"
	"time"
)

var (
	createSchema = `
CREATE TABLE test_users (
	id int(11) unsigned NOT NULL AUTO_INCREMENT,
	name  varchar(100) NOT NULL DEFAULT '',
	email  varchar(100) NOT NULL DEFAULT '',
	created_at datetime NOT NULL,
	updated_at datetime NOT NULL,
  	PRIMARY KEY (id)
)ENGINE=InnoDB CHARSET=utf8;
`

	dropSchema = `
	drop table test_users
`
)

type Users struct {
	Id        int       `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func TestConnect(t *testing.T) {
	configs := make(map[string]*Config)

	dsn := os.Getenv("MYSQL_TEST_DSN")

	if dsn == "" {
		dsn = "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True&loc=Asia%2FShanghai"
	}

	configs["default"] = &Config{
		Enable: true,
		Driver: "mysql",
		Dsn:    dsn,
	}

	Connect(configs)
	db := DB()

	if db.DriverName() != "mysql" {
		t.Fatalf("sqlx database connection error")
	}

	_, err := db.Exec(createSchema)

	if err != nil {
		t.Fatalf("create schema error:%s", err)
	}

	defer func() {
		db.Exec(dropSchema)
	}()

	result, err := db.Exec("insert into test_users(name,email,created_at,updated_at) values(?,?,?,?)", "test", "test@mail.com", time.Now(), time.Now())

	if err != nil {
		t.Fatalf("insert error:%s", err)
	}

	id, err := result.LastInsertId()

	if err != nil {
		t.Fatalf("LastInsertId error:%s", err)
	}

	_, err = db.Exec("update test_users set name = ? where id = ?", "test2", id)

	if err != nil {
		t.Fatalf("update error:%s", err)
	}

	user := Users{}
	DB().Get(&user, "select * from test_users where id = ?", id)

	if user.Name != "test2" {
		t.Fatalf("select error:%#v", user)
	}

	_, err = db.Exec("delete from test_users where id = ?", id)

	if err != nil {
		t.Fatalf("delete error:%s", err)
	}
}
