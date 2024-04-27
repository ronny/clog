package clog

import "github.com/ronny/clog/trace"

// Standard JSON log fields as per
// https://cloud.google.com/logging/docs/structured-logging#structured_logging_special_fields
const (
	TimestampKey   = "timestamp"
	MessageKey     = "message"
	SeverityKey    = "severity"
	HttpRequestKey = "httpRequest"
)

// Extended JSON log fields as per
// https://cloud.google.com/logging/docs/structured-logging#structured_logging_special_fields
const (
	LabelsKey         = "logging.googleapis.com/labels"
	OperationKey      = "logging.googleapis.com/operation"
	SourceLocationKey = "logging.googleapis.com/sourceLocation"
	TraceKey          = trace.TraceKey
	SpanIDKey         = trace.SpanIDKey
	TraceSampledKey   = trace.TraceSampledKey
)
