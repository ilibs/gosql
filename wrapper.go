package gosql

import (
	"bufio"
	"bytes"
	"context"
	"database/sql"
	"os"
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
}

var (
	defaultWrapper = Use(Default)
)

type Wrapper struct {
	database string
	tx       *sqlx.Tx
	logging  bool
}

func (w *Wrapper) db() ISqlx {
	if w.tx != nil {
		return w.tx.Unsafe()
	}

	return DB(w.database).Unsafe()
}

func ShowSql() *Wrapper  {
	w := Use(Default)
	w.logging = true
	return w
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

	return w.db().Queryx(query, args...)
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

	return w.db().QueryRowx(query, args...)
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

	return w.db().Get(dest, query, args...)
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

	return w.db().Select(dest, query, args...)
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
			tx.Rollback()
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
			tx.Rollback()
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
	return Table(t, w.tx)
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
