package clog_test

import (
	"fmt"
	"log/slog"
	"testing"

	"github.com/ronny/clog"
	"github.com/stretchr/testify/assert"
)

func TestLevelName(t *testing.T) {
	testCases := []struct {
		desc     string
		level    slog.Level
		expected string
	}{
		{
			desc:     "clog.LevelDefault -> DEFAULT",
			level:    clog.LevelDefault,
			expected: "DEFAULT",
		},
		{
			desc:     "clog.LevelDebug -> DEBUG",
			level:    clog.LevelDebug,
			expected: "DEBUG",
		},
		{
			desc:     "clog.LevelInfo -> INFO",
			level:    clog.LevelInfo,
			expected: "INFO",
		},
		{
			desc:     "clog.LevelNotice -> NOTICE",
			level:    clog.LevelNotice,
			expected: "NOTICE",
		},
		{
			desc:     "clog.LevelWarning -> WARNING",
			level:    clog.LevelWarning,
			expected: "WARNING",
		},
		{
			desc:     "slog.LevelWarn -> WARNING", // ⚠️ Warn vs Warning
			level:    slog.LevelWarn,
			expected: "WARNING",
		},
		{
			desc:     "clog.LevelError -> ERROR",
			level:    clog.LevelError,
			expected: "ERROR",
		},
		{
			desc:     "clog.LevelCritical -> CRITICAL",
			level:    clog.LevelCritical,
			expected: "CRITICAL",
		},
		{
			desc:     "clog.LevelAlert -> ALERT",
			level:    clog.LevelAlert,
			expected: "ALERT",
		},
		{
			desc:     "clog.LevelEmergency -> EMERGENCY",
			level:    clog.LevelEmergency,
			expected: "EMERGENCY",
		},
	}
	for i, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("%d: %s", i, tc.desc), func(t *testing.T) {
			t.Parallel()

			actual := clog.LevelName(tc.level)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
