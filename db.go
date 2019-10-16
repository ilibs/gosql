package gosql

import (
	"bufio"
	"bytes"
	"context"
	"database/sql"
	"log"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type ISqlx interface {
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
	QueryRowx(query string, args ...interface{}) *sqlx.Row
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
	Exec(query string, args ...interface{}) (sql.Result, error)
	Rebind(query string) string
	DriverName() string
}

type BuilderChainFunc func(b *ModelStruct)

type DB struct {
	database    *sqlx.DB
	tx          *sqlx.Tx
	logging     bool
	RelationMap map[string]BuilderChainFunc
}

// return database instance, if it is a transaction, the transaction priority is higher
func (w *DB) db() ISqlx {
	if w.tx != nil {
		return w.tx.Unsafe()
	}

	return w.database.Unsafe()
}

// ShowSql single show sql log
func ShowSql() *DB {
	w := Use(defaultLink)
	w.logging = true
	return w
}

func (w *DB) argsIn(query string, args []interface{}) (string, []interface{}, error) {
	newArgs := make([]interface{}, 0)
	newQuery, newArgs, err := sqlx.In(query, args...)

	if err != nil {
		return query, args, err
	}

	return newQuery, newArgs, nil
}

//DriverName wrapper sqlx.DriverName
func (w *DB) DriverName() string {
	if w.tx != nil {
		return w.tx.DriverName()
	}

	return w.database.DriverName()
}

//Rebind wrapper sqlx.Rebind
func (w *DB) Rebind(query string) string {
	return w.db().Rebind(query)
}

//Exec wrapper sqlx.Exec
func (w *DB) Exec(query string, args ...interface{}) (result sql.Result, err error) {
	defer func(start time.Time) {
		logger.Log(&QueryStatus{
			Query: query,
			Args:  args,
			Err:   err,
			Start: start,
			End:   time.Now(),
		}, w.logging)

	}(time.Now())

	return w.db().Exec(query, args...)
}

//Queryx wrapper sqlx.Queryx
func (w *DB) Queryx(query string, args ...interface{}) (rows *sqlx.Rows, err error) {
	defer func(start time.Time) {
		logger.Log(&QueryStatus{
			Query: query,
			Args:  args,
			Err:   err,
			Start: start,
			End:   time.Now(),
		}, w.logging)
	}(time.Now())

	query, newArgs, err := w.argsIn(query, args)
	if err != nil {
		return nil, err
	}

	return w.db().Queryx(query, newArgs...)
}

//QueryRowx wrapper sqlx.QueryRowx
func (w *DB) QueryRowx(query string, args ...interface{}) (rows *sqlx.Row) {
	defer func(start time.Time) {
		logger.Log(&QueryStatus{
			Query: query,
			Args:  args,
			Err:   rows.Err(),
			Start: start,
			End:   time.Now(),
		}, w.logging)
	}(time.Now())

	query, newArgs, _ := w.argsIn(query, args)

	return w.db().QueryRowx(query, newArgs...)
}

//Get wrapper sqlx.Get
func (w *DB) Get(dest interface{}, query string, args ...interface{}) (err error) {
	defer func(start time.Time) {
		logger.Log(&QueryStatus{
			Query: query,
			Args:  args,
			Err:   err,
			Start: start,
			End:   time.Now(),
		}, w.logging)
	}(time.Now())

	hook := NewHook(w)
	refVal := reflect.ValueOf(dest)
	hook.callMethod("BeforeFind", refVal)

	query, newArgs, err := w.argsIn(query, args)
	if err != nil {
		return err
	}

	err = w.db().Get(dest, query, newArgs...)
	if err != nil {
		return err
	}

	if reflect.Indirect(refVal).Kind() == reflect.Struct {
		// relation data fill
		err = RelationOne(dest, w.RelationMap)
	}

	if err != nil {
		return err
	}

	hook.callMethod("AfterFind", refVal)
	if hook.HasError() > 0 {
		return hook.Error()
	}

	return nil
}

func indirectType(v reflect.Type) reflect.Type {
	if v.Kind() != reflect.Ptr {
		return v
	}
	return v.Elem()
}

//Select wrapper sqlx.Select
func (w *DB) Select(dest interface{}, query string, args ...interface{}) (err error) {
	defer func(start time.Time) {
		logger.Log(&QueryStatus{
			Query: query,
			Args:  args,
			Err:   err,
			Start: start,
			End:   time.Now(),
		}, w.logging)
	}(time.Now())

	query, newArgs, err := w.argsIn(query, args)
	if err != nil {
		return err
	}

	err = w.db().Select(dest, query, newArgs...)
	if err != nil {
		return err
	}

	t := indirectType(reflect.TypeOf(dest))
	if t.Kind() == reflect.Slice {
		if indirectType(t.Elem()).Kind() == reflect.Struct {
			// relation data fill
			err = RelationAll(dest, w.RelationMap)
		}
	}

	if err != nil {
		return err
	}

	return nil
}

//Txx the transaction with context
func (w *DB) Txx(ctx context.Context, fn func(ctx context.Context, tx *DB) error) (err error) {
	tx, err := w.database.BeginTxx(ctx, nil)

	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				log.Printf("gosql rollback error:%s", err)
			}
		}
	}()

	err = fn(ctx, &DB{tx: tx})
	if err == nil {
		err = tx.Commit()
	}
	return
}

