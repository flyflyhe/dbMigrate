package main

import (
	"flag"
	"github.com/dbMigrate/v2/config"
	"github.com/dbMigrate/v2/internal/db"
	"github.com/dbMigrate/v2/internal/db/connection"
	"github.com/dbMigrate/v2/internal/scripts"
	"github.com/dbMigrate/v2/pkg/logging"
	"github.com/dbMigrate/v2/pkg/tasks"
	"os"
)

var (
	configPath  string
	filterPath  string
	convertPath string
	f           string
)

func main() {
	flag.StringVar(&configPath, "c", "config.yaml", "配置文件")
	flag.StringVar(&filterPath, "filter", "scripts/filter.lua", "过滤文件")
	flag.StringVar(&convertPath, "convert", "scripts/convert.lua", "ddl转换文件")
	flag.StringVar(&f, "f", "compare", "任务执行函数")

	flag.Parse()
	config.InitConfig(configPath)
	logging.InitLogger(config.GetApp().Log)

	if _, err := os.Stat(filterPath); err != nil {
		logging.Logger.Sugar().Info("未设置过滤文件")
	} else {
		scripts.LoadFromFile(filterPath)
	}

	if _, err := os.Stat(convertPath); err != nil {
		logging.Logger.Sugar().Info("未设置ddl转换文件")
	} else {
		scripts.LoadFromFileConvert(convertPath)
	}

	source, err := connection.InitDb(config.GetApp().DbConfig.Mysql0)
	if err != nil {
		panic(err)
	}

	dst, err := connection.InitDb(config.GetApp().DbConfig.Mysql1)
	if err != nil {
		panic(err)
	}

	task := &tasks.Task{
		Source:         &db.Wrapper{DB: source},
		SourceDatabase: config.GetApp().DbConfig.Mysql0.Database,
		Dst:            &db.Wrapper{DB: dst},
		DstDatabase:    config.GetApp().DbConfig.Mysql1.Database,
	}

	if f == "task" {
		if err := task.Start(); err != nil {
			logging.Logger.Sugar().Error(err)
		} else {
			logging.Logger.Sugar().Info("同步完成")
		}
	} else if f == "compare" {
		if err := task.Compare(); err != nil {
			logging.Logger.Sugar().Error(err)
		} else {
			logging.Logger.Sugar().Info("同步完成")
		}
	}
}
