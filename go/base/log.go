package base

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Log struct {
	Log      *zap.Logger
	LogLevel zap.AtomicLevel
}

func NewDefaultLogger() Log {
	logLevel := zap.NewAtomicLevelAt(zap.InfoLevel)
	log := zap.New(
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
			zapcore.Lock(os.Stdout),
			logLevel,
		),
		zap.AddStacktrace(zap.ErrorLevel),
	)
	return Log{
		Log:      log,
		LogLevel: logLevel,
	}
}

func (l *Log) SetLevel(level zapcore.Level) {
	l.LogLevel.SetLevel(level)
}

func (l *Log) Sync() {
	_ = l.Log.Sync()
}

func (l *Log) Debug(msg string, fields ...zap.Field) {
	l.Log.Debug(msg, fields...)
}

func (l *Log) Debugf(template string, args ...interface{}) {
	l.Log.Sugar().Debugf(template, args...)
}

func (l *Log) Info(msg string, fields ...zap.Field) {
	l.Log.Info(msg, fields...)
}

func (l *Log) Infof(template string, args ...interface{}) {
	l.Log.Sugar().Infof(template, args...)
}

func (l *Log) Warn(msg string, fields ...zap.Field) {
	l.Log.Warn(msg, fields...)
}

func (l *Log) Warnf(template string, args ...interface{}) {
	l.Log.Sugar().Warnf(template, args...)
}

func (l *Log) Error(msg string, fields ...zap.Field) {
	l.Log.Error(msg, fields...)
}

func (l *Log) Errorf(template string, args ...interface{}) {
	l.Log.Sugar().Errorf(template, args...)
}

func (l *Log) Panic(msg string, fields ...zap.Field) {
	l.Log.Panic(msg, fields...)
}

func (l *Log) Panicf(template string, args ...interface{}) {
	l.Log.Sugar().Panicf(template, args...)
}

func (l *Log) Fatal(msg string, fields ...zap.Field) {
	l.Log.Fatal(msg, fields...)
}

func (l *Log) Fatalf(template string, args ...interface{}) {
	l.Log.Sugar().Fatalf(template, args...)
}
