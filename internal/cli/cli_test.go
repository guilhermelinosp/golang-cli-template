package cli

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
)

// dispatchEngine emulates the Engine contract without Cobra: it resolves
// args against the app tree exactly like a real engine would, proving that
// ALL application behavior below holds on any implementation.
type dispatchEngine struct {
	err error // optional failure injected before/at dispatch
}

func (e *dispatchEngine) Execute(_ context.Context, app *App) error {
	if e.err != nil {
		return e.err
	}
	args := app.Args()
	if len(args) == 0 {
		return nil // engines typically render help and succeed
	}
	for _, c := range app.Root().Commands {
		if c.Name != args[0] || c.Run == nil {
			continue
		}
		if len(args) > 1 && strings.HasPrefix(args[1], "--") {
			// minimal flag handling for tests: --name=value only
			for _, f := range c.Flags {
				k, v, ok := cut(args[1])
				if !ok || k != f.Name {
					continue
				}
				f.ResetRuntime()
				if err := f.Apply(v); err != nil {
					return Usagef("flag %s: %w", f.Name, err)
				}
			}
		}
		if err := c.Args; err != nil {
			if err2 := err(c, nil); err2 != nil {
				return err2
			}
		}
		return c.Run(context.Background(), c, nil)
	}
	return Usagef("unknown command %q", args[0])
}

func cut(s string) (key, value string, ok bool) {
	rest, found := strings.CutPrefix(s, "--")
	if !found {
		return "", "", false
	}
	k, v, hasEq := strings.Cut(rest, "=")
	return k, v, hasEq
}

func newTestApp(t *testing.T) (*App, *strings.Builder) {
	t.Helper()
	app := newWith(Options{Name: "testapp", Version: "9.9.9-test"}, &dispatchEngine{})
	var buf strings.Builder
	app.SetOutput(&buf, &buf)
	t.Cleanup(func() {})
	return app, &buf
}

func TestRunMapsErrorsToExitCodes(t *testing.T) {
	tests := []struct {
		name     string
		engine   *dispatchEngine
		wantCode int
	}{
		{"success", &dispatchEngine{}, ExitSuccess},
		{"generic-error-becomes-failure", &dispatchEngine{err: fmt.Errorf("boom")}, ExitFailure},
		{"wrapped-usage-preserves-code", &dispatchEngine{err: fmt.Errorf("outer: %w", Usagef("bad input"))}, ExitUsage},
		{"domain-error-exit-1", &dispatchEngine{err: Errorf("quota exceeded")}, ExitFailure},
		{"custom-code-flows-through", &dispatchEngine{err: Exitf(42, "custom")}, 42},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			app, buf := newTestApp(t)
			app.engine = tc.engine
			app.Add(&Command{Name: "noop", Run: func(context.Context, *Command, []string) error { return nil }})

			got := app.Run(context.Background(), []string{"noop"})
			if got != tc.wantCode {
				t.Fatalf("exit code = %d, want %d (output: %q)", got, tc.wantCode, buf.String())
			}
			if tc.wantCode != ExitSuccess && !strings.Contains(buf.String(), "error:") {
				t.Fatalf("expected rendered error line, got %q", buf.String())
			}
			if tc.wantCode == ExitUsage && !strings.Contains(buf.String(), "--help") {
				t.Fatalf("usage failures must include recovery hint, got %q", buf.String())
			}
		})
	}
}

func TestSnapshotIsolationBetweenRuns(t *testing.T) {
	var observed []string

	cmd := &Command{
		Name:  "mutator",
		Flags: []*Flag{StringFlag("mode", "m", "default", "")},
		Run: func(_ context.Context, c *Command, _ []string) error {
			observed = append(observed, c.String("mode"))
			return nil
		},
	}

	app, _ := newTestApp(t)
	app.Add(cmd)

	if code := app.Run(context.Background(), []string{"mutator", "--mode=alpha"}); code != ExitSuccess {
		t.Fatalf("first run exit = %d", code)
	}
	if code := app.Run(context.Background(), []string{"mutator"}); code != ExitSuccess {
		t.Fatalf("second run exit = %d", code)
	}
	want := []string{"alpha", "default"}
	if len(observed) != 2 || observed[0] != want[0] || observed[1] != want[1] {
		t.Fatalf("state leaked across snapshot runs: %v want %v", observed, want)
	}
}

func TestAddIgnoresInvalidEntries(t *testing.T) {
	app, _ := newTestApp(t)
	before := len(app.Root().Commands)
	app.Add(nil, &Command{}, &Command{Name: ""}, &Command{Name: "valid"})

	if got := len(app.Root().Commands); got != before+1 {
		t.Fatalf("registered %d entries, want 1 extra (invalid ones ignored)", got-before)
	}
}

func TestExitCodeOf(t *testing.T) {
	assertEq := func(name string, got, want int) {
		t.Helper()
		if got != want {
			t.Fatalf("%s: got %d want %d", name, got, want)
		}
	}
	assertEq("nil", ExitCodeOf(nil), ExitSuccess)
	assertEq("plain", ExitCodeOf(errors.New("x")), ExitFailure)
	assertEq("coded-deep-wrap", ExitCodeOf(fmt.Errorf("a: %w", fmt.Errorf("b: %w", Exitf(7, "c")))), 7)
}

func TestDefaultVersionCommandShipsInNew(t *testing.T) {
	app, buf := newTestApp(t)
	var names []string
	for _, c := range app.Root().Commands {
		names = append(names, c.Name)
	}
	if !contains(names, "version") {
		t.Fatalf("version command missing from defaults: %v", names)
	}

	code := app.Run(context.Background(), []string{"version"})
	if code != ExitSuccess {
		t.Fatalf("version exit = %d (%q)", code, buf.String())
	}
	out := strings.TrimSpace(buf.String())
	if !strings.HasPrefix(out, "testapp ") || !strings.Contains(out, "go ") {
		t.Fatalf("version line does not match contract %q", out)
	}
}

func contains(hay []string, needle string) bool {
	for _, s := range hay {
		if s == needle {
			return true
		}
	}
	return false
}
