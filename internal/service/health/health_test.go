package health

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
	"time"
)

func newTestService() *Service {
	return New(time.Now().Add(-90*time.Second), "9.9.9-test", slog.Default())
}

func TestCheckReportsOKWithUptime(t *testing.T) {
	got := newTestService().Check(context.Background())
	if got.Status != StatusOK {
		t.Fatalf("status=%s want ok", got.Status)
	}
	if got.UptimeText != "1m30s" {
		t.Fatalf("uptime=%q want 1m30s", got.UptimeText)
	}
	if _, err := time.Parse(time.RFC3339, got.CheckedAt); err != nil {
		t.Fatalf("checked_at not RFC3339: %v", err)
	}
}

func TestCheckRespectsCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	got := newTestService().Check(ctx)
	if got.Status != StatusFail {
		t.Fatalf("canceled context must fail closed; got %+v", got)
	}
}

func TestReportJSONShape(t *testing.T) {
	rep := newTestService().Check(context.Background())
	b, _ := json.Marshal(rep)
	for _, key := range []string{"status", "version", "uptime", "checked_at"} {
		if !strings.Contains(string(b), `"`+key+`"`) {
			t.Fatalf("missing key %q in %s", key, b)
		}
	}
}
