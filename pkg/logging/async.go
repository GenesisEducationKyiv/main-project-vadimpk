package logging

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type asyncLogger struct {
	logger *zap.SugaredLogger
}

var _ Logger = (*asyncLogger)(nil)

type Syncer interface {
	Write(p []byte) (n int, err error)
	Sync() error
	Close() error
}

func NewAsyncLogger(syncer Syncer, level string) *asyncLogger {
	var l zapcore.Level
	l, err := zapcore.ParseLevel(level)

	if err != nil {
		l = zap.InfoLevel
	}

	config := zap.NewProductionEncoderConfig()
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(config),
		zapcore.AddSync(syncer),
		l,
	)
	logger := zap.New(core)

	return &asyncLogger{
		logger: logger.Sugar(),
	}
}

func (l *asyncLogger) Named(name string) Logger {
	return &asyncLogger{
		logger: l.logger.Named(name),
	}
}

// With - returns a new logger with parameters.
func (l *asyncLogger) With(args ...interface{}) Logger {
	return &asyncLogger{
		logger: l.logger.With(args...),
	}
}

func (l *asyncLogger) WithContext(ctx context.Context) Logger {
	return l.With("RequestID", ctx.Value("RequestID"))
}

// Debug - logs in debug level.
func (l *asyncLogger) Debug(message string, args ...interface{}) {
	l.logger.Debugw(message, args...)
}

// Info - logs in info level.
func (l *asyncLogger) Info(message string, args ...interface{}) {
	l.logger.Infow(message, args...)
}

// Warn - logs in warn level.
func (l *asyncLogger) Warn(message string, args ...interface{}) {
	l.logger.Warnw(message, args...)
}

// Error - logs in error level.
func (l *asyncLogger) Error(message string, args ...interface{}) {
	l.logger.Errorw(message, args...)
}

// Fatal - logs and exits program with status 1.
func (l *asyncLogger) Fatal(message string, args ...interface{}) {
	l.logger.Fatalw(message, args...)
	os.Exit(1)
}

func (l *asyncLogger) Unwrap() *zap.Logger {
	return l.logger.Desugar()
}
