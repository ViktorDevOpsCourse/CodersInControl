package logger

import (
	"context"
	"github.com/sirupsen/logrus"
	"os"
)

type contextKey string

const (
	LoggerContextKey contextKey = "ctx-logger"
)

var def *logrus.Logger
var defaultContext context.Context = nil

func GetDefaultLogger() *logrus.Logger {
	if def == nil {
		def = logrus.New()
	}

	def.SetLevel(logrus.DebugLevel)
	def.SetFormatter(initFormatterFromOption("json"))
	def.SetOutput(os.Stdout)
	return def
}

func ToContext(ctx context.Context, l *logrus.Logger) context.Context {
	return context.WithValue(ctx, LoggerContextKey, l)
}

func FromContext(ctx context.Context) *logrus.Logger {
	if l, ok := ctx.Value(LoggerContextKey).(*logrus.Logger); ok {
		return l
	}

	return GetDefaultLogger()
}

func FromDefaultContext() *logrus.Logger {
	if defaultContext == nil {
		defaultContext = context.Background()
	}
	if l, ok := defaultContext.Value(LoggerContextKey).(*logrus.Logger); ok {
		return l
	}

	return GetDefaultLogger()
}

// initFormatterFromOption copied from github.com/bcmi-labs/logger.go
func initFormatterFromOption(format string) logrus.Formatter {
	switch format {
	case "json":
		return &logrus.JSONFormatter{}
	case "text":
		return &logrus.TextFormatter{}
	default:
		return &logrus.JSONFormatter{}
	}
}
