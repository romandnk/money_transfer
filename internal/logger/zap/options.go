package zap_logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
)

type Option func(*zap.Config)

func Level(level string) Option {
	return func(c *zap.Config) {
		var lvl zapcore.Level
		switch strings.ToLower(level) {
		case "debug":
			lvl = zapcore.DebugLevel
			c.Level = zap.NewAtomicLevelAt(lvl)
		case "info":
			lvl = zapcore.InfoLevel
			c.Level = zap.NewAtomicLevelAt(lvl)
		case "warn":
			lvl = zapcore.WarnLevel
			c.Level = zap.NewAtomicLevelAt(lvl)
		case "error":
			lvl = zapcore.ErrorLevel
			c.Level = zap.NewAtomicLevelAt(lvl)
		case "dpanic":
			lvl = zapcore.DPanicLevel
			c.Level = zap.NewAtomicLevelAt(lvl)
		case "panic":
			lvl = zapcore.PanicLevel
			c.Level = zap.NewAtomicLevelAt(lvl)
		case "fatal":
			lvl = zapcore.FatalLevel
			c.Level = zap.NewAtomicLevelAt(lvl)
		}
	}
}

func Encoding(encoding string) Option {
	return func(c *zap.Config) {
		c.Encoding = encoding
	}
}

func OutputPaths(outputPaths []string) Option {
	return func(c *zap.Config) {
		c.OutputPaths = outputPaths
	}
}

func ErrorOutputPaths(errorOutputPaths []string) Option {
	return func(c *zap.Config) {
		c.ErrorOutputPaths = errorOutputPaths
	}
}
