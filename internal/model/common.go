package model

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/wonli/gormt/config"
	"github.com/wonli/gormt/internal/cnf"
	"github.com/wonli/gormt/mybigcamel"
)

// getCamelName Big Hump or Capital Letter.大驼峰或者首字母大写
func getCamelName(name string) string {
	return mybigcamel.Marshal(strings.ToLower(name))
}

// getTypeName Type acquisition filtering.类型获取过滤
func getTypeName(name string, isNull bool) string {
	// 优先匹配自定义类型
	selfDefineTypeMqlDicMap := config.GetSelfTypeDefine()
	if v, ok := selfDefineTypeMqlDicMap[name]; ok {
		return fixNullToPorint(v, isNull)
	}

	// Fuzzy Regular Matching.模糊正则匹配自定义类型
	for selfKey, selfVal := range selfDefineTypeMqlDicMap {
		if ok, _ := regexp.MatchString(selfKey, name); ok {
			return fixNullToPorint(selfVal, isNull)
		}
	}

	// Precise matching first.先精确匹配
	if v, ok := cnf.TypeMysqlDicMp[name]; ok {
		return fixNullToPorint(v, isNull)
	}

	// Fuzzy Regular Matching.模糊正则匹配
	for _, l := range cnf.TypeMysqlMatchList {
		if ok, _ := regexp.MatchString(l.Key, name); ok {
			return fixNullToPorint(l.Value, isNull)
		}
	}

	panic(fmt.Sprintf("type (%v) not match in any way.maybe need to add on (https://github.com/xxjwxc/gormt/blob/master/data/view/cnf/def.go)", name))
}

// 过滤null point 类型
func fixNullToPorint(name string, isNull bool) string {
	if isNull && config.GetIsNullToPoint() {
		if strings.HasPrefix(name, "uint") {
			return "*" + name
		}
		if strings.HasPrefix(name, "int") {
			return "*" + name
		}
		if strings.HasPrefix(name, "float") {
			return "*" + name
		}
		if strings.HasPrefix(name, "date") {
			return "*" + name
		}
		if strings.HasPrefix(name, "time") {
			return "*" + name
		}
		if strings.HasPrefix(name, "bool") {
			return "*" + name
		}
		if strings.HasPrefix(name, "string") {
			return "*" + name
		}
	}
	if isNull && config.GetIsNullToSqlNull() {

		if strings.HasPrefix(name, "uint") {
			return "sql.NullInt64"
		}
		if strings.HasPrefix(name, "int") {
			return "sql.NullInt32"
		}
		if strings.HasPrefix(name, "float") {
			return "sql.NullFloat64"
		}
		if strings.HasPrefix(name, "date") {
			return "sql.NullTime"
		}
		if strings.HasPrefix(name, "time") {
			return "sql.NullTime"
		}
		if strings.HasPrefix(name, "bool") {
			return "sql.NullBool"
		}
		if strings.HasPrefix(name, "string") {
			return "sql.NullString"
		}
	}

	return name
}

func getUninStr(left, middle, right string) string {
	re := left
	if len(right) > 0 {
		re = left + middle + right
	}
	return re
}
