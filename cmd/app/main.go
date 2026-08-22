// Package main provides a minimal CLI entry point for golang-cli-template.
// Replace this with your actual CLI commands (e.g., using cobra, urfave/cli, or stdlib flag).
package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var name string
	flag.StringVar(&name, "name", "world", "name to greet")
	flag.Parse()

	_, _ = fmt.Fprintf(os.Stdout, "hello %s\n", name)
}
