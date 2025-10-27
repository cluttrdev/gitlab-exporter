package main

import (
	"fmt"
	"os"

	"go.cluttr.dev/gitlab-exporter/cmd"
)

var version string

func main() {
	cmd.Version = version

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
