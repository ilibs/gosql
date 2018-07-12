package gosql

import (
	"database/sql"
	"log"
	"reflect"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
)

var (
	mapper = reflectx.NewMapper("db")
	//Insert database automatically updates fields
	AUTO_CREATE_TIME_FIELDS = []string{
		"create_time",
		"create_at",
		"created_at",
		"update_time",
		"update_at",
		"updated_at",
	}
	//Update database automatically updates fields
	AUTO_UPDATE_TIME_FIELDS = []string{
		"update_time",
		"update_at",
		"updated_at",
	}

	logger = &defaultLogger{}
)

//Model interface
type IModel interface {
	TableName() string
	DbName() string
	PK() string
}

type ISqlx interface {
	QueryRowx(query string, args ...interface{}) *sqlx.Row
	Get(dest interface{}, query string, args ...interface{}) error
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
	Select(dest interface{}, query string, args ...interface{}) error
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type Builder struct {
	model    interface{}
	tx       *sqlx.Tx
	database string
	SQLBuilder
}

func Model(model interface{}, tx ...*sqlx.Tx) *Builder {
	value := reflect.ValueOf(model)
	if value.Kind() != reflect.Ptr {
		log.Fatalf("model argument must pass a pointer, not a value %#v", model)
	}

	if value.IsNil() {
		log.Fatalf("model argument cannot be nil pointer passed")
	}

	var txx *sqlx.Tx

	if tx != nil {
		txx = tx[0]
	}

	return &Builder{
		model: model,
		tx:    txx,
	}
}

func (b *Builder) db() ISqlx {
	if b.tx != nil {
		return b.tx.Unsafe()
	}

	return DB(b.database).Unsafe()
}

func (b *Builder) initModel() {
	if m, ok := b.model.(IModel); ok {
		b.database = m.DbName()
		b.table = m.TableName()
	} else {
		value := reflect.ValueOf(b.model)
		tp := reflect.Indirect(value).Type()
		if tp.Kind() != reflect.Slice {
			log.Fatalf("model argument must slice, but get %#v", b.model)
		}

		if m, ok := reflect.Indirect(reflect.New(tp.Elem())).Interface().(IModel); ok {
			b.database = m.DbName()
			b.table = m.TableName()
		} else {
			log.Fatalf("model argument must implementation IModel interface or slice []IModel and pointer,but get %#v", b.model)
		}
	}
}

func (b *Builder) Where(str string, args ...interface{}) *Builder {
	b.SQLBuilder.Where(str, args...)
	return b
}

func (b *Builder) Limit(i int) *Builder {
	b.limit = i
	return b
}

func (b *Builder) Offset(i int) *Builder {
	b.offset = i
	return b
}

func (b *Builder) OrderBy(str string) *Builder {
	b.order = str
	return b
}

func (b *Builder) Get() (err error) {
	b.initModel()

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

	err = b.db().QueryRowx(query, b.args...).StructScan(b.model)
	return err
}

func (b *Builder) All() (err error) {
	b.initModel()

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

	err = b.db().Select(b.model, query, b.args...)
	return err
}

func (b *Builder) Create() (lastInsertId int64, err error) {
	b.initModel()

	rv := reflect.Indirect(reflect.ValueOf(b.model))
	fields := mapper.FieldMap(rv)
	structAutoTime(fields, AUTO_CREATE_TIME_FIELDS)
	m := structToMap(fields)

	query := b.insertString(m)
	result, err := exec(b.db(), query, b.args...)

	if err != nil {
		return 0, err
	}

	lastInsertId, err = result.LastInsertId()

	return lastInsertId, err
}

//gosql.Model(&User{Status:0}).Where("id = ?",1).Update("status")
func (b *Builder) Update(zeroValues ...string) (affected int64, err error) {
	b.initModel()

	uv := reflect.Indirect(reflect.ValueOf(b.model))
	fields := mapper.FieldMap(uv)
	structAutoTime(fields, AUTO_UPDATE_TIME_FIELDS)
	m := zeroValueFilter(fields, zeroValues)

	query := b.updateString(m)
	result, err := exec(b.db(), query, b.args...)

	if err != nil {
		return 0, err
	}

	affected, err = result.RowsAffected()

	return affected, err
}

//gosql.Model(&User{}).Delete()
func (b *Builder) Delete() (affected int64, err error) {
	b.initModel()

	query := b.deleteString()
	result, err := exec(b.db(), query, b.args...)
	if err != nil {
		return 0, err
	}
	affected, err = result.RowsAffected()
	return affected, err
}

//gosql.Model(&User{}).Where("status = 0").Count()
func (b *Builder) Count() (num int64, err error) {
	b.initModel()

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
	err = b.db().Get(&num, query, b.args...)
	return num, err
}
