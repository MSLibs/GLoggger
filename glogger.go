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
	std = CreateLog(GLoggerConfig{Level: 0, OutputPath: "./logs"})
)

const (
	RequestID  string = "requestId"
	PlatformID string = "platformId"
	UserFlag   string = "userflag"
	Duration   string = "duration"
	Size       string = "size"
	UserAgent  string = "userAgent"
	Referer    string = "referer"
	Method     string = "method"
	Url        string = "url"
	ServerIP   string = "serverip"
	SourceIP   string = "sourceip"
)

type GLogger struct {
	log     *zap.Logger
	sugar   *zap.SugaredLogger
	context *context.Context
}

type FormatTemplateWithor interface {
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

func WithInfo(ctx *context.Context, msg string, fields ...zap.Field) {
	std.WithInfo(ctx, msg, fields...)
}

func WithError(ctx *context.Context, msg string, fields ...zap.Field) {
	std.WithError(ctx, msg, fields...)
}

func WithDebug(ctx *context.Context, msg string, fields ...zap.Field) {
	std.WithDebug(ctx, msg, fields...)
}

func (log GLogger) Info(msg string, fields ...zap.Field) {
	fields = log.appendFields(fields...)
	log.log.Info(msg, fields...)
}

func (log GLogger) WithInfo(ctx *context.Context, msg string, fields ...zap.Field) {
	fields = log.SetContext(ctx).appendFields(fields...)
	log.log.Info(msg, fields...)
}

func (log GLogger) appendFields(fields ...zap.Field) []zap.Field {
	if nil == log.context {
		context := context.Background()
		log.context = &context
	}
	fields2 := defaultFields(*log.context)
	fields = append(fields, fields2...)
	return fields
}

func (log GLogger) SetContext(ctx *context.Context) GLogger {
	log.context = ctx
	return log
}

func defaultFields(ctx context.Context) []zap.Field {
	// payload := read(ctx)
	// var slice = writeFields(payload).([]zap.Field)
	// var fields []zap.Field
	// copy(fields, slice[:])
	start, ok := ctx.Value(Duration).(time.Time)
	var requestID, userflag, platformID, duration, userAgent, referer, method, url, serverip, sourceip = "", "", "", "", "", "", "", "", "", ""
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
	if s, ok := ctx.Value(UserAgent).(string); ok {
		userAgent = s
	}
	if s, ok := ctx.Value(Referer).(string); ok {
		referer = s
	}
	if s, ok := ctx.Value(Method).(string); ok {
		method = s
	}
	if s, ok := ctx.Value(Url).(string); ok {
		url = s
	}
	if s, ok := ctx.Value(ServerIP).(string); ok {
		serverip = s
	}
	if s, ok := ctx.Value(SourceIP).(string); ok {
		sourceip = s
	}

	fileds := []zapcore.Field{
		zap.String(RequestID, requestID),
		zap.String(UserFlag, userflag),
		zap.String(PlatformID, platformID),
		zap.String(Duration, duration),
		zap.Int64(Size, size),
		zap.String(UserAgent, userAgent),
		zap.String(Referer, referer),
		zap.String(Method, method),
		zap.String(Url, url),
		zap.String(SourceIP, sourceip),
		zap.String(ServerIP, serverip),
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

func (log GLogger) WithError(ctx *context.Context, msg string, fields ...zap.Field) {
	fields = log.SetContext(ctx).appendFields(fields...)
	log.log.Error(msg, fields...)
}

func (log GLogger) Debug(msg string, fields ...zap.Field) {
	fields = log.appendFields(fields...)
	log.log.Debug(msg, fields...)
}

func (log GLogger) WithDebug(ctx *context.Context, msg string, fields ...zap.Field) {
	fields = log.SetContext(ctx).appendFields(fields...)
	log.log.Debug(msg, fields...)
}

func (log GLogger) Infof(msg string, args ...interface{}) {
	a := log.buildFormatTemplateWithor()
	if a != nil {
		log.sugar.With(a...).Infof(msg, args...)
	} else {
		log.sugar.Infof(msg, args...)
	}
}

func (log GLogger) WithInfof(ctx *context.Context, msg string, args ...interface{}) {
	a := log.SetContext(ctx).buildFormatTemplateWithor()
	if a != nil {
		log.sugar.With(a...).Infof(msg, args...)
	} else {
		log.sugar.Infof(msg, args...)
	}
}

//TODO 待优化，既然格式一定，是不是可以暂时直接写死
func (log GLogger) Withf(args []interface{}) *zap.SugaredLogger {
	log.sugar.With(args...)
	for _, v := range args {
		log.sugar.With(v)
	}
	return log.sugar
}
func (log GLogger) Errorf(msg string, args ...interface{}) {
	a := log.defaultLogData()
	if a != nil {
		log.sugar.With(a).Errorf(msg, args...)
	} else {
		log.sugar.Errorf(msg, args...)
	}
}
func (log GLogger) WithErrorf(ctx *context.Context, msg string, args ...interface{}) {
	a := log.SetContext(ctx).defaultLogData()
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
func (log GLogger) WithWarnf(ctx *context.Context, msg string, args ...interface{}) {
	a := log.SetContext(ctx).defaultLogData()
	if a != nil {
		log.sugar.With(a).Warnf(msg, args...)
	} else {
		log.sugar.Warnf(msg, args...)
	}
}

func (log GLogger) buildFormatTemplateWithor() []interface{} {
	if log.context == nil {
		context := context.Background()
		log.context = &context
	}
	ctx := *log.context
	payload := read(ctx)
	fields := writeFields(payload)
	return sweetenFields(fields)
}

func sweetenFields(fields []zap.Field) []interface{} {
	slices := make([]interface{}, 0, len(fields))
	for i := 0; i < len(fields); {
		slices = append(slices, fields[i])
		i++
		continue
	}
	return slices
}

func (log GLogger) defaultLogData() interface{} {
	if log.context == nil {
		context := context.Background()
		log.context = &context
	}
	args := read(*log.context)
	return args
}

func read(ctx context.Context) (payload LogPayload) {
	//TODO：重复代码，用类来代替
	if s, ok := ctx.Value(RequestID).(string); ok {
		payload.RequestID = s
	}
	if s, ok := ctx.Value(UserFlag).(string); ok {
		payload.UserFlag = s
	}
	if s, ok := ctx.Value(PlatformID).(string); ok {
		payload.PlatformID = s
	}
	if start, ok := ctx.Value(Duration).(time.Time); ok {
		payload.Duration = strconv.FormatInt(time.Since(start).Milliseconds(), 10)
	}
	if s, ok := ctx.Value(Size).(int64); ok {
		payload.Size = s
	}
	if s, ok := ctx.Value(UserAgent).(string); ok {
		payload.UserAgent = s
	}
	if s, ok := ctx.Value(Referer).(string); ok {
		payload.Referer = s
	}
	if s, ok := ctx.Value(Method).(string); ok {
		payload.Method = s
	}
	if s, ok := ctx.Value(Url).(string); ok {
		payload.Url = s
	}
	if s, ok := ctx.Value(ServerIP).(string); ok {
		payload.ServerIP = s
	}
	if s, ok := ctx.Value(SourceIP).(string); ok {
		payload.SourceIP = s
	}
	return
}
func writeFields(payload LogPayload) []zap.Field {
	return []zap.Field{
		zap.String(RequestID, payload.RequestID),
		zap.String(UserFlag, payload.UserFlag),
		zap.String(PlatformID, payload.PlatformID),
		zap.String(Duration, payload.Duration),
		zap.Int64(Size, payload.Size),
		zap.String(UserAgent, payload.UserAgent),
		zap.String(Referer, payload.Referer),
		zap.String(Method, payload.Method),
		zap.String(Url, payload.Url),
		zap.String(SourceIP, payload.SourceIP),
		zap.String(ServerIP, payload.ServerIP),
	}
}

var config zap.Config

func CreateLog(gconfig GLoggerConfig) GLogger {
	initDefaultConfig(gconfig)
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

func initDefaultConfig(gconfig GLoggerConfig) {
	registerEncoder()
	level := gconfig.Level
	outputs := []string{"stdout"}
	if gconfig.OutputPath == "" {
		gconfig.OutputPath = "./tmp/logs"
	}
	outputs = append(outputs, gconfig.OutputPath)

	config = zap.Config{
		Level:       zap.NewAtomicLevelAt(level),
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
		OutputPaths:      outputs,
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

type LogPayload struct {
	RequestID  string
	UserFlag   string
	PlatformID string
	Duration   string
	Size       int64
	UserAgent  string
	Referer    string
	Method     string
	Url        string
	ServerIP   string
	SourceIP   string
}
