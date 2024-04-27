# clog

[![Go Reference](https://pkg.go.dev/badge/github.com/ronny/clog.svg)](https://pkg.go.dev/github.com/ronny/clog)

A lightweight [`log/slog.JSONHandler`](https://pkg.go.dev/log/slog#JSONHandler)
wrapper that adapts the fields to the [Google Cloud Logging structured log
format](https://cloud.google.com/logging/docs/structured-logging#structured_logging_special_fields).

The intended use case is Cloud Run, but it should work in similar environments
where logs are emitted to stdout/stderr and automatically picked up by Cloud
Logging (e.g. App Engine, Cloud Functions, GKE).

## Features

- Lightweight. The handler merely reformats/renames the structured JSON log
  fields. It's still [`log/slog.JSONHandler`](https://pkg.go.dev/log/slog#JSONHandler)
  under the hood. It does NOT send logs to Cloud Logging directly (e.g. using
  the Cloud SDK).

- Tracing. A tracing middleware is provided to automatically extract tracing
  information from `traceparent` or `X-Cloud-Trace-Context` HTTP request header,
  and attaches it to the request context. The Handler automatically includes any
  tracing information as log attributes.

- Custom levels as supported by Google Cloud Logging, e.g. CRITICAL and NOTICE.

## Usage

```go
import (
	"log/slog"

	"cloud.google.com/go/compute/metadata"
	"github.com/ronny/clog"
)

func main() {
	projectID, err := metadata.ProjectID()
	if err != nil {
		panic(err)
	}

	logger := slog.New(
		clog.NewHandler(os.Stderr, clog.HandlerOptions{
			Level: clog.LevelInfo,
			GoogleProjectID: projectID,
		}),
	)
	slog.SetDefault(logger)

	mux := http.NewServeMux()
	mux.Handle("POST /warn", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// ⚠️ The log will have tracing attrs since we're using
		// trace.Middleware below (assuming trace header is in the request):
		// "logging.googleapis.com/trace"
		// "logging.googleapis.com/spanId"
		// "logging.googleapis.com/traceSampled"
		slog.WarnContext(ctx, "flux capacitor is too warm",
			"mycount", 42,
		)
	}))
	mux.Handle("POST /critical", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// ⚠️ Custom level CRITICAL
		slog.Log(ctx, clog.LevelCritical, "flux capacitor is on fire")
	}))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("listening on port %s", port)

	// ⚠️ `trace.Middleware` to make tracing information available in ctx in mux
	// handlers.
	if err := http.ListenAndServe(":"+port, trace.Middleware(mux)); err != nil {
		log.Fatal(err)
	}
}
```

## Credits and acknowledgements

Thank you to Remko Tronçon for doing most of the hard work in
https://github.com/remko/cloudrun-slog which is the basis for this library.
