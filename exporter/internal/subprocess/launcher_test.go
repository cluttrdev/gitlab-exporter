package subprocess

import (
	"context"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// buildTestHelper compiles the test helper binary and returns its path.
func buildTestHelper(t *testing.T) string {
	t.Helper()

	helperDir := filepath.Join("testdata", "testhelper")
	binary := filepath.Join(t.TempDir(), "testhelper")

	cmd := exec.Command("go", "build", "-o", binary, ".")
	cmd.Dir = helperDir
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build test helper: %v\n%s", err, output)
	}

	return binary
}

// --- Unit Tests (no subprocess execution) ---

func TestNewLauncher_Defaults(t *testing.T) {
	// Use explicit command to avoid PATH lookup failure
	tempBinary := filepath.Join(t.TempDir(), "fake-binary")
	if err := os.WriteFile(tempBinary, []byte(""), 0755); err != nil {
		t.Fatalf("create fake binary: %v", err)
	}

	launcher, err := NewLauncher(LauncherConfig{
		RecorderType: "test",
		Command:      tempBinary,
	})
	if err != nil {
		t.Fatalf("NewLauncher error: %v", err)
	}

	// Check defaults
	if launcher.config.SocketTimeout != DefaultSocketTimeout {
		t.Errorf("SocketTimeout = %v, want %v", launcher.config.SocketTimeout, DefaultSocketTimeout)
	}
	if launcher.config.ShutdownTimeout != DefaultShutdownTimeout {
		t.Errorf("ShutdownTimeout = %v, want %v", launcher.config.ShutdownTimeout, DefaultShutdownTimeout)
	}
	if launcher.config.MaxRestarts != DefaultMaxRestarts {
		t.Errorf("MaxRestarts = %d, want %d", launcher.config.MaxRestarts, DefaultMaxRestarts)
	}
}

func TestNewLauncher_SocketPathGenerated(t *testing.T) {
	tempBinary := filepath.Join(t.TempDir(), "fake-binary")
	if err := os.WriteFile(tempBinary, []byte(""), 0755); err != nil {
		t.Fatalf("create fake binary: %v", err)
	}

	launcher, err := NewLauncher(LauncherConfig{
		RecorderType: "mytype",
		Command:      tempBinary,
	})
	if err != nil {
		t.Fatalf("NewLauncher error: %v", err)
	}

	socketPath := launcher.SocketPath()

	// Should contain the recorder type
	if !strings.Contains(socketPath, "mytype") {
		t.Errorf("SocketPath %q should contain recorder type 'mytype'", socketPath)
	}

	// Should be in /tmp
	if !strings.HasPrefix(socketPath, "/tmp/") {
		t.Errorf("SocketPath %q should start with /tmp/", socketPath)
	}

	// Should end with .sock
	if !strings.HasSuffix(socketPath, ".sock") {
		t.Errorf("SocketPath %q should end with .sock", socketPath)
	}
}

func TestNewLauncher_ExplicitSocketPath(t *testing.T) {
	tempBinary := filepath.Join(t.TempDir(), "fake-binary")
	if err := os.WriteFile(tempBinary, []byte(""), 0755); err != nil {
		t.Fatalf("create fake binary: %v", err)
	}

	customPath := "/custom/path/to/socket.sock"
	launcher, err := NewLauncher(LauncherConfig{
		RecorderType: "test",
		Command:      tempBinary,
		SocketPath:   customPath,
	})
	if err != nil {
		t.Fatalf("NewLauncher error: %v", err)
	}

	if launcher.SocketPath() != customPath {
		t.Errorf("SocketPath = %q, want %q", launcher.SocketPath(), customPath)
	}
}

