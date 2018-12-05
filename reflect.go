package gosql

import (
	"reflect"

	"github.com/jmoiron/sqlx/reflectx"
)

type ReflectMapper struct {
	mapper *reflectx.Mapper
}

func NewReflectMapper(tagName string) *ReflectMapper {
	return &ReflectMapper{
		mapper: reflectx.NewMapper(tagName),
	}
}

// FieldByName returns a field by its mapped name as a reflect.Value.
// Panics if v's Kind is not Struct or v is not Indirectable to a struct Kind.
// Returns zero Value if the name is not found.
func (r *ReflectMapper) FieldByName(v reflect.Value, name string) reflect.Value {
	return r.mapper.FieldByName(v, name)
}

// FieldMap returns the mapper's mapping of field names to reflect values.  Panics
// if v's Kind is not Struct, or v is not Indirectable to a struct kind.
func (r *ReflectMapper) FieldMap(v reflect.Value) map[string]reflect.Value {
	v = reflect.Indirect(v)

	ret := map[string]reflect.Value{}
	tm := r.mapper.TypeMap(v.Type())
	for tagName, fi := range tm.Names {
		//fmt.Println(tagName,fi.Parent.Zero.Kind(),fi.Parent.Field.Anonymous)
		if (fi.Parent.Zero.Kind() == reflect.Struct || (fi.Zero.Kind() == reflect.Ptr && fi.Zero.Type().Elem().Kind() == reflect.Struct)) && !fi.Parent.Field.Anonymous {
			continue
		}
		ret[tagName] = reflectx.FieldByIndexes(v, fi.Index)
	}

	return ret
}
