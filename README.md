# clog

A lightweight `log/slog.JSONHandler` wrapper that adapts the JSON log fields
to the Google Cloud Logging structured log format.

The handler merely reformats/renames the structured JSON log fields. It's still `JSONandler` under the hood. It does NOT send logs to Cloud Logging directly (e.g. using the Cloud SDK).

The intended use case is Cloud Run, but it should work in similar environments (e.g. App Engine, Cloud Functions).

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