func TestNewLauncher_ExplicitCommand(t *testing.T) {
	tempBinary := filepath.Join(t.TempDir(), "my-custom-recorder")
	if err := os.WriteFile(tempBinary, []byte(""), 0755); err != nil {
		t.Fatalf("create fake binary: %v", err)
	}

	launcher, err := NewLauncher(LauncherConfig{
		RecorderType: "test",
		Command:      tempBinary,
	})
	if err != nil {
		t.Fatalf("NewLauncher error: %v", err)
	}

	if launcher.config.Command != tempBinary {
		t.Errorf("Command = %q, want %q", launcher.config.Command, tempBinary)
	}
}

func TestNewLauncher_CommandNotFound(t *testing.T) {
	_, err := NewLauncher(LauncherConfig{
		RecorderType: "nonexistent-recorder-type-12345",
		// No explicit Command, so it will try PATH lookup
	})

	if err == nil {
		t.Error("expected error for missing binary")
	}

	if !strings.Contains(err.Error(), "not found in PATH") {
		t.Errorf("error should mention 'not found in PATH', got: %v", err)
	}
}

func TestNewLauncher_CustomTimeouts(t *testing.T) {
	tempBinary := filepath.Join(t.TempDir(), "fake-binary")
	if err := os.WriteFile(tempBinary, []byte(""), 0755); err != nil {
		t.Fatalf("create fake binary: %v", err)
	}

	launcher, err := NewLauncher(LauncherConfig{
		RecorderType:    "test",
		Command:         tempBinary,
		SocketTimeout:   5 * time.Second,
		ShutdownTimeout: 2 * time.Second,
		MaxRestarts:     10,
	})
	if err != nil {
		t.Fatalf("NewLauncher error: %v", err)
	}

	if launcher.config.SocketTimeout != 5*time.Second {
		t.Errorf("SocketTimeout = %v, want 5s", launcher.config.SocketTimeout)
	}
	if launcher.config.ShutdownTimeout != 2*time.Second {
		t.Errorf("ShutdownTimeout = %v, want 2s", launcher.config.ShutdownTimeout)
	}
	if launcher.config.MaxRestarts != 10 {
		t.Errorf("MaxRestarts = %d, want 10", launcher.config.MaxRestarts)
	}
}

func TestLauncher_SocketPath(t *testing.T) {
	tempBinary := filepath.Join(t.TempDir(), "fake-binary")
	if err := os.WriteFile(tempBinary, []byte(""), 0755); err != nil {
		t.Fatalf("create fake binary: %v", err)
	}

	socketPath := filepath.Join(t.TempDir(), "test.sock")
	launcher, err := NewLauncher(LauncherConfig{
		RecorderType: "test",
		Command:      tempBinary,
		SocketPath:   socketPath,
	})
	if err != nil {
		t.Fatalf("NewLauncher error: %v", err)
	}

	if got := launcher.SocketPath(); got != socketPath {
		t.Errorf("SocketPath() = %q, want %q", got, socketPath)
	}
}

func TestLauncher_IsRunning_Initial(t *testing.T) {
	tempBinary := filepath.Join(t.TempDir(), "fake-binary")
	if err := os.WriteFile(tempBinary, []byte(""), 0755); err != nil {
		t.Fatalf("create fake binary: %v", err)
	}

	launcher, err := NewLauncher(LauncherConfig{
		RecorderType: "test",
		Command:      tempBinary,
	})
	if err != nil {
		t.Fatalf("NewLauncher error: %v", err)
	}

	if launcher.IsRunning() {
		t.Error("IsRunning() should be false before Start()")
	}
}

func TestLauncher_StopNotRunning(t *testing.T) {
	tempBinary := filepath.Join(t.TempDir(), "fake-binary")
	if err := os.WriteFile(tempBinary, []byte(""), 0755); err != nil {
		t.Fatalf("create fake binary: %v", err)
	}

	launcher, err := NewLauncher(LauncherConfig{
		RecorderType: "test",
		Command:      tempBinary,
	})
	if err != nil {
		t.Fatalf("NewLauncher error: %v", err)
	}

	err = launcher.Stop(t.Context())
	if err != nil {
		t.Errorf("Stop() on non-running launcher should not error, got: %v", err)
	}
}

