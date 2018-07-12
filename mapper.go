package gosql

import (
	"github.com/jmoiron/sqlx"
)

type Mapper struct {
	database string
	tx       *sqlx.Tx
	SQLBuilder
}

func (m *Mapper) db() ISqlx {
	return &Wrapper{
		database: m.database,
		tx:       m.tx,
	}
}

//Where
func (m *Mapper) Where(str string, args ...interface{}) *Mapper {
	m.SQLBuilder.Where(str, args...)
	return m
}

//Update data from to map[string]interface
func (m *Mapper) Update(data map[string]interface{}) (affected int64, err error) {
	result, err := exec(m.db(), m.updateString(data), m.args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

//Create data from to map[string]interface
func (m *Mapper) Create(data map[string]interface{}) (lastInsertId int64, err error) {
	result, err := exec(m.db(), m.insertString(data), m.args...)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

//Delete data from to map[string]interface
func (m *Mapper) Delete() (affected int64, err error) {
	result, err := exec(m.db(), m.deleteString(), m.args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

//Count data from to map[string]interface
func (m *Mapper) Count() (num int64, err error) {
	err = m.db().Get(&num, m.countString(), m.args...)
	return num, err
}

// Table select table name
func Table(t string, tx ...*sqlx.Tx) *Mapper {
	var txx *sqlx.Tx

	if tx != nil {
		txx = tx[0]
	}

	return &Mapper{database: Default, SQLBuilder: SQLBuilder{table: t}, tx: txx}
}
