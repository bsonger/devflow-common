package logging

import (
	"context"
	"os"
	"strings"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/bsonger/devflow-common/model"
)

//
// =======================
// 全局 Logger（兜底）
// =======================
//

var Logger *zap.Logger

//
// =======================
// context key（私有，防冲突）
// =======================
//

type loggerKeyType struct{}

var loggerKey = loggerKeyType{}

//
// =======================
// 初始化 Zap Logger
// =======================
//

func InitZapLogger(config *model.LogConfig) {
	if config == nil {
		panic("InitZapLogger: log config is nil")
	}

	var cfg zap.Config
	if strings.ToLower(config.Format) == "json" {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
	}

	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	cfg.EncoderConfig.StacktraceKey = "stacktrace"

	// 默认 INFO
	level := zapcore.InfoLevel
	if config.Level != "" {
		_ = level.Set(strings.ToLower(config.Level))
	}
	cfg.Level = zap.NewAtomicLevelAt(level)

	cfg.DisableStacktrace = false
	cfg.DisableCaller = false
	cfg.Development = strings.ToLower(config.Format) != "json"

	logger, err := cfg.Build(
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if err != nil {
		panic(err)
	}

	Logger = withEnvFields(logger)
}

//
// =======================
// 请求级 Logger 注入（核心）
// =======================
//
// ⚠️ 只在「请求入口」调用一次
//

func InjectLogger(ctx context.Context, base *zap.Logger) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if base == nil {
		base = Logger
	}

	log := base

	span := trace.SpanFromContext(ctx)
	if sc := span.SpanContext(); sc.IsValid() {
		log = log.With(
			zap.String("trace_id", sc.TraceID().String()),
			zap.String("span_id", sc.SpanID().String()),
		)
	}

	return context.WithValue(ctx, loggerKey, log)
}

//
// =======================
// 从 context 取 logger（高频调用）
// =======================
//

func LoggerFromContext(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return Logger
	}
	if l, ok := ctx.Value(loggerKey).(*zap.Logger); ok {
		return l
	}
	return Logger
}

//
// =======================
// 向后兼容（可逐步废弃）
// =======================
//

func LoggerWithContext(ctx context.Context) *zap.Logger {
	return LoggerFromContext(ctx)
}

//
// =======================
// ZapAdapter（给 Pyroscope / 第三方库用）
// =======================
//

type ZapAdapter struct {
	logger *zap.Logger
	sugar  *zap.SugaredLogger
}

func NewZapAdapter(logger *zap.Logger) *ZapAdapter {
	if logger == nil {
		logger = Logger
	}
	return &ZapAdapter{
		logger: logger,
		sugar:  logger.Sugar(),
	}
}

func (z *ZapAdapter) Infof(msg string, args ...interface{}) {
	z.sugar.Infof(msg, args...)
}

func (z *ZapAdapter) Debugf(msg string, args ...interface{}) {
	z.sugar.Debugf(msg, args...)
}

func (z *ZapAdapter) Errorf(msg string, args ...interface{}) {
	z.sugar.Errorf(msg, args...)
}

func withEnvFields(l *zap.Logger) *zap.Logger {
	fields := []zap.Field{
		envField("service", "SERVICE_NAME"),
		envField("version", "SERVICE_VERSION"),
		envField("env", "ENV"),
		envField("pod", "POD_NAME"),
		envField("namespace", "POD_NAMESPACE"),
		envField("node", "NODE_NAME"),
		envField("cluster", "CLUSTER_NAME"),
	}

	out := l
	for _, f := range fields {
		if f.Key != "" {
			out = out.With(f)
		}
	}
	return out
}

func envField(key, envKey string) zap.Field {
	if v := os.Getenv(envKey); v != "" {
		return zap.String(key, v)
	}
	return zap.Field{}
}
