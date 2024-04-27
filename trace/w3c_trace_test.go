package trace_test

import (
	"fmt"
	"testing"

	"github.com/ronny/clog/trace"
	"github.com/stretchr/testify/assert"
)

func TestParseW3CTraceParent(t *testing.T) {
	testCases := []struct {
		desc        string
		input       string
		expected    *trace.W3CTraceParent
		expectedErr error
	}{
		{
			desc:        "empty string -> ErrUnparseable",
			input:       "",
			expectedErr: trace.ErrUnparseable,
		},
		{
			desc:  "valid looking -> expected",
			input: "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-00",
			expected: &trace.W3CTraceParent{
				Version:    0,
				TraceID:    "4bf92f3577b34da6a3ce929d0e0e4736",
				ParentID:   "00f067aa0ba902b7",
				TraceFlags: 0,
			},
		},
		{
			desc:        "valid looking but wrong number of bytes -> ErrUnparseable",
			input:       "00-abc-def-01",
			expectedErr: trace.ErrUnparseable,
		},
	}
	for i, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("%d: %s", i, tc.desc), func(t *testing.T) {
			tp, err := trace.ParseW3CTraceParent(tc.input)
			if tc.expectedErr != nil {
				assert.NotNil(t, err)
				assert.ErrorIs(t, err, tc.expectedErr)
				return
			}
			assert.Nil(t, err)
			assert.Equal(t, tc.expected.GetVersion(), tp.GetVersion())
			assert.Equal(t, tc.expected.GetTraceID(), tp.GetTraceID())
			assert.Equal(t, tc.expected.GetParentID(), tp.GetParentID())
			assert.Equal(t, tc.expected.GetSpanID(), tp.GetSpanID())
			assert.Equal(t, tc.expected.GetTraceFlags(), tp.GetTraceFlags())
			assert.Equal(t, tc.expected.Sampled(), tp.Sampled())
		})
	}
}

func TestW3CTraceContext_Sampled(t *testing.T) {
	testCases := []struct {
		desc     string
		trace    *trace.W3CTraceParent
		expected bool
	}{
		{
			desc:     "nil trace -> false",
			trace:    nil,
			expected: false,
		},
		{
			desc: "0x00 -> false",
			trace: &trace.W3CTraceParent{
				TraceFlags: 0,
			},
			expected: false,
		},
		{
			desc: "0x01 -> true",
			trace: &trace.W3CTraceParent{
				TraceFlags: 1,
			},
			expected: true,
		},
		{
			desc: "0xf3 -> true",
			trace: &trace.W3CTraceParent{
				TraceFlags: 0xf3,
			},
			expected: true,
		},
		{
			desc: "0xa0 -> false",
			trace: &trace.W3CTraceParent{
				TraceFlags: 0xa0,
			},
			expected: false,
		},
	}

	for i, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("%d: %s", i, tc.desc), func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.trace.Sampled())
		})
	}
}
