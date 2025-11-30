package subprocess

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
)

const (
	// DefaultSocketTimeout is how long to wait for the subprocess socket to become available
	DefaultSocketTimeout = 10 * time.Second

	// DefaultShutdownTimeout is how long to wait for graceful shutdown before SIGKILL
	DefaultShutdownTimeout = 5 * time.Second

	// DefaultMaxRestarts is the default number of restart attempts (0 = no restarts)
	DefaultMaxRestarts = 3

	// socketPollInterval is how often to check if the socket is available
	socketPollInterval = 100 * time.Millisecond

	// restartBackoffBase is the base delay between restart attempts
	restartBackoffBase = 100 * time.Millisecond

	// restartBackoffMax is the maximum delay between restart attempts
	restartBackoffMax = 5 * time.Second
)

// LauncherConfig holds configuration for the subprocess launcher
type LauncherConfig struct {
	// RecorderType is the type of recorder (e.g., "sqlite", "clickhouse")
	RecorderType string

	// Command is the explicit path to the recorder binary (optional)
	// If empty, will look for gitlab-exporter-{type}-recorder in PATH
	Command string

	// Args are additional arguments to pass to the subprocess
	Args []string

	// ConfigPath is the path to the config file to pass to the subprocess
	ConfigPath string

	// SocketPath is the Unix socket path for IPC (auto-generated if empty)
	SocketPath string

	// SocketTimeout is how long to wait for the socket to become available
	SocketTimeout time.Duration

	// ShutdownTimeout is how long to wait for graceful shutdown
	ShutdownTimeout time.Duration

	// MaxRestarts is the maximum number of restart attempts (0 = no restarts, -1 = unlimited)
	MaxRestarts int
}

// Launcher manages a subprocess recorder
type Launcher struct {
	config LauncherConfig

	mu            sync.RWMutex
	cmd           *exec.Cmd
	ctx           context.Context    // parent context for subprocess lifecycle
	processCancel context.CancelFunc // cancels the subprocess context
	restartCount  int
	running       bool

	// For process monitoring and clean shutdown
	stopCh chan struct{} // signals monitor to stop
	doneCh chan struct{} // signals monitor has stopped
}

// NewLauncher creates a new subprocess launcher
func NewLauncher(cfg LauncherConfig) (*Launcher, error) {
	// Set defaults
	if cfg.SocketTimeout == 0 {
		cfg.SocketTimeout = DefaultSocketTimeout
	}
	if cfg.ShutdownTimeout == 0 {
		cfg.ShutdownTimeout = DefaultShutdownTimeout
	}
	if cfg.MaxRestarts == 0 {
		cfg.MaxRestarts = DefaultMaxRestarts
	}

	// Generate socket path if not provided
	if cfg.SocketPath == "" {
		cfg.SocketPath = fmt.Sprintf("/tmp/gitlab-exporter-%s-%d.sock", cfg.RecorderType, os.Getpid())
	}

	// Resolve command if not explicitly provided
	if cfg.Command == "" {
		binaryName := fmt.Sprintf("gitlab-exporter-%s-recorder", cfg.RecorderType)
		path, err := exec.LookPath(binaryName)
		if err != nil {
			return nil, fmt.Errorf("recorder binary %q not found in PATH: %w", binaryName, err)
		}
		cfg.Command = path
	}

	return &Launcher{
		config: cfg,
		stopCh: make(chan struct{}),
		doneCh: make(chan struct{}),
	}, nil
}

// Start launches the recorder subprocess and begins monitoring it
func (l *Launcher) Start(ctx context.Context) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.running {
		return fmt.Errorf("already running")
	}

	// Reset state for fresh start
	l.ctx = ctx
	l.stopCh = make(chan struct{})
	l.doneCh = make(chan struct{})
	l.restartCount = 0

	// Clean up any stale socket file
	if err := l.cleanupSocket(); err != nil {
		return fmt.Errorf("cleanup stale socket: %w", err)
	}

	// Start the subprocess
	if err := l.startProcess(ctx); err != nil {
		return fmt.Errorf("start subprocess: %w", err)
	}

	// Wait for socket to become available
	if err := l.waitForSocket(ctx); err != nil {
		l.killProcess()
		return fmt.Errorf("wait for socket: %w", err)
	}

	l.running = true

	// Start monitoring goroutine
	go l.monitor(ctx)

	return nil
}

