package logging

import (
	"github.com/dbMigrate/v2/config"
	"go.uber.org/zap"
)

var Logger *zap.Logger

func InitLogger(logConfig config.LogConfig) {
	initZap(logConfig)
	Logger = zapLog
}
