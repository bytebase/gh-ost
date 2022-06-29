package base

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	Logger   *zap.Logger
	LogLevel zap.AtomicLevel
}

func NewDefaultLogger() Logger {
	logLevel := zap.NewAtomicLevelAt(zap.InfoLevel)
	log := zap.New(
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
			zapcore.Lock(os.Stdout),
			logLevel,
		),
		zap.AddStacktrace(zap.ErrorLevel),
	)
	return Logger{
		Logger:   log,
		LogLevel: logLevel,
	}
}

func (l *Logger) SetLevel(level zapcore.Level) {
	l.LogLevel.SetLevel(level)
}

func (l *Logger) Sync() {
	_ = l.Logger.Sync()
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.Logger.Debug(msg, fields...)
}

func (l *Logger) Debugf(template string, args ...interface{}) {
	l.Logger.Sugar().Debugf(template, args...)
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.Logger.Info(msg, fields...)
}

func (l *Logger) Infof(template string, args ...interface{}) {
	l.Logger.Sugar().Infof(template, args...)
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.Logger.Warn(msg, fields...)
}

func (l *Logger) Warnf(template string, args ...interface{}) {
	l.Logger.Sugar().Warnf(template, args...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.Logger.Error(msg, fields...)
}

func (l *Logger) Errorf(template string, args ...interface{}) {
	l.Logger.Sugar().Errorf(template, args...)
}

func (l *Logger) Errore(err error) {
	l.Logger.Error(err.Error())
}

func (l *Logger) Panic(msg string, fields ...zap.Field) {
	l.Logger.Panic(msg, fields...)
}

func (l *Logger) Panicf(template string, args ...interface{}) {
	l.Logger.Sugar().Panicf(template, args...)
}

func (l *Logger) Panice(err error) {
	l.Logger.Panic(err.Error())
}

func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.Logger.Fatal(msg, fields...)
}

func (l *Logger) Fatalf(template string, args ...interface{}) {
	l.Logger.Sugar().Fatalf(template, args...)
}
func (l *Logger) Fatale(err error) {
	l.Logger.Fatal(err.Error())
}
