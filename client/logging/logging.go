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
	var cfg zap.Config

	if config.Format == "json" {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
	}

	// 设置日志级别
	// 默认 INFO
	level := zapcore.InfoLevel

	// 根据字符串设置日志级别（大小写不敏感）
	if config.Level != "" {
		if err := level.Set(strings.ToLower(config.Level)); err != nil {
			// 配置错误时兜底 INFO，并打一次标准日志
			level = zapcore.InfoLevel
		}
	}

	_ = level.Set(model.C.Log.Level)
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
