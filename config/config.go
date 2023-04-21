package config

import (
	"github.com/spf13/viper"
	"log"
)

var (
	app Config //全局配置
)

type Config struct {
	Log      LogConfig
	DbConfig DbConfig
}

type LogConfig struct {
	Debug bool
	Info  struct {
		Filename string
	}
	Error struct {
		Filename string
	}
}

type DbConfig struct {
	Mysql0 MysqlConfig
	Mysql1 MysqlConfig
}

type MysqlConfig struct {
	Host     string
	User     string
	Port     int
	Password string
	Database string
}

func InitConfig(taxConfigFile string) {
	log.Printf(taxConfigFile)
	viper.SetConfigFile(taxConfigFile)

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&app); err != nil {
		panic(err)
	}
}

func GetApp() Config {
	return app
}
