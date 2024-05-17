package config

import (
	"fmt"
	"log"
	"strings"

	"gorm.io/gorm"
)

// Config custom config struct
type Config struct {
	DBConfig          *DBConfig
	PkgName           string            // 生成文件包名
	OutDir            string            // 保存路径
	OutFileName       string            // 保存文件名
	DbTag             string            // 数据库标签（gorm）
	IsJsonTag         bool              // 是否使用json标签
	IsJsonTagPkHidden bool              // json标记是否隐藏主键
	IsNullToPoint     bool              // null
	IsNullToSqlNull   bool              // sql null
	TablePrefix       string            // 表前缀
	StripTablePrefix  bool              // 移除表前缀
	SelfTypeDef       map[string]string // 自定义类型
	TableNames        string            // 表名（多个表名用","隔开）
}

type DBConfig struct {
	Gorm     *gorm.DB
	Database string //数据库名字
}

func GetDBConfig() *DBConfig {
	return conf.DBConfig
}

// GetOutDir Get Output Directory.获取输出目录
func GetOutDir() string {
	if len(conf.OutDir) == 0 {
		conf.OutDir = "./model"
	}

	return conf.OutDir
}

// GetIsJsonTag json tag.json标记
func GetIsJsonTag() bool {
	return conf.IsJsonTag
}

// GetIsWebTagPkHidden web tag是否隐藏主键
func GetIsWebTagPkHidden() bool {
	return conf.IsJsonTagPkHidden
}

// GetDBTag get database tag.
func GetDBTag() string {
	if conf.DbTag != "gorm" && conf.DbTag != "db" {
		conf.DbTag = "gorm"
	}

	return conf.DbTag
}

// GetIsNullToPoint get if with null to porint in sturct
func GetIsNullToPoint() bool {
	return conf.IsNullToPoint
}

func GetIsNullToSqlNull() bool {
	return conf.IsNullToSqlNull
}

// GetTablePrefix get table prefix
func GetTablePrefix() string {
	return conf.TablePrefix
}

func GetPkgName() string {
	return conf.PkgName
}

func StripTablePrefix() bool {
	return conf.StripTablePrefix
}

// GetSelfTypeDefine 获取自定义字段映射
func GetSelfTypeDefine() map[string]string {
	return conf.SelfTypeDef
}

// GetOutFileName 获取输出文件名
func GetOutFileName() string {
	return conf.OutFileName
}

// GetTableNames get format tableNames by config. 获取格式化后设置的表名
func GetTableNames() string {
	var sb strings.Builder
	if conf.TableNames != "" {
		tableNames := conf.TableNames
		tableNames = strings.TrimLeft(tableNames, ",")
		tableNames = strings.TrimRight(tableNames, ",")
		if tableNames == "" {
			return ""
		}

		arr := strings.Split(conf.TableNames, ",")
		if len(arr) == 0 {
			log.Println("tableNames is vailed, genmodel will by default global")
			return ""
		}

		for i, val := range arr {
			sb.WriteString(fmt.Sprintf("'%s'", val))
			if i != len(arr)-1 {
				sb.WriteString(",")
			}
		}
	}

	return sb.String()
}

// GetOriginTableNames get origin tableNames. 获取原始的设置的表名
func GetOriginTableNames() string {
	return conf.TableNames
}
