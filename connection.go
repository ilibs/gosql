package gosql

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

// defaultLink set database default link name
var defaultLink = "default"

// If database fatal exit
var FatalExit = true
var dbService = make(map[string]*sqlx.DB, 0)

// DB gets the specified database engine,
// or the default DB if no name is specified.
func Sqlx(name ...string) *sqlx.DB {
	dbName := defaultLink
	if name != nil {
		dbName = name[0]
	}

	engine, ok := dbService[dbName]
	if !ok {
		panic(fmt.Sprintf("[db] the database link `%s` is not configured", dbName))
	}
	return engine
}

// List gets the list of database engines
func List() map[string]*sqlx.DB {
	return dbService
}

type Options struct {
	maxOpenConns int
	maxIdleConns int
	maxLifetime  int
}

type Option func(*Options)

func WithMaxOpenConns(i int) Option {
	return func(options *Options) {
		options.maxOpenConns = i
	}
}

func WithMaxIdleConns(i int) Option {
	return func(options *Options) {
		options.maxIdleConns = i
	}
}

func WithMaxLifetimes(i int) Option {
	return func(options *Options) {
		options.maxLifetime = i
	}
}

// Open gosql.DB with sqlx
func Open(driver, dbSource string, opts ...Option) (*DB, error) {

	var options Options
	for _, opt := range opts {
		opt(&options)
	}

	db, err := sqlx.Connect(driver, dbSource)
	if err != nil {
		return nil, err
	}

	if options.maxOpenConns > 0 {
		db.SetMaxOpenConns(options.maxOpenConns)
	}

	if options.maxIdleConns > 0 {
		db.SetMaxIdleConns(options.maxIdleConns)
	}

	if options.maxLifetime > 0 {
		db.SetConnMaxLifetime(time.Duration(options.maxLifetime) * time.Second)
	}

	return &DB{database: db}, nil
}

// OpenWithDB open gosql.DB with sql.DB
func OpenWithDB(driver string, db *sql.DB) *DB {
	return &DB{database: sqlx.NewDb(db, driver)}
}

// Connect database
func Connect(configs map[string]*Config) (err error) {

	var errs []string
	defer func() {
		if len(errs) > 0 {
			err = errors.New("[db] " + strings.Join(errs, "\n"))
			if FatalExit {
				log.Fatal(err)
			}
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
		if conf.MaxLifetime > 0 {
			sess.SetConnMaxLifetime(time.Duration(conf.MaxLifetime) * time.Second)
		}

		if db, ok := dbService[key]; ok {
			_ = db.Close()
		}

		dbService[key] = sess
	}
	return
}
