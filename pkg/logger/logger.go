package logger

import (
	"os"

	"github.com/IndigoCloud6/go-web-template/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger

// Init initializes the logger
func Init(cfg *config.LoggerConfig) error {
	var level zapcore.Level
	switch cfg.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var encoder zapcore.Encoder
	if cfg.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	var writer zapcore.WriteSyncer
	if cfg.Output == "file" {
		// Set default values if not configured
		maxSize := cfg.MaxSize
		if maxSize == 0 {
			maxSize = 100
		}
		maxBackups := cfg.MaxBackups
		if maxBackups == 0 {
			maxBackups = 3
		}
		maxAge := cfg.MaxAge
		if maxAge == 0 {
			maxAge = 28
		}

		lumberJackLogger := &lumberjack.Logger{
			Filename:   cfg.FilePath,
			MaxSize:    maxSize,    // MB
			MaxBackups: maxBackups,
			MaxAge:     maxAge,     // days
			Compress:   cfg.Compress,
		}
		writer = zapcore.AddSync(lumberJackLogger)
	} else {
		writer = zapcore.AddSync(os.Stdout)
	}

	core := zapcore.NewCore(encoder, writer, level)
	Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return nil
}

// Info logs an info message
func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

// Debug logs a debug message
func Debug(msg string, fields ...zap.Field) {
	Logger.Debug(msg, fields...)
}

// Warn logs a warning message
func Warn(msg string, fields ...zap.Field) {
	Logger.Warn(msg, fields...)
}

// Error logs an error message
func Error(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}

// Fatal logs a fatal message and exits
func Fatal(msg string, fields ...zap.Field) {
	Logger.Fatal(msg, fields...)
}
