package trace

import (
	"fmt"
	"regexp"
)

var _ Trace = (*CloudTraceContext)(nil)

// TRACE_ID/SPAN_ID();o=OPTIONS
// Example: 70e0091f6f5d4643bb4eca9d81320c76/97123319527522;o=1
var cloudTraceContextRegex = regexp.MustCompile(`^(?P<TraceID>[0-9a-fA-F]{32})/(?P<SpanID>[0-9]+);o=(?P<Options>.+)$`)

// Deprecated: use [ParseW3CTraceParent] instead.
func ParseCloudTraceContext(s string) (*CloudTraceContext, error) {
	match := cloudTraceContextRegex.FindStringSubmatch(s)
	result := make(map[string]string)
	for i, name := range cloudTraceContextRegex.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}
	t := &CloudTraceContext{
		TraceID: result["TraceID"],
		SpanID:  result["SpanID"],
		Options: result["Options"],
	}
	err := t.Validate()
	if err != nil {
		return nil, err
	}
	return t, nil
}

// Based on https://cloud.google.com/trace/docs/trace-context#http-requests.
// See also https://cloud.google.com/run/docs/trace for Cloud Run specific info.
//
// Deprecated: use [W3CTraceParent] instead.
type CloudTraceContext struct {
	TraceID string
	SpanID  string
	Options string
}

func (t *CloudTraceContext) Validate() error {
	if t == nil {
		return fmt.Errorf("%w: nil CloudTraceContext", ErrInvalidTrace)
	}
	if len(t.TraceID) != 32 {
		return fmt.Errorf("%w: invalid TraceID", ErrInvalidTrace)
	}
	if t.SpanID == "" {
		return fmt.Errorf("%w: invalid SpanID", ErrInvalidTrace)
	}
	if t.Options == "" {
		return fmt.Errorf("%w: invalid Options", ErrInvalidTrace)
	}
	return nil
}

func (t *CloudTraceContext) GetTraceID() string {
	if t == nil {
		return ""
	}
	return t.TraceID
}

func (t *CloudTraceContext) GetSpanID() string {
	if t == nil {
		return ""
	}
	return t.SpanID
}

func (t *CloudTraceContext) GetOptions() string {
	if t == nil {
		return ""
	}
	return t.Options
}

func (t *CloudTraceContext) Sampled() bool {
	if t == nil {
		return false
	}
	return t.Options == "1"
}
