package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/ivan-jorge001/gormt"
	"github.com/ivan-jorge001/gormt/config"
)

func main() {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&interpolateParams=True",
		"root",
		"abc123456",
		"127.0.0.1",
		3306,
		"qmd",
	)

	gormOptions := &gorm.Config{
		PrepareStmt:    false,
		NamingStrategy: schema.NamingStrategy{SingularTable: true}, // 全局禁用表名复数
	}

	orm, err := gorm.Open(mysql.Open(dsn), gormOptions)
	if err != nil {
		log.Fatalf("生成失败")
	}

	conf := config.Config{
		DBConfig: &config.DBConfig{
			Gorm:     orm,
			Database: "qmd",
		},

		PkgName:          "model", //包名称
		OutDir:           "./examples/model",
		DbTag:            "gorm",
		IsJsonTag:        true,
		IsNullToSqlNull:  true,
		TablePrefix:      "q_",
		StripTablePrefix: true,
		OutFileName:      "model",
	}

	gormt.ExecuteConfig(&conf)
}