func TestLauncher_CleanupSocket(t *testing.T) {
	tempDir := t.TempDir()
	socketPath := filepath.Join(tempDir, "test.sock")

	// Create a dummy socket file
	if err := os.WriteFile(socketPath, []byte(""), 0644); err != nil {
		t.Fatalf("create dummy socket: %v", err)
	}

	tempBinary := filepath.Join(tempDir, "fake-binary")
	if err := os.WriteFile(tempBinary, []byte(""), 0755); err != nil {
		t.Fatalf("create fake binary: %v", err)
	}

	launcher, err := NewLauncher(LauncherConfig{
		RecorderType: "test",
		Command:      tempBinary,
		SocketPath:   socketPath,
	})
	if err != nil {
		t.Fatalf("NewLauncher error: %v", err)
	}

	// Call cleanupSocket directly
	if err := launcher.cleanupSocket(); err != nil {
		t.Errorf("cleanupSocket error: %v", err)
	}

	// Verify socket file was removed
	if _, err := os.Stat(socketPath); !os.IsNotExist(err) {
		t.Error("socket file should have been removed")
	}
}

func TestLauncher_CleanupSocket_NonExistent(t *testing.T) {
	tempDir := t.TempDir()
	socketPath := filepath.Join(tempDir, "nonexistent.sock")

	tempBinary := filepath.Join(tempDir, "fake-binary")
	if err := os.WriteFile(tempBinary, []byte(""), 0755); err != nil {
		t.Fatalf("create fake binary: %v", err)
	}

	launcher, err := NewLauncher(LauncherConfig{
		RecorderType: "test",
		Command:      tempBinary,
		SocketPath:   socketPath,
	})
	if err != nil {
		t.Fatalf("NewLauncher error: %v", err)
	}

	// Should not error when socket doesn't exist
	if err := launcher.cleanupSocket(); err != nil {
		t.Errorf("cleanupSocket should not error for non-existent socket: %v", err)
	}
}

func TestLauncher_CleanupSocket_WithUnixPrefix(t *testing.T) {
	tempDir := t.TempDir()
	socketPath := filepath.Join(tempDir, "test.sock")

	// Create a dummy socket file
	if err := os.WriteFile(socketPath, []byte(""), 0644); err != nil {
		t.Fatalf("create dummy socket: %v", err)
	}

	tempBinary := filepath.Join(tempDir, "fake-binary")
	if err := os.WriteFile(tempBinary, []byte(""), 0755); err != nil {
		t.Fatalf("create fake binary: %v", err)
	}

	// Use unix:// prefix in socket path
	launcher, err := NewLauncher(LauncherConfig{
		RecorderType: "test",
		Command:      tempBinary,
		SocketPath:   "unix://" + socketPath,
	})
	if err != nil {
		t.Fatalf("NewLauncher error: %v", err)
	}

	// cleanupSocket should handle the unix:// prefix
	if err := launcher.cleanupSocket(); err != nil {
		t.Errorf("cleanupSocket error: %v", err)
	}

	// Verify socket file was removed
	if _, err := os.Stat(socketPath); !os.IsNotExist(err) {
		t.Error("socket file should have been removed")
	}
}

// --- Integration Tests (with subprocess execution) ---

func TestLauncher_StartStop(t *testing.T) {
	binary := buildTestHelper(t)
	socketPath := filepath.Join(t.TempDir(), "test.sock")

	launcher, err := NewLauncher(LauncherConfig{
		RecorderType:  "test",
		Command:       binary,
		SocketPath:    socketPath,
		SocketTimeout: 2 * time.Second,
	})
	if err != nil {
		t.Fatalf("NewLauncher error: %v", err)
	}

	ctx := context.Background()

	// Start
	if err := launcher.Start(ctx); err != nil {
		t.Fatalf("Start error: %v", err)
	}

	// Verify running
	if !launcher.IsRunning() {
		t.Error("IsRunning() should be true after Start()")
	}

	// Verify socket exists
	conn, err := net.DialTimeout("unix", socketPath, time.Second)
	if err != nil {
		t.Errorf("should be able to connect to socket: %v", err)
	} else {
		_ = conn.Close()
	}

	// Stop
	if err := launcher.Stop(ctx); err != nil {
		t.Errorf("Stop error: %v", err)
	}

	// Verify not running
	if launcher.IsRunning() {
		t.Error("IsRunning() should be false after Stop()")
	}
}

