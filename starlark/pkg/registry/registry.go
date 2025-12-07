package registry

import (
	"fmt"
	"sync"

	"go.starlark.net/starlark"
)

// Module represents a Starlark module that can be registered
type Module interface {
	// Metadata returns the module's metadata
	Metadata() ModuleMetadata

	// Build constructs the Starlark module value
	// This is called each time a new script execution environment is created
	Build() starlark.Value
}

// Registry manages all registered modules
type Registry struct {
	modules map[string]Module
	mu      sync.RWMutex
}

// DefaultRegistry is the global registry instance
var DefaultRegistry = NewRegistry()

// NewRegistry creates a new Registry instance
func NewRegistry() *Registry {
	return &Registry{
		modules: make(map[string]Module),
	}
}

// Register adds a module to the registry
func (r *Registry) Register(module Module) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	meta := module.Metadata()
	if meta.Name == "" {
		return fmt.Errorf("module name cannot be empty")
	}

	if _, exists := r.modules[meta.Name]; exists {
		return fmt.Errorf("module %s already registered", meta.Name)
	}

	r.modules[meta.Name] = module
	return nil
}

// Get retrieves a module by name
func (r *Registry) Get(name string) (Module, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	module, ok := r.modules[name]
	if !ok {
		return nil, fmt.Errorf("module %s not found", name)
	}

	return module, nil
}

// All returns all registered modules
func (r *Registry) All() []Module {
	r.mu.RLock()
	defer r.mu.RUnlock()

	modules := make([]Module, 0, len(r.modules))
	for _, mod := range r.modules {
		modules = append(modules, mod)
	}
	return modules
}

// BuildGlobals creates a Starlark globals dict from all registered modules
func (r *Registry) BuildGlobals() starlark.StringDict {
	r.mu.RLock()
	defer r.mu.RUnlock()

	globals := make(starlark.StringDict)
	for name, module := range r.modules {
		globals[name] = module.Build()
	}

	return globals
}
