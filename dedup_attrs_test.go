package clog_test

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/ronny/clog"
	"github.com/stretchr/testify/assert"
)

func TestHandler_DuplicateAttrs(t *testing.T) {
	timestamp := time.Date(2024, time.May, 29, 12, 34, 56, 0, time.UTC)

	testCases := []struct {
		desc     string
		ctx      context.Context
		record   func() slog.Record
		expected string
	}{
		{
			desc: "without attrs",
			ctx:  context.Background(),
			record: func() slog.Record {
				return slog.NewRecord(timestamp, slog.LevelInfo, "hello", 0)
			},
			expected: `{"time":"2024-05-29T12:34:56Z","severity":"INFO","message":"hello"}`,
		},
		{
			desc: "without duplicate attrs",
			ctx:  context.Background(),
			record: func() slog.Record {
				r := slog.NewRecord(timestamp, slog.LevelInfo, "hello", 0)
				r.AddAttrs(
					slog.Attr{Key: "a", Value: slog.StringValue("one")},
				)
				return r
			},
			expected: `{"time":"2024-05-29T12:34:56Z","severity":"INFO","message":"hello","a":"one"}`,
		},
		{
			desc: "overwrites duplicate attrs",
			ctx:  context.Background(),
			record: func() slog.Record {
				r := slog.NewRecord(timestamp, slog.LevelInfo, "hello", 0)
				r.AddAttrs(
					slog.Attr{Key: "a", Value: slog.StringValue("one")},
					slog.Attr{Key: "a", Value: slog.StringValue("two")},
				)
				return r
			},
			expected: `{"time":"2024-05-29T12:34:56Z","severity":"INFO","message":"hello","a":"two"}`,
		},
	}
	for i, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("%d: %s", i, tc.desc), func(t *testing.T) {
			t.Parallel()

			buf := &bytes.Buffer{}

			handler, err := clog.NewHandler(buf, clog.HandlerOptions{
				Level:           clog.LevelInfo,
				GoogleProjectID: "my-project-id",
			})
			if !assert.Nil(t, err) {
				return
			}

			err = handler.Handle(tc.ctx, tc.record())
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tc.expected, strings.TrimSpace(buf.String()))
		})
	}
}
