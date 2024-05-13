package gensqlite

import "regexp"

// genColumns show full columns
type genColumns struct {
	Name    string `gorm:"column:name"`
	Type    string `gorm:"column:type"`
	Pk      int    `gorm:"column:pk"`
	NotNull int    `gorm:"column:notnull"`
}

var noteRegex = regexp.MustCompile(`^\[@gorm\s(\S+)+\]`)
var foreignKeyRegex = regexp.MustCompile(`^\[@fk\s(\S+)+\]`)
