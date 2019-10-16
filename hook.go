package gosql

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
)

type Hook struct {
	db   *DB
	Errs []error
}

func NewHook(db *DB) *Hook {
	return &Hook{
		db: db,
	}
}

func (h *Hook) callMethod(methodName string, reflectValue reflect.Value) {
	// Only get address from non-pointer
	if reflectValue.CanAddr() && reflectValue.Kind() != reflect.Ptr {
		reflectValue = reflectValue.Addr()
	}

	if methodValue := reflectValue.MethodByName(methodName); methodValue.IsValid() {
		switch method := methodValue.Interface().(type) {
		case func():
			method()
		case func() error:
			h.Err(method())
		case func(db *DB):
			method(h.db)
		case func(db *DB) error:
			h.Err(method(h.db))
		default:
			log.Fatal(fmt.Errorf("unsupported function %v", methodName))
		}
	}
}

// Err add error
func (h *Hook) Err(err error) error {
	if err != nil {
		h.Errs = append(h.Errs, err)
	}
	return err
}

// HasError has errors
func (h *Hook) HasError() int {
	return len(h.Errs)
}

// Error format happened errors
func (h *Hook) Error() error {
	var errs = make([]string, 0)
	for _, e := range h.Errs {
		errs = append(errs, e.Error())
	}
	return errors.New(strings.Join(errs, "; "))
}
