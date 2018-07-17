package gosql

import (
	"fmt"
	"strings"
)

type SQLBuilder struct {
	table  string
	where  string
	order  string
	limit  string
	offset string
	// Extra args to be substituted in the *where* clause
	args []interface{}
}

func (s *SQLBuilder) limitFormat() string {
	if s.limit != "" {
		return fmt.Sprintf("LIMIT %s", s.limit)
	}
	return ""
}

func (s *SQLBuilder) offsetFormat() string {
	if s.offset != "" {
		return fmt.Sprintf("OFFSET %s", s.offset)
	}
	return ""
}

func (s *SQLBuilder) orderFormat() string {
	if s.order != "" {
		return fmt.Sprintf("ORDER BY %s", s.order)
	}
	return ""
}

//queryString Assemble the query statement
func (s *SQLBuilder) queryString() string {
	query := fmt.Sprintf("SELECT * FROM %s %s %s %s %s", s.table, s.where, s.orderFormat(), s.limitFormat(), s.offsetFormat())
	query = strings.TrimRight(query, " ")
	query = query + ";"

	return query
}

//countString Assemble the count statement
func (s *SQLBuilder) countString() string {
	query := fmt.Sprintf("SELECT count(*) FROM %s %s %s %s", s.table, s.where, s.limitFormat(), s.offsetFormat())
	query = strings.TrimRight(query, " ")
	query = query + ";"

	return query
}

//insertString Assemble the insert statement
func (s *SQLBuilder) insertString(params map[string]interface{}) string {
	var cols, vals []string
	for _, k := range sortedParamKeys(params) {
		cols = append(cols, fmt.Sprintf("`%s`", k))
		vals = append(vals, "?")
		s.args = append(s.args, params[k])
	}

	return fmt.Sprintf("INSERT INTO %s (%s) VALUES(%s);", s.table, strings.Join(cols, ","), strings.Join(vals, ","))
}

//updateString Assemble the update statement
func (s *SQLBuilder) updateString(params map[string]interface{}) string {
	var updateFields []string
	args := make([]interface{},0)
	for _, k := range sortedParamKeys(params) {
		updateFields = append(updateFields, fmt.Sprintf("%s=?", fmt.Sprintf("`%s`", k)))
		args = append(args, params[k])
	}
	args = append(args,s.args...)
	s.args = args

	query := fmt.Sprintf("UPDATE %s SET %s %s", s.table, strings.Join(updateFields, ","), s.where)
	query = strings.TrimRight(query, " ")
	query = query + ";"
	return query
}

//deleteString Assemble the delete statement
func (s *SQLBuilder) deleteString() string {
	query := fmt.Sprintf("DELETE FROM %s %s", s.table, s.where)
	query = strings.TrimRight(query, " ")
	query = query + ";"
	return query
}

func (s *SQLBuilder) Where(str string, args ...interface{}) {
	if s.where != "" {
		s.where = fmt.Sprintf("%s AND (%s)", s.where, str)
	} else {
		s.where = fmt.Sprintf("WHERE (%s)", str)
	}

	// NB this assumes that args are only supplied for where clauses
	// this may be an incorrect assumption!
	if args != nil {
		if s.args == nil {
			s.args = args
		} else {
			s.args = append(s.args, args...)
		}
	}
}
