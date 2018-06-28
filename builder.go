package gosql

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
)

var (
	mapper                  = reflectx.NewMapper("db")
	AUTO_CREATE_TIME_FIELDS = []string{
		"create_time",
		"create_at",
		"created_at",
		"update_time",
		"update_at",
		"updated_at",
	}
	AUTO_UPDATE_TIME_FIELDS = []string{
		"update_time",
		"update_at",
		"updated_at",
	}

	logger = &defaultLogger{}
)

type IModel interface {
	TableName() string
	DbName() string
	PK() string
}

type IBuilder interface {
	//条件
	Where(s string, args ...interface{}) IBuilder
	//分页
	Limit(i int) IBuilder
	//分页
	Offset(i int) IBuilder
	//排序
	OrderBy(s string) IBuilder
	//查询一条
	Get() error
	//查询多条
	All() error
	//统计
	Count() (int64, error)
	//创建
	Create() (int64, error)
	//更新
	Update(zeroValues ...string) (int64, error)
	//删除
	Delete() (int64, error)
}

var _ IBuilder = (*Builder)(nil)

type Builder struct {
	model  interface{}
	db     *sqlx.DB
	table  string
	sql    string
	where  string
	order  string
	offset string
	limit  string

	// Extra args to be substituted in the *where* clause
	args []interface{}
}

func Model(model interface{}) *Builder {
	value := reflect.ValueOf(model)
	if value.Kind() != reflect.Ptr {
		log.Fatalf("model argument must pass a pointer, not a value %#v", model)
	}

	if value.IsNil() {
		log.Fatalf("model argument cannot be nil pointer passed")
	}

	builder := &Builder{
		model: model,
	}

	if m, ok := model.(IModel); ok {
		builder.db = DB(m.DbName())
		builder.table = m.TableName()
	} else {
		tp := reflect.Indirect(value).Type()
		if tp.Kind() != reflect.Slice {
			log.Fatalf("model argument must slice, but get %#v", model)
		}

		if m, ok := reflect.Indirect(reflect.New(tp.Elem())).Interface().(IModel); ok {
			builder.db = DB(m.DbName())
			builder.table = m.TableName()
		} else {
			log.Fatalf("model argument must implementation IModel interface or slice []IModel and pointer,but get %#v", model)
		}
	}

	return builder
}

func (b *Builder) Where(str string, args ...interface{}) IBuilder {
	if len(b.where) > 0 {
		b.where = fmt.Sprintf("%s AND (%s)", b.where, str)
	} else {
		b.where = fmt.Sprintf("WHERE (%s)", str)
	}

	// NB this assumes that args are only supplied for where clauses
	// this may be an incorrect assumption!
	if args != nil {
		if b.args == nil {
			b.args = args
		} else {
			b.args = append(b.args, args...)
		}
	}
	return b
}

func (b *Builder) Limit(i int) IBuilder {
	b.limit = fmt.Sprintf("LIMIT %d", i)
	return b
}

func (b *Builder) Offset(i int) IBuilder {
	b.offset = fmt.Sprintf("OFFSET %d", i)
	return b
}

func (b *Builder) OrderBy(str string) IBuilder {
	b.order = fmt.Sprintf("ORDER BY %s", str)
	return b
}

func (b *Builder) queryString() string {
	b.sql = fmt.Sprintf("SELECT * FROM %s %s %s %s %s", b.table, b.where, b.order, b.limit, b.offset)
	b.sql = strings.TrimRight(b.sql, " ")
	b.sql = b.sql + ";"

	return b.sql
}

func (b *Builder) countString() string {
	b.sql = fmt.Sprintf("SELECT count(*) FROM %s %s", b.table, b.where)
	b.sql = strings.TrimRight(b.sql, " ")
	b.sql = b.sql + ";"

	return b.sql
}

func (b *Builder) insertString(params map[string]interface{}) string {
	var cols, vals []string
	for _, k := range sortedParamKeys(params) {
		cols = append(cols, fmt.Sprintf("`%s`", k))
		vals = append(vals, "?")
	}

	b.sql = fmt.Sprintf("INSERT INTO %s (%s) VALUES(%s);", b.table, strings.Join(cols, ","), strings.Join(vals, ","))
	return b.sql
}

func (b *Builder) updateString(params map[string]interface{}) string {
	var updateFields []string
	for _, k := range sortedParamKeys(params) {
		updateFields = append(updateFields, fmt.Sprintf("%s=?", fmt.Sprintf("`%s`", k)))
	}

	b.sql = fmt.Sprintf("UPDATE %s SET %s %s", b.table, strings.Join(updateFields, ","), b.where)
	b.sql = strings.TrimRight(b.sql, " ")
	b.sql = b.sql + ";"
	return b.sql
}

func (b *Builder) deleteString() string {
	b.sql = fmt.Sprintf("DELETE FROM %s %s", b.table, b.where)
	b.sql = strings.TrimRight(b.sql, " ")
	b.sql = b.sql + ";"
	return b.sql
}

func (b *Builder) Get() (err error) {
	query := b.queryString()
	defer func(start time.Time) {
		logger.Log(&QueryStatus{
			Query: query,
			Args:  b.args,
			Err:   err,
			Start: start,
			End:   time.Now(),
		})
	}(time.Now())

	err = b.db.QueryRowx(query, b.args...).StructScan(b.model)
	return err
}