func TestLauncher_StartAlreadyRunning(t *testing.T) {
	binary := buildTestHelper(t)
	socketPath := filepath.Join(t.TempDir(), "test.sock")

	launcher, err := NewLauncher(LauncherConfig{
		RecorderType:  "test",
		Command:       binary,
		SocketPath:    socketPath,
		SocketTimeout: 2 * time.Second,
	})
	if err != nil {
		t.Fatalf("NewLauncher error: %v", err)
	}

	ctx := context.Background()

	// Start first time
	if err := launcher.Start(ctx); err != nil {
		t.Fatalf("Start error: %v", err)
	}
	defer func() { _ = launcher.Stop(ctx) }()

	// Start second time should error
	err = launcher.Start(ctx)
	if err == nil {
		t.Error("expected error when starting already running launcher")
	}
	if !strings.Contains(err.Error(), "already running") {
		t.Errorf("error should mention 'already running', got: %v", err)
	}
}

func TestLauncher_IsRunning_AfterStart(t *testing.T) {
	binary := buildTestHelper(t)
	socketPath := filepath.Join(t.TempDir(), "test.sock")

	launcher, err := NewLauncher(LauncherConfig{
		RecorderType:  "test",
		Command:       binary,
		SocketPath:    socketPath,
		SocketTimeout: 2 * time.Second,
	})
	if err != nil {
		t.Fatalf("NewLauncher error: %v", err)
	}

	ctx := context.Background()

	if launcher.IsRunning() {
		t.Error("IsRunning() should be false before Start()")
	}

	if err := launcher.Start(ctx); err != nil {
		t.Fatalf("Start error: %v", err)
	}

	if !launcher.IsRunning() {
		t.Error("IsRunning() should be true after Start()")
	}

	if err := launcher.Stop(ctx); err != nil {
		t.Errorf("Stop error: %v", err)
	}

	if launcher.IsRunning() {
		t.Error("IsRunning() should be false after Stop()")
	}
}

func TestLauncher_SocketTimeout(t *testing.T) {
	binary := buildTestHelper(t)
	socketPath := filepath.Join(t.TempDir(), "test.sock")

	launcher, err := NewLauncher(LauncherConfig{
		RecorderType:  "test",
		Command:       binary,
		SocketPath:    socketPath,
		SocketTimeout: 200 * time.Millisecond,    // Short timeout
		Args:          []string{"--delay", "5s"}, // Subprocess delays socket creation
	})
	if err != nil {
		t.Fatalf("NewLauncher error: %v", err)
	}

	ctx := context.Background()

	start := time.Now()
	err = launcher.Start(ctx)
	elapsed := time.Since(start)

	if err == nil {
		_ = launcher.Stop(ctx)
		t.Fatal("expected timeout error")
	}

	// Should timeout around 200ms, not wait for the full 5s delay
	if elapsed > time.Second {
		t.Errorf("Start took too long (%v), should have timed out around 200ms", elapsed)
	}

	if !strings.Contains(err.Error(), "wait for socket") {
		t.Errorf("error should mention 'wait for socket', got: %v", err)
	}
}

