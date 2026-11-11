package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Debug(msg string, field ...zap.Field)
	Info(msg string, field ...zap.Field)
	Warn(msg string, field ...zap.Field)
	Error(msg string, field ...zap.Field)
	Fatal(msg string, field ...zap.Field)
	Sync() error
}

type log struct {
	logger *zap.Logger
}

func New(lvl string) (Logger, error) {
	zapConfig := zap.NewProductionConfig()
	var level zapcore.Level
	if err := level.Set(lvl); err != nil {
		return nil, err
	}

	zapConfig.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	zapConfig.Level = zap.NewAtomicLevelAt(level)

	logger, err := zapConfig.Build(
		zap.WithCaller(true),
		zap.AddStacktrace(zapcore.ErrorLevel),
		zap.IncreaseLevel(level),
		zap.AddCallerSkip(1),
	)
	if err != nil {
		return nil, err
	}
	zap.ReplaceGlobals(logger)

	return &log{logger: logger}, nil
}

func (l *log) Debug(msg string, field ...zap.Field) {
	l.logger.Debug(msg, field...)
}

func (l *log) Info(msg string, field ...zap.Field) {
	l.logger.Info(msg, field...)
}

func (l *log) Warn(msg string, field ...zap.Field) {
	l.logger.Warn(msg, field...)
}

func (l *log) Error(msg string, field ...zap.Field) {
	l.logger.Error(msg, field...)
}

func (l *log) Fatal(msg string, field ...zap.Field) {
	l.logger.Fatal(msg, field...)
}

func (l *log) Sync() error {
	return l.logger.Sync()
}
