package logger

import (
	"context"
	"google.golang.org/grpc/metadata"
	"log"
	"strings"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ctxKey struct{}

var attachedLoggerKey = &ctxKey{}

var globalLogger *zap.SugaredLogger

func fromContext(ctx context.Context) *zap.SugaredLogger {
	var result *zap.SugaredLogger
	if attachedLogger, ok := ctx.Value(attachedLoggerKey).(*zap.SugaredLogger); ok {
		result = attachedLogger
	} else {
		result = globalLogger
	}

	jaegerSpan := opentracing.SpanFromContext(ctx)
	if jaegerSpan != nil {
		if spanCtx, ok := opentracing.SpanFromContext(ctx).Context().(jaeger.SpanContext); ok {
			result = result.With("trace-id", spanCtx.TraceID())
		}
	}

	return result
}

// ErrorKV is a helper function to create an error key/value pair.
func ErrorKV(ctx context.Context, message string, kvs ...interface{}) {
	fromContext(ctx).Errorw(message, kvs...)
}

// InfoKV is a helper function to create an info key/value pair.
func InfoKV(ctx context.Context, message string, kvs ...interface{}) {
	fromContext(ctx).Infow(message, kvs...)
}

// DebugKV is a helper function to create a debug key/value pair.
func DebugKV(ctx context.Context, message string, kvs ...interface{}) {
	fromContext(ctx).Debugw(message, kvs...)
}

// AttachLogger attaches a logger to the context.
func AttachLogger(ctx context.Context, logger *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, attachedLoggerKey, logger)
}

// CloneWithLevel clones the logger with the given level.
func CloneWithLevel(ctx context.Context, newLevel zapcore.Level) *zap.SugaredLogger {
	return fromContext(ctx).
		Desugar().
		WithOptions(WithLevel(newLevel)).
		Sugar()
}

// SetLogger sets the global logger.
func SetLogger(newLogger *zap.SugaredLogger) {
	globalLogger = newLogger
}

// LogLevelFromContext returns the log level from the context.
func LogLevelFromContext(ctx context.Context) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		levels := md.Get("log-level")
		if len(levels) > 0 {
			if parsedLevel, ok := parseLevel(levels[0]); ok {
				newLogger := CloneWithLevel(ctx, parsedLevel)
				ctx = AttachLogger(ctx, newLogger)
			}
		}
	}
	return ctx
}

func parseLevel(str string) (zapcore.Level, bool) {
	switch strings.ToLower(str) {
	case "debug":
		return zapcore.DebugLevel, true
	case "info":
		return zapcore.InfoLevel, true
	case "warn":
		return zapcore.WarnLevel, true
	case "error":
		return zapcore.ErrorLevel, true
	default:
		return zapcore.DebugLevel, false
	}
}

func init() {
	notSugaredLogger, err := zap.NewProduction()
	if err != nil {
		log.Panic(err)
	}

	globalLogger = notSugaredLogger.Sugar()
}
