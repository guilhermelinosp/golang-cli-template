// Package build holds build metadata that is injected at link time.
//
// Do not hardcode values here. GoReleaser (or plain go build with -ldflags)
// sets these variables, e.g.:
//
//	go build -ldflags "-X myapp/internal/build.Version=1.2.3 \
//	                  -X myapp/internal/build.Commit=abc1234 \
//	                  -X myapp/internal/build.Date=2026-01-02T15:04:05Z"
//
// When the binary is built without ldflags injection (go run, make build
// without VCS metadata), Version falls back to "dev" so it is always obvious
// which binary is running.
package build

import "runtime"

// Overridable build metadata. All fields default to safe fallbacks so the CLI
// never prints empty output when it is compiled manually.
var (
	// Version is the semantic version of the binary (e.g. "1.4.2").
	Version = "dev"
	// Commit is the full or short git commit SHA the binary was built from.
	Commit = "none"
	// Date is the RFC3339 build timestamp (e.g. "2026-08-26T10:00:00Z").
	Date = "unknown"
)

// Info is an immutable snapshot of the current build metadata.
type Info struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	Date      string `json:"date"`
	GoVersion string `json:"go_version"`
}

// Get returns the current build metadata snapshot.
func Get() Info {
	return Info{
		Version:   Version,
		Commit:    Commit,
		Date:      Date,
		GoVersion: runtime.Version(),
	}
}

// String renders a single-line, script-friendly summary:
//
//	myapp 1.4.2 (commit abc1234, built 2026-08-26T10:00:00Z, go go1.27.0)
func (i Info) String() string {
	return i.Version + " (commit " + i.Commit + ", built " + i.Date + ", go " + i.GoVersion + ")"
}
