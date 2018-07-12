package gosql

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type Mapper struct {
	database string
	tx       *sqlx.Tx
	table    string
	where    string
	args     map[string]interface{}
}

type INamedSqlx interface {
	NamedExec(query string, arg interface{}) (sql.Result, error)
	PrepareNamed(query string) (*sqlx.NamedStmt, error)
}

func (m *Mapper) db() INamedSqlx {
	if m.tx != nil {
		return m.tx
	}

	return DB(m.database)
}

//insertString Assemble the insert statement
func (m *Mapper) insertString(params map[string]interface{}) string {
	var cols, vals []string
	for k := range params {
		cols = append(cols, fmt.Sprintf("`%s`", k))
		vals = append(vals, fmt.Sprintf(":%s", k))
	}

	return fmt.Sprintf("INSERT INTO %s (%s) VALUES(%s);", m.table, strings.Join(cols, ","), strings.Join(vals, ","))
}

//updateString Assemble the update statement
func (m *Mapper) updateString(params map[string]interface{}) string {
	var updateFields []string
	for k := range params {
		updateFields = append(updateFields, fmt.Sprintf("%s=:%s", fmt.Sprintf("`%s`", k), k))
	}

	query := fmt.Sprintf("UPDATE %s SET %s %s", m.table, strings.Join(updateFields, ","), m.where)
	query = strings.TrimRight(query, " ")
	query = query + ";"
	return query
}

//deleteString Assemble the delete statement
func (m *Mapper) deleteString() string {
	query := fmt.Sprintf("DELETE FROM %s %s", m.table, m.where)
	query = strings.TrimRight(query, " ")
	query = query + ";"
	return query
}

//countString Assemble the count statement
func (m *Mapper) countString() string {
	query := fmt.Sprintf("SELECT count(*) FROM %s %s", m.table, m.where)
	query = strings.TrimRight(query, " ")
	query = query + ";"

	return query
}

func (m *Mapper) Where(str string, args map[string]interface{}) *Mapper {
	if len(m.where) > 0 {
		m.where = fmt.Sprintf("%s AND (%s)", m.where, str)
	} else {
		m.where = fmt.Sprintf("WHERE (%s)", str)
	}

	if args != nil {
		if m.args == nil {
			m.args = args
		} else {
			for k, v := range args {
				m.args[k] = v
			}
		}
	}
	return m
}

func (m *Mapper) exec(query string, data map[string]interface{}) (sql.Result, error) {
	reqult, err := m.db().NamedExec(query, data)

	defer func(start time.Time) {
		logger.Log(&QueryStatus{
			Query: query,
			Args:  data,
			Err:   err,
			Start: start,
			End:   time.Now(),
		})
	}(time.Now())

	return reqult, err
}

//Update data using map[string]interface
func (m *Mapper) Update(data map[string]interface{}) (affected int64, err error) {
	query := m.updateString(data)
	for k, v := range m.args {
		data[k] = v
	}

	result, err := m.exec(query, data)

	if err != nil {
		return 0, err
	}

	affected, err = result.RowsAffected()
	return affected, err
}

//Create data using map[string]interface
func (m *Mapper) Create(data map[string]interface{}) (lastInsertId int64, err error) {
	query := m.insertString(data)
	result, err := m.exec(query, data)

	if err != nil {
		return 0, err
	}

	lastInsertId, err = result.LastInsertId()

	return lastInsertId, err
}

//Delete data using map[string]interface
func (m *Mapper) Delete() (affected int64, err error) {
	query := m.deleteString()
	result, err := m.exec(query, m.args)
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
	nstmt, err := m.db().PrepareNamed(query)
	err = nstmt.Get(&num, m.args)

	return num, err
}

// Table is default database new Mapper
func Table(t string, tx ...*sqlx.Tx) *Mapper {
	var txx *sqlx.Tx

	if tx != nil {
		txx = tx[0]
	}

	return &Mapper{database: Default, table: t, tx: txx}
}
