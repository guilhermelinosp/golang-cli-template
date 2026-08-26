package cli

// runSnapshot produces an isolated working copy for one invocation:
//   - fresh command tree (deep copy) so Flag.changed/value mutations from a
//     previous Run can never leak into the next — hermetic by construction;
//   - invocation-scoped streams wired recursively;
//   - same engine and options.
//
// This is what makes repeated app.Run(...) calls safely independent in tests.
func (a *App) runSnapshot() *App {
	snap := &App{
		opts:   a.opts,
		engine: a.engine,
		out:    a.out,
		errOut: a.errOut,
		root:   cloneCommand(a.root),
		args:   a.args,
	}
	snap.root.wireOutput(snap.out)
	return snap
}

// cloneCommand deep-copies a command tree. Flags are copied structurally with
// runtime slots RESET to declaration defaults, so every snapshot starts
// pristine; declaration order is preserved deterministically for stable help.
func cloneCommand(c *Command) *Command {
	if c == nil {
		return nil
	}
	cp := *c

	if c.Aliases != nil {
		cp.Aliases = append([]string(nil), c.Aliases...)
	}
	if c.Flags != nil {
		cp.Flags = make([]*Flag, len(c.Flags))
		for i, f := range c.Flags {
			fc := *f
			fc.value = nil // drop last run's parsed state...
			fc.changed = false
			cp.Flags[i] = &fc
		}
	}
	if c.Commands != nil {
		cp.Commands = make([]*Command, len(c.Commands))
		for i, sub := range c.Commands {
			cp.Commands[i] = cloneCommand(sub)
		}
	}
	return &cp
}
