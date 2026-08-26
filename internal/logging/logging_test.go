package logging

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
)

func TestParseLevelDefensive(t *testing.T) {
	tests := map[string]slog.Level{
		"DEBUG": slog.LevelDebug,
		"debug": slog.LevelDebug,
		"warn":  slog.LevelWarn,
		"WARN ": slog.LevelWarn,
		"error": slog.LevelError,
		"":      slog.LevelInfo,
		"bogus": slog.LevelInfo, // unknown values must never silence the app
	}
	for in, want := range tests {
		if got := parseLevel(in); got != want {
			t.Fatalf("parseLevel(%q) = %v want %v", in, got, want)
		}
	}
}

func TestJSONFormatCarriesServiceAndVersion(t *testing.T) {
	var buf bytes.Buffer
	logger := New(Options{Service: "svc", Version: "2.0.0", Level: "info", Format: "json", Out: &buf})
	logger.Info("hello", "k", "v")

	var rec map[string]any
	if err := json.Unmarshal(buf.Bytes(), &rec); err != nil {
		t.Fatalf("not valid JSON log: %v (%q)", err, buf.String())
	}
	for _, key := range []string{"service", "version", "msg", "k"} {
		if _, ok := rec[key]; !ok {
			t.Fatalf("missing key %q in %v", key, rec)
		}
	}
}

func TestTextFormatOmitsDevVersionNoise(t *testing.T) {
	var buf bytes.Buffer
	logger := New(Options{Service: "svc", Version: "dev", Format: "text", Level: "debug", Out: &buf})
	logger.Debug("x")
	if strings.Contains(buf.String(), "version=dev") {
		t.Fatalf("dev noise leaked into logs: %q", buf.String())
	}
}
