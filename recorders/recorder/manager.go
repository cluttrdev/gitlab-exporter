package recorder

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// Manager manages the lifecycle of multiple recorder instances
type Manager struct {
	recorders []Recorder
	mu        sync.RWMutex
}

// NewManager creates a new recorder manager
func NewManager() *Manager {
	return &Manager{
		recorders: make([]Recorder, 0),
	}
}

// Add adds a recorder to be managed
func (m *Manager) Add(r Recorder) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.recorders = append(m.recorders, r)
}

// StartAll starts all managed recorders
func (m *Manager) StartAll(ctx context.Context) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var errs error
	for _, r := range m.recorders {
		if err := r.Start(ctx); err != nil {
			errs = errors.Join(errs, fmt.Errorf("failed to start recorder %q: %w", r.Name(), err))
		}
	}

	return errs
}

// StopAll stops all managed recorders
func (m *Manager) StopAll(ctx context.Context) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var errs error
	for _, r := range m.recorders {
		if err := r.Stop(ctx); err != nil {
			errs = errors.Join(errs, fmt.Errorf("failed to stop recorder %q: %w", r.Name(), err))
		}
	}

	return errs
}

// HealthCheckAll checks the health of all managed recorders
func (m *Manager) HealthCheckAll(ctx context.Context) []HealthStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()

	statuses := make([]HealthStatus, 0, len(m.recorders))
	now := time.Now()

	for _, r := range m.recorders {
		status := HealthStatus{
			Name:      r.Name(),
			Healthy:   true,
			CheckedAt: now,
		}

		if err := r.CheckHealth(ctx); err != nil {
			status.Healthy = false
			status.Message = err.Error()
		}

		statuses = append(statuses, status)
	}

	return statuses
}

// List returns all managed recorders
func (m *Manager) List() []Recorder {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Return a copy to prevent external modification
	recorders := make([]Recorder, len(m.recorders))
	copy(recorders, m.recorders)
	return recorders
}

// Count returns the number of managed recorders
func (m *Manager) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.recorders)
}
