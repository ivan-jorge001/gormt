package model

import (
	"fmt"
	"log"
	"strings"

	"github.com/wonli/gormt/config"
	"github.com/wonli/gormt/internal/genstruct"
	"github.com/wonli/gormt/mybigcamel"
)

type GenModel struct {
	info DBInfo
	pkg  *genstruct.GenPackage
}

// Generate build code string.生成代码
func Generate(info DBInfo) (out []GenOutInfo, m GenModel) {
	m = GenModel{
		info: info,
	}

	var stt GenOutInfo
	stt.FileCtx = m.generate()
	stt.FileName = info.DbName + ".go"

	if name := config.GetOutFileName(); len(name) > 0 {
		stt.FileName = name + ".go"
	}

	out = append(out, stt)
	return
}

// getTableNameWithPrefix get table name with prefix
func getTableNameWithPrefix(tableName string) string {
	tablePrefix := config.GetTablePrefix()
	if tablePrefix == "" {
		return tableName
	}

	if config.StripTablePrefix() {
		trimPrefix := strings.TrimPrefix(tablePrefix, "-")
		trimPrefix = strings.TrimPrefix(tablePrefix, "_")
		tableName = strings.TrimPrefix(tableName, trimPrefix)

		return tableName
	}

	if strings.HasPrefix(tablePrefix, "-") {
		trimPrefix := strings.TrimPrefix(tablePrefix, "-")
		tableName = strings.TrimPrefix(tableName, trimPrefix)
	} else {
		tableName = tablePrefix + tableName
	}

	return tableName
}

// GetPackage gen struct on table
func (m *GenModel) GetPackage() genstruct.GenPackage {
	if m.pkg == nil {
		var pkg genstruct.GenPackage
		pkg.SetPackage(m.info.PackageName)

		for _, tab := range m.info.TabList {
			var sct genstruct.GenStruct
			sct.SetTableName(tab.Name)
			tab.Name = getTableNameWithPrefix(tab.Name)
			sct.SetStructName(getCamelName(tab.Name))
			sct.SetNotes(tab.Notes)
			sct.AddElement(m.genTableElement(tab.Em)...)
			sct.SetCreatTableStr(tab.SQLBuildStr)
			pkg.AddStruct(sct)

			log.Printf("Gen table %s. [OK]", tab.Name)
		}

		m.pkg = &pkg
	}

	return *m.pkg
}

func (m *GenModel) generate() string {
	m.pkg = nil
	m.GetPackage()
	return m.pkg.Generate()
}

