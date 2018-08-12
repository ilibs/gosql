package gosql

import (
	"fmt"
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
	var txx *sqlx.Tx

	if tx != nil {
		txx = tx[0]
	}

	return &Builder{
		model:   model,
		wrapper: &Wrapper{tx: txx},
	}
}

func (b *Builder) ShowSQL() *Builder {
	b.wrapper.logging = true
	return b
}

func (b *Builder) db() ISqlx {
	return b.wrapper
}

func (b *Builder) initModel() {
	if m, ok := b.model.(IModel); ok {
		b.wrapper.database = m.DbName()
		b.table = m.TableName()
	} else {
		value := reflect.ValueOf(b.model)
		if value.Kind() != reflect.Ptr {
			log.Fatalf("model argument must pass a pointer, not a value %#v", b.model)
		}

		if value.IsNil() {
			log.Fatalf("model argument cannot be nil pointer passed")
		}

		tp := reflect.Indirect(value).Type()
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

func (b *Builder) reflectModel(autoTime []string) map[string]reflect.Value {
	uv := reflect.Indirect(reflect.ValueOf(b.model))
	fields := mapper.FieldMap(uv)
	if autoTime != nil {
		structAutoTime(fields, autoTime)
	}
	return fields
}

//All get data row from to Struct
func (b *Builder) Get(zeroValues ...string) (err error) {
	b.initModel()
	m := zeroValueFilter(b.reflectModel(nil), zeroValues)
	//If where is empty, the primary key where condition is generated automatically
	b.generateWhere(m)

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

	fields := b.reflectModel(AUTO_CREATE_TIME_FIELDS)
	m := structToMap(fields)

	result, err := b.db().Exec(b.insertString(m), b.args...)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (b *Builder) generateWhere(m map[string]interface{}) {
	for k, v := range m {
		b.Where(fmt.Sprintf("%s=?", k), v)
	}
}

func (b *Builder) generateWhereForPK(m map[string]interface{}) {
	pk := b.model.(IModel).PK()
	pval, has := m[pk]
	if b.where == "" && has {
		b.Where(fmt.Sprintf("%s=?", pk), pval)
		delete(m, pk)
	}
}

//gosql.Model(&User{Id:1,Status:0}).Update("status")
func (b *Builder) Update(zeroValues ...string) (affected int64, err error) {
	b.initModel()

	fields := b.reflectModel(AUTO_UPDATE_TIME_FIELDS)
	m := zeroValueFilter(fields, zeroValues)

	//If where is empty, the primary key where condition is generated automatically
	b.generateWhereForPK(m)

	result, err := b.db().Exec(b.updateString(m), b.args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

//gosql.Model(&User{Id:1}).Delete()
func (b *Builder) Delete(zeroValues ...string) (affected int64, err error) {
	b.initModel()

	m := zeroValueFilter(b.reflectModel(nil), zeroValues)
	//If where is empty, the primary key where condition is generated automatically
	b.generateWhere(m)

	result, err := b.db().Exec(b.deleteString(), b.args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

//gosql.Model(&User{}).Where("status = 0").Count()
func (b *Builder) Count(zeroValues ...string) (num int64, err error) {
	b.initModel()

	m := zeroValueFilter(b.reflectModel(nil), zeroValues)
	//If where is empty, the primary key where condition is generated automatically
	b.generateWhere(m)

	err = b.db().Get(&num, b.countString(), b.args...)
	return num, err
}
