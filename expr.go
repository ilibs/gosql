package gosql

// SQL expression
type expr struct {
	expr string
	args []interface{}
}

// Expr generate raw SQL expression, for example:
//     gosql.Table("user").Update(map[string]interface{}{"price", gorm.Expr("price * ? + ?", 2, 100)})
func Expr(expression string, args ...interface{}) *expr {
	return &expr{expr: expression, args: args}
}
