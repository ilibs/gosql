package gosql

import (
	"context"
	"database/sql"
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

type Wrapper struct {
	database string
	tx       *sqlx.Tx
}

func (w *Wrapper) db() ISqlx {
	if w.tx != nil {
		return w.tx.Unsafe()
	}

	return DB(w.database).Unsafe()
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
		})

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
		})
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
		})
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
		})
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
		})
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

//Table new Mapper in Use database
func (w *Wrapper) Table(t string) *Mapper {
	return &Mapper{database: w.database, SQLBuilder: SQLBuilder{table: t}}
}

//Use is change database
func Use(db string) *Wrapper {
	return &Wrapper{database: db}
}

//Exec default database
func Exec(query string, args ...interface{}) (sql.Result, error) {
	return (&Wrapper{database: Default}).Exec(query, args...)
}

//Queryx default database
func Queryx(query string, args ...interface{}) (*sqlx.Rows, error) {
	return (&Wrapper{database: Default}).Queryx(query, args...)
}

//QueryRowx default database
func QueryRowx(query string, args ...interface{}) *sqlx.Row {
	return (&Wrapper{database: Default}).QueryRowx(query, args...)
}

//Txx default database the transaction with context
func Txx(ctx context.Context, fn func(ctx context.Context, tx *sqlx.Tx) error) error {
	return (&Wrapper{database: Default}).Txx(ctx, fn)
}

//Tx default database the transaction
func Tx(fn func(tx *sqlx.Tx) error) error {
	return (&Wrapper{database: Default}).Tx(fn)
}

//Get default database
func Get(dest interface{}, query string, args ...interface{}) error {
	return (&Wrapper{database: Default}).Get(dest, query, args...)
}

//Select default database
func Select(dest interface{}, query string, args ...interface{}) error {
	return (&Wrapper{database: Default}).Select(dest, query, args...)
}
