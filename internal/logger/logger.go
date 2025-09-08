package logger

import (
	"go.uber.org/zap"
)

var Logger *zap.Logger

func InitLogger() {
	cfg := zap.NewDevelopmentConfig()

	var err error
	Logger, err = cfg.Build()
	if err != nil {
		panic("failed to initialize logger " + err.Error())
	}
}
