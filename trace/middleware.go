package trace

import "net/http"

func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// https://cloud.google.com/trace/docs/trace-context#http-requests
		var t Trace

		traceparent := r.Header.Get("traceparent")
		if traceparent != "" {
			// Intentionally ignoring err since there's nothing useful we can do it.
			t, _ = ParseW3CTraceParent(traceparent)
		}

		if t == nil {
			cloudTraceContext := r.Header.Get("X-Cloud-Trace-Context")
			if cloudTraceContext != "" {
				// Intentionally ignoring err since there's nothing useful we can do it.
				t, _ = ParseCloudTraceContext(cloudTraceContext)
			}
		}

		if t != nil {
			r = r.WithContext(NewContext(r.Context(), t))
		}

		h.ServeHTTP(w, r)
	})
}
