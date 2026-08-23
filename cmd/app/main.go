// Package main provides a minimal CLI entry point for golang-cli-template.
// Replace this with your actual CLI commands (e.g., using cobra, urfave/cli, or stdlib flag).
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/guilhermelinosp/hellnet-lib-telemetry/telemetry"
)

func main() {
	ops, err := telemetry.New(telemetry.Options{
		ServiceName: "golang-cli-template",
		Enabled:     os.Getenv("HELLNET_TELEMETRY_ENABLED") == "true",
	})
	if err != nil {
		log.Fatalf("failed to init telemetry: %v", err)
	}
	defer func() { _ = ops.Shutdown() }()

	var name string
	flag.StringVar(&name, "name", "world", "name to greet")
	flag.Parse()

	ops.Log().Info("greeting", "name", name)
	_, _ = fmt.Fprintf(os.Stdout, "hello %s\n", name)
}
