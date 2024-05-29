package clog_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/ronny/clog"
	"github.com/stretchr/testify/assert"
)

func TestNewHandler(t *testing.T) {
	t.Parallel()

	buf := &strings.Builder{}

	_, err := clog.NewHandler(buf, clog.HandlerOptions{})
	if !assert.NotNil(t, err) {
		return
	}
	assert.ErrorIs(t, err, clog.ErrInvalidHandlerOptions)
	assert.ErrorContains(t, err, "missing GoogleProjectID")

	handler, err := clog.NewHandler(buf, clog.HandlerOptions{
		AddSource:       true,
		Level:           clog.LevelInfo,
		GoogleProjectID: "my-project-id",
	})
	if !assert.Nil(t, err) {
		return
	}

	logger := slog.New(handler)

	logger.Debug("this is debug")

	logger.Warn("this is warn",
		"foo", "bar",
		"level", 42,
		clog.TraceKey, "fish",
	)

	logger.Log(context.Background(), clog.LevelCritical, "this is critical",
		clog.TraceKey, "fish",
		clog.SpanIDKey, "banana",
	)

	// fmt.Printf("%q\n", buf.String())

	lines := strings.Split(buf.String(), "\n")
	// the log lines are \n terminated, so the last line will always be empty since we split on \n
	lines = lines[:len(lines)-1]

	// 2 because the debug one shouldn't be written, since handler level is set to info
	assert.Equal(t, 2, len(lines))

	var warnEntry Entry
	err = json.Unmarshal([]byte(lines[0]), &warnEntry)
	if err != nil {
		t.Fatal(err)
	}

	ts, hasTime := warnEntry.GetString(slog.TimeKey)
	assert.True(t, hasTime)
	assert.NotEmpty(t, ts)

	sev, hasSev := warnEntry.GetString(clog.SeverityKey)
	assert.True(t, hasSev)
	assert.Equal(t, clog.LevelName(clog.LevelWarning), sev)

	source, hasSource := warnEntry.GetMap(clog.SourceLocationKey)
	assert.True(t, hasSource)
	assert.NotEmpty(t, source["function"])
	assert.NotEmpty(t, source["file"])
	assert.NotEmpty(t, source["line"])

	msg, hasMsg := warnEntry.GetString(clog.MessageKey)
	assert.True(t, hasMsg)
	assert.Equal(t, "this is warn", msg)

	foo, hasFoo := warnEntry.GetString("foo")
	assert.True(t, hasFoo)
	assert.Equal(t, "bar", foo)

	// custom "level" vs severity
	level, hasLevel := warnEntry.GetAny("level")
	assert.True(t, hasLevel)
	assert.Equal(t, float64(42), level) // float64 because json

	trace, hasTrace := warnEntry.GetString(clog.TraceKey)
	assert.True(t, hasTrace)
	assert.Equal(t, "fish", trace)

	var critEntry Entry
	err = json.Unmarshal([]byte(lines[1]), &critEntry)
	if err != nil {
		t.Fatal(err)
	}

	ts, hasTime = critEntry.GetString(slog.TimeKey)
	assert.True(t, hasTime)
	assert.NotEmpty(t, ts)

	sev, hasSev = critEntry.GetString(clog.SeverityKey)
	assert.True(t, hasSev)
	assert.Equal(t, clog.LevelName(clog.LevelCritical), sev)

	source, hasSource = critEntry.GetMap(clog.SourceLocationKey)
	assert.True(t, hasSource)
	assert.NotEmpty(t, source["function"])
	assert.NotEmpty(t, source["file"])
	assert.NotEmpty(t, source["line"])

	msg, hasMsg = critEntry.GetString(clog.MessageKey)
	assert.True(t, hasMsg)
	assert.Equal(t, "this is critical", msg)

	trace, hasTrace = critEntry.GetString(clog.TraceKey)
	assert.True(t, hasTrace)
	assert.Equal(t, "fish", trace)

	spanID, hasSpanID := critEntry.GetString(clog.SpanIDKey)
	assert.True(t, hasSpanID)
	assert.Equal(t, "banana", spanID)
}

var benchmarkJSONHandler_HandleResult error

func BenchmarkJSONHandler_Handle(b *testing.B) {
	buf := &bytes.Buffer{}
	h := slog.NewJSONHandler(buf, &slog.HandlerOptions{Level: slog.LevelInfo})

	ctx := context.Background()

	b.ResetTimer()
	var err error
	for n := 0; n < b.N; n++ {
		record := slog.NewRecord(time.Now(), slog.LevelInfo, "hello", 0)
		record.AddAttrs(
			slog.Attr{Key: "a", Value: slog.StringValue("one")},
			slog.Attr{Key: "a", Value: slog.StringValue("two")},
		)
		err = h.Handle(ctx, record)
	}

	// Always store the result to a package level variable
	// so the compiler cannot eliminate the Benchmark itself.
	benchmarkJSONHandler_HandleResult = err
}

var benchmarkHandler_HandleResult error

func BenchmarkHandler_Handle(b *testing.B) {
	buf := &bytes.Buffer{}
	h, err := clog.NewHandler(buf, clog.HandlerOptions{
		Level:           slog.LevelInfo,
		GoogleProjectID: "example",
	})
	if err != nil {
		b.Fatal(err)
	}

	ctx := context.Background()

	b.ResetTimer()
	var e error
	for n := 0; n < b.N; n++ {
		record := slog.NewRecord(time.Now(), slog.LevelInfo, "hello", 0)
		record.AddAttrs(
			slog.Attr{Key: "a", Value: slog.StringValue("one")},
			slog.Attr{Key: "a", Value: slog.StringValue("two")},
		)
		e = h.Handle(ctx, record)
	}

	// Always store the result to a package level variable
	// so the compiler cannot eliminate the Benchmark itself.
	benchmarkHandler_HandleResult = e
}
