package genmssql

import "regexp"

// TableDescription 表及表注释
type TableDescription struct {
	Name  string `gorm:"column:name"`  // 表名
	Value string `gorm:"column:value"` // 表注释
}

type ColumnKeys struct {
	ID     int    `gorm:"column:id"`
	Name   string `gorm:"column:name"`   // 列名
	Pk     int    `gorm:"column:pk"`     // 是否主键
	Type   string `gorm:"column:tp"`     // 类型
	Length int    `gorm:"column:len"`    // 长度
	Isnull int    `gorm:"column:isnull"` // 是否为空
	Desc   string `gorm:"column:des"`    // 列注释
}

var noteRegex = regexp.MustCompile(`^\[@gorm\s(\S+)+\]`)
var foreignKeyRegex = regexp.MustCompile(`^\[@fk\s(\S+)+\]`)
