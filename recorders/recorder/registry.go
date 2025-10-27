package recorder

import (
	"fmt"
	"sync"
)

type Registry struct {
	factories map[string]Factory
	mu        sync.RWMutex
}

func NewRegistry() *Registry {
	return &Registry{
		factories: make(map[string]Factory),
	}
}

// Register registers a recorder factory with the given type name.
func (r *Registry) Register(typeName string, factory Factory) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if factory == nil {
		return fmt.Errorf("recorder factory is nil")
	}

	if _, exists := r.factories[typeName]; exists {
		return fmt.Errorf("recorder: Register called twice for type %q", typeName)
	}

	r.factories[typeName] = factory
	return nil
}

// Get retrieves a recorder factory by type name.
// Returns nil if the type is not registered.
func (r *Registry) Get(typeName string) Factory {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.factories[typeName]
}

// List returns all registered recorder type names.
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	types := make([]string, 0, len(r.factories))
	for typeName := range r.factories {
		types = append(types, typeName)
	}
	return types
}

// Create creates a new recorder instance of the specified type.
// Returns an error if the type is not registered.
func (r *Registry) Create(typeName string) (Recorder, error) {
	factory := r.Get(typeName)
	if factory == nil {
		return nil, fmt.Errorf("recorder: unknown type %q (available: %v)", typeName, r.List())
	}

	return factory(), nil
}
