package recorder

import (
	"context"
	"testing"

	"go.cluttr.dev/gitlab-exporter/protobuf/servicepb"
)

// mockRecorder is a simple mock implementation for testing
type mockRecorder struct {
	servicepb.UnimplementedGitLabExporterServer
	name        string
	initialized bool
	started     bool
	stopped     bool
}

func (m *mockRecorder) Name() string {
	return m.name
}

func (m *mockRecorder) Initialize(ctx context.Context, config []byte) error {
	m.initialized = true
	return nil
}

func (m *mockRecorder) Start(ctx context.Context) error {
	m.started = true
	return nil
}

func (m *mockRecorder) Stop(ctx context.Context) error {
	m.stopped = true
	return nil
}

func (m *mockRecorder) CheckHealth(ctx context.Context) error {
	return nil
}

func TestRegister(t *testing.T) {
	registry := NewRegistry()

	factory := func() Recorder {
		return &mockRecorder{name: "test"}
	}

	err := registry.Register("test", factory)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if registry.Get("test") == nil {
		t.Error("expected factory to be registered")
	}
}

func TestRegisterPanicsOnNil(t *testing.T) {
	registry := NewRegistry()

	err := registry.Register("nil-test", nil)
	if err == nil {
		t.Error("expected panic for nil factory")
	}
}

func TestRegisterPanicsOnDuplicate(t *testing.T) {
	var (
		registry = NewRegistry()
		err      error
	)

	factory := func() Recorder {
		return &mockRecorder{name: "test"}
	}

	err = registry.Register("duplicate", factory)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = registry.Register("duplicate", factory)
	if err == nil {
		t.Error("expected panic for duplicate registration")
	}
}

func TestGet(t *testing.T) {
	registry := NewRegistry()

	factory := func() Recorder {
		return &mockRecorder{name: "test"}
	}

	registry.Register("test", factory)

	got := registry.Get("test")
	if got == nil {
		t.Fatal("expected factory to be found")
	}

	rec := got()
	if rec.Name() != "test" {
		t.Errorf("expected name %q, got %q", "test", rec.Name())
	}
}

func TestGetUnknown(t *testing.T) {
	registry := NewRegistry()

	got := registry.Get("nonexistent")
	if got != nil {
		t.Error("expected nil for unknown type")
	}
}

func TestList(t *testing.T) {
	registry := NewRegistry()

	factory := func() Recorder {
		return &mockRecorder{name: "test"}
	}

	registry.Register("type1", factory)
	registry.Register("type2", factory)

	types := registry.List()
	if len(types) != 2 {
		t.Errorf("expected 2 types, got %d", len(types))
	}

	hasType1 := false
	hasType2 := false
	for _, typ := range types {
		if typ == "type1" {
			hasType1 = true
		}
		if typ == "type2" {
			hasType2 = true
		}
	}

	if !hasType1 || !hasType2 {
		t.Error("expected both type1 and type2 in list")
	}
}

func TestCreate(t *testing.T) {
	registry := NewRegistry()

	factory := func() Recorder {
		return &mockRecorder{name: "test-recorder"}
	}

	registry.Register("test", factory)

	rec, err := registry.Create("test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rec.Name() != "test-recorder" {
		t.Errorf("expected name %q, got %q", "test-recorder", rec.Name())
	}
}

func TestCreateUnknown(t *testing.T) {
	registry := NewRegistry()

	_, err := registry.Create("unknown-type")
	if err == nil {
		t.Error("expected error for unknown type")
	}
}
