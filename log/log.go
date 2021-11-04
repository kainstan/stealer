package log

import (
	"go.uber.org/zap"
	"stealer/configs"
)

var (
	AppLog *zap.SugaredLogger
)

// Setup initialize the log instance
func init() {
	AppLog = Init(configs.AppConfig.LogFile)
	AppLog.Info("Init log module")
}

func Info(args ...interface{}) {
	AppLog.Info(args)
}

func Warn(args ...interface{}) {
	AppLog.Warn(args)
}

func Error(args ...interface{}) {
	AppLog.Error(args)
}

