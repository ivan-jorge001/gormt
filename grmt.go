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
	switch config.GetDbInfo().Type {
	case 0:
		mdb = genmysql.GetModel()
	case 1:
		mdb = gensqlite.GetModel()
	case 2:
		mdb = genmssql.GetModel()
	}

	if mdb == nil {
		log.Printf("Check db_info.type (0:mysql , 1:sqlite , 2:mssql)")
		return
	}

	pkg := mdb.GenModel()
	list, _ := model.Generate(pkg)

	workDir, _ := os.Getwd()
	baseDir := filepath.Join(workDir, config.GetOutDir())

	for _, v := range list {
		path := filepath.Join(baseDir, v.FileName)
		tools.WriteFile(path, []string{v.FileCtx}, true)

		_, _ = exec.Command("fieldalignment", "-fix", path).Output()
		_, _ = exec.Command("goimports", "-l", "-w", path).Output()
		_, _ = exec.Command("gofmt", "-l", "-w", path).Output()
	}
}
