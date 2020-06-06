package logger

import (
	"context"
	"os"
	"reflect"
	"testing"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func TestInitGlobalLogging(t *testing.T) {
	InitGlobalLogging()
	if StdLog.Level != defaultLogLevel {
		t.Errorf("expected the default log level to be %v, but got %v",
			defaultLogLevel, StdLog.Level)
	}
}

func TestNewLoggerExpectStdErrLoggerForEmptyFilename(t *testing.T) {
	stdErrLogger := newLogger("", defaultLogLevel, 5, 5)
	if !reflect.DeepEqual(stdErrLogger.Out, os.Stderr) {
		t.Errorf("expected the logger to be redirected to stderr")
	}
}

func TestNewLoggerExpectLumberjackRotatingLogger(t *testing.T) {
	lumberjackLogger := newLogger("dummyfilename", defaultLogLevel, 5, 5)
	if _, ok := lumberjackLogger.Out.(*lumberjack.Logger); !ok {
		t.Errorf("expected the type assertion to be of lumberjack.Logger type")
	}
}

func TestGetFromContextIfLoggerPresentInContext(t *testing.T) {
	ctx := context.Background()
	expectedLogger := &logrus.Logger{}
	logCtx := context.WithValue(ctx, ctxLoggerKey, expectedLogger)
	resultLogger := GetFromContext(logCtx)
	if !reflect.DeepEqual(resultLogger, expectedLogger) {
		t.Errorf("expected %v, but got %v", expectedLogger, resultLogger)
	}
}

func TestGetFromContextExpectStandartLoggerIfLoggerNotInContext(t *testing.T) {
	ctx := context.Background()
	resultLogger := GetFromContext(ctx)
	if !reflect.DeepEqual(resultLogger, StdLog) {
		t.Errorf("expected %v, but got %v", StdLog, resultLogger)
	}
}

func TestPutLoggerInContext(t *testing.T) {
	ctx := context.Background()
	logCtx := PutLoggerInContext(ctx)
	if value, ok := logCtx.Value(ctxLoggerKey).(logrus.FieldLogger); !ok || value != StdLog {
		t.Errorf("unexpected error")
	}
}
