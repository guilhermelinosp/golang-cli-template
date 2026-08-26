package cli

import (
	"errors"
	"fmt"
)

// Exit codes used by the template. Apps may define additional codes for
// their own domain semantics, but should stay in this range style.
const (
	// ExitSuccess means the command completed without error.
	ExitSuccess = 0
	// ExitFailure is the generic runtime/domain failure code.
	ExitFailure = 1
	// ExitUsage is returned for flag/argument misuse (mirrors BSD/GNU tools).
	ExitUsage = 2
)

// ExitCoder is implemented by errors that know the process exit code they
// should produce. Commands may return plain errors; the framework wraps them
// with [ExitFailure] automatically.
type ExitCoder interface {
	error
	ExitCode() int
}

// exitError couples an exit code with a formatted error (wrap-friendly via
// %w) for centralized rendering decisions downstream.
type exitError struct {
	code int
	err  error
}

func (e *exitError) Error() string { return e.err.Error() }
func (e *exitError) ExitCode() int { return e.code }
func (e *exitError) Unwrap() error { return e.err }

// Errorf returns a domain error with [ExitFailure] that is displayed to the
// user verbatim — no stack traces, no internal details. Use it when user
// input or business rules fail:
//
//	return cli.Errorf("unknown environment %q", env)
//
// Wrap an existing error with %w to keep it available for logging:
//
//	return cli.Errorf("read config: %w", err)
func Errorf(format string, args ...any) error {
	return &exitError{code: ExitFailure, err: fmt.Errorf(format, args...)}
}

// Exitf behaves like [Errorf] but with an explicit process exit code.
func Exitf(code int, format string, args ...any) error {
	return &exitError{code: code, err: fmt.Errorf(format, args...)}
}

// Usagef returns a usage error ([ExitUsage]) shown next to command help.
// The framework already produces usage errors for malformed flags and bad
// argument counts; use this only for custom input validation.
func Usagef(format string, args ...any) error {
	return &exitError{code: ExitUsage, err: fmt.Errorf(format, args...)}
}

// ExitCodeOf reports the process exit code for err: 0 when err is nil, the
// code carried by any wrapped [ExitCoder], or [ExitFailure] otherwise.
func ExitCodeOf(err error) int {
	if err == nil {
		return ExitSuccess
	}
	var ec ExitCoder
	if errors.As(err, &ec) {
		return ec.ExitCode()
	}
	return ExitFailure
}
