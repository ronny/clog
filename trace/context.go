package trace

import "context"

type ctxKeyType int

const ctxKey ctxKeyType = iota

// NewContext returns a new context derived from ctx with trace attached.
func NewContext(ctx context.Context, t Trace) context.Context {
	return context.WithValue(ctx, ctxKey, t)
}

// FromContext extracts any [Trace] information from ctx and
// returns it if found, otherwise returns nil.
func FromContext(ctx context.Context) Trace {
	t, ok := ctx.Value(ctxKey).(Trace)
	if !ok {
		return nil
	}
	return t
}
