package logging

import (
	"context"
	"strings"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/bsonger/devflow-common/model"
)

var Logger *zap.Logger

func InitZapLogger(ctx context.Context, config *model.LogConfig) {
	if config == nil {
		panic("InitZapLogger: log config is nil")
	}

	var cfg zap.Config
	if config.Format == "json" {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
	}

	// 默认 INFO
	level := zapcore.InfoLevel

	// 根据字符串设置日志级别
	if config.Level != "" {
		if err := level.Set(strings.ToLower(config.Level)); err != nil {
			level = zapcore.InfoLevel
		}
	}

	cfg.Level = zap.NewAtomicLevelAt(level)

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	Logger = logger
}

func LoggerWithContext(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return Logger
	}
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		return Logger
	}
	traceID := span.SpanContext().TraceID().String()
	return Logger.With(zap.String("trace_id", traceID))
}