func TestLauncher_ProcessExitedPrematurely(t *testing.T) {
	binary := buildTestHelper(t)
	socketPath := filepath.Join(t.TempDir(), "test.sock")

	launcher, err := NewLauncher(LauncherConfig{
		RecorderType:  "test",
		Command:       binary,
		SocketPath:    socketPath,
		SocketTimeout: 2 * time.Second,
		Args:          []string{"--exit-early"}, // Subprocess exits immediately
	})
	if err != nil {
		t.Fatalf("NewLauncher error: %v", err)
	}

	ctx := context.Background()

	err = launcher.Start(ctx)
	if err == nil {
		_ = launcher.Stop(ctx)
		t.Fatal("expected error when subprocess exits prematurely")
	}

	if !strings.Contains(err.Error(), "wait for socket") {
		t.Errorf("error should mention 'wait for socket', got: %v", err)
	}
}

func TestLauncher_ForceKill(t *testing.T) {
	binary := buildTestHelper(t)
	socketPath := filepath.Join(t.TempDir(), "test.sock")

	launcher, err := NewLauncher(LauncherConfig{
		RecorderType:    "test",
		Command:         binary,
		SocketPath:      socketPath,
		SocketTimeout:   2 * time.Second,
		ShutdownTimeout: 200 * time.Millisecond,       // Short shutdown timeout
		Args:            []string{"--ignore-signals"}, // Subprocess ignores SIGTERM
	})
	if err != nil {
		t.Fatalf("NewLauncher error: %v", err)
	}

	ctx := context.Background()

	if err := launcher.Start(ctx); err != nil {
		t.Fatalf("Start error: %v", err)
	}

	// Stop should force kill after shutdown timeout
	start := time.Now()
	err = launcher.Stop(ctx)
	elapsed := time.Since(start)

	if err != nil {
		t.Errorf("Stop error: %v", err)
	}

	// Should take around shutdown timeout (200ms), not hang forever
	if elapsed > 2*time.Second {
		t.Errorf("Stop took too long (%v), should have force-killed around 200ms", elapsed)
	}

	if launcher.IsRunning() {
		t.Error("IsRunning() should be false after Stop()")
	}
}

func TestLauncher_ConfigPath(t *testing.T) {
	binary := buildTestHelper(t)
	tempDir := t.TempDir()
	socketPath := filepath.Join(tempDir, "test.sock")
	configPath := filepath.Join(tempDir, "config.yaml")

	// Create a dummy config file
	if err := os.WriteFile(configPath, []byte("test: true"), 0644); err != nil {
		t.Fatalf("create config file: %v", err)
	}

	launcher, err := NewLauncher(LauncherConfig{
		RecorderType:  "test",
		Command:       binary,
		SocketPath:    socketPath,
		SocketTimeout: 2 * time.Second,
		ConfigPath:    configPath,
	})
	if err != nil {
		t.Fatalf("NewLauncher error: %v", err)
	}

	ctx := context.Background()

	// Start should pass config path to subprocess
	if err := launcher.Start(ctx); err != nil {
		t.Fatalf("Start error: %v", err)
	}
	defer func() { _ = launcher.Stop(ctx) }()

	// Verify socket was created (subprocess started successfully)
	conn, err := net.DialTimeout("unix", socketPath, time.Second)
	if err != nil {
		t.Errorf("should be able to connect to socket: %v", err)
	} else {
		_ = conn.Close()
	}
}

func TestLauncher_ContextCancellation(t *testing.T) {
	binary := buildTestHelper(t)
	socketPath := filepath.Join(t.TempDir(), "test.sock")

	launcher, err := NewLauncher(LauncherConfig{
		RecorderType:    "test",
		Command:         binary,
		SocketPath:      socketPath,
		SocketTimeout:   2 * time.Second,
		ShutdownTimeout: 2 * time.Second,
		Args:            []string{"--ignore-signals"},
	})
	if err != nil {
		t.Fatalf("NewLauncher error: %v", err)
	}

	ctx := context.Background()

	if err := launcher.Start(ctx); err != nil {
		t.Fatalf("Start error: %v", err)
	}

	// Use cancelled context for Stop
	cancelCtx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	start := time.Now()
	err = launcher.Stop(cancelCtx)
	elapsed := time.Since(start)

	if err != nil {
		t.Errorf("Stop error: %v", err)
	}

	// Should force kill quickly due to cancelled context
	if elapsed > time.Second {
		t.Errorf("Stop took too long (%v) with cancelled context", elapsed)
	}

	if launcher.IsRunning() {
		t.Error("IsRunning() should be false after Stop()")
	}
}

