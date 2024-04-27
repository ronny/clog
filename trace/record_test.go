package trace_test

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/ronny/clog/trace"
	"github.com/stretchr/testify/assert"
)

func TestNewRecord(t *testing.T) {
	now := time.Now()

	record := slog.NewRecord(
		now,
		slog.LevelInfo,
		"hello fellow kids",
		0,
	)

	tr, err := trace.ParseW3CTraceParent(
		"00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
	)
	if err != nil {
		t.Fatal(err)
	}

	ctx := trace.NewContext(context.Background(), tr)

	projectID := "your-project-id"

	returnedRecord := trace.NewRecord(ctx, record, projectID)

	hasTrace := false
	hasSpan := false
	hasSampled := false

	returnedRecord.Attrs(func(attr slog.Attr) bool {
		switch attr.Key {
		case trace.TraceKey:
			hasTrace = true
			assert.Equal(t, slog.KindString, attr.Value.Kind())
			assert.Equal(t,
				"projects/your-project-id/traces/4bf92f3577b34da6a3ce929d0e0e4736",
				attr.Value.String(),
			)
		case trace.SpanIDKey:
			hasSpan = true
			assert.Equal(t, slog.KindString, attr.Value.Kind())
			assert.Equal(t, "00f067aa0ba902b7", attr.Value.String())
		case trace.TraceSampledKey:
			hasSampled = true
			assert.Equal(t, slog.KindBool, attr.Value.Kind())
			assert.Equal(t, true, attr.Value.Bool())
		}
		return true
	})

	assert.True(t, hasTrace)
	assert.True(t, hasSpan)
	assert.True(t, hasSampled)
}

var benchmarkResult slog.Record

func BenchmarkNewRecord(b *testing.B) {
	now := time.Now()

	record := slog.NewRecord(
		now,
		slog.LevelInfo,
		"hello fellow kids",
		0,
	)

	tr, err := trace.ParseW3CTraceParent(
		"00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
	)
	if err != nil {
		b.Fatal(err)
	}

	ctx := trace.NewContext(context.Background(), tr)

	projectID := "your-project-id"

	b.ResetTimer()

	var r slog.Record
	for n := 0; n < b.N; n++ {
		// Always store the result to avoid the compiler eliminating the function
		// call.
		r = trace.NewRecord(ctx, record, projectID)
	}

	// Always store the result to a package level variable
	// so the compiler cannot eliminate the Benchmark itself.
	benchmarkResult = r
}
