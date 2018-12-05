package main

import (
	"encoding/json"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ilibs/gosql"
	"github.com/ilibs/gosql/example/models"
)

func connect()  {
	configs := make(map[string]*gosql.Config)

	configs["default"] = &gosql.Config{
		Enable:  true,
		Driver:  "mysql",
		Dsn:     "root:123456@tcp(127.0.0.1:3306)/mykids?charset=utf8&parseTime=True&loc=Asia%2FShanghai",
		ShowSql: true,
	}
	gosql.Connect(configs)
}

func main()  {
	connect()

	//models.MomentGetList()
	m,err := models.MomentGetList()
	b , _ :=json.MarshalIndent(m,"","	")

	fmt.Println(string(b),err)
}
