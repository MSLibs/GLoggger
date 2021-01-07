package glogger

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/MSLibs/glogger/core/encoder"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// std is the name of the standard logger in stdlib `log`
	std = CreateLog()
)

const (
	RequestID  string = "requestId"
	PlatformID string = "platformId"
	UserFlag   string = "userFlag"
	Duration   string = "duration"
	Size       string = "size"
)

type GLogger struct {
	log   *zap.Logger
	sugar *zap.SugaredLogger
}

func Info(msg string, fields ...zap.Field) {
	std.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	std.Error(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	std.Debug(msg, fields...)
}

func (log GLogger) Info(msg string, fields ...zap.Field) {
	fields = log.appendFields(fields...)
	log.log.Info(msg, fields...)
}

func (log GLogger) appendFields(fields ...zap.Field) []zap.Field {
	ctx := context.Background()
	if nil == ctx {
		return nil
	}
	fields2 := defaultFields(ctx)
	fields = append(fields, fields2...)
	return fields
}

func defaultFields(ctx context.Context) []zap.Field {
	start, ok := ctx.Value(Duration).(time.Time)
	var requestID, userflag, platformID, duration = "", "", "", ""
	if s, ok := ctx.Value(RequestID).(string); ok {
		requestID = s
	}
	if s, ok := ctx.Value(UserFlag).(string); ok {
		userflag = s
	}
	if s, ok := ctx.Value(PlatformID).(string); ok {
		platformID = s
	}
	if ok {
		duration = strconv.FormatInt(time.Since(start).Milliseconds(), 10)
	}
	var size int64 = -1
	if s, ok := ctx.Value(Size).(int64); ok {
		size = s
	}

	fileds := []zapcore.Field{
		zap.String(RequestID, requestID),
		zap.String(UserFlag, userflag),
		zap.String(PlatformID, platformID),
		zap.String(Duration, duration),
		zap.Int64(Size, size),
	}
	return fileds
}

func (log GLogger) With(fields ...zap.Field) GLogger {
	log.log.With(fields...)
	return log
}

func (log GLogger) Error(msg string, fields ...zap.Field) {
	fields = log.appendFields(fields...)
	log.log.Error(msg, fields...)
}

func (log GLogger) Debug(msg string, fields ...zap.Field) {
	fields = log.appendFields(fields...)
	log.log.Debug(msg, fields...)
}

func (log GLogger) Infof(msg string, args ...interface{}) {
	a := log.defaultLogData()
	if a != nil {
		log.sugar.With(a).Infof(msg, args...)
	} else {
		log.sugar.Infof(msg, args...)
	}
}
func (log GLogger) Errorf(msg string, args ...interface{}) {
	a := log.defaultLogData()
	if a != nil {
		log.sugar.With(a).Errorf(msg, args...)
	} else {
		log.sugar.Errorf(msg, args...)
	}
}
func (log GLogger) Warnf(msg string, args ...interface{}) {
	a := log.defaultLogData()
	if a != nil {
		log.sugar.With(a).Warnf(msg, args...)
	} else {
		log.sugar.Warnf(msg, args...)
	}
}

func (log GLogger) defaultLogData() interface{} {
	ctx := context.Background()
	if nil == ctx {
		return nil
	}

	start, ok := ctx.Value(Duration).(time.Time)
	var duration = ""
	if ok {
		duration = strconv.FormatInt(time.Since(start).Milliseconds(), 10)
	}
	args := struct {
		RequestID  string
		UserFlag   string
		PlatformID string
		Duration   string
		Size       int64
	}{
		ctx.Value(RequestID).(string),
		ctx.Value(UserFlag).(string),
		ctx.Value(PlatformID).(string),
		duration,
		ctx.Value(Size).(int64),
	}

	return args
}

var config zap.Config

func CreateLog() GLogger {
	initDefaultConfig()
	logger, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		logger.Error("logger construction falied")
		panic(err)
	}
	zap.ReplaceGlobals(logger)
	defer logger.Sync()
	logger.Info("logger construction succeeded")
	return GLogger{
		log:   logger,
		sugar: logger.Sugar(),
	}
}

func initDefaultConfig() {
	registerEncoder()
	config = zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Encoding:    "kvpare",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "t",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "trace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     formatEncodeTime,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout", "./tmp/logs"},
		ErrorOutputPaths: []string{"stderr"},
		// InitialFields: map[string]interface{}{
		// 	"requestId":  "",
		// 	"userflag":   "",
		// 	"platformId": "",
		// },
	}
}

func formatEncodeTime(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second()))
}

func registerEncoder() {
	zap.RegisterEncoder("kvpare", func(c zapcore.EncoderConfig) (zapcore.Encoder, error) {
		return encoder.NewKVEncoder(c), nil
	})
}
