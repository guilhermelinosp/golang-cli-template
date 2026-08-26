package config

import (
	"testing"
)

func TestLoadDefaultsWhenEnvEmpty(t *testing.T) {
	got := Load("NOPE")
	if got != Defaults() {
		t.Fatalf("got %+v want defaults %+v", got, Defaults())
	}
}

func TestLoadReadsPrefixedEnv(t *testing.T) {
	t.Setenv("MYAPP_LOG_LEVEL", "debug")
	t.Setenv("MYAPP_LOG_FORMAT", "JSON")

	got := Load("MYAPP")
	if got.LogLevel != "debug" || got.LogFormat != "json" {
		t.Fatalf("env overrides ignored: %+v", got)
	}
}

func TestEnvKeyNormalizesPrefix(t *testing.T) {
	cases := map[string]string{
		"my-app":   "MY_APP_LOG_LEVEL",
		"myapp":    "MYAPP_LOG_LEVEL",
		"_my_app_": "MY_APP_LOG_LEVEL",
	}
	for prefix, want := range cases {
		if got := envKey(prefix, "LOG_LEVEL"); got != want {
			t.Fatalf("prefix %q: got %s want %s", prefix, got, want)
		}
	}
}
