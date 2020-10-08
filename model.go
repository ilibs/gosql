package gosql

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
)

var (
	mapper = NewReflectMapper("db")
	// Insert database automatically updates fields
	AUTO_CREATE_TIME_FIELDS = []string{
		"create_time",
		"create_at",
		"created_at",
		"update_time",
		"update_at",
		"updated_at",
	}
	// Update database automatically updates fields
	AUTO_UPDATE_TIME_FIELDS = []string{
		"update_time",
		"update_at",
		"updated_at",
	}
)

// Model interface
type IModel interface {
	TableName() string
	PK() string
}

type Builder struct {
	model             interface{}
	modelReflectValue reflect.Value
	modelEntity       IModel
	db                *DB
	SQLBuilder
}

// Model construct SQL from Struct
func Model(model interface{}) *Builder {
	return &Builder{
		model: model,
		db:    &DB{database: Sqlx(defaultLink)},
	}
}

// ShowSQL output single sql
func (b *Builder) ShowSQL() *Builder {
	b.db.logging = true
	return b
}

func (b *Builder) initModel() {
	if m, ok := b.model.(IModel); ok {
		b.modelEntity = m
		if b.table == "" {
			b.table = m.TableName()
		}
		b.modelReflectValue = reflect.ValueOf(m)
		b.dialect = newDialect(b.db.DriverName())
	} else {
		value := reflect.ValueOf(b.model)
		if value.Kind() != reflect.Ptr {
			log.Panicf("model argument must pass a pointer, not a value %#v", b.model)
		}

		if value.IsNil() {
			log.Panicf("model argument cannot be nil pointer passed")
		}

		tp := reflect.Indirect(value).Type()

		// If b.model is *interface{} have to do a second Elem
		//
		// For example,
		// var m interface{}
		// mm := make([]*Model,0)
		// mm = append(mm, &Model{Id:1})
		// m = mm
		// reflect.Indirect(reflect.ValueOf(&m)).Elem().Type().Kind() == reflect.Slice

		if tp.Kind() == reflect.Interface {
			tp = reflect.Indirect(value).Elem().Type()
		}

		if tp.Kind() != reflect.Slice {
			log.Panicf("model argument must slice, but get %#v", b.model)
		}

		tpEl := tp.Elem()

		// Compatible with []*Struct or []Struct
		if tpEl.Kind() == reflect.Ptr {
			tpEl = tpEl.Elem()
		}

		if m, ok := reflect.New(tpEl).Interface().(IModel); ok {
			b.modelEntity = m
			if b.table == "" {
				b.table = m.TableName()
			}
			b.modelReflectValue = reflect.ValueOf(m)
			b.dialect = newDialect(b.db.DriverName())
		} else {
			log.Panicf("model argument must implementation IModel interface or slice []IModel and pointer,but get %#v", b.model)
		}
	}
}

// Hint is set TDDL "/*+TDDL:slave()*/"
func (b *Builder) Hint(hint string) *Builder {
	b.hint = hint
	return b
}

// ForceIndex
func (b *Builder) ForceIndex(i string) *Builder {
	b.forceIndex = i
	return b
}

// Where for example Where("id = ? and name = ?",1,"test")
func (b *Builder) Where(str string, args ...interface{}) *Builder {
	b.SQLBuilder.Where(str, args...)
	return b
}

// Select filter column
func (b *Builder) Select(fields string) *Builder {
	b.fields = fields
	return b
}

// Limit
func (b *Builder) Limit(i int) *Builder {
	b.limit = strconv.Itoa(i)
	return b
}

// Offset
func (b *Builder) Offset(i int) *Builder {
	b.offset = strconv.Itoa(i)
	return b
}

// OrderBy for example "id desc"
func (b *Builder) OrderBy(str string) *Builder {
	b.order = str
	return b
}

func (b *Builder) reflectModel(autoTime []string) map[string]reflect.Value {
	fields := mapper.FieldMap(b.modelReflectValue)
	if autoTime != nil {
		structAutoTime(fields, autoTime)
	}
	return fields
}

// Relation association table builder handle
func (b *Builder) Relation(fieldName string, fn BuilderChainFunc) *Builder {
	if b.db.RelationMap == nil {
		b.db.RelationMap = make(map[string]BuilderChainFunc)
	}
	b.db.RelationMap[fieldName] = fn
	return b
}

