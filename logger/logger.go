package logger

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	logFilePath              = "/Users/hhalil/Documents/Personal/system-programming/src/gohil/gohil.log"
	defaultLogLevel          = logrus.InfoLevel
	defaultLogFileMaxSizeMB  = 10
	defaultLogFileMaxBackups = 5
	ctxLoggerKey             = 0
)

// StdLog holds the global default logger.
var StdLog = logrus.StandardLogger()

// InitGlobalLogging initializes the default logger.
func InitGlobalLogging() {
	StdLog = newLogger(logFilePath, defaultLogLevel, defaultLogFileMaxSizeMB, defaultLogFileMaxBackups)
}

// newLogger creates a new logger instance.
func newLogger(logFilePath string, level logrus.Level, mb int, backups int) *logrus.Logger {
	logger := logrus.New()

	logger.SetLevel(level)
	logger.SetFormatter(&gohilLogFormatter{})

	if logFilePath == "" {
		logger.SetOutput(os.Stderr)
	} else {
		// Setup rotating log file writer.
		fileRotateLogWriter := lumberjack.Logger{
			Filename:   logFilePath,
			MaxSize:    mb,
			MaxBackups: backups,
			MaxAge:     0, // Rotate based on size only.
			Compress:   true}

		logger.SetOutput(&fileRotateLogWriter)
	}

	return logger
}

// GetFromContext returns a context log entry on the logger instance
// if it is present in context as a value.
func GetFromContext(ctx context.Context) logrus.FieldLogger {
	if fl, ok := ctx.Value(ctxLoggerKey).(logrus.FieldLogger); ok {
		return fl
	}
	return StdLog
}

// PutLoggerInContext adds the logger to the context.
func PutLoggerInContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxLoggerKey, StdLog)
}
