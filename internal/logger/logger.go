package logger

import (
	"go.uber.org/zap"
)

var logger *zap.Logger

func InitLogger() {
	cfg := zap.NewDevelopmentConfig()

	var err error
	logger, err = cfg.Build()
	if err != nil {
		panic("failed to initialize logger " + err.Error())
	}
}

func Info(msg string, fields ...any) {
	logger.Sugar().Infow(msg, fields...)
}

func Warn(msg string, fields ...any) {
	logger.Sugar().Warnw(msg, fields...)
}

func Error(msg string, fields ...any) {
	logger.Sugar().Errorw(msg, fields...)
}

func Debug(msg string, fields ...any) {
	logger.Sugar().Debugw(msg, fields...)
}
