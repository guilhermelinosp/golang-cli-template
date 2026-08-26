package build

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestFallbacksAreSafeWithoutLdflags(t *testing.T) {
	info := Get()
	if info.Version == "" || info.Commit == "" || info.Date == "" || info.GoVersion == "" {
		t.Fatalf("empty metadata would break UX contracts: %+v", info)
	}
}

func TestStringMatchesScriptContract(t *testing.T) {
	Version, Commit, Date = "1.0.0", "abc123", "2026-01-01T00:00:00Z"
	defer func() { Version, Commit, Date = "dev", "none", "unknown" }()

	got := Get().String()
	for _, frag := range []string{"1.0.0", "commit abc123", "built 2026-01-01T00:00:00Z", "go "} {
		if !strings.Contains(got, frag) {
			t.Fatalf("%q missing fragment %q", got, frag)
		}
	}
}

func TestInfoJSONRoundTrip(t *testing.T) {
	Version = "2.3.4"
	data, err := json.Marshal(Get())
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var back Info
	if err := json.Unmarshal(data, &back); err != nil {
		t.Fatalf("unmarshal: %v (%s)", err, data)
	}
	if back.Version != "2.3.4" {
		t.Fatalf("round trip lost version: %+v", back)
	}
}
