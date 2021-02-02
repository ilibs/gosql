package gosql

import (
	"testing"
)

func TestGetDialect(t *testing.T) {
	RegisterDialect("mysql", &mysqlDialect{})
	d := newDialect("mysql")
	if d.GetName() != "mysql" {
		t.Fatal("get dialect not mysql")
	}
}
