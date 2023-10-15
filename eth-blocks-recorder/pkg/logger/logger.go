package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
)

var Log *zap.SugaredLogger

// Initialize initializes the logger to be used across the application.
func Initialize(logLevel string) error {
	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(parseLogLevel(logLevel)),
		Development: false,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := cfg.Build()
	if err != nil {
		return err
	}

	Log = logger.Sugar()
	return nil
}

// parseLogLevel converts a level string to a Zap log level.
func parseLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "fatal":
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
}

// Error logs error messages with key-value pairs.
func Error(message string, keysAndValues ...interface{}) {
	if Log == nil {
		err := Initialize("error") // default to "error" level if not initialized
		if err != nil {
			panic(err)
		}
	}

	Log.Errorw(message, keysAndValues...)
}

// Sync flushes the logger's buffers and performs cleanup.
// Generally, it's better to defer this function in main() after initializing the logger.
func Sync() {
	if err := Log.Sync(); err != nil {
		log.Printf("Error syncing logger: %v", err)
	}
}
