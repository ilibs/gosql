package gosql

import (
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var Default = "default"
var dbService = make(map[string]*sqlx.DB, 0)

// DB gets the specified database engine,
// or the default DB if no name is specified.
func DB(name ...string) *sqlx.DB {
	dbName := Default
	if name != nil {
		dbName = name[0]
	}

	engine, ok := dbService[dbName]
	if !ok {
		log.Fatalf("[db] the database link `%s` is not configured", dbName)
	}
	return engine
}

// List gets the list of database engines
func List() map[string]*sqlx.DB {
	return dbService
}

func Connect(configs map[string]*Config) {

	var errs []string
	defer func() {
		if len(errs) > 0 {
			panic("[db] " + strings.Join(errs, "\n"))
		}
	}()

	for key, conf := range configs {
		if !conf.Enable {
			continue
		}

		sess, err := sqlx.Connect(conf.Driver, conf.Dsn)

		if err != nil {
			errs = append(errs, err.Error())
			continue
		}
		log.Println("[db] connect:" + key)

		if conf.ShowSql {
			logger.SetLogging(true)
		}

		sess.SetMaxOpenConns(conf.MaxOpenConns)
		sess.SetMaxIdleConns(conf.MaxIdleConns)

		dbService[key] = sess
	}
}
