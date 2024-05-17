package gormt

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/wonli/gormt/config"
	"github.com/wonli/gormt/internal/model"
	"github.com/wonli/gormt/internal/model/genmssql"
	"github.com/wonli/gormt/internal/model/genmysql"
	"github.com/wonli/gormt/internal/model/gensqlite"
	"github.com/wonli/gormt/tools"
)

// ExecuteConfig exe the cmd
func ExecuteConfig(conf *config.Config) {
	if conf != nil {
		config.InitConfig(conf)
	}

	run()
}

func run() {
	var mdb model.IModel
	t := config.GetDBConfig().Gorm.Dialector.Name()
	switch t {
	case "mysql":
		mdb = genmysql.GetModel()
	case "sqlite":
		mdb = gensqlite.GetModel()
	case "mssql":
		mdb = genmssql.GetModel()
	}

	if mdb == nil {
		log.Printf("Check DBConfig.Type (mysql,sqlite,mssql)")
		return
	}

	pkg := mdb.GenModel()
	list, _ := model.Generate(pkg)

	workDir, _ := os.Getwd()
	baseDir := filepath.Join(workDir, config.GetOutDir())

	for _, v := range list {
		path := filepath.Join(baseDir, v.FileName)
		tools.WriteFile(path, []string{v.FileCtx}, true)

		_, _ = exec.Command("goimports", "-l", "-w", path).Output()
		_, _ = exec.Command("gofmt", "-l", "-w", path).Output()
	}
}
