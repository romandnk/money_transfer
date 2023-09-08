package zap_logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

const (
	defaultLevel    zapcore.Level = 0
	defaultEncoding string        = "json"
)

var (
	defaultOutputPath      = []string{"stdout"}
	defaultErrorOutputPath = []string{"stdout"}
)

type Logger struct {
	Log *zap.Logger
}

func New(opts ...Option) (*Logger, error) {
	cfg := &zap.Config{
		Level:            zap.NewAtomicLevelAt(defaultLevel),
		Encoding:         defaultEncoding,
		OutputPaths:      defaultOutputPath,
		ErrorOutputPaths: defaultErrorOutputPath,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "msg",
			LevelKey:   "lvl",
			TimeKey:    "ts",
			EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(t.Format(time.DateTime))
			},
			EncodeLevel: func(lvl zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(lvl.String())
			},
		},
	}

	for _, opt := range opts {
		opt(cfg)
	}

	logger, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("error building zap logging with config: %w", err)
	}

	log := &Logger{
		Log: logger,
	}

	return log, nil
}

func (l *Logger) Info(msg string, fields ...any) {
	zapFields := make([]zap.Field, 0, len(fields))

	for _, field := range fields {
		switch f := field.(type) {
		case zap.Field:
			zapFields = append(zapFields, f)
		default:
			return
		}
	}

	l.Log.Info(msg, zapFields...)
}

func (l *Logger) Error(msg string, fields ...any) {
	zapFields := make([]zap.Field, 0, len(fields))

	for _, field := range fields {
		switch f := field.(type) {
		case zap.Field:
			zapFields = append(zapFields, f)
		default:
			return
		}
	}

	l.Log.Error(msg, zapFields...)
}
