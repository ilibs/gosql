package gosql

import (
	"context"
	"errors"
	"log"
	"reflect"
	"strings"
)

type Hook struct {
	db   *DB
	Errs []error
	ctx  context.Context
}

func NewHook(ctx context.Context, db *DB) *Hook {
	return &Hook{
		db:  db,
		ctx: ctx,
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
		case func(ctx context.Context):
			method(h.ctx)
		case func(ctx context.Context) error:
			h.Err(method(h.ctx))
		case func(ctx context.Context, db *DB):
			method(h.ctx, h.db)
		case func(ctx context.Context, db *DB) error:
			h.Err(method(h.ctx, h.db))
		default:
			log.Panicf("unsupported function %v", methodName)
		}
	}
}

// Err add error
func (h *Hook) Err(err error) {
	if err != nil {
		h.Errs = append(h.Errs, err)
	}
}

// HasError has errors
func (h *Hook) HasError() bool {
	return len(h.Errs) > 0
}

// Error format happened errors
func (h *Hook) Error() error {
	var errs = make([]string, 0)
	for _, e := range h.Errs {
		errs = append(errs, e.Error())
	}
	return errors.New(strings.Join(errs, "; "))
}
