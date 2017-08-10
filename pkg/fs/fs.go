package fs

import (
	"io"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// Filesystem is an abstraction over the native filesystem that allows us to
// create mock implementations for better testing.
type Filesystem interface {
	// Create takes a path, creates the file and then returns a File back that
	// can be used. This returns an error if the file can not be created in
	// some way.
	Create(string) (File, error)

	// Open takes a path, opens a potential file and then returns a File if
	// that file exists, otherwise it returns an error if the file wasn't found.
	Open(string) (File, error)

	// Rename takes a current destination path and a new destination path and will
	// rename the a File if it exists, otherwise it returns an error if the file
	// wasn't found.
	Rename(string, string) error

	// Exists takes a path and checks to see if the potential file exists or
	// not.
	// Note: If there is an error trying to read that file, it will return false
	// even if the file already exists.
	Exists(string) bool

	// Remote takes a path, removes a potential file, if no file doesn't exist it
	// will return not found.
	Remove(string) error

	// Walk over a specific directory and will return an error if there was an
	// error whilst walking.
	Walk(string, filepath.WalkFunc) error
}

// File is an abstraction for reading, writing and also closing a file. These
// interfaces already exist, it's just a matter of composing them to be more
// usable by other components.
type File interface {
	io.Reader
	io.Writer
	io.Closer

	// Name returns the name of the file
	Name() string

	// Size returns the size of the file
	Size() int64

	// Sync attempts to sync the file with the underlying storage or errors if it
	// can't not succeed.
	Sync() error
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

// Config encapsulates the requirements for generating a Filesystem
type Config struct {
	name         string
	remoteConfig *RemoteConfig
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

// With adds a type of filesystem to use for the configuration.
func With(name string) Option {
	return func(config *Config) error {
		config.name = name
		return nil
	}
}

// WithConfig adds a remote filesystem config to the configuration
func WithConfig(remoteConfig *RemoteConfig) Option {
	return func(config *Config) error {
		config.remoteConfig = remoteConfig
		return nil
	}
}

// New creates a filesystem from a configuration or returns error if on failure.
func New(config *Config) (fsys Filesystem, err error) {
	switch strings.ToLower(config.name) {
	case "remote":
		fsys, err = NewRemoteFilesystem(config.remoteConfig)
		if err != nil {
			err = errors.Wrap(err, "remote filesystem")
			return
		}
	case "local":
		fsys = NewLocalFilesystem()
	case "virtual":
		fsys = NewVirtualFilesystem()
	case "nop":
		fsys = NewNopFilesystem()
	default:
		err = errors.Errorf("unexpected fs type %q", config.name)
	}
	return
}
