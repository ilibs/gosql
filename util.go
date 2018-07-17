package gosql

import (
	"bufio"
	"bytes"
	"database/sql"
	"io"
	"reflect"
	"sort"
	"strings"
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

// Import SQL DDL from io.Reader
func Import(db ISqlx, r io.Reader) ([]sql.Result, error) {
	var results []sql.Result
	scanner := bufio.NewScanner(r)

	semiColSpliter := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		if i := bytes.IndexByte(data, ';'); i >= 0 {
			return i + 1, data[0:i], nil
		}
		// If we're at EOF, we have a final, non-terminated line. Return it.
		if atEOF {
			return len(data), data, nil
		}
		// Request more data.
		return 0, nil, nil
	}

	scanner.Split(semiColSpliter)

	for scanner.Scan() {
		query := strings.Trim(scanner.Text(), " \t\n\r")
		if len(query) > 0 {
			result, err := db.Exec(query)
			results = append(results, result)
			if err != nil {
				return nil, err
			}
		}
	}

	return results, nil
}
