package cli

import (
	"fmt"
	"strings"
	"testing"
)

func TestRegisterEngineDuplicatePanics(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil || !strings.Contains(fmt.Sprint(r), "registered twice") {
			t.Fatalf("expected duplicate-registration panic, got %v", r)
		}
	}()
	RegisterEngine("dup-test", func() Engine { return nil })
	RegisterEngine("dup-test", func() Engine { return nil })
}

func TestRegisterEngineInvalidPanics(t *testing.T) {
	for _, tc := range []struct {
		name   string
		f      EngineFactory
		expect string
	}{
		{"", func() Engine { return nil }, "invalid engine"},
		{"blank-factory", nil, "invalid engine"},
	} {
		func() {
			defer func() {
				r := recover()
				if r == nil || !strings.Contains(fmt.Sprint(r), tc.expect) {
					t.Fatalf("%q: expected panic %q got %v", tc.name, tc.expect, r)
				}
			}()
			RegisterEngine(tc.name, tc.f)
		}()
	}
}

func TestUnknownEngineFailsFast(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil || !strings.Contains(fmt.Sprint(r), `no engine registered`) {
			t.Fatalf("expected missing-engine panic, got %v", r)
		}
	}()
	saved := registry["cobra"]
	delete(registry, "cobra")
	defer func() { registry["cobra"] = saved }()
	defaultEngine()
}
