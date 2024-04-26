# clog

[![Go Reference](https://pkg.go.dev/badge/github.com/ronny/clog.svg)](https://pkg.go.dev/github.com/ronny/clog)

A lightweight [`log/slog.JSONHandler`](https://pkg.go.dev/log/slog#JSONHandler)
wrapper that adapts the fields to the [Google Cloud Logging structured log
format](https://cloud.google.com/logging/docs/structured-logging#structured_logging_special_fields).

The handler merely reformats/renames the structured JSON log fields. It's
still `JSONHandler` under the hood. It does NOT send logs to Cloud Logging
directly (e.g. using the Cloud SDK).

The intended use case is Cloud Run, but it should work in similar environments
where logs are emitted to stdout/stderr and automatically picked up by Cloud
Logging (e.g. App Engine, Cloud Functions, GKE).

## Usage

```go
import (
	"log/slog"

	"github.com/ronny/clog"
)

func main() {
	logger := slog.New(
		clog.NewHandler(os.Stderr, clog.HandlerOptions{
			Level: clog.LevelInfo,
		}),
	)
	logger.Warn("flux capacitor is too warm", "tempCelsius", 42)
	logger.Log(ctx, clog.LevelCritical, "flux capacitor is on fire")
}
```

## Credits and acknowledgements

Thank you to Remko Tronçon for doing all the hard work in
https://github.com/remko/cloudrun-slog which is the basis for this library.