// Stop gracefully shuts down the subprocess
func (l *Launcher) Stop(_ context.Context) error {
	l.mu.Lock()
	if !l.running {
		l.mu.Unlock()
		return nil
	}

	// Signal monitor goroutine to stop (prevents restart attempts)
	close(l.stopCh)

	// Cancel the subprocess context - this triggers graceful shutdown
	// via cmd.Cancel (SIGTERM) followed by SIGKILL after WaitDelay
	if l.processCancel != nil {
		l.processCancel()
	}
	l.mu.Unlock()

	// Wait for monitor to finish (it will exit after process dies)
	<-l.doneCh

	l.mu.Lock()
	defer l.mu.Unlock()

	// Clean up socket file
	_ = l.cleanupSocket()

	l.running = false
	l.cmd = nil

	return nil
}

// IsRunning returns whether the subprocess is running
func (l *Launcher) IsRunning() bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.running
}

// SocketPath returns the Unix socket path
func (l *Launcher) SocketPath() string {
	return l.config.SocketPath
}

// RestartCount returns the number of times the subprocess has been restarted
func (l *Launcher) RestartCount() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.restartCount
}

// monitor watches the subprocess and restarts it if it exits unexpectedly
func (l *Launcher) monitor(ctx context.Context) {
	defer close(l.doneCh)

	for {
		// Wait for process to exit
		waitCh := make(chan error, 1)
		go func() {
			if l.cmd != nil && l.cmd.Process != nil {
				waitCh <- l.cmd.Wait()
			} else {
				waitCh <- fmt.Errorf("no process to wait for")
			}
		}()

		select {
		case <-l.stopCh:
			// Intentional stop requested - don't restart
			return

		case <-ctx.Done():
			// Parent context cancelled - don't restart
			return

		case err := <-waitCh:
			// Process exited - check if we should restart
			l.mu.Lock()

			// Check if stop was requested or context cancelled while we were waiting
			select {
			case <-l.stopCh:
				l.mu.Unlock()
				return
			case <-ctx.Done():
				l.mu.Unlock()
				return
			default:
			}

			// Log the exit
			exitCode := -1
			if l.cmd != nil && l.cmd.ProcessState != nil {
				exitCode = l.cmd.ProcessState.ExitCode()
			}
			slog.Warn("Subprocess exited unexpectedly",
				"recorder_type", l.config.RecorderType,
				"exit_code", exitCode,
				"error", err,
				"restart_count", l.restartCount,
			)

			// Check if we can restart
			if !l.canRestart() {
				slog.Error("Max restarts exceeded, giving up",
					"recorder_type", l.config.RecorderType,
					"max_restarts", l.config.MaxRestarts,
					"restart_count", l.restartCount,
				)
				l.running = false
				l.mu.Unlock()
				return
			}

			// Attempt restart with backoff
			l.restartCount++
			backoff := l.calculateBackoff()
			l.mu.Unlock()

			slog.Info("Attempting to restart subprocess",
				"recorder_type", l.config.RecorderType,
				"restart_count", l.restartCount,
				"backoff", backoff,
			)

			// Wait for backoff period (can be interrupted by stop or context cancellation)
			select {
			case <-l.stopCh:
				return
			case <-ctx.Done():
				return
			case <-time.After(backoff):
			}

			// Perform restart
			l.mu.Lock()
			select {
			case <-l.stopCh:
				l.mu.Unlock()
				return
			case <-ctx.Done():
				l.mu.Unlock()
				return
			default:
			}

			if err := l.restart(ctx); err != nil {
				slog.Error("Failed to restart subprocess",
					"recorder_type", l.config.RecorderType,
					"error", err,
				)
				l.running = false
				l.mu.Unlock()
				return
			}
			l.mu.Unlock()

			slog.Info("Subprocess restarted successfully",
				"recorder_type", l.config.RecorderType,
				"restart_count", l.restartCount,
			)
		}
	}
}

