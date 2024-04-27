package trace

type Trace interface {
	GetTraceID() string
	GetSpanID() string
	Sampled() bool
}
