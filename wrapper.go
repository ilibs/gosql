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
}

var (
	defaultWrapper = Use(defaultLink)
)

type BuilderChainFunc func(b *Builder)

type Wrapper struct {
	database    string
	tx          *sqlx.Tx
	logging     bool
	RelationMap map[string]BuilderChainFunc
}

func (w *Wrapper) db() ISqlx {
	if w.tx != nil {
		return w.tx.Unsafe()
	}

	return DB(w.database).Unsafe()
}

func ShowSql() *Wrapper {
	w := Use(defaultLink)
	w.logging = true
	return w
}

func (w *Wrapper) argsIn(query string, args []interface{}) (string, []interface{}, error) {
	newArgs := make([]interface{}, 0)
	newQuery, newArgs, err := sqlx.In(query, args...)

	if err != nil {
		return query, args, err
	}

	return newQuery, newArgs, nil
}

//Rebind wrapper sqlx.Rebind
func (w *Wrapper) Rebind(query string) string {
	return w.db().Rebind(query)
}

//Exec wrapper sqlx.Exec
func (w *Wrapper) Exec(query string, args ...interface{}) (result sql.Result, err error) {
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
func (w *Wrapper) Queryx(query string, args ...interface{}) (rows *sqlx.Rows, err error) {
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
func (w *Wrapper) QueryRowx(query string, args ...interface{}) (rows *sqlx.Row) {
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
func (w *Wrapper) Get(dest interface{}, query string, args ...interface{}) (err error) {
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
func (w *Wrapper) Select(dest interface{}, query string, args ...interface{}) (err error) {
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
func (w *Wrapper) Txx(ctx context.Context, fn func(ctx context.Context, tx *sqlx.Tx) error) (err error) {
	db := DB(w.database)
	tx, err := db.BeginTxx(ctx, nil)
	tx = tx.Unsafe()
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

	err = fn(ctx, tx)
	if err == nil {
		err = tx.Commit()
	}
	return
}

//Tx the transaction
func (w *Wrapper) Tx(fn func(tx *sqlx.Tx) error) (err error) {
	db := DB(w.database)
	tx, err := db.Beginx()
	tx = tx.Unsafe()
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

	err = fn(tx)
	if err == nil {
		err = tx.Commit()
	}
	return
}

//Table database handler from to table name
//for example gosql.Use("db2").Table("users")
func (w *Wrapper) Table(t string) *Mapper {
	return &Mapper{wrapper: w, SQLBuilder: SQLBuilder{table: t}}
}

//Import SQL DDL from sql file
func (w *Wrapper) Import(f string) ([]sql.Result, error) {
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
func (w *Wrapper) Relation(name string, fn BuilderChainFunc) *Wrapper {
	if w.RelationMap == nil {
		w.RelationMap = make(map[string]BuilderChainFunc)
	}
	w.RelationMap[name] = fn
	return w
}

//Use is change database
func Use(db string) *Wrapper {
	return &Wrapper{database: db}
}

//WithTx use the specified transaction session
func WithTx(tx *sqlx.Tx) *Wrapper {
	return &Wrapper{tx: tx}
}

//Exec default database
func Exec(query string, args ...interface{}) (sql.Result, error) {
	return defaultWrapper.Exec(query, args...)
}

//Queryx default database
func Queryx(query string, args ...interface{}) (*sqlx.Rows, error) {
	return defaultWrapper.Queryx(query, args...)
}

//QueryRowx default database
func QueryRowx(query string, args ...interface{}) *sqlx.Row {
	return defaultWrapper.QueryRowx(query, args...)
}

//Txx default database the transaction with context
func Txx(ctx context.Context, fn func(ctx context.Context, tx *sqlx.Tx) error) error {
	return defaultWrapper.Txx(ctx, fn)
}

//Tx default database the transaction
func Tx(fn func(tx *sqlx.Tx) error) error {
	return defaultWrapper.Tx(fn)
}

//Get default database
func Get(dest interface{}, query string, args ...interface{}) error {
	return defaultWrapper.Get(dest, query, args...)
}

//Select default database
func Select(dest interface{}, query string, args ...interface{}) error {
	return defaultWrapper.Select(dest, query, args...)
}

// Import SQL DDL from io.Reader
func Import(f string) ([]sql.Result, error) {
	return defaultWrapper.Import(f)
}

// Relation association table builder handle
func Relation(name string, fn BuilderChainFunc) *Wrapper {
	w := Use(defaultLink)
	w.RelationMap = make(map[string]BuilderChainFunc)
	w.RelationMap[name] = fn
	return w
}

// SetDefaultLink set default link name
func SetDefaultLink(db string) {
	defaultLink = db
	defaultWrapper = Use(defaultLink)
}