// canRestart returns true if restart is allowed based on MaxRestarts config
func (l *Launcher) canRestart() bool {
	if l.config.MaxRestarts < 0 {
		// Unlimited restarts
		return true
	}
	return l.restartCount < l.config.MaxRestarts
}

// calculateBackoff returns the backoff duration for the current restart attempt
func (l *Launcher) calculateBackoff() time.Duration {
	// Exponential backoff: base * 2^(restartCount-1), capped at max
	backoff := restartBackoffBase
	for i := 1; i < l.restartCount && backoff < restartBackoffMax; i++ {
		backoff *= 2
	}
	if backoff > restartBackoffMax {
		backoff = restartBackoffMax
	}
	return backoff
}

// restart performs the actual restart of the subprocess (must be called with lock held)
func (l *Launcher) restart(ctx context.Context) error {
	// Clean up old process
	l.killProcess()
	if err := l.cleanupSocket(); err != nil {
		return fmt.Errorf("cleanup socket: %w", err)
	}

	// Start new process
	if err := l.startProcess(ctx); err != nil {
		return fmt.Errorf("start process: %w", err)
	}

	// Wait for socket
	if err := l.waitForSocket(ctx); err != nil {
		l.killProcess()
		return fmt.Errorf("wait for socket: %w", err)
	}

	return nil
}

// startProcess starts the recorder subprocess
func (l *Launcher) startProcess(ctx context.Context) error {
	args := []string{
		"--address", "unix://" + l.config.SocketPath,
	}

	if l.config.ConfigPath != "" {
		args = append(args, "--config", l.config.ConfigPath)
	}

	args = append(args, l.config.Args...)

	// Derive from caller's context so either parent cancellation or Stop() triggers shutdown
	processCtx, processCancel := context.WithCancel(ctx)
	l.processCancel = processCancel

	l.cmd = exec.CommandContext(processCtx, l.config.Command, args...)

	// Configure graceful shutdown: SIGTERM first, then SIGKILL after WaitDelay
	l.cmd.Cancel = func() error {
		return l.cmd.Process.Signal(syscall.SIGTERM)
	}
	l.cmd.WaitDelay = l.config.ShutdownTimeout

	// Inherit stdout/stderr for logging
	l.cmd.Stdout = os.Stdout
	l.cmd.Stderr = os.Stderr

	if err := l.cmd.Start(); err != nil {
		processCancel() // Clean up context on failure
		return fmt.Errorf("exec %s: %w", l.config.Command, err)
	}

	return nil
}

// waitForSocket polls until the socket becomes available or timeout
func (l *Launcher) waitForSocket(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, l.config.SocketTimeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Check if process is still alive
		if l.cmd.ProcessState != nil && l.cmd.ProcessState.Exited() {
			return fmt.Errorf("subprocess exited prematurely with code %d", l.cmd.ProcessState.ExitCode())
		}

		// Try to connect to the socket
		conn, err := net.DialTimeout("unix", l.config.SocketPath, socketPollInterval)
		if err == nil {
			_ = conn.Close()
			return nil
		}

		time.Sleep(socketPollInterval)
	}
}

// killProcess forcefully kills the subprocess
func (l *Launcher) killProcess() {
	if l.cmd != nil && l.cmd.Process != nil {
		_ = l.cmd.Process.Kill()
		_, _ = l.cmd.Process.Wait() // Reap the zombie
	}
}

// cleanupSocket removes the socket file if it exists
func (l *Launcher) cleanupSocket() error {
	// Handle unix:// prefix if present
	socketPath := strings.TrimPrefix(l.config.SocketPath, "unix://")

	// Ensure parent directory exists
	dir := filepath.Dir(socketPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil // Directory doesn't exist, nothing to clean
	}

	// Remove socket file if it exists
	if err := os.Remove(socketPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove socket file %s: %w", socketPath, err)
	}

	return nil
}
