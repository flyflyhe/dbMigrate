package main

import (
	"flag"
	"github.com/dbMigrate/v2/config"
	"github.com/dbMigrate/v2/internal/db"
	"github.com/dbMigrate/v2/internal/db/connection"
	"github.com/dbMigrate/v2/internal/filter"
	"github.com/dbMigrate/v2/pkg/logging"
	"github.com/dbMigrate/v2/pkg/tasks"
	"os"
)

var (
	configPath string
	filterPath string
)

func main() {
	flag.StringVar(&configPath, "c", "config.yaml", "配置文件")
	flag.StringVar(&filterPath, "filter", "filters/filter.lua", "过滤文件")
	flag.Parse()
	config.InitConfig(configPath)
	logging.InitLogger(config.GetApp().Log)

	if _, err := os.Stat(filterPath); err != nil {
		logging.Logger.Sugar().Info("未设置过滤文件")
	} else {
		filter.LoadFromFile(filterPath)
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

	if err := task.Start(); err != nil {
		logging.Logger.Sugar().Error(err)
	} else {
		logging.Logger.Sugar().Info("同步完成")
	}
}