func (b *Builder) All() (err error) {
	query := b.queryString()
	defer func(start time.Time) {
		logger.Log(&QueryStatus{
			Query: b.queryString(),
			Args:  b.args,
			Err:   err,
			Start: start,
			End:   time.Now(),
		})
	}(time.Now())

	err = b.db.Select(b.model, query, b.args...)
	return err
}

func (b *Builder) exec(query string, args ...interface{}) (sql.Result, error) {
	reqult, err := b.db.Exec(query, args...)

	defer func(start time.Time) {
		logger.Log(&QueryStatus{
			Query: query,
			Args:  args,
			Err:   err,
			Start: start,
			End:   time.Now(),
		})
	}(time.Now())
	return reqult, err
}

func (b *Builder) Create() (lastInsertId int64, err error) {
	rv := reflect.Indirect(reflect.ValueOf(b.model))
	fields := mapper.FieldMap(rv)
	structAutoTime(fields, AUTO_CREATE_TIME_FIELDS)
	m := structToMap(fields)

	query := b.insertString(m)
	args := sortedMap(m)
	result, err := b.exec(query, args...)

	if err != nil {
		return 0, err
	}

	lastInsertId, err = result.LastInsertId()

	return lastInsertId, err
}

//gosql.Model(&User{Status:0}).Where("id = ?",1).Update("status")
func (b *Builder) Update(zeroValues ...string) (affected int64, err error) {
	uv := reflect.Indirect(reflect.ValueOf(b.model))
	fields := mapper.FieldMap(uv)
	structAutoTime(fields, AUTO_UPDATE_TIME_FIELDS)
	m := structToMap(fields)
	m = zeroValueFilter(fields, zeroValues)

	query := b.updateString(m)
	args := append(sortedMap(m), b.args...)
	result, err := b.exec(query, args...)

	if err != nil {
		return 0, err
	}

	affected, err = result.RowsAffected()

	return affected, err
}

//gosql.Model(&User{}).Delete()
func (b *Builder) Delete() (affected int64, err error) {
	query := b.deleteString()
	result, err := b.exec(query, b.args...)
	if err != nil {
		return 0, err
	}
	affected, err = result.RowsAffected()
	return affected, err
}

//gosql.Model(&User{}).Where("status = 0").Count()
func (b *Builder) Count() (num int64, err error) {
	var id int
	query := b.countString()

	defer func(start time.Time) {
		logger.Log(&QueryStatus{
			Query: query,
			Args:  b.args,
			Err:   err,
			Start: start,
			End:   time.Now(),
		})
	}(time.Now())
	err = b.db.Get(&id, query, b.args...)

	return int64(id), err
}

func inSlice(k string, s []string) bool {
	for _, v := range s {
		if k == v {
			return true
		}
	}
	return false
}

func zeroValueFilter(fields map[string]reflect.Value, zv []string) map[string]interface{} {
	m := make(map[string]interface{})

	for k, v := range fields {
		v = reflect.Indirect(v)
		if v.IsValid() && !inSlice(k, zv) {
			t, ok := v.Interface().(time.Time)
			if ok && t.IsZero() {
				continue
			}

			switch v.Interface().(type) {
			case int, int8, int16, int32, int64:
				c := v.Int()
				if c != 0 {
					m[k] = c
				}
			case uint, uint8, uint16, uint32, uint64:
				c := v.Uint()
				if c != 0 {
					m[k] = c
				}
			case float32, float64:
				c := v.Float()
				if c != 0.0 {
					m[k] = c
				}
			case bool:
				c := v.Bool()
				if c != false {
					m[k] = c
				}
			case string:
				c := v.String()
				if c != "" {
					m[k] = c
				}
			default:
				m[k] = v.Interface()
			}
		} else {
			m[k] = v.Interface()
		}
	}

	return m
}

func structAutoTime(fields map[string]reflect.Value, f []string) {
	for k, v := range fields {
		v = reflect.Indirect(v)
		if v.IsValid() && inSlice(k, f) {
			switch v.Type().Kind() {
			case reflect.String:
				v.SetString(time.Now().Format("2006-01-02 15:04:05"))
			case reflect.Struct:
				v.Set(reflect.ValueOf(time.Now()))
			}
		}
	}
}

func structToMap(fields map[string]reflect.Value) map[string]interface{} {
	m := make(map[string]interface{})
	for k, v := range fields {
		v = reflect.Indirect(v)
		m[k] = v.Interface()
	}
	return m
}

func sortedMap(m map[string]interface{}) []interface{} {
	var vals []interface{}

	for _, v := range sortedParamKeys(m) {
		vals = append(vals, m[v])
	}
	return vals
}

// Sorts the param names given - map iteration order is explicitly random in Go
// but we need params in a defined order to avoid unexpected results.
func sortedParamKeys(params map[string]interface{}) []string {
	sortedKeys := make([]string, len(params))
	i := 0
	for k := range params {
		sortedKeys[i] = k
		i++
	}
	sort.Strings(sortedKeys)

	return sortedKeys
}
