package main

import (
	"github.com/wonli/gormt"
	"github.com/wonli/gormt/config"
)

func main() {
	dbInfo := config.DBInfo{
		Host:     "127.0.0.1",
		Port:     3306,
		Username: "root",
		Password: "123456",
		Database: "test",
		Type:     0,
	}

	conf := config.Config{
		DBInfo:           dbInfo,
		PkgName:          "schema", //包名称
		OutDir:           "./examples/model",
		DbTag:            "gorm",
		IsJsonTag:        true,
		IsNullToSqlNull:  true,
		TablePrefix:      "q_",
		StripTablePrefix: true,
		OutFileName:      "schema",
	}

	gormt.ExecuteConfig(&conf)
}
