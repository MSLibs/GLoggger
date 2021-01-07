package glogger

import (
	"context"

	"go.uber.org/zap"
)

func WithContext(ctx *context.Context) GLogger {
	logger := zap.L()
	return GLogger{
		logger,
		logger.Sugar(),
		ctx,
	}
}
