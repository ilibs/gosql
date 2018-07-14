package gosql

import (
	"log"
	"reflect"
	"strconv"

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

type Builder struct {
	model   interface{}
	SQLBuilder
	wrapper *Wrapper
}

// Model construct SQL from Struct
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
		model:   model,
		wrapper: &Wrapper{tx: txx},
	}
}

func (b *Builder) db() ISqlx {
	return b.wrapper
}

func (b *Builder) initModel() {
	if m, ok := b.model.(IModel); ok {
		b.wrapper.database = m.DbName()
		b.table = m.TableName()
	} else {
		tp := reflect.Indirect(reflect.ValueOf(b.model)).Type()
		if tp.Kind() != reflect.Slice {
			log.Fatalf("model argument must slice, but get %#v", b.model)
		}

		if m, ok := reflect.Indirect(reflect.New(tp.Elem())).Interface().(IModel); ok {
			b.wrapper.database = m.DbName()
			b.table = m.TableName()
		} else {
			log.Fatalf("model argument must implementation IModel interface or slice []IModel and pointer,but get %#v", b.model)
		}
	}
}

//Where for example Where("id = ? and name = ?",1,"test")
func (b *Builder) Where(str string, args ...interface{}) *Builder {
	b.SQLBuilder.Where(str, args...)
	return b
}

//Limit
func (b *Builder) Limit(i int) *Builder {
	b.limit = strconv.Itoa(i)
	return b
}

//Offset
func (b *Builder) Offset(i int) *Builder {
	b.offset = strconv.Itoa(i)
	return b
}

//OrderBy for example "id desc"
func (b *Builder) OrderBy(str string) *Builder {
	b.order = str
	return b
}

//All get data row from to Struct
func (b *Builder) Get() (err error) {
	b.initModel()
	return b.db().Get(b.model, b.queryString(), b.args...)
}

//All get data rows from to Struct
func (b *Builder) All() (err error) {
	b.initModel()
	return b.db().Select(b.model, b.queryString(), b.args...)
}

//Create data from to Struct
func (b *Builder) Create() (lastInsertId int64, err error) {
	b.initModel()

	rv := reflect.Indirect(reflect.ValueOf(b.model))
	fields := mapper.FieldMap(rv)
	structAutoTime(fields, AUTO_CREATE_TIME_FIELDS)
	m := structToMap(fields)

	result, err := b.db().Exec(b.insertString(m), b.args...)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

//gosql.Model(&User{Status:0}).Where("id = ?",1).Update("status")
func (b *Builder) Update(zeroValues ...string) (affected int64, err error) {
	b.initModel()

	uv := reflect.Indirect(reflect.ValueOf(b.model))
	fields := mapper.FieldMap(uv)
	structAutoTime(fields, AUTO_UPDATE_TIME_FIELDS)
	m := zeroValueFilter(fields, zeroValues)

	result, err := b.db().Exec(b.updateString(m), b.args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

//gosql.Model(&User{}).Delete()
func (b *Builder) Delete() (affected int64, err error) {
	b.initModel()
	result, err := b.db().Exec(b.deleteString(), b.args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

//gosql.Model(&User{}).Where("status = 0").Count()
func (b *Builder) Count() (num int64, err error) {
	b.initModel()
	err = b.db().Get(&num, b.countString(), b.args...)
	return num, err
}
