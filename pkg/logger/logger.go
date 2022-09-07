package logger

import "go.uber.org/zap"

// Logger ...
type Logger interface {
	Debugf(format string, a ...any)
	Infof(format string, a ...any)
	Warnf(format string, a ...any)
	Errorf(format string, a ...any)
	Fatalf(format string, a ...any)

	Debug(message string, fields ...zap.Field)
	Info(message string, fields ...zap.Field)
	Warn(message string, fields ...zap.Field)
	Error(message string, fields ...zap.Field)
	Fatal(message string, fields ...zap.Field)
	Sync() error
}
