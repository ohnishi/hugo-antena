package env_test

import (
	"fmt"
	"testing"

	"go.uber.org/zap/zapcore"

	"github.com/ohnishi/antena/backend/common/env"
)

func TestExportLogLevelFromText(t *testing.T) {
	cases := []struct {
		logLevelText     string
		expectedLogLevel zapcore.Level
		expectedFound    bool
	}{
		{
			logLevelText:  "xxx",
			expectedFound: false,
		},

		{
			logLevelText:  "",
			expectedFound: false,
		},

		{
			logLevelText:     "info",
			expectedLogLevel: zapcore.InfoLevel,
			expectedFound:    true,
		},

		{
			logLevelText:     "WARN",
			expectedLogLevel: zapcore.WarnLevel,
			expectedFound:    true,
		},

		{
			logLevelText:     "error",
			expectedLogLevel: zapcore.ErrorLevel,
			expectedFound:    true,
		},
	}

	for idx, eachCase := range cases {
		t.Run(fmt.Sprintf("case %d", idx), func(t *testing.T) {
			actualLogLevel, found := env.ExportLogLevelFromText(eachCase.logLevelText)
			if eachCase.expectedFound {
				if !found {
					t.Error("`found` must be true")
				}
				if eachCase.expectedLogLevel != *actualLogLevel {
					t.Errorf("invalid log level. expected %+v but actual was %+v",
						eachCase.expectedLogLevel, actualLogLevel)
				}
			} else {
				if found {
					t.Error("`found` must be false")
				}
				if actualLogLevel != nil {
					t.Error("log level must be nil")
				}
			}
		})
	}

}
