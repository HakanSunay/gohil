package logger

import (
	"errors"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

const rfc3339Milli = "2006-01-02T15:04:05.999Z07:00"

// gohilLogFormatter represents a log line formatter.
type gohilLogFormatter struct{}

// Format formats log entry to form a log line.
// Line format <UTC RFC3339 formatted time> <LEVEL> <msg>
func (gf *gohilLogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	if entry == nil {
		return nil, errors.New("logger entry is nil")
	}

	if entry.Buffer == nil {
		return nil, errors.New("logger empty entry buffer")
	}

	time := entry.Time.UTC().Format(rfc3339Milli)
	level := strings.ToUpper(entry.Level.String())

	logLineData := entry.Buffer
	logLineData.WriteString(
		fmt.Sprintf("%s %s %s\n", time, level, entry.Message))

	return logLineData.Bytes(), nil
}
