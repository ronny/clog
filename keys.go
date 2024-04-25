package clog

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
	SpanIDKey         = "logging.googleapis.com/spanId"
	TraceKey          = "logging.googleapis.com/trace"
	TraceSampledKey   = "logging.googleapis.com/traceSampled"
)
