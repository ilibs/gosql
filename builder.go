package gosql

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
)

var (
	mapper        = NewReflectMapper("db")
	foreignMapper = reflectx.NewMapper("foreign_key")

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
	model             interface{}
	modelReflectValue reflect.Value
	SQLBuilder
	modelEntity IModel
	wrapper     *Wrapper
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
		b.modelEntity = m
		b.wrapper.database = m.DbName()
		b.table = m.TableName()
		b.modelReflectValue = reflect.ValueOf(m)
	} else {
		value := reflect.ValueOf(b.model)
		if value.Kind() != reflect.Ptr {
			log.Fatalf("model argument must pass a pointer, not a value %#v", b.model)
		}

		if value.IsNil() {
			log.Fatalf("model argument cannot be nil pointer passed")
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
			log.Fatalf("model argument must slice, but get %#v", b.model)
		}

		tpEl := tp.Elem()

		//Compatible with []*Struct or []Struct
		if tpEl.Kind() == reflect.Ptr {
			tpEl = tpEl.Elem()
		}

		if m, ok := reflect.New(tpEl).Interface().(IModel); ok {
			b.modelEntity = m
			b.wrapper.database = m.DbName()
			b.table = m.TableName()
			b.modelReflectValue = reflect.ValueOf(m)
		} else {
			log.Fatalf("model argument must implementation IModel interface or slice []IModel and pointer,but get %#v", b.model)
		}
	}
}

//Hint is set TDDL "/*+TDDL:slave()*/"
func (b *Builder) Hint(hint string) *Builder {
	b.hint = hint
	return b
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
	fields := mapper.FieldMap(b.modelReflectValue)
	if autoTime != nil {
		structAutoTime(fields, autoTime)
	}
	return fields
}

func (b *Builder) relation(refVal reflect.Value) error {
	t := reflect.Indirect(refVal).Type()
	for i := 0; i < t.NumField(); i++ {
		val := t.Field(i).Tag.Get("relation")
		name := t.Field(i).Name
		if val != "" && val != "-" {

			relations := strings.Split(val, ",")

			var foreignModel reflect.Value
			if t.Field(i).Type.Kind() == reflect.Slice {
				foreignModel = reflect.New(t.Field(i).Type)
				fi := foreignModel.Interface()
				err := Model(fi).Where(fmt.Sprintf("%s=?", relations[1]), mapper.FieldByName(refVal, relations[0]).Interface()).All()
				if err != nil {
					return err
				}
				if reflect.Indirect(foreignModel).Len() == 0 {
					reflect.Indirect(refVal).FieldByName(name).Set(reflect.MakeSlice(t.Field(i).Type, 0, 0))
				} else {
					reflect.Indirect(refVal).FieldByName(name).Set(foreignModel.Elem())
				}

			} else {
				foreignModel = reflect.New(t.Field(i).Type.Elem())
				fi := foreignModel.Interface()
				err := Model(fi).Where(fmt.Sprintf("%s=?", relations[1]), mapper.FieldByName(refVal, relations[0]).Interface()).Get()
				if err != nil {
					return err
				}

				reflect.Indirect(refVal).FieldByName(name).Set(foreignModel)
			}
		}
	}
	return nil
}

func (b *Builder) relationAll(refVal reflect.Value) error {
	t := reflect.Indirect(refVal).Type()
	for i := 0; i < t.NumField(); i++ {
		val := t.Field(i).Tag.Get("relation")
		name := t.Field(i).Name
		if val != "" && val != "-" {

			relations := strings.Split(val, ",")

			var foreignModel reflect.Value
			if t.Field(i).Type.Kind() == reflect.Slice {
				foreignModel = reflect.New(t.Field(i).Type)
				fi := foreignModel.Interface()
				err := Model(fi).Where(fmt.Sprintf("%s=?", relations[1]), mapper.FieldByName(refVal, relations[0]).Interface()).All()
				if err != nil {
					return err
				}
				if reflect.Indirect(foreignModel).Len() == 0 {
					reflect.Indirect(refVal).FieldByName(name).Set(reflect.MakeSlice(t.Field(i).Type, 0, 0))
				} else {
					reflect.Indirect(refVal).FieldByName(name).Set(foreignModel.Elem())
				}

			} else {
				foreignModel = reflect.New(t.Field(i).Type.Elem())
				fi := foreignModel.Interface()
				err := Model(fi).Where(fmt.Sprintf("%s=?", relations[1]), mapper.FieldByName(refVal, relations[0]).Interface()).Get()
				if err != nil {
					return err
				}

				reflect.Indirect(refVal).FieldByName(name).Set(foreignModel)
			}
		}
	}
	return nil
}

//All get data row from to Struct
func (b *Builder) Get(zeroValues ...string) (err error) {
	b.initModel()
	hook := NewHook(b.wrapper)
	hook.callMethod("BeforeFind", b.modelReflectValue)
	m := zeroValueFilter(b.reflectModel(nil), zeroValues)
	//If where is empty, the primary key where condition is generated automatically
	b.generateWhere(m)

	err = b.db().Get(b.model, b.queryString(), b.args...)

	if err != nil {
		return err
	}

	//auto relation fill
	err = b.relation(b.modelReflectValue)
	if err == nil {
		hook.callMethod("AfterFind", b.modelReflectValue)
		if hook.HasError() > 0 {
			return hook.Error()
		}
	}

	return err
}

