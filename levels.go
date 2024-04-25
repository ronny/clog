package clog

import (
	"fmt"
	"log/slog"
)

// The levels' name and order match
// https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#LogSeverity
// but the internal int values differ.
//
// The internal int values are changed so that they can be used with
// [log/slog.Log] for logging at a custom level, for example:
//
//	slog.Log(ctx, clog.LevelCritical, "flux capacitor exploded")
const (
	LevelDefault   slog.Level = -5
	LevelDebug     slog.Level = slog.LevelDebug // -4
	LevelInfo      slog.Level = slog.LevelInfo  // 0
	LevelNotice    slog.Level = 2
	LevelWarning   slog.Level = slog.LevelWarn  // 4 ⚠️ Warning vs Warn
	LevelError     slog.Level = slog.LevelError // 8
	LevelCritical  slog.Level = 9
	LevelAlert     slog.Level = 10
	LevelEmergency slog.Level = 11
)

// LevelName returns the Google Cloud Logging name (in uppercase) for the level
// as per
// https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#LogSeverity
//
// If the level is between named values, then an integer is appended to the
// uppercased name.
//
// Examples:
//
//	LevelName(LevelWarning) => "WARNING"
//	LevelName(LevelNotice+1) => "NOTICE+1"
func LevelName(l slog.Level) string {
	str := func(base string, val slog.Level) string {
		if val == 0 {
			return base
		}
		return fmt.Sprintf("%s%+d", base, val)
	}

	switch l {
	case LevelDefault:
		return str("DEFAULT", l-LevelDefault)
	case LevelNotice:
		return str("NOTICE", l-LevelNotice)
	case LevelWarning:
		return str("WARNING", l-LevelWarning) // ⚠️ not "WARN"
	case LevelCritical:
		return str("CRITICAL", l-LevelCritical)
	case LevelAlert:
		return str("ALERT", l-LevelAlert)
	case LevelEmergency:
		return str("EMERGENCY", l-LevelEmergency)
	default:
		return l.String()
	}
}
