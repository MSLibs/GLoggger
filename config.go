package glogger

import (
	"go.uber.org/zap/zapcore"
)

type GLoggerConfig struct {
	OutputPath string
	Level      zapcore.Level
}

var _ GLoggerConfig = GLoggerConfig{}
