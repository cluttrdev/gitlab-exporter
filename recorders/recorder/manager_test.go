package recorder

import (
	"context"
	"errors"
	"testing"
)

type mockRecorderWithError struct {
	mockRecorder
	startErr  error
	stopErr   error
	healthErr error
}

func (m *mockRecorderWithError) Start(ctx context.Context) error {
	if m.startErr != nil {
		return m.startErr
	}
	return m.mockRecorder.Start(ctx)
}

func (m *mockRecorderWithError) Stop(ctx context.Context) error {
	if m.stopErr != nil {
		return m.stopErr
	}
	return m.mockRecorder.Stop(ctx)
}

func (m *mockRecorderWithError) CheckHealth(ctx context.Context) error {
	if m.healthErr != nil {
		return m.healthErr
	}
	return m.mockRecorder.CheckHealth(ctx)
}

func TestNewManager(t *testing.T) {
	m := NewManager()
	if m == nil {
		t.Fatal("expected non-nil manager")
	}

	if m.Count() != 0 {
		t.Errorf("expected 0 recorders, got %d", m.Count())
	}
}

func TestManagerAdd(t *testing.T) {
	m := NewManager()
	rec := &mockRecorder{name: "test"}

	m.Add(rec)

	if m.Count() != 1 {
		t.Errorf("expected 1 recorder, got %d", m.Count())
	}
}

func TestManagerStartAll(t *testing.T) {
	m := NewManager()
	rec1 := &mockRecorder{name: "rec1"}
	rec2 := &mockRecorder{name: "rec2"}

	m.Add(rec1)
	m.Add(rec2)

	ctx := context.Background()
	err := m.StartAll(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !rec1.started {
		t.Error("rec1 should be started")
	}
	if !rec2.started {
		t.Error("rec2 should be started")
	}
}

func TestManagerStartAllWithError(t *testing.T) {
	m := NewManager()
	rec1 := &mockRecorder{name: "rec1"}
	rec2 := &mockRecorderWithError{
		mockRecorder: mockRecorder{name: "rec2"},
		startErr:     errors.New("start failed"),
	}

	m.Add(rec1)
	m.Add(rec2)

	ctx := context.Background()
	err := m.StartAll(ctx)
	if err == nil {
		t.Error("expected error")
	}

	// rec1 should still have started
	if !rec1.started {
		t.Error("rec1 should be started despite rec2 error")
	}
}

func TestManagerStopAll(t *testing.T) {
	m := NewManager()
	rec1 := &mockRecorder{name: "rec1"}
	rec2 := &mockRecorder{name: "rec2"}

	m.Add(rec1)
	m.Add(rec2)

	ctx := context.Background()
	err := m.StopAll(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !rec1.stopped {
		t.Error("rec1 should be stopped")
	}
	if !rec2.stopped {
		t.Error("rec2 should be stopped")
	}
}

func TestManagerStopAllWithError(t *testing.T) {
	m := NewManager()
	rec1 := &mockRecorder{name: "rec1"}
	rec2 := &mockRecorderWithError{
		mockRecorder: mockRecorder{name: "rec2"},
		stopErr:      errors.New("stop failed"),
	}

	m.Add(rec1)
	m.Add(rec2)

	ctx := context.Background()
	err := m.StopAll(ctx)
	if err == nil {
		t.Error("expected error")
	}

	// rec1 should still have stopped
	if !rec1.stopped {
		t.Error("rec1 should be stopped despite rec2 error")
	}
}

func TestManagerHealthCheckAll(t *testing.T) {
	m := NewManager()
	rec1 := &mockRecorder{name: "rec1"}
	rec2 := &mockRecorderWithError{
		mockRecorder: mockRecorder{name: "rec2"},
		healthErr:    errors.New("unhealthy"),
	}

	m.Add(rec1)
	m.Add(rec2)

	ctx := context.Background()
	statuses := m.HealthCheckAll(ctx)

	if len(statuses) != 2 {
		t.Fatalf("expected 2 statuses, got %d", len(statuses))
	}

	// rec1 should be healthy
	if statuses[0].Name != "rec1" {
		t.Errorf("expected name %q, got %q", "rec1", statuses[0].Name)
	}
	if !statuses[0].Healthy {
		t.Error("rec1 should be healthy")
	}

	// rec2 should be unhealthy
	if statuses[1].Name != "rec2" {
		t.Errorf("expected name %q, got %q", "rec2", statuses[1].Name)
	}
	if statuses[1].Healthy {
		t.Error("rec2 should be unhealthy")
	}
	if statuses[1].Message != "unhealthy" {
		t.Errorf("expected message %q, got %q", "unhealthy", statuses[1].Message)
	}
}

func TestManagerList(t *testing.T) {
	m := NewManager()
	rec1 := &mockRecorder{name: "rec1"}
	rec2 := &mockRecorder{name: "rec2"}

	m.Add(rec1)
	m.Add(rec2)

	recorders := m.List()
	if len(recorders) != 2 {
		t.Fatalf("expected 2 recorders, got %d", len(recorders))
	}

	// Verify it returns a copy
	recorders[0] = nil
	if m.Count() != 2 {
		t.Error("manager's recorders should not be affected by modifications to returned slice")
	}
}

func TestManagerCount(t *testing.T) {
	m := NewManager()

	if m.Count() != 0 {
		t.Errorf("expected count 0, got %d", m.Count())
	}

	m.Add(&mockRecorder{name: "rec1"})
	if m.Count() != 1 {
		t.Errorf("expected count 1, got %d", m.Count())
	}

	m.Add(&mockRecorder{name: "rec2"})
	if m.Count() != 2 {
		t.Errorf("expected count 2, got %d", m.Count())
	}
}
