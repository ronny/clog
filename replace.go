package clog

import "log/slog"

// ReplaceAttr replaces attr keys and values to fit Google Cloud Logging's
// expected structure as per:
// - https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry
// - https://cloud.google.com/functions/docs/monitoring/logging
// - https://cloud.google.com/functions/docs/monitoring/logging
func ReplaceAttr(groups []string, attr slog.Attr) slog.Attr {
	// Based on https://github.com/remko/cloudrun-slog
	switch attr.Key {
	// This is not strictly necessary, as Cloud Logging supports "time" as well.
	// TODO: make configurable.
	// case slog.TimeKey:
	// 	attr.Key = TimestampKey
	case slog.MessageKey:
		attr.Key = MessageKey
	case slog.SourceKey:
		attr.Key = SourceLocationKey
	case slog.LevelKey:
		level, ok := attr.Value.Any().(slog.Level)
		if ok {
			attr.Key = SeverityKey
			attr.Value = slog.StringValue(LevelName(level))
		}
	}
	return attr
}
