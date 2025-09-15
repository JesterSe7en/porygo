// Package logger provides logging functionality for the scrapego tool. It wraps zap logger
// to provide a simple interface for structured logging.
package logger

import (
	"go.uber.org/zap"
)

var logger *zap.Logger

func InitLogger() {
	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = []string{"scrapego.log"}
	cfg.ErrorOutputPaths = []string{"scrapego.log"}

	var err error
	logger, err = cfg.Build()
	if err != nil {
		panic("failed to initialize logger " + err.Error())
	}
}

func Info(msg string, args ...any) {
	logger.Sugar().Infof(msg, args...)
}

func Warn(msg string, args ...any) {
	logger.Sugar().Warnf(msg, args...)
}

func Error(msg string, args ...any) {
	logger.Sugar().Errorf(msg, args...)
}

func Debug(msg string, args ...any) {
	logger.Sugar().Debugf(msg, args...)
}
