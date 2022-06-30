package base

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	logger   *zap.SugaredLogger
	logLevel zap.AtomicLevel
}

func NewZapLogger() zapLogger {
	logLevel := zap.NewAtomicLevelAt(zap.InfoLevel)
	log := zap.New(
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
			zapcore.Lock(os.Stdout),
			logLevel,
		),
		zap.AddStacktrace(zap.ErrorLevel),
	).Sugar()
	return zapLogger{
		logger:   log,
		logLevel: logLevel,
	}
}

func (l *zapLogger) SetLevel(level zapcore.Level) {
	l.logLevel.SetLevel(level)
}

func (l *zapLogger) Sync() {
	_ = l.logger.Sync()
}

func (l *zapLogger) Debug(args ...interface{}) {
	l.logger.Debug(args)
}

func (l *zapLogger) Debugf(template string, args ...interface{}) {
	l.logger.Debugf(template, args...)
}

func (l *zapLogger) Info(args ...interface{}) {
	l.logger.Info(args)
}

func (l *zapLogger) Infof(template string, args ...interface{}) {
	l.logger.Infof(template, args...)
}

func (l *zapLogger) Warning(args ...interface{}) {
	l.logger.Warn(args)
}

func (l *zapLogger) Warningf(template string, args ...interface{}) {
	l.logger.Warnf(template, args...)
}

func (l *zapLogger) Error(args ...interface{}) {
	l.logger.Error(args)
}

func (l *zapLogger) Errorf(template string, args ...interface{}) error {
	err := fmt.Errorf(template, args...)
	l.logger.Error(err)
	return err
}

func (l *zapLogger) Errore(err error) error {
	if err != nil {
		l.logger.Error(err.Error())
	}
	return err
}

func (l *zapLogger) Panic(args ...interface{}) {
	l.logger.Panic(args)
}

func (l *zapLogger) Panicf(template string, args ...interface{}) error {
	err := fmt.Errorf(template, args...)
	l.logger.Panic(err)
	return err
}

func (l *zapLogger) Panice(err error) {
	l.logger.Panic(err.Error())
}

// panic in Fatal
func (l *zapLogger) Fatal(args ...interface{}) {
	l.Panic(args)
}

// panic in Fatalf
func (l *zapLogger) Fatalf(template string, args ...interface{}) error {
	return l.Panicf(template, args...)
}

// panic in Fatale
func (l *zapLogger) Fatale(err error) {
	l.Panice(err)
}
