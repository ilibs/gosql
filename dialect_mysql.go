package gosql

import (
	"fmt"
)

type mysqlDialect struct {
	commonDialect
}

func init() {
	RegisterDialect("mysql", &mysqlDialect{})
}

func (mysqlDialect) GetName() string {
	return "mysql"
}

func (mysqlDialect) Quote(key string) string {
	return fmt.Sprintf("`%s`", key)
}
