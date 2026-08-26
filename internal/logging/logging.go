// Package logging standardizes structured logs through log/slog.
//
// One constructor, no framework: level and format come from configuration,
// every line carries service identity so derived CLIs correlate cleanly in
// aggregators. Handlers stay swappable if a project later needs richer
// telemetry sinks.
package logging

import (
	"io"
	"log/slog"
	"os"
	"strings"
)

// Options shape logger construction. Keep them intentionally tiny.
type Options struct {
	// Service identifies the binary emitting logs.
	Service string
	// Version stamps build metadata on each record when non-empty.
	Version string
	// Level is one of debug, info, warn, error (anything else → info).
	Level string
	// Format is "text" (human) or "json" (machine); anything else → text.
	Format string
	// Out overrides the sink (json→stdout and text→stderr by default).
	Out io.Writer
}

// New builds the process logger. Returns *slog.Logger — never a wrapper
// type — so consumers keep the full standard API surface.
func New(opts Options) *slog.Logger {
	lvl := parseLevel(opts.Level)

	sink := opts.Out
	if sink == nil && strings.EqualFold(opts.Format, "json") {
		sink = os.Stdout // machine-readable records belong on stdout pipelines
	}
	if sink == nil {
		sink = os.Stderr // human-readable diagnostics stay off result streams
	}

	var handler slog.Handler
	hopts := &slog.HandlerOptions{Level: lvl}
	if strings.EqualFold(opts.Format, "json") {
		handler = slog.NewJSONHandler(sink, hopts)
	} else {
		handler = slog.NewTextHandler(sink, hopts)
	}

	log := slog.New(handler).With(
		slog.String("service", opts.Service),
	)
	if opts.Version != "" && opts.Version != "dev" {
		log = log.With(slog.String("version", opts.Version))
	}
	return log
}

// parseLevel maps textual levels defensively; unknown values stay at info so
// misconfiguration never silences an entire application.
func parseLevel(s string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