// All get data row from to Struct
func (b *Builder) Get(zeroValues ...string) (err error) {
	b.initModel()
	m := zeroValueFilter(b.reflectModel(nil), zeroValues)
	// If where is empty, the primary key where condition is generated automatically
	b.generateWhere(m)

	return b.db.Get(b.model, b.queryString(), b.args...)
}

// All get data rows from to Struct
func (b *Builder) All() (err error) {
	b.initModel()
	return b.db.Select(b.model, b.queryString(), b.args...)
}

// Create data from to Struct
func (b *Builder) Create() (lastInsertId int64, err error) {
	b.initModel()
	hook := NewHook(b.db)
	hook.callMethod("BeforeChange", b.modelReflectValue)
	hook.callMethod("BeforeCreate", b.modelReflectValue)
	if hook.HasError() > 0 {
		return 0, hook.Error()
	}

	fields := b.reflectModel(AUTO_CREATE_TIME_FIELDS)
	m := structToMap(fields)

	result, err := b.db.Exec(b.insertString(m), b.args...)
	if err != nil {
		return 0, err
	}

	hook.callMethod("AfterCreate", b.modelReflectValue)
	hook.callMethod("AfterChange", b.modelReflectValue)

	if hook.HasError() > 0 {
		return 0, hook.Error()
	}

	lastId, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	if v, ok := fields[b.modelEntity.PK()]; ok {
		fillPrimaryKey(v, lastId)
	}

	return lastId, err
}

func (b *Builder) generateWhere(m map[string]interface{}) {
	for k, v := range m {
		b.Where(fmt.Sprintf("%s=?", k), v)
	}
}

func (b *Builder) generateWhereForPK(m map[string]interface{}) {
	pk := b.modelEntity.PK()
	pval, has := m[pk]
	if b.where == "" && has {
		b.Where(fmt.Sprintf("%s=?", pk), pval)
		delete(m, pk)
	}
}

// gosql.Model(&User{Id:1,Status:0}).Update("status")
func (b *Builder) Update(zeroValues ...string) (affected int64, err error) {
	b.initModel()
	hook := NewHook(b.db)
	hook.callMethod("BeforeChange", b.modelReflectValue)
	hook.callMethod("BeforeUpdate", b.modelReflectValue)
	if hook.HasError() > 0 {
		return 0, hook.Error()
	}

	fields := b.reflectModel(AUTO_UPDATE_TIME_FIELDS)
	m := zeroValueFilter(fields, zeroValues)

	// If where is empty, the primary key where condition is generated automatically
	b.generateWhereForPK(m)

	result, err := b.db.Exec(b.updateString(m), b.args...)
	if err != nil {
		return 0, err
	}

	hook.callMethod("AfterUpdate", b.modelReflectValue)
	hook.callMethod("AfterChange", b.modelReflectValue)

	if hook.HasError() > 0 {
		return 0, hook.Error()
	}

	return result.RowsAffected()
}

// gosql.Model(&User{Id:1}).Delete()
func (b *Builder) Delete(zeroValues ...string) (affected int64, err error) {
	b.initModel()
	hook := NewHook(b.db)
	hook.callMethod("BeforeChange", b.modelReflectValue)
	hook.callMethod("BeforeDelete", b.modelReflectValue)
	if hook.HasError() > 0 {
		return 0, hook.Error()
	}

	m := zeroValueFilter(b.reflectModel(nil), zeroValues)
	// If where is empty, the primary key where condition is generated automatically
	b.generateWhere(m)

	result, err := b.db.Exec(b.deleteString(), b.args...)
	if err != nil {
		return 0, err
	}

	hook.callMethod("AfterDelete", b.modelReflectValue)
	hook.callMethod("AfterChange", b.modelReflectValue)

	if hook.HasError() > 0 {
		return 0, hook.Error()
	}

	return result.RowsAffected()
}

// gosql.Model(&User{}).Where("status = 0").Count()
func (b *Builder) Count(zeroValues ...string) (num int64, err error) {
	b.initModel()

	m := zeroValueFilter(b.reflectModel(nil), zeroValues)
	// If where is empty, the primary key where condition is generated automatically
	b.generateWhere(m)

	err = b.db.Get(&num, b.countString(), b.args...)
	return num, err
}

// gosql.Model(&users).Table("user_1").Where("status = 0").OrderBy("id DESC").All()
func (b *Builder) Table(tableName string) *Builder {
	b.table = tableName
	return b
}