// --- Restart Tests ---

func TestLauncher_RestartCount_Initial(t *testing.T) {
	tempBinary := filepath.Join(t.TempDir(), "fake-binary")
	if err := os.WriteFile(tempBinary, []byte(""), 0755); err != nil {
		t.Fatalf("create fake binary: %v", err)
	}

	launcher, err := NewLauncher(LauncherConfig{
		RecorderType: "test",
		Command:      tempBinary,
	})
	if err != nil {
		t.Fatalf("NewLauncher error: %v", err)
	}

	if launcher.RestartCount() != 0 {
		t.Errorf("RestartCount() = %d, want 0", launcher.RestartCount())
	}
}

func TestLauncher_RestartOnCrash(t *testing.T) {
	binary := buildTestHelper(t)
	socketPath := filepath.Join(t.TempDir(), "test.sock")

	launcher, err := NewLauncher(LauncherConfig{
		RecorderType:  "test",
		Command:       binary,
		SocketPath:    socketPath,
		SocketTimeout: 2 * time.Second,
		MaxRestarts:   3,
		Args:          []string{"--exit-after", "200ms"}, // Crash after 200ms
	})
	if err != nil {
		t.Fatalf("NewLauncher error: %v", err)
	}

	ctx := context.Background()

	if err := launcher.Start(ctx); err != nil {
		t.Fatalf("Start error: %v", err)
	}

	// Wait for crash and restart to occur
	time.Sleep(600 * time.Millisecond)

	// Should have restarted at least once
	if launcher.RestartCount() < 1 {
		t.Errorf("RestartCount() = %d, want >= 1", launcher.RestartCount())
	}

	// Should still be running (was restarted)
	if !launcher.IsRunning() {
		t.Error("IsRunning() should be true after restart")
	}

	// Stop the launcher
	if err := launcher.Stop(ctx); err != nil {
		t.Errorf("Stop error: %v", err)
	}
}

func TestLauncher_MaxRestartsExceeded(t *testing.T) {
	binary := buildTestHelper(t)
	socketPath := filepath.Join(t.TempDir(), "test.sock")

	launcher, err := NewLauncher(LauncherConfig{
		RecorderType:  "test",
		Command:       binary,
		SocketPath:    socketPath,
		SocketTimeout: 500 * time.Millisecond,
		MaxRestarts:   2,                                 // Allow only 2 restarts
		Args:          []string{"--exit-after", "100ms"}, // Crash quickly
	})
	if err != nil {
		t.Fatalf("NewLauncher error: %v", err)
	}

	ctx := context.Background()

	if err := launcher.Start(ctx); err != nil {
		t.Fatalf("Start error: %v", err)
	}

	// Wait for multiple crashes to exceed max restarts
	// With 100ms crash time + backoff, should exceed 2 restarts within 2 seconds
	time.Sleep(2 * time.Second)

	// Should have stopped after max restarts exceeded
	if launcher.IsRunning() {
		// Try to stop anyway for cleanup
		_ = launcher.Stop(ctx)
		t.Error("IsRunning() should be false after max restarts exceeded")
	}

	// RestartCount should be at MaxRestarts
	if launcher.RestartCount() < 2 {
		t.Errorf("RestartCount() = %d, want >= 2", launcher.RestartCount())
	}
}

