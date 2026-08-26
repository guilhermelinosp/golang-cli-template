package cli

import (
	"fmt"
	"strconv"
	"time"
)

// FlagType enumerates the supported flag value kinds.
type FlagType int

// Supported flag types. Adding a new type touches only this file plus the
// engine adapter; do not leak engine-specific builders into the API.
const (
	FlagString FlagType = iota
	FlagBool
	FlagInt
	FlagFloat64
	FlagDuration
)

// Flag describes a command-line flag declaratively. Declare flags in
// [Command.Flags]; read their parsed values inside [Command.Run] through the
// typed accessors on [Command] (e.g. [Command.String], [Command.Bool]).
type Flag struct {
	// Name is the long form, used as --name. Required.
	Name string
	// Shorthand is an optional one-letter alias, used as -n.
	Shorthand string
	// Usage is the help text shown next to the flag.
	Usage string
	// Default is the value applied when the flag is absent.
	Default any
	// Type discriminates how Default and input are interpreted.
	Type FlagType
	// Hidden removes the flag from help output (it stays functional).
	Hidden bool

	value   any // last parsed value
	changed bool
}

// String renders a compact help-oriented representation of the flag.
func (f *Flag) String() string {
	if f == nil {
		return ""
	}
	sh := " "
	if f.Shorthand != "" {
		sh = f.Shorthand
	}
	return fmt.Sprintf("-%s, --%s", sh, f.Name)
}

// StringFlag declares a --name string flag.
func StringFlag(name, shorthand, def, usage string) *Flag {
	return &Flag{Name: name, Shorthand: shorthand, Usage: usage, Default: def, Type: FlagString}
}

// BoolFlag declares a boolean switch usable without an explicit value:
//
//	myapp health --json
func BoolFlag(name, shorthand string, def bool, usage string) *Flag {
	return &Flag{Name: name, Shorthand: shorthand, Usage: usage, Default: def, Type: FlagBool}
}

// IntFlag declares a base-10 integer flag (0x/0b/0o prefixes accepted).
func IntFlag(name, shorthand string, def int, usage string) *Flag {
	return &Flag{Name: name, Shorthand: shorthand, Usage: usage, Default: def, Type: FlagInt}
}

// Float64Flag declares a 64-bit floating point flag.
func Float64Flag(name, shorthand string, def float64, usage string) *Flag {
	return &Flag{Name: name, Shorthand: shorthand, Usage: usage, Default: def, Type: FlagFloat64}
}

// DurationFlag declares a time.Duration flag accepting "300ms", "2h45m", etc.
func DurationFlag(name, shorthand string, def time.Duration, usage string) *Flag {
	return &Flag{Name: name, Shorthand: shorthand, Usage: usage, Default: def, Type: FlagDuration}
}

// Effective returns the runtime value honoring explicit sets over defaults.
// Engines call it to render help annotations; application code should prefer
// the typed accessors on [Command].
func (f *Flag) Effective() any {
	if f.changed && f.value != nil {
		return f.value
	}
	return f.Default
}

// Apply converts raw textual input into the stored typed value. Conversion
// lives here — not in the engine — so parse errors stay uniform regardless
// of which implementation powers the CLI.
func (f *Flag) Apply(raw string) error {
	switch f.Type {
	case FlagString:
		f.value = raw
	case FlagBool:
		if raw == "" {
			raw = "true"
		}
		b, err := strconv.ParseBool(raw)
		if err != nil {
			return fmt.Errorf("invalid boolean %q for --%s", raw, f.Name)
		}
		f.value = b
	case FlagInt:
		n, err := strconv.ParseInt(raw, 0, strconv.IntSize)
		if err != nil {
			return fmt.Errorf("invalid integer %q for --%s", raw, f.Name)
		}
		f.value = int(n)
	case FlagFloat64:
		fl, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return fmt.Errorf("invalid float %q for --%s", raw, f.Name)
		}
		f.value = fl
	case FlagDuration:
		d, err := time.ParseDuration(raw)
		if err != nil {
			return fmt.Errorf("invalid duration %q for --%s", raw, f.Name)
		}
		f.value = d
	default:
		return fmt.Errorf("unsupported flag type for --%s", f.Name)
	}
	f.changed = true
	return nil
}

// --- typed accessors for Run implementations -------------------------------

func (c *Command) lookup(name string, want FlagType) *Flag {
	for _, f := range c.Flags {
		if f.Name == name {
			if f.Type != want {
				panic(fmt.Sprintf("cli: flag %q of %q has unexpected type (want %v)", name, c.Name, want))
			}
			return f
		}
	}
	panic(fmt.Sprintf("cli: %q does not declare a flag named %q (add it to Command.Flags)", c.Name, name))
}

// String returns the value of the string flag name.
func (c *Command) String(name string) string { return c.lookup(name, FlagString).Effective().(string) }

// Bool returns the value of the boolean flag name.
func (c *Command) Bool(name string) bool { return c.lookup(name, FlagBool).Effective().(bool) }

// Int returns the value of the int flag name.
func (c *Command) Int(name string) int { return c.lookup(name, FlagInt).Effective().(int) }

// Float64 returns the value of the float64 flag name.
func (c *Command) Float64(name string) float64 {
	return c.lookup(name, FlagFloat64).Effective().(float64)
}

// Duration returns the value of the duration flag name.
func (c *Command) Duration(name string) time.Duration {
	return c.lookup(name, FlagDuration).Effective().(time.Duration)
}

// Changed reports whether flag name was explicitly provided on the command
// line rather than falling back to its default.
func (c *Command) Changed(name string) bool {
	for _, f := range c.Flags {
		if f.Name == name {
			return f.changed
		}
	}
	return false
}

// ResetRuntime clears parsed state so a Flag struct can serve another
// invocation from pristine defaults. Engine adapters call it at build time.
func (f *Flag) ResetRuntime() {
	f.value = nil
	f.changed = false
}
