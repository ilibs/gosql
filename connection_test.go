package gosql

import (
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestMain(m *testing.M) {
	configs := make(map[string]*Config)

	dsn1 := os.Getenv("MYSQL_TEST_DSN1")

	if dsn1 == "" {
		dsn1 = "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True&loc=Asia%2FShanghai"
	}

	dsn2 := os.Getenv("MYSQL_TEST_DSN2")

	if dsn2 == "" {
		dsn2 = "root:123456@tcp(127.0.0.1:3306)/test2?charset=utf8&parseTime=True&loc=Asia%2FShanghai"
	}

	configs["default"] = &Config{
		Enable:  true,
		Driver:  "mysql",
		Dsn:     dsn1,
		ShowSql: true,
	}

	configs["db2"] = &Config{
		Enable:  true,
		Driver:  "mysql",
		Dsn:     dsn2,
		ShowSql: true,
	}

	_ = Connect(configs)

	m.Run()
}

func TestConnect(t *testing.T) {
	db := Sqlx()

	if db.DriverName() != "mysql" {
		t.Fatalf("sqlx database connection error")
	}
}
