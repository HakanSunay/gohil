package logger

import (
	"bytes"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestGohilLogFormatterFormat(t *testing.T) {
	tests := []struct {
		entry        *logrus.Entry
		wantedSuffix string
		wantErr      bool
	}{
		{
			entry:   nil,
			wantErr: true,
		},
		{
			entry:   &logrus.Entry{Buffer: nil},
			wantErr: true,
		},
		{
			entry:        &logrus.Entry{Level: logrus.DebugLevel, Message: "Message 1", Buffer: &bytes.Buffer{}},
			wantedSuffix: "DEBUG Message 1\n",
			wantErr:      false,
		},
		{
			entry:        &logrus.Entry{Level: logrus.InfoLevel, Message: "Message 2", Buffer: &bytes.Buffer{}},
			wantedSuffix: "INFO Message 2\n",
			wantErr:      false,
		},
		{
			entry:        &logrus.Entry{Level: logrus.WarnLevel, Message: "Message 3", Buffer: &bytes.Buffer{}},
			wantedSuffix: "WARNING Message 3\n",
			wantErr:      false,
		},
		{
			entry:        &logrus.Entry{Level: logrus.ErrorLevel, Message: "Message 4", Buffer: &bytes.Buffer{}},
			wantedSuffix: "ERROR Message 4\n",
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.wantedSuffix, func(t *testing.T) {
			gf := &gohilLogFormatter{}
			got, err := gf.Format(tt.entry)
			if tt.wantErr && err == nil {
				t.Errorf("Format() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// test if the resulting log line contains the wanted suffix

			stringLog := string(got)
			if !strings.HasSuffix(stringLog, tt.wantedSuffix) {
				t.Errorf("Format() got = %v, wantedSuffix %v", got, tt.wantedSuffix)
			}
		})
	}
}
