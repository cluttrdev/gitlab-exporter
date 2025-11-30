// Package main provides a minimal test helper that mimics a recorder subprocess.
// It creates a Unix socket and waits for signals, with flags to simulate various behaviors.
package main

import (
	"flag"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	address := flag.String("address", "", "Unix socket address (unix:///path/to/socket)")
	_ = flag.String("config", "", "Config file path (accepted but ignored)")
	delay := flag.Duration("delay", 0, "Delay before creating socket")
	exitEarly := flag.Bool("exit-early", false, "Exit immediately before creating socket")
	exitAfter := flag.Duration("exit-after", 0, "Exit after this duration (simulates crash)")
	ignoreSignals := flag.Bool("ignore-signals", false, "Ignore SIGTERM (for force kill testing)")
	flag.Parse()

	if *exitEarly {
		os.Exit(1)
	}

	if *delay > 0 {
		time.Sleep(*delay)
	}

	socketPath := strings.TrimPrefix(*address, "unix://")
	if socketPath == "" {
		os.Exit(1)
	}

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		os.Exit(1)
	}
	defer listener.Close()

	sigCh := make(chan os.Signal, 1)
	if !*ignoreSignals {
		signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
	}

	if *exitAfter > 0 {
		// Exit after specified duration (simulates a crash)
		select {
		case <-sigCh:
			return
		case <-time.After(*exitAfter):
			os.Exit(1)
		}
	} else if *ignoreSignals {
		// Block forever when ignoring signals (will be killed)
		select {}
	} else {
		// Wait for signal
		<-sigCh
	}
}