// genTableElement Get table columns and comments.获取表列及注释
func (m *GenModel) genTableElement(cols []ColumnsInfo) (el []genstruct.GenElement) {
	tagGorm := config.GetDBTag()
	for _, v := range cols {
		var tmp genstruct.GenElement
		var isPK bool
		if strings.EqualFold(v.Type, "gorm.Model") { // gorm model
			tmp.SetType(v.Type) //
		} else {
			tmp.SetName(getCamelName(v.Name))
			tmp.SetNotes(v.Notes)
			tmp.SetType(getTypeName(v.Type, v.IsNull))
			// 是否输出gorm标签
			if len(tagGorm) > 0 {
				if strings.EqualFold(v.Extra, "auto_increment") {
					tmp.AddTag(tagGorm, "autoIncrement:true")
				}

				for _, v1 := range v.Index {
					switch v1.Key {
					// case ColumnsKeyDefault:
					case ColumnsKeyPrimary: // primary key.主键
						tmp.AddTag(tagGorm, "primaryKey")
						isPK = true
					case ColumnsKeyUnique: // unique key.唯一索引
						tmp.AddTag(tagGorm, "unique")
					case ColumnsKeyIndex: // index key.复合索引
						uninStr := getUninStr("index", ":", v1.KeyName)
						// 兼容 gorm 本身 sort 标签
						if v1.KeyName == "sort" {
							uninStr = "index"
						}
						if v1.KeyType == "FULLTEXT" {
							uninStr += ",class:FULLTEXT"
						}
						tmp.AddTag(tagGorm, uninStr)
					case ColumnsKeyUniqueIndex: // unique index key.唯一复合索引
						tmp.AddTag(tagGorm, getUninStr("uniqueIndex", ":", v1.KeyName))
					}
				}

			}
		}

		if len(v.Name) > 0 {
			// 是否输出gorm标签
			if len(tagGorm) > 0 {
				tmp.AddTag(tagGorm, "column:"+v.Name)
				tmp.AddTag(tagGorm, "type:"+v.Type)
				if !v.IsNull {
					tmp.AddTag(tagGorm, "not null")
				} else if v.IsNull && !config.GetIsNullToPoint() {
					// 当该字段默认值为null，并且结构不用指针类型时，添加default:null的tag
					tmp.AddTag(tagGorm, "default:null")
				}
				// default tag
				if len(v.Gormt) > 0 {
					tmp.AddTag(tagGorm, v.Gormt)
				}
				if len(v.Notes) > 0 {
					tmp.AddTag(tagGorm, fmt.Sprintf("comment:%v", v.Notes))
				}
			}

			// json tag
			if config.GetIsJsonTag() {
				if isPK && config.GetIsWebTagPkHidden() {
					tmp.AddTag("json", "-")
				} else {
					tmp.AddTag("json", mybigcamel.UnSmallMarshal(mybigcamel.Marshal(v.Name)))
				}
			}
		}

		tmp.ColumnName = v.Name // 列名
		el = append(el, tmp)
	}

	return
}

// genForeignKey Get information about foreign key of table column.获取表列外键相关信息
func (m *GenModel) genForeignKey(col ColumnsInfo) (fklist []genstruct.GenElement) {
	tagGorm := config.GetDBTag()
	for _, v := range col.ForeignKeyList {
		isMulti, isFind, notes := m.getColumnsKeyMulti(v.TableName, v.ColumnName)
		if isFind {
			var tmp genstruct.GenElement
			tmp.SetNotes(notes)
			if isMulti {
				tmp.SetName(getCamelName(v.TableName) + "List")
				tmp.SetType("[]" + getCamelName(v.TableName))
			} else {
				tmp.SetName(getCamelName(v.TableName))
				tmp.SetType(getCamelName(v.TableName))
			}

			tmp.AddTag(tagGorm, "joinForeignKey:"+col.Name)
			tmp.AddTag(tagGorm, "foreignKey:"+v.ColumnName)
			tmp.AddTag(tagGorm, "references:"+getCamelName(col.Name))

			// json tag
			if config.GetIsJsonTag() {
				tmp.AddTag("json", mybigcamel.UnSmallMarshal(mybigcamel.Marshal(v.TableName))+"List")
			}

			fklist = append(fklist, tmp)
		}
	}

	return
}

func (m *GenModel) getColumnsKeyMulti(tableName, col string) (isMulti bool, isFind bool, notes string) {
	var haveGomod bool
	for _, v := range m.info.TabList {
		if strings.EqualFold(v.Name, tableName) {
			for _, v1 := range v.Em {
				if strings.EqualFold(v1.Name, col) {
					for _, v2 := range v1.Index {
						switch v2.Key {
						case ColumnsKeyPrimary, ColumnsKeyUnique, ColumnsKeyUniqueIndex:
							{
								if !v2.Multi { // 唯一索引
									return false, true, v.Notes
								}
							}
						}
					}
					return true, true, v.Notes
				} else if strings.EqualFold(v1.Type, "gorm.Model") {
					haveGomod = true
					notes = v.Notes
				}
			}
			break
		}
	}

	// default gorm.Model
	if haveGomod {
		if strings.EqualFold(col, "id") {
			return false, true, notes
		}

		if strings.EqualFold(col, "created_at") ||
			strings.EqualFold(col, "updated_at") ||
			strings.EqualFold(col, "deleted_at") {
			return true, true, notes
		}
	}

	return false, false, ""
}
