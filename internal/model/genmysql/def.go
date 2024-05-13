package genmysql

import "regexp"

type keys struct {
	NonUnique  int    `gorm:"column:Non_unique"`
	KeyName    string `gorm:"column:Key_name"`
	ColumnName string `gorm:"column:Column_name"`
	IndexType  string `gorm:"column:Index_type"`
}

// genColumns show full columns
type genColumns struct {
	Field   string  `gorm:"column:Field"`
	Type    string  `gorm:"column:Type"`
	Key     string  `gorm:"column:Key"`
	Desc    string  `gorm:"column:Comment"`
	Null    string  `gorm:"column:Null"`
	Extra   string  `gorm:"Extra"`
	Default *string `gorm:"column:Default"`
}

var noteRegex = regexp.MustCompile(`^\[@gorm\s(\S+)+\]`)
var foreignKeyRegex = regexp.MustCompile(`^\[@fk\s(\S+)+\]`)
