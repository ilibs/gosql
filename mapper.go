package gosql

import (
	"github.com/jmoiron/sqlx"
)

type Mapper struct {
	wrapper *Wrapper
	SQLBuilder
}

func (m *Mapper) ShowSQL() *Mapper {
	m.wrapper.logging = true
	return m
}

func (m *Mapper) db() ISqlx {
	return m.wrapper
}

//Where
func (m *Mapper) Where(str string, args ...interface{}) *Mapper {
	m.SQLBuilder.Where(str, args...)
	return m
}

//Update data from to map[string]interface
func (m *Mapper) Update(data map[string]interface{}) (affected int64, err error) {
	result, err := m.db().Exec(m.updateString(data), m.args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

//Create data from to map[string]interface
func (m *Mapper) Create(data map[string]interface{}) (lastInsertId int64, err error) {
	result, err := m.db().Exec(m.insertString(data), m.args...)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

//Delete data from to map[string]interface
func (m *Mapper) Delete() (affected int64, err error) {
	result, err := m.db().Exec(m.deleteString(), m.args...)
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

	return &Mapper{wrapper: &Wrapper{database: Default, tx: txx}, SQLBuilder: SQLBuilder{table: t}}
}
