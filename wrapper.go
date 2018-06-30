package gosql

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type Wrapper struct {
	database string
}

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

	return DB(w.database).Exec(query, args...)
}

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

	return DB(w.database).Queryx(query, args...)

}

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

	return DB(w.database).QueryRowx(query, args...)
}

func (w *Wrapper) Tx(ctx context.Context, fn func(ctx context.Context, tx *sqlx.Tx) error) (err error) {
	db := DB(w.database)
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	err = fn(ctx, tx)
	if err != nil {
		err = tx.Commit()
	}
	return
}

//change database
func Use(db string) *Wrapper {
	return &Wrapper{db}
}

//default database Exec
func Exec(query string, args ...interface{}) (sql.Result, error) {
	return (&Wrapper{Default}).Exec(query, args...)
}

//default database Queryx
func Queryx(query string, args ...interface{}) (*sqlx.Rows, error) {
	return (&Wrapper{Default}).Queryx(query, args...)

}

//default database QueryRowx
func QueryRowx(query string, args ...interface{}) *sqlx.Row {
	return (&Wrapper{Default}).QueryRowx(query, args...)
}

func Tx(ctx context.Context, fn func(ctx context.Context,tx *sqlx.Tx) error) (err error) {
	return (&Wrapper{Default}).Tx(ctx, fn)
}
