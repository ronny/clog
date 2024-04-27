package clog

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"

	"github.com/ronny/clog/trace"
)

var ErrInvalidHandlerOptions = errors.New("invalid HandlerOptions")

type HandlerOptions struct {
	AddSource       bool
	Level           slog.Leveler
	ReplaceAttr     func(groups []string, a slog.Attr) slog.Attr
	GoogleProjectID string
}

var _ slog.Handler = (*Handler)(nil)

// Handler is a [log/slog.JSONHandler] preconfigured for Google Cloud Logging.
type Handler struct {
	opts    HandlerOptions
	handler slog.Handler
}

func NewHandler(w io.Writer, opts HandlerOptions) (*Handler, error) {
	opts.ReplaceAttr = ReplaceAttr
	if opts.GoogleProjectID == "" {
		return nil, fmt.Errorf("%w: missing GoogleProjectID", ErrInvalidHandlerOptions)
	}
	return &Handler{
		opts: opts,
		handler: slog.NewJSONHandler(w, &slog.HandlerOptions{
			AddSource:   opts.AddSource,
			Level:       opts.Level,
			ReplaceAttr: opts.ReplaceAttr,
		}),
	}, nil
}

// Handle implements [log/slog.Handler].
func (h *Handler) Handle(ctx context.Context, record slog.Record) error {
	return h.handler.Handle(ctx,
		trace.NewRecord(ctx, record, h.opts.GoogleProjectID),
	)
}

// Enabled implements [log/slog.Handler].
func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

// WithAttrs implements [log/slog.Handler].
func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	h.handler = h.handler.WithAttrs(attrs)
	return h
}

// WithGroup implements [log/slog.Handler].
func (h *Handler) WithGroup(name string) slog.Handler {
	h.handler = h.handler.WithGroup(name)
	return h
}
