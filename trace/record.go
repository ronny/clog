package trace

import (
	"context"
	"log/slog"
	"strings"
)

const (
	TraceKey        = "logging.googleapis.com/trace"
	SpanIDKey       = "logging.googleapis.com/spanId"
	TraceSampledKey = "logging.googleapis.com/traceSampled"
)

// NewRecord extracts trace information from ctx, then returns a clone of record
// with the trace attrs as per
// https://cloud.google.com/logging/docs/structured-logging#structured_logging_special_fields
// added to it.
func NewRecord(ctx context.Context, record slog.Record, projectID string) slog.Record {
	t := FromContext(ctx)
	if t == nil {
		return record
	}

	record = record.Clone()

	traceValue := strings.Builder{}
	traceValue.WriteString("projects/")
	traceValue.WriteString(projectID)
	traceValue.WriteString("/traces/")
	traceValue.WriteString(t.GetTraceID())

	record.Add(TraceKey, slog.StringValue(traceValue.String()))
	record.Add(SpanIDKey, slog.StringValue(t.GetSpanID()))
	record.Add(TraceSampledKey, slog.BoolValue(t.Sampled()))

	return record
}
