package logger

import (
	"fmt"
	"github.com/abdivasiyev/project_template/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/url"
	"time"

	"go.uber.org/zap/zapcore"

	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var Module = fx.Provide(New)

type zapLogger struct {
	zap *zap.SugaredLogger
}

type Params struct {
	fx.In
	Config config.Config
}

// DebugLevel shows all messages which is written in log
// InfoLevel shows info messages
// WarnLevel shows warning messages
// ErrorLevel shows error messages
// PanicLevel shows panic messages and stops service
// FatalLevel shows fatal messages and stops service
const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
	PanicLevel = "panic"
	FatalLevel = "fatal"
)

var levelMap = map[string]zapcore.Level{
	DebugLevel: zapcore.DebugLevel,
	InfoLevel:  zapcore.InfoLevel,
	WarnLevel:  zapcore.WarnLevel,
	ErrorLevel: zapcore.ErrorLevel,
	PanicLevel: zapcore.PanicLevel,
	FatalLevel: zapcore.FatalLevel,
}

type lumberjackSink struct {
	*lumberjack.Logger
}

func (lumberjackSink) Sync() error {
	return nil
}

// New ...
func New(params Params) Logger {
	var (
		logLevel    = params.Config.GetString(config.LogLevelKey)
		namespace   = params.Config.GetString(config.NamespaceKey)
		environment = params.Config.GetString(config.EnvironmentKey)

		globalLevel = levelMap[logLevel]
		encoderCfg  = func(environment, timeFormat string) zapcore.EncoderConfig {
			var cfg zapcore.EncoderConfig

			if environment == config.Production {
				cfg = zap.NewProductionEncoderConfig()
			} else {
				cfg = zap.NewDevelopmentEncoderConfig()
			}

			cfg.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(t.Format(timeFormat))
			}
			return cfg
		}(environment, config.DateTimeFormat)
		logFile = config.RootDir() + "/logs/" + namespace + ".log"
		ll      = lumberjack.Logger{
			Filename:   logFile,
			MaxSize:    1024, // MB
			MaxBackups: 30,
			MaxAge:     90, // days
			Compress:   true,
		}
	)

	if err := zap.RegisterSink("lumberjack", func(*url.URL) (zap.Sink, error) {
		return lumberjackSink{
			Logger: &ll,
		}, nil
	}); err != nil {
		panic(err)
	}

	loggerConfig := zap.Config{
		Level:             zap.NewAtomicLevelAt(globalLevel),
		Development:       environment != config.Production,
		Encoding:          "json",
		EncoderConfig:     encoderCfg,
		OutputPaths:       []string{fmt.Sprintf("lumberjack:%s", logFile), "stderr"},
		ErrorOutputPaths:  []string{fmt.Sprintf("lumberjack:%s", logFile), "stderr"},
		DisableStacktrace: true,
	}

	zapLog, err := loggerConfig.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}

	zap.ReplaceGlobals(zapLog)()

	return &zapLogger{
		zap: zapLog.Named(namespace).Sugar(),
	}
}

func (l *zapLogger) Debug(message string, fields ...zap.Field) {
	l.zap.Desugar().Debug(message, fields...)
}

func (l *zapLogger) Info(message string, fields ...zap.Field) {
	l.zap.Desugar().Info(message, fields...)
}

func (l *zapLogger) Warn(message string, fields ...zap.Field) {
	l.zap.Desugar().Warn(message, fields...)
}

func (l *zapLogger) Error(message string, fields ...zap.Field) {
	l.zap.Desugar().Error(message, fields...)
}

func (l *zapLogger) Fatal(message string, fields ...zap.Field) {
	l.zap.Desugar().Fatal(message, fields...)
}

func (l *zapLogger) Debugf(format string, a ...any) {
	l.zap.Debugf(format, a...)
}

func (l *zapLogger) Infof(format string, a ...any) {
	l.zap.Infof(format, a...)
}

func (l *zapLogger) Warnf(format string, a ...any) {
	l.zap.Warnf(format, a...)
}

func (l *zapLogger) Errorf(format string, a ...any) {
	l.zap.Errorf(format, a...)
}

func (l *zapLogger) Fatalf(format string, a ...any) {
	l.zap.Fatalf(format, a...)
}

func (l *zapLogger) Sync() error {
	return l.zap.Sync()
}
