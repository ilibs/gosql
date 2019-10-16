package gosql

import (
	"fmt"
)

// Dialect interface contains behaviors that differ across SQL database
type Dialect interface {
	// GetName get dialect's name
	GetName() string

	// Quote quotes field name to avoid SQL parsing exceptions by using a reserved word as a field name
	Quote(key string) string
}

type commonDialect struct {
}

func (commonDialect) GetName() string {
	return "common"
}

func (commonDialect) Quote(key string) string {
	return fmt.Sprintf(`"%s"`, key)
}

var dialectsMap = map[string]Dialect{}

// RegisterDialect register new dialect
func RegisterDialect(name string, dialect Dialect) {
	dialectsMap[name] = dialect
}

// GetDialect gets the dialect for the specified dialect name
func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectsMap[name]
	return
}

func mustGetDialect(name string) Dialect {
	if dialect, ok := dialectsMap[name]; ok {
		return dialect
	}
	panic(fmt.Sprintf("`%v` is not officially supported", name))
	return nil
}

func newDialect(name string) Dialect {
	if value, ok := GetDialect(name); ok {
		return value
	}

	fmt.Printf("`%v` is not officially supported, running under compatibility mode.\n", name)
	return &commonDialect{}
}
