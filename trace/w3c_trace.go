package trace

import (
	"encoding/hex"
	"fmt"
	"strings"
)

var _ Trace = (*W3CTraceParent)(nil)

func ParseW3CTraceParent(s string) (*W3CTraceParent, error) {
	parts := strings.Split(s, "-")
	if len(parts) != 4 {
		return nil, fmt.Errorf("%w: incorrect number of parts %d, expected 4", ErrUnparseable, len(parts))
	}

	versionBytes, err := hex.DecodeString(parts[0])
	if err != nil {
		return nil, fmt.Errorf("%w: Version: hex.DecodeString: %w", ErrUnparseable, err)
	}
	if len(versionBytes) != 1 {
		return nil, fmt.Errorf("%w: Version: expected 1 byte, got %d", ErrUnparseable, len(versionBytes))
	}

	traceIDBytes, err := hex.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("%w: TraceID: hex.DecodeString: %w", ErrUnparseable, err)
	}
	if len(traceIDBytes) != 16 {
		return nil, fmt.Errorf("%w: TraceID: expected 16 bytes, got %d", ErrUnparseable, len(traceIDBytes))
	}

	parentIDBytes, err := hex.DecodeString(parts[2])
	if err != nil {
		return nil, fmt.Errorf("%w: ParentID: hex.DecodeString: %w", ErrUnparseable, err)
	}
	if len(parentIDBytes) != 8 {
		return nil, fmt.Errorf("%w: ParentID: expected 8 bytes, got %d", ErrUnparseable, len(parentIDBytes))
	}

	flagsBytes, err := hex.DecodeString(parts[3])
	if err != nil {
		return nil, fmt.Errorf("%w: TraceFlags: hex.DecodeString: %w", ErrUnparseable, err)
	}
	if len(flagsBytes) != 1 {
		return nil, fmt.Errorf("%w: TraceFlags: expected 1 byte, got %d", ErrUnparseable, len(flagsBytes))
	}

	return &W3CTraceParent{
		Version:    versionBytes[0],
		TraceID:    parts[1],
		ParentID:   parts[2],
		TraceFlags: flagsBytes[0],
	}, nil
}

// Based on https://www.w3.org/TR/trace-context/#traceparent-header
type W3CTraceParent struct {
	Version    byte
	TraceID    string
	ParentID   string
	TraceFlags byte
}

func (t *W3CTraceParent) GetVersion() byte {
	if t == nil {
		return 0
	}
	return t.Version
}

func (t *W3CTraceParent) GetTraceID() string {
	if t == nil {
		return ""
	}
	return t.TraceID
}

func (t *W3CTraceParent) GetParentID() string {
	if t == nil {
		return ""
	}
	return t.ParentID
}

func (t *W3CTraceParent) GetSpanID() string {
	return t.GetParentID()
}

func (t *W3CTraceParent) GetTraceFlags() byte {
	if t == nil {
		return 0
	}
	return t.TraceFlags
}

// https://www.w3.org/TR/trace-context/#sampled-flag
const W3CFlagSampled byte = 1

func (t *W3CTraceParent) Sampled() bool {
	if t == nil {
		return false
	}

	// Need to mask because TraceFlags is a bit field, we need to check
	// whether the least significant bit (the rightmost one) is on or not,
	// regardless of the other bits. Another way to check is to test if the
	// number is odd or not, but the mask way works with any bit position.
	//
	// https://www.w3.org/TR/trace-context/#trace-flags
	return t.TraceFlags&W3CFlagSampled == W3CFlagSampled
}
