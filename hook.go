package gosql

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
)

type Hook struct {
	wrapper *Wrapper
	Errs    []error
}

func NewHook(wrapper *Wrapper) *Hook {
	return &Hook{
		wrapper: wrapper,
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
		case func(*sqlx.Tx):
			method(h.wrapper.tx)
		case func(*sqlx.Tx) error:
			h.Err(method(h.wrapper.tx))
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
