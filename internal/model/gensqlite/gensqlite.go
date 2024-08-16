package gensqlite

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"gorm.io/gorm"

	"github.com/ivan-jorge001/gormt/config"
	"github.com/ivan-jorge001/gormt/internal/model"
	"github.com/ivan-jorge001/gormt/tools"
)

// SQLiteModel mysql model from IModel
var SQLiteModel sqliteModel

type sqliteModel struct {
}

// GenModel get model.DBInfo info.获取数据库相关属性
func (m *sqliteModel) GenModel() model.DBInfo {
	db := config.GetDBConfig().Gorm
	defer func() {
		sqldb, _ := db.DB()
		_ = sqldb.Close()
	}()

	var dbInfo model.DBInfo
	m.getPackageInfo(db, &dbInfo)
	dbInfo.PackageName = m.GetPkgName()
	dbInfo.DbName = m.GetDbName()
	return dbInfo
}

// GetDbName get database name.获取数据库名字
func (m *sqliteModel) GetDbName() string {
	dir := config.GetDBConfig().Database
	dir = strings.Replace(dir, "\\", "/", -1)
	if len(dir) > 0 {
		if dir[len(dir)-1] == '/' {
			dir = dir[:(len(dir) - 1)]
		}
	}
	var dbName string
	list := strings.Split(dir, "/")
	if len(list) > 0 {
		dbName = list[len(list)-1]
	}
	list = strings.Split(dbName, ".")
	if len(list) > 0 {
		dbName = list[0]
	}

	if len(dbName) == 0 || dbName == "." {
		panic(fmt.Sprintf("%v : db host config err.must file dir", dbName))
	}

	return dbName
}

// GetTableNames get table name.获取指定的表名
func (m *sqliteModel) GetTableNames() string {
	return config.GetTableNames()
}

// GetPkgName package names through config outdir configuration.通过config outdir 配置获取包名
func (m *sqliteModel) GetPkgName() string {
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

func (m *sqliteModel) getPackageInfo(orm *gorm.DB, info *model.DBInfo) {
	tables := m.getTables(orm) // get table and notes
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
func (m *sqliteModel) getTableElement(orm *gorm.DB, tab string) (el []model.ColumnsInfo) {
	var list []genColumns
	// Get table annotations.获取表注释
	orm.Raw(fmt.Sprintf("PRAGMA table_info(%v)", assemblyTable(tab))).Scan(&list)

	for _, v := range list {
		var tmp model.ColumnsInfo
		tmp.Name = v.Name
		tmp.Type = v.Type
		FixNotes(&tmp, "")
		if v.Pk == 1 { // 主键
			tmp.Index = append(tmp.Index, model.KList{
				Key:   model.ColumnsKeyPrimary,
				Multi: false,
			})
		}

		tmp.IsNull = v.NotNull != 1

		el = append(el, tmp)
	}
	return
}

// getTables Get columns and comments.获取表列及注释
func (m *sqliteModel) getTables(orm *gorm.DB) map[string]string {
	tbDesc := make(map[string]string)

	// Get column names.获取列名
	var tables []string

	rows, err := orm.Raw("SELECT name FROM sqlite_master WHERE type='table'").Rows()
	if err != nil {
		log.Printf("%s", err)
		return tbDesc
	}

	for rows.Next() {
		var table string
		rows.Scan(&table)
		if !strings.EqualFold(table, "sqlite_sequence") { // 剔除系统默认
			tables = append(tables, table)
			tbDesc[table] = ""
		}
	}

	rows.Close()
	return tbDesc
}

func assemblyTable(name string) string {
	return "'" + name + "'"
}
