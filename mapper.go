package gosql

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type Mapper struct {
	database string
	tx       *sqlx.Tx
	SQLBuilder
}

func (m *Mapper) db() ISqlx {
	if m.tx != nil {
		return m.tx
	}

	return DB(m.database)
}

func (m *Mapper) Where(str string, args ...interface{}) *Mapper {
	m.SQLBuilder.Where(str, args...)
	return m
}

//Update data using map[string]interface
func (m *Mapper) Update(data map[string]interface{}) (affected int64, err error) {
	query := m.updateString(data)
	result, err := exec(m.db(), query, m.args...)

	if err != nil {
		return 0, err
	}

	affected, err = result.RowsAffected()
	return affected, err
}

//Create data using map[string]interface
func (m *Mapper) Create(data map[string]interface{}) (lastInsertId int64, err error) {
	query := m.insertString(data)
	result, err := exec(m.db(), query, m.args...)

	if err != nil {
		return 0, err
	}

	lastInsertId, err = result.LastInsertId()

	return lastInsertId, err
}

//Delete data using map[string]interface
func (m *Mapper) Delete() (affected int64, err error) {
	query := m.deleteString()
	result, err := exec(m.db(), query, m.args...)
	if err != nil {
		return 0, err
	}
	affected, err = result.RowsAffected()
	return affected, err
}

func (m *Mapper) Count() (num int64, err error) {
	query := m.countString()

	defer func(start time.Time) {
		logger.Log(&QueryStatus{
			Query: query,
			Args:  m.args,
			Err:   err,
			Start: start,
			End:   time.Now(),
		})
	}(time.Now())
	err = m.db().Get(&num, query, m.args...)
	return num, err
}

// Table is default database new Mapper
func Table(t string, tx ...*sqlx.Tx) *Mapper {
	var txx *sqlx.Tx

	if tx != nil {
		txx = tx[0]
	}

	return &Mapper{database: Default, SQLBuilder: SQLBuilder{table: t}, tx: txx}
}
