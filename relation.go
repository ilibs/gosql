package gosql

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func eachField(t reflect.Type, fn func(field reflect.StructField, val string, name string, relations []string, connection string) error) error {
	for i := 0; i < t.NumField(); i++ {
		val := t.Field(i).Tag.Get("relation")
		connection := t.Field(i).Tag.Get("connection")
		name := t.Field(i).Name
		field := t.Field(i)

		if val != "" && val != "-" {
			relations := strings.Split(val, ",")
			if len(relations) != 2 {
				return errors.New(fmt.Sprintf("relation tag error, length must 2,but get %v", relations))
			}

			err := fn(field, val, name, relations, connection)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func newModel(value reflect.Value, connection string) *ModelStruct {
	var m *ModelStruct
	if connection != "" {
		m = Use(connection).Model(value.Interface())
	} else {
		m = Model(value.Interface())
	}

	return m
}

// RelationOne is get the associated relational data for a single piece of data
func RelationOne(data interface{}, chains map[string]BuilderChainFunc) error {
	refVal := reflect.Indirect(reflect.ValueOf(data))
	t := refVal.Type()

	return eachField(t, func(field reflect.StructField, val string, name string, relations []string, connection string) error {
		var foreignModel reflect.Value
		// if field type is slice then one-to-many ,eg: []*Struct
		if field.Type.Kind() == reflect.Slice {
			foreignModel = reflect.New(field.Type)
			m := newModel(foreignModel, connection)

			if chainFn, ok := chains[name]; ok {
				chainFn(m)
			}

			// batch get field values
			// Since the structure is slice, there is no need to new Value
			err := m.Where(fmt.Sprintf("%s=?", relations[1]), mapper.FieldByName(refVal, relations[0]).Interface()).All()
			if err != nil {
				return err
			}

			if reflect.Indirect(foreignModel).Len() == 0 {
				// If relation data is empty, must set empty slice
				// Otherwise, the JSON result will be null instead of []
				refVal.FieldByName(name).Set(reflect.MakeSlice(field.Type, 0, 0))
			} else {
				refVal.FieldByName(name).Set(foreignModel.Elem())
			}

		} else {
			// If field type is struct the one-to-one,eg: *Struct
			foreignModel = reflect.New(field.Type.Elem())
			m := newModel(foreignModel, connection)

			if chainFn, ok := chains[name]; ok {
				chainFn(m)
			}

			err := m.Where(fmt.Sprintf("%s=?", relations[1]), mapper.FieldByName(refVal, relations[0]).Interface()).Get()
			// If one-to-one NoRows is not an error that needs to be terminated
			if err != nil && err != sql.ErrNoRows {
				return err
			}

			if err == nil {
				refVal.FieldByName(name).Set(foreignModel)
			}
		}
		return nil
	})
}

// RelationAll is gets the associated relational data for multiple pieces of data
func RelationAll(data interface{}, chains map[string]BuilderChainFunc) error {
	refVal := reflect.Indirect(reflect.ValueOf(data))

	l := refVal.Len()

	if l == 0 {
		return nil
	}

	// get the struct field in slice
	t := reflect.Indirect(refVal.Index(0)).Type()

	return eachField(t, func(field reflect.StructField, val string, name string, relations []string, connection string) error {
		relVals := make([]interface{}, 0)
		relValsMap := make(map[interface{}]interface{}, 0)

		// get relation field values and unique
		for j := 0; j < l; j++ {
			v := mapper.FieldByName(refVal.Index(j), relations[0]).Interface()
			relValsMap[v] = nil
		}

		for k, _ := range relValsMap {
			relVals = append(relVals, k)
		}

		var foreignModel reflect.Value
		// if field type is slice then one to many ,eg: []*Struct
		if field.Type.Kind() == reflect.Slice {
			foreignModel = reflect.New(field.Type)
			m := newModel(foreignModel, connection)

			if chainFn, ok := chains[name]; ok {
				chainFn(m)
			}

			// batch get field values
			// Since the structure is slice, there is no need to new Value
			err := m.Where(fmt.Sprintf("%s in(?)", relations[1]), relVals).All()
			if err != nil {
				return err
			}

			fmap := make(map[interface{}]reflect.Value)

			// Combine relation data as a one-to-many relation
			// For example, if there are multiple images under an article
			// we use the article ID to associate the images, map[1][]*Images
			for n := 0; n < reflect.Indirect(foreignModel).Len(); n++ {
				val := reflect.Indirect(foreignModel).Index(n)
				fid := mapper.FieldByName(val, relations[1])
				if _, has := fmap[fid.Interface()]; !has {
					fmap[fid.Interface()] = reflect.New(reflect.SliceOf(field.Type.Elem())).Elem()
				}
				fmap[fid.Interface()] = reflect.Append(fmap[fid.Interface()], val)
			}

			// Set the result to the model
			for j := 0; j < l; j++ {
				fid := mapper.FieldByName(refVal.Index(j), relations[0])
				if value, has := fmap[fid.Interface()]; has {
					reflect.Indirect(refVal.Index(j)).FieldByName(name).Set(value)
				} else {
					// If relation data is empty, must set empty slice
					// Otherwise, the JSON result will be null instead of []
					reflect.Indirect(refVal.Index(j)).FieldByName(name).Set(reflect.MakeSlice(field.Type, 0, 0))
				}
			}
		} else {
			// If field type is struct the one to one,eg: *Struct
			foreignModel = reflect.New(field.Type.Elem())

			// Batch get field values, but must new slice []*Struct
			fi := reflect.New(reflect.SliceOf(foreignModel.Type()))
			m := newModel(fi, connection)

			if chainFn, ok := chains[name]; ok {
				chainFn(m)
			}

			err := m.Where(fmt.Sprintf("%s in(?)", relations[1]), relVals).All()
			if err != nil {
				return err
			}

			// Combine relation data as a one-to-one relation
			fmap := make(map[interface{}]reflect.Value)
			for n := 0; n < reflect.Indirect(fi).Len(); n++ {
				val := reflect.Indirect(fi).Index(n)
				fid := mapper.FieldByName(val, relations[1])
				fmap[fid.Interface()] = val
			}

			// Set the result to the model
			for j := 0; j < l; j++ {
				fid := mapper.FieldByName(refVal.Index(j), relations[0])
				if value, has := fmap[fid.Interface()]; has {
					reflect.Indirect(refVal.Index(j)).FieldByName(name).Set(value)
				}
			}
		}

		return nil
	})
}
