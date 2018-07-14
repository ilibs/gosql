package gosql

import (
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestMain(m *testing.M) {
	configs := make(map[string]*Config)

	dsn := os.Getenv("MYSQL_TEST_DSN")

	if dsn == "" {
		dsn = "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True&loc=Asia%2FShanghai"
	}

	configs["default"] = &Config{
		Enable:  true,
		Driver:  "mysql",
		Dsn:     dsn,
		ShowSql: true,
	}

	Connect(configs)
	m.Run()
}

func TestConnect(t *testing.T) {
	db := DB()

	if db.DriverName() != "mysql" {
		t.Fatalf("sqlx database connection error")
	}
}
