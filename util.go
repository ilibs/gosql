package gosql

import (
	"reflect"
	"sort"
	"time"
)

//inSlice
func inSlice(k string, s []string) bool {
	for _, v := range s {
		if k == v {
			return true
		}
	}
	return false
}

//zeroValueFilter filter zero value and keep the specified zero value
func zeroValueFilter(fields map[string]reflect.Value, zv []string) map[string]interface{} {
	m := make(map[string]interface{})

	for k, v := range fields {
		v = reflect.Indirect(v)
		if v.IsValid() && !inSlice(k, zv) {
			if v.Type().Kind() == reflect.Struct {
				if t, ok := v.Interface().(time.Time); ok {
					if t.IsZero() {
						continue
					}
				} else {
					valid := v.FieldByName("Valid").Interface()
					if va, ok := valid.(bool); ok && !va {
						continue
					}
				}
			}

			switch v.Interface().(type) {
			case int, int8, int16, int32, int64:
				if c := v.Int(); c != 0 {
					m[k] = c
				}
			case uint, uint8, uint16, uint32, uint64:
				if c := v.Uint(); c != 0 {
					m[k] = c
				}
			case float32, float64:
				if c := v.Float(); c != 0.0 {
					m[k] = c
				}
			case bool:
				if c := v.Bool(); c != false {
					m[k] = c
				}
			case string:
				if c := v.String(); c != "" {
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

// structAutoTime  auto set created_at updated_at
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

// structToMap
func structToMap(fields map[string]reflect.Value) map[string]interface{} {
	m := make(map[string]interface{})
	for k, v := range fields {
		v = reflect.Indirect(v)
		m[k] = v.Interface()
	}
	return m
}

// sortedParamKeys Sorts the param names given - map iteration order is explicitly random in Go
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
