// Copyright (c) 2025 Alexander Chan
// SPDX-License-Identifier: MIT

// Package logger provides logging functionality for the scrapego tool. It wraps zap logger
// to provide a simple interface for structured logging.
package logger

import (
	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

func InitLogger(filename string, verbose bool, debug bool) {
	cfg := zap.NewDevelopmentConfig()
	if filename == "" {
		cfg.OutputPaths = []string{"stderr"}
		cfg.ErrorOutputPaths = []string{"stderr"}
	} else {
		cfg.OutputPaths = []string{filename}
		cfg.ErrorOutputPaths = []string{filename}
	}

	if debug {
		// debug + info + warn + error
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else if verbose {
		// info + warn + error
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	} else {
		// error only
		cfg.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	}

	var err error
	l, err := cfg.Build()
	if err != nil {
		panic("failed to initialize logger " + err.Error())
	}

	logger = l.Sugar()
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
