// Package providers defines the interface for opnlab data providers.
package providers

import (
	"errors"
)

// ErrActionNotFound is returned when an action is not found.
var ErrActionNotFound = errors.New("action not found")

// Action represents an executable action on a provider.
type Action struct {
	Name        string
	Description string
	Execute     func(params map[string]string) error
}

// Provider is the interface that all data providers must implement.
type Provider interface {
	// Name returns the provider's unique name.
	Name() string

	// Collect gathers data from the provider.
	// Returns a map of data and any error encountered.
	Collect() (map[string]interface{}, error)

	// Actions returns a map of available actions for this provider.
	Actions() map[string]Action
}

// BaseProvider provides a base implementation of Provider.
type BaseProvider struct {
	name string
}

// NewBaseProvider creates a new BaseProvider with the given name.
func NewBaseProvider(name string) *BaseProvider {
	return &BaseProvider{name: name}
}

// Name returns the provider's name.
func (p *BaseProvider) Name() string {
	return p.name
}

// Actions returns an empty map of actions.
func (p *BaseProvider) Actions() map[string]Action {
	return make(map[string]Action)
}
