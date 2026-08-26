// Package health is the template's reference domain service. It exists to
// demonstrate the intended layering — commands stay thin, services hold
// logic — and doubles as the default `health` command payload.
package health

import (
	"context"
	"log/slog"
	"time"
)

// Status values reported by checks.
const (
	StatusOK   = "ok"
	StatusFail = "fail"
)

// Service answers liveness questions. Real projects would ping databases,
// queues or downstream APIs here; the template keeps a deterministic check
// so builds and CI never flake.
type Service struct {
	// StartAt anchors process uptime measurement; wire time.Now() in main.
	StartAt time.Time

	// version stamps reports for traceability (injected from internal/build).
	version string

	// logger is an explicitly injected dependency — no locator magic.
	logger *slog.Logger
}

// New wires the service. Explicit constructor injection replaces any need
// for dependency-injection frameworks.
func New(startAt time.Time, version string, logger *slog.Logger) *Service {
	return &Service{StartAt: startAt, version: version, logger: logger}
}

// Report is the machine-readable outcome of a Check.
type Report struct {
	Status     string        `json:"status"`
	Version    string        `json:"version"`
	Uptime     time.Duration `json:"-"`
	UptimeText string        `json:"uptime"`
	CheckedAt  string        `json:"checked_at"`
}

// Check runs all probes under ctx so long-circuit dependencies can be cut
// short on shutdown signals.
func (s *Service) Check(ctx context.Context) Report {
	if err := ctx.Err(); err != nil {
		s.logger.DebugContext(ctx, "health check aborted", slog.String("reason", err.Error()))
		return Report{Status: StatusFail, Version: s.version}
	}
	up := time.Since(s.StartAt)
	report := Report{
		Status:     StatusOK,
		Version:    s.version,
		UptimeText: up.Round(time.Second).String(),
		CheckedAt:  time.Now().UTC().Format(time.RFC3339),
	}
	s.logger.DebugContext(ctx, "health check passed",
		slog.String("status", report.Status),
		slog.Duration("uptime", up),
	)
	return report
}
