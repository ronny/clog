package clog

import (
	"io"
	"log/slog"
)

type HandlerOptions = slog.HandlerOptions

// NewHandler returns a [log/slog.JSONHandler] pre-configured for Google Cloud
// Logging.
//
// Basically it replaces `opts.ReplaceAttr` with [ReplaceAttr].
func NewHandler(w io.Writer, opts *HandlerOptions) *slog.JSONHandler {
	opts.ReplaceAttr = ReplaceAttr
	return slog.NewJSONHandler(w, opts)
}
