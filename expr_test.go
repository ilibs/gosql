package gosql

import (
	"testing"
)

func TestExpr(t *testing.T) {
	b := &SQLBuilder{
		table:   "users",
		dialect: mustGetDialect("mysql"),
	}

	q := b.updateString(map[string]interface{}{
		"id":    2,
		"count": Expr("count + ?", 1),
	})

	//fmt.Println(q, b.args)

	if q != "UPDATE `users` SET `count`=count + ?,`id`=?;" {
		t.Error("Expr error,get:", q)
	}

	if len(b.args) != 2 {
		t.Error("Expr args count error")
	}
}
