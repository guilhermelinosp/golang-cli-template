package cli

import (
	"context"
	"io"
)

// RunFunc is the business logic of a command. It receives the executing
// [Command] (for flag access) and remaining positional arguments.
//
// Returning an error is all a command ever does: the framework owns process
// exit codes, user-facing messages and internal logging (see errors.go).
type RunFunc func(ctx context.Context, c *Command, args []string) error

// ArgsFunc validates positional arguments before Run executes. Prebuilt
// policies: [NoArgs], [ExactArgs], [RangeArgs].
type ArgsFunc func(c *Command, args []string) error

// Command models one CLI command — root, subcommand or nested group —
// independently of any concrete CLI engine.
//
// Declare commands as plain values; build them with literals in wiring code:
//
//	greet := &cli.Command{
//	    Name:  "greet",
//	    Short: "Greet someone",
//	    Flags: []*cli.Flag{cli.StringFlag("name", "n", "world", "who to greet")},
//	    Run: func(ctx context.Context, c *cli.Command, args []string) error {
//	        fmt.Fprintf(out, "hello %s\n", c.String("name"))
//	        return nil
//	    },
//	}
type Command struct {
	// Name is the invocation word, e.g. "health" for `myapp health`.
	Name string
	// Aliases are additional names that resolve to this command.
	Aliases []string
	// Short is the one-line description shown in listings and help headers.
	Short string
	// Long is the extended help body; falls back to Short when empty.
	Long string
	// Example renders an indented examples block in help output.
	Example string
	// Hidden removes the command from help and completion suggestions while
	// keeping it invocable — useful for maintenance commands.
	Hidden bool

	// Args optionally validates positional arguments prior to Run.
	Args ArgsFunc
	// Flags declares the flags accepted by this command only.
	Flags []*Flag
	// Commands declares nested subcommands (groups included).
	Commands []*Command
	// Run holds the business logic; leave nil for pure grouping parents.
	Run RunFunc

	// out resolves per-invocation stdout (bound by App snapshots only).
	out io.Writer
}

// NoArgs rejects any positional argument: `myapp version extra`.
func NoArgs(c *Command, args []string) error {
	if len(args) > 0 {
		return Usagef("unknown argument %q for %q", args[0], c.Name)
	}
	return nil
}

// ExactArgs requires exactly n positional arguments.
func ExactArgs(n int) ArgsFunc {
	return func(c *Command, args []string) error {
		if len(args) != n {
			return Usagef("%q requires exactly %d argument(s), got %d", c.Name, n, len(args))
		}
		return nil
	}
}

// RangeArgs requires between min and max positional arguments.
func RangeArgs(min, max int) ArgsFunc {
	return func(c *Command, args []string) error {
		if len(args) < min || len(args) > max {
			return Usagef("%q accepts %d to %d argument(s), got %d", c.Name, min, max, len(args))
		}
		return nil
	}
}

// UsageLine returns the canonical `myapp command [flags]` hint line used by
// help and error rendering on every engine.
func (c *Command) UsageLine(path string) string {
	line := path + " " + c.Name
	if len(c.Commands) > 0 {
		line += " <command>"
	}
	return line + " [flags]"
}
