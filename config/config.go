package config

import (
	"fmt"
	"log"
	"strings"
)

// Config custom config struct
type Config struct {
	DBInfo            DBInfo            `yaml:"db_info"`
	PkgName           string            `yaml:"pkg_name"`              // 生成文件包名
	OutDir            string            `yaml:"out_dir"`               // 保存路径
	OutFileName       string            `yaml:"out_file_name"`         // 保存文件名
	DbTag             string            `yaml:"db_tag"`                // 数据库标签（gorm）
	IsJsonTag         bool              `yaml:"is_json_tag"`           // 是否使用json标签
	IsJsonTagPkHidden bool              `yaml:"is_json_tag_pk_hidden"` // json标记是否隐藏主键
	IsNullToPoint     bool              `yaml:"is_null_to_point"`      // null
	IsNullToSqlNull   bool              `yaml:"is_null_to_sql_null"`   // sql null
	TablePrefix       string            `yaml:"table_prefix"`          // 表前缀
	StripTablePrefix  bool              `yaml:"strip_table_prefix"`    // 移除表前缀
	SelfTypeDef       map[string]string `yaml:"self_type_define"`      // 自定义类型
	TableNames        string            `yaml:"table_names"`           // 表名（多个表名用","隔开）
}

// DBInfo mysql database information. mysql 数据库信息
type DBInfo struct {
	Host     string `validate:"required"` // Host. 地址
	Port     int    // Port 端口号
	Username string // Username 用户名
	Password string // Password 密码
	Database string // Database 数据库名
	Type     int    // 数据库类型: 0:mysql , 1:sqlite , 2:mssql
}

// GetDbInfo Get configuration information .获取数据配置信息
func GetDbInfo() DBInfo {
	return conf.DBInfo
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