//Tx the transaction
func (w *DB) Tx(fn func(w *DB) error) (err error) {
	tx, err := w.database.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				log.Printf("gosql rollback error:%s", err)
			}
		}
	}()
	err = fn(&DB{tx: tx})
	if err == nil {
		err = tx.Commit()
	}
	return
}

// Table database handler from to table name
// for example:
// gosql.Use("db2").Table("users")
func (w *DB) Table(t string) *Mapper {
	return &Mapper{db: w, SQLBuilder: SQLBuilder{table: t, dialect: newDialect(w.DriverName())}}
}

//Model database handler from to struct
//for example:
// gosql.Use("db2").Model(&users{})
func (w *DB) Model(m interface{}) *ModelStruct {
	return &ModelStruct{model: m, db: w, SQLBuilder: SQLBuilder{dialect: newDialect(w.DriverName())}}
}

//Import SQL DDL from sql file
func (w *DB) Import(f string) ([]sql.Result, error) {
	file, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var results []sql.Result
	scanner := bufio.NewScanner(file)

	semiColSpliter := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		if i := bytes.IndexByte(data, ';'); i >= 0 {
			return i + 1, data[0:i], nil
		}
		// If we're at EOF, we have a final, non-terminated line. Return it.
		if atEOF {
			return len(data), data, nil
		}
		// Request more data.
		return 0, nil, nil
	}

	scanner.Split(semiColSpliter)

	for scanner.Scan() {
		query := strings.Trim(scanner.Text(), " \t\n\r")
		if len(query) > 0 {
			result, err := w.db().Exec(query)
			results = append(results, result)
			if err != nil {
				return nil, err
			}
		}
	}

	return results, nil
}

// Relation association table builder handle
func (w *DB) Relation(name string, fn BuilderChainFunc) *DB {
	if w.RelationMap == nil {
		w.RelationMap = make(map[string]BuilderChainFunc)
	}
	w.RelationMap[name] = fn
	return w
}

//Use is change database
func Use(db string) *DB {
	return &DB{database: Sqlx(db)}
}

//Exec default database
func Exec(query string, args ...interface{}) (sql.Result, error) {
	return Use(defaultLink).Exec(query, args...)
}

//Queryx default database
func Queryx(query string, args ...interface{}) (*sqlx.Rows, error) {
	return Use(defaultLink).Queryx(query, args...)
}

//QueryRowx default database
func QueryRowx(query string, args ...interface{}) *sqlx.Row {
	return Use(defaultLink).QueryRowx(query, args...)
}

//Txx default database the transaction with context
func Txx(ctx context.Context, fn func(ctx context.Context, tx *DB) error) error {
	return Use(defaultLink).Txx(ctx, fn)
}

//Tx default database the transaction
func Tx(fn func(tx *DB) error) error {
	return Use(defaultLink).Tx(fn)
}

//Get default database
func Get(dest interface{}, query string, args ...interface{}) error {
	return Use(defaultLink).Get(dest, query, args...)
}

//Select default database
func Select(dest interface{}, query string, args ...interface{}) error {
	return Use(defaultLink).Select(dest, query, args...)
}

// Import SQL DDL from io.Reader
func Import(f string) ([]sql.Result, error) {
	return Use(defaultLink).Import(f)
}

// Relation association table builder handle
func Relation(name string, fn BuilderChainFunc) *DB {
	w := Use(defaultLink)
	w.RelationMap = make(map[string]BuilderChainFunc)
	w.RelationMap[name] = fn
	return w
}

// SetDefaultLink set default link name
func SetDefaultLink(db string) {
	defaultLink = db
}
