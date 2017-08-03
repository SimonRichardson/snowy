package store

import (
	"strings"

	"github.com/trussle/snowy/pkg/uuid"
	"github.com/go-kit/kit/log"
)

// Store represents a API over a persistent store.
type Store interface {

	// Get returns a stored document from the datastore, minus the actual content
	Get(resourceID uuid.UUID) (Entity, error)

	// Put inserts a entity with in the datastore.
	Put(Entity) error

	// Run manages the store, keeping the store reliable.
	Run() error

	// Stop closes the store and prevents any new actions running on the
	// underlying datastore.
	Stop()
}

// Config encapsulates the requirements for generating a Filesystem
type Config struct {
	name       string
	realConfig *RealConfig
}

// Option defines a option for generating a filesystem Config
type Option func(*Config) error

// Build ingests configuration options to then yield a Config and return an
// error if it fails during setup.
func Build(opts ...Option) (*Config, error) {
	var config Config
	for _, opt := range opts {
		err := opt(&config)
		if err != nil {
			return nil, err
		}
	}
	return &config, nil
}

// With adds a type of store to use for the configuration.
func With(name string) Option {
	return func(config *Config) error {
		config.name = name
		return nil
	}
}

// WithConfig adds a remote store config to the configuration
func WithConfig(realConfig *RealConfig) Option {
	return func(config *Config) error {
		config.realConfig = realConfig
		return nil
	}
}

// New creates a store from a configuration or returns error if on failure.
func New(config *Config, logger log.Logger) (store Store, err error) {
	switch strings.ToLower(config.name) {
	case "real":
		store = NewRealStore(config.realConfig, logger)
	case "virtual":
		store = NewVirtualStore()
	case "nop":
		store = NewNopStore()
	}
	return
}

type notFound interface {
	NotFound() bool
}

type errNotFound struct {
	err error
}

func (e errNotFound) Error() string {
	return e.err.Error()
}

func (e errNotFound) NotFound() bool {
	return true
}

// ErrNotFound tests to see if the error passed is a not found error or not.
func ErrNotFound(err error) bool {
	if err != nil {
		if _, ok := err.(notFound); ok {
			return true
		}
	}
	return false
}
