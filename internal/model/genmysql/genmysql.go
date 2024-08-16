package genmysql

import (
	"database/sql"
	"fmt"
	"log"
	"sort"
	"strings"

	"gorm.io/gorm"

	"github.com/ivan-jorge001/gormt/config"
	"github.com/ivan-jorge001/gormt/internal/model"
	"github.com/ivan-jorge001/gormt/tools"
)

// MySQLModel mysql model from IModel
var MySQLModel mysqlModel

type mysqlModel struct {
}

func (m *mysqlModel) GenModel() model.DBInfo {
	orm := config.GetDBConfig().Gorm
	var dbInfo model.DBInfo
	m.getPackageInfo(orm, &dbInfo)
	dbInfo.PackageName = m.GetPkgName()
	dbInfo.DbName = m.GetDbName()
	return dbInfo
}

// GetDbName get database name.获取数据库名字
func (m *mysqlModel) GetDbName() string {
	return config.GetDBConfig().Database
}

// GetTableNames get table name.获取格式化后指定的表名
func (m *mysqlModel) GetTableNames() string {
	return config.GetTableNames()
}

// GetOriginTableNames get table name.获取原始指定的表名
func (m *mysqlModel) GetOriginTableNames() string {
	return config.GetOriginTableNames()
}

// GetPkgName package names through config outdir configuration.通过config outdir 配置获取包名
func (m *mysqlModel) GetPkgName() string {
	setPkgName := config.GetPkgName()
	if setPkgName != "" {
		return setPkgName
	}

	dir := config.GetOutDir()
	dir = strings.Replace(dir, "\\", "/", -1)
	if len(dir) > 0 {
		if dir[len(dir)-1] == '/' {
			dir = dir[:(len(dir) - 1)]
		}
	}
	var pkgName string
	list := strings.Split(dir, "/")
	if len(list) > 0 {
		pkgName = list[len(list)-1]
	}

	if len(pkgName) == 0 || pkgName == "." {
		list = strings.Split(tools.GetModelPath(), "/")
		if len(list) > 0 {
			pkgName = list[len(list)-1]
		}
	}

	return pkgName
}

func (m *mysqlModel) getPackageInfo(orm *gorm.DB, info *model.DBInfo) {
	tables := m.getTables(orm)
	for tabName, notes := range tables {
		var tab model.TabInfo
		tab.Name = tabName
		tab.Notes = notes
		tab.Em = m.getTableElement(orm, tabName)

		info.TabList = append(info.TabList, tab)
	}

	sort.Slice(info.TabList, func(i, j int) bool {
		return info.TabList[i].Name < info.TabList[j].Name
	})
}

// getTableElement Get table columns and comments.获取表列及注释
func (m *mysqlModel) getTableElement(orm *gorm.DB, tab string) (el []model.ColumnsInfo) {
	keyNameCount := make(map[string]int)
	KeyColumnMp := make(map[string][]keys)

	var Keys []keys
	orm.Raw("show keys from " + assemblyTable(tab)).Scan(&Keys)
	for _, v := range Keys {
		keyNameCount[v.KeyName]++
		KeyColumnMp[v.ColumnName] = append(KeyColumnMp[v.ColumnName], v)
	}

	var list []genColumns
	orm.Raw("show FULL COLUMNS from " + assemblyTable(tab)).Scan(&list)

	for _, v := range list {
		var tmp model.ColumnsInfo
		tmp.Name = v.Field
		tmp.Type = v.Type
		tmp.Extra = v.Extra
		FixNotes(&tmp, v.Desc)

		if v.Default != nil {
			if *v.Default == "" {
				tmp.Gormt = "default:''"
			} else {
				tmp.Gormt = fmt.Sprintf("default:%s", *v.Default)
			}
		}

		// keys
		if keyList, ok := KeyColumnMp[v.Field]; ok {
			for _, v := range keyList {
				if v.NonUnique == 0 {
					if strings.EqualFold(v.KeyName, "PRIMARY") {
						tmp.Index = append(tmp.Index, model.KList{
							Key:     model.ColumnsKeyPrimary,
							Multi:   keyNameCount[v.KeyName] > 1,
							KeyType: v.IndexType,
						})
					} else {
						if keyNameCount[v.KeyName] > 1 {
							tmp.Index = append(tmp.Index, model.KList{
								Key:     model.ColumnsKeyUniqueIndex,
								Multi:   keyNameCount[v.KeyName] > 1,
								KeyName: v.KeyName,
								KeyType: v.IndexType,
							})
						} else {
							tmp.Index = append(tmp.Index, model.KList{
								Key:     model.ColumnsKeyUnique,
								Multi:   keyNameCount[v.KeyName] > 1,
								KeyName: v.KeyName,
								KeyType: v.IndexType,
							})
						}
					}
				} else { // mut
					tmp.Index = append(tmp.Index, model.KList{
						Key:     model.ColumnsKeyIndex,
						Multi:   true,
						KeyName: v.KeyName,
						KeyType: v.IndexType,
					})
				}
			}
		}

		tmp.IsNull = strings.EqualFold(v.Null, "YES")
		el = append(el, tmp)
	}
	return
}

// getTables Get columns and comments.获取表列及注释
func (m *mysqlModel) getTables(orm *gorm.DB) map[string]string {
	tbDesc := make(map[string]string)
	var tables []string
	if m.GetOriginTableNames() != "" {
		arr := strings.Split(m.GetOriginTableNames(), ",")
		if len(arr) != 0 {
			for _, val := range arr {
				tbDesc[val] = ""
			}
		}
	} else {
		rows, err := orm.Raw("show tables").Rows()
		if err != nil {
			log.Printf("%s", err)
			return tbDesc
		}

		for rows.Next() {
			var table string
			rows.Scan(&table)
			tables = append(tables, table)
			tbDesc[table] = ""
		}
		rows.Close()
	}

	// Get table annotations.获取表注释
	var err error
	var rows1 *sql.Rows
	if m.GetTableNames() != "" {
		sql1 := fmt.Sprintf(
			"SELECT TABLE_NAME,TABLE_COMMENT FROM information_schema.TABLES WHERE table_schema= '%s' and TABLE_NAME IN(%s)",
			m.GetDbName(),
			m.GetTableNames(),
		)
		rows1, err = orm.Raw(sql1).Rows()
	} else {
		sql1 := fmt.Sprintf(
			"SELECT TABLE_NAME,TABLE_COMMENT FROM information_schema.TABLES WHERE table_schema= '%s'",
			m.GetDbName(),
		)

		rows1, err = orm.Raw(sql1).Rows()
	}

	if err != nil {
		log.Printf("%s", err)
		return tbDesc
	}

	for rows1.Next() {
		var table, desc string
		rows1.Scan(&table, &desc)
		tbDesc[table] = desc
	}
	rows1.Close()

	return tbDesc
}

func assemblyTable(name string) string {
	return "`" + name + "`"
}