//All get data rows from to Struct
func (b *Builder) All() (err error) {
	b.initModel()

	err = b.db().Select(b.model, b.queryString(), b.args...)
	if err != nil {
		return err
	}

	refVal := reflect.ValueOf(b.model)

	l := reflect.Indirect(refVal).Len()
	t := reflect.Indirect(refVal).Index(0).Elem().Type()

	for i := 0; i < t.NumField(); i++ {
		relVals := make([]interface{}, 0)
		val := t.Field(i).Tag.Get("relation")
		name := t.Field(i).Name
		if val != "" && val != "-" {
			relations := strings.Split(val, ",")
			for j := 0; j < l; j++ {
				relVals = append(relVals, mapper.FieldByName(reflect.Indirect(refVal).Index(j), relations[0]).Interface())
			}

			var foreignModel reflect.Value
			if t.Field(i).Type.Kind() == reflect.Slice {
				foreignModel = reflect.New(t.Field(i).Type)
				fi := foreignModel.Interface()
				err := Model(fi).Where(fmt.Sprintf("%s in(?)", relations[1]), relVals).All()
				if err != nil {
					return err
				}

				fmap := make(map[interface{}]reflect.Value)

				for n := 0; n < reflect.Indirect(foreignModel).Len(); n++ {
					fid := mapper.FieldByName(reflect.Indirect(refVal).Index(n), relations[0])
					fmap[fid.Interface()] = reflect.New(reflect.SliceOf(t.Field(i).Type.Elem())).Elem()
					fmap[fid.Interface()] = reflect.Append(fmap[fid.Interface()],reflect.Indirect(foreignModel).Index(n))
				}

				for j := 0; j < l; j++ {
					pRefVal := mapper.FieldByName(reflect.Indirect(refVal).Index(j), relations[0])
					if pVal, has := fmap[pRefVal.Interface()]; has {
						reflect.Indirect(reflect.Indirect(refVal).Index(j)).FieldByName(name).Set(pVal)
					} else {
						reflect.Indirect(reflect.Indirect(refVal).Index(j)).FieldByName(name).Set(reflect.MakeSlice(t.Field(i).Type, 0, 0))
					}
				}

			} else {
				foreignModel = reflect.New(t.Field(i).Type.Elem())
				fi := reflect.New(reflect.SliceOf(foreignModel.Type()))
				err := Model(fi.Interface()).Where(fmt.Sprintf("%s in(?)", relations[1]), relVals).All()
				if err != nil {
					return err
				}

				fmap := make(map[interface{}]reflect.Value)
				for n := 0; n < reflect.Indirect(fi).Len(); n++ {
					fmap[mapper.FieldByName(reflect.Indirect(refVal).Index(n), relations[0]).Interface()] = reflect.Indirect(fi).Index(n)
				}

				for j := 0; j < l; j++ {
					pRefVal := mapper.FieldByName(reflect.Indirect(refVal).Index(j), relations[0])
					if pVal, has := fmap[pRefVal.Interface()]; has {
						reflect.Indirect(reflect.Indirect(refVal).Index(j)).FieldByName(name).Set(pVal)
					}
				}
			}
		}
	}
	return nil
}

//Create data from to Struct
func (b *Builder) Create() (lastInsertId int64, err error) {
	b.initModel()
	hook := NewHook(b.wrapper)
	hook.callMethod("BeforeChange", b.modelReflectValue)
	hook.callMethod("BeforeCreate", b.modelReflectValue)
	if hook.HasError() > 0 {
		return 0, hook.Error()
	}

	fields := b.reflectModel(AUTO_CREATE_TIME_FIELDS)
	m := structToMap(fields)

	result, err := b.db().Exec(b.insertString(m), b.args...)
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
	hook := NewHook(b.wrapper)
	hook.callMethod("BeforeChange", b.modelReflectValue)
	hook.callMethod("BeforeUpdate", b.modelReflectValue)
	if hook.HasError() > 0 {
		return 0, hook.Error()
	}

	fields := b.reflectModel(AUTO_UPDATE_TIME_FIELDS)
	m := zeroValueFilter(fields, zeroValues)

	//If where is empty, the primary key where condition is generated automatically
	b.generateWhereForPK(m)

	result, err := b.db().Exec(b.updateString(m), b.args...)
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

//gosql.Model(&User{Id:1}).Delete()
func (b *Builder) Delete(zeroValues ...string) (affected int64, err error) {
	b.initModel()
	hook := NewHook(b.wrapper)
	hook.callMethod("BeforeChange", b.modelReflectValue)
	hook.callMethod("BeforeDelete", b.modelReflectValue)
	if hook.HasError() > 0 {
		return 0, hook.Error()
	}

	m := zeroValueFilter(b.reflectModel(nil), zeroValues)
	//If where is empty, the primary key where condition is generated automatically
	b.generateWhere(m)

	result, err := b.db().Exec(b.deleteString(), b.args...)
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

//gosql.Model(&User{}).Where("status = 0").Count()
func (b *Builder) Count(zeroValues ...string) (num int64, err error) {
	b.initModel()

	m := zeroValueFilter(b.reflectModel(nil), zeroValues)
	//If where is empty, the primary key where condition is generated automatically
	b.generateWhere(m)

	err = b.db().Get(&num, b.countString(), b.args...)
	return num, err
}
