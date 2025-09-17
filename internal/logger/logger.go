// Package logger provides logging functionality for the scrapego tool. It wraps zap logger
// to provide a simple interface for structured logging.
package logger

import (
	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

func InitLogger(enabled bool) {
	if enabled {
		cfg := zap.NewDevelopmentConfig()
		cfg.OutputPaths = []string{"scrapego.log"}
		cfg.ErrorOutputPaths = []string{"scrapego.log"}

		var err error
		l, err := cfg.Build()
		if err != nil {
			panic("failed to initialize logger " + err.Error())
		}

		logger = l.Sugar()
	} else {
		logger = zap.NewNop().Sugar()
	}
}

func Info(msg string, args ...any) {
	logger.Infof(msg, args...)
}

func Warn(msg string, args ...any) {
	logger.Warnf(msg, args...)
}

func Error(msg string, args ...any) {
	logger.Errorf(msg, args...)
}

func Debug(msg string, args ...any) {
	logger.Debugf(msg, args...)
}

func Sync() {
	_ = logger.Sync()
}