func TestLauncher_NoRestarts(t *testing.T) {
	binary := buildTestHelper(t)
	socketPath := filepath.Join(t.TempDir(), "test.sock")

	launcher, err := NewLauncher(LauncherConfig{
		RecorderType:  "test",
		Command:       binary,
		SocketPath:    socketPath,
		SocketTimeout: 500 * time.Millisecond,
		MaxRestarts:   0, // No restarts allowed (but default is 3, so set explicitly)
		Args:          []string{"--exit-after", "100ms"},
	})
	if err != nil {
		t.Fatalf("NewLauncher error: %v", err)
	}

	// Override MaxRestarts to 0 (disable restarts)
	launcher.config.MaxRestarts = 0

	ctx := context.Background()

	if err := launcher.Start(ctx); err != nil {
		t.Fatalf("Start error: %v", err)
	}

	// Wait for crash
	time.Sleep(300 * time.Millisecond)

	// Should have stopped (no restarts allowed)
	if launcher.IsRunning() {
		_ = launcher.Stop(ctx)
		t.Error("IsRunning() should be false when MaxRestarts=0")
	}

	if launcher.RestartCount() != 0 {
		t.Errorf("RestartCount() = %d, want 0", launcher.RestartCount())
	}
}

func TestLauncher_StopDuringRestart(t *testing.T) {
	binary := buildTestHelper(t)
	socketPath := filepath.Join(t.TempDir(), "test.sock")

	launcher, err := NewLauncher(LauncherConfig{
		RecorderType:  "test",
		Command:       binary,
		SocketPath:    socketPath,
		SocketTimeout: 2 * time.Second,
		MaxRestarts:   10, // Allow many restarts
		Args:          []string{"--exit-after", "100ms"},
	})
	if err != nil {
		t.Fatalf("NewLauncher error: %v", err)
	}

	ctx := context.Background()

	if err := launcher.Start(ctx); err != nil {
		t.Fatalf("Start error: %v", err)
	}

	// Wait for first crash to start restart process
	time.Sleep(200 * time.Millisecond)

	// Stop during restart cycle
	start := time.Now()
	if err := launcher.Stop(ctx); err != nil {
		t.Errorf("Stop error: %v", err)
	}
	elapsed := time.Since(start)

	// Should stop quickly (not wait for all restarts)
	if elapsed > 2*time.Second {
		t.Errorf("Stop took too long (%v), should interrupt restart cycle", elapsed)
	}

	if launcher.IsRunning() {
		t.Error("IsRunning() should be false after Stop()")
	}
}

func TestLauncher_RestartPreservesSocket(t *testing.T) {
	binary := buildTestHelper(t)
	socketPath := filepath.Join(t.TempDir(), "test.sock")

	launcher, err := NewLauncher(LauncherConfig{
		RecorderType:  "test",
		Command:       binary,
		SocketPath:    socketPath,
		SocketTimeout: 2 * time.Second,
		MaxRestarts:   3,
		Args:          []string{"--exit-after", "300ms"}, // Longer run time
	})
	if err != nil {
		t.Fatalf("NewLauncher error: %v", err)
	}

	ctx := context.Background()

	if err := launcher.Start(ctx); err != nil {
		t.Fatalf("Start error: %v", err)
	}

	// Wait for at least one restart to complete
	// 300ms (crash) + 100ms (backoff) + 500ms (socket timeout margin) = ~900ms
	time.Sleep(700 * time.Millisecond)

	// Verify restart happened
	if launcher.RestartCount() < 1 {
		t.Errorf("expected at least one restart, got %d", launcher.RestartCount())
	}

	// Poll for socket availability (it may be briefly unavailable during restart)
	var conn net.Conn
	for i := 0; i < 10; i++ {
		conn, err = net.DialTimeout("unix", socketPath, 200*time.Millisecond)
		if err == nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	if err != nil {
		t.Errorf("should be able to connect to socket after restart: %v", err)
	} else {
		_ = conn.Close()
	}

	if err := launcher.Stop(ctx); err != nil {
		t.Errorf("Stop error: %v", err)
	}
}
