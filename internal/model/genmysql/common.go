package genmysql

import (
	"log"
	"strings"

	"github.com/wonli/gormt/internal/model"
)

// GetModel get model interface. 获取model接口
func GetModel() model.IModel {
	return &MySQLModel
}

// FixNotes 分析元素表注释
func FixNotes(em *model.ColumnsInfo, note string) {
	b0 := FixElementTag(em, note)        // gorm
	b1 := FixForeignKeyTag(em, em.Notes) // 外键
	if !b0 && b1 {                       // 补偿
		FixElementTag(em, em.Notes) // gorm
	}
}

// FixElementTag 分析元素表注释
func FixElementTag(em *model.ColumnsInfo, note string) bool {
	matches := noteRegex.FindStringSubmatch(note)
	if len(matches) < 2 {
		em.Notes = note
		return false
	}

	log.Printf("get one gorm tag:(%v) ==> (%v)", em.BaseInfo.Name, matches[1])
	em.Notes = note[len(matches[0]):]
	em.Gormt = matches[1]
	return true
}

// FixForeignKeyTag 分析元素表注释(外键)
func FixForeignKeyTag(em *model.ColumnsInfo, note string) bool {
	matches := foreignKeyRegex.FindStringSubmatch(note) // foreign key 外键
	if len(matches) < 2 {
		em.Notes = note
		return false
	}
	em.Notes = note[len(matches[0]):]

	// foreign key 外键
	tmp := strings.Split(matches[1], ".")
	if len(tmp) > 0 {
		log.Printf("get one foreign key:(%v) ==> (%v)", em.BaseInfo.Name, matches[1])
		em.ForeignKeyList = append(em.ForeignKeyList, model.ForeignKey{
			TableName:  tmp[0],
			ColumnName: tmp[1],
		})
	}

	return true
}
