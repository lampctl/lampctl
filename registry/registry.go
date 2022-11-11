package registry

import (
	"errors"
)

var errInvalidProvider = errors.New("invalid provider specified")

// Registry maintains a list of providers and provides access to them.
type Registry struct {
	providers map[string]Provider
}

// New creates and initializes a new Registry instance.
func New() *Registry {
	return &Registry{}
}

// Register adds a provider to the registry.
func (r *Registry) Register(provider Provider) {
	r.providers[provider.ID()] = provider
}

// Providers returns a slice of all registered providers.
func (r *Registry) Providers() []Provider {
	providers := []Provider{}
	for _, p := range r.providers {
		providers = append(providers, p)
	}
	return providers
}

// GetProvider retrieves the specified provider by its ID.
func (r *Registry) GetProvider(id string) (Provider, error) {
	p, ok := r.providers[id]
	if !ok {
		return nil, errInvalidProvider
	}
	return p, nil
}

// Close frees all providers and resources used by the registry.
func (r *Registry) Close() {
	for _, v := range r.providers {
		v.Close()
	}
}
