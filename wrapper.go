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

	return DB(w.database).Exec(query, args...)
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

	return DB(w.database).Queryx(query, args...)

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

	return DB(w.database).QueryRowx(query, args...)
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

	return DB(w.database).Get(dest, query, args...)
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

	return DB(w.database).Select(dest, query, args...)
}

//Tx the transaction
func (w *Wrapper) Tx(ctx context.Context, fn func(ctx context.Context, tx *sqlx.Tx) error) (err error) {
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

	err = fn(ctx, tx)
	if err != nil {
		err = tx.Commit()
	}
	return
}

//Use is change database
func Use(db string) *Wrapper {
	return &Wrapper{db}
}

//Exec default database
func Exec(query string, args ...interface{}) (sql.Result, error) {
	return (&Wrapper{Default}).Exec(query, args...)
}

//Queryx default database
func Queryx(query string, args ...interface{}) (*sqlx.Rows, error) {
	return (&Wrapper{Default}).Queryx(query, args...)
}

//QueryRowx default database
func QueryRowx(query string, args ...interface{}) *sqlx.Row {
	return (&Wrapper{Default}).QueryRowx(query, args...)
}

//Tx default database the transaction
func Tx(ctx context.Context, fn func(ctx context.Context, tx *sqlx.Tx) error) error {
	return (&Wrapper{Default}).Tx(ctx, fn)
}

//Get default database
func Get(dest interface{}, query string, args ...interface{}) error {
	return (&Wrapper{Default}).Get(dest, query, args...)
}

//Select default database
func Select(dest interface{}, query string, args ...interface{}) error {
	return (&Wrapper{Default}).Select(dest, query, args...)
}
