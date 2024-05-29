package clog

import "log/slog"

// Returns a new slog.Record with its attrs having the same key deduped.
// Keys are case sensitive.
func dedupAttrs(record slog.Record) slog.Record {
	attrMap := map[string]slog.Attr{}

	record.Attrs(func(a slog.Attr) bool {
		attrMap[a.Key] = a
		return true
	})

	if len(attrMap) == 0 {
		return record
	}

	r := slog.NewRecord(record.Time, record.Level, record.Message, record.PC)

	for _, attr := range attrMap {
		// We could have used an extra array to call AddAttrs() once only, but we
		// want to avoid further allocs.
		r.AddAttrs(attr)
	}

	return r
}
