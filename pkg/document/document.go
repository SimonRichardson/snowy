package document

import (
	"time"

	"github.com/pkg/errors"
	"github.com/trussle/snowy/pkg/uuid"
)

// Document encapsulates all values that are required to represent a document of
// the system.
type Document interface {

	// ID returns the id of the document resource, this is the unique identifier
	// for each document.
	ID() string

	// ResourceID returns the id associated with the document resource.
	ResourceID() uuid.UUID

	// AuthorID returns the id associated with the document resource.
	AuthorID() uuid.UUID

	// Name returns the name of the document
	Name() string

	// Tags returns the associated tags that categorize the document.
	Tags() []string

	CreatedOn() time.Time

	DeletedOn() time.Time
}

// Option defines a option for generating a document
type Option func(Document) error

// Build ingests configuration options to then yield a Document and returns a
// error if it fails during setup.
func Build(opts ...Option) (Document, error) {
	var doc realDocument
	for _, opt := range opts {
		err := opt(&doc)
		if err != nil {
			return nil, err
		}
	}
	return &doc, nil
}

// WithID adds a ID to the document
func WithID(id string) Option {
	return real(func(doc *realDocument) error {
		doc.id = id
		return nil
	})
}

// WithName adds a Name to the document
func WithName(name string) Option {
	return real(func(doc *realDocument) error {
		doc.name = name
		return nil
	})
}

// WithResourceID adds a ResourceID to the document
func WithResourceID(resourceID uuid.UUID) Option {
	return real(func(doc *realDocument) error {
		doc.resourceID = resourceID
		return nil
	})
}

// WithAuthorID adds a AuthorID to the document
func WithAuthorID(authorID uuid.UUID) Option {
	return real(func(doc *realDocument) error {
		doc.authorID = authorID
		return nil
	})
}

// WithTags adds a Tags to the document
func WithTags(tags []string) Option {
	return real(func(doc *realDocument) error {
		doc.tags = tags
		return nil
	})
}

// WithCreatedOn adds a CreatedOn to the document
func WithCreatedOn(createdOn time.Time) Option {
	return real(func(doc *realDocument) error {
		doc.createdOn = createdOn
		return nil
	})
}

// WithDeletedOn adds a DeletedOn to the document
func WithDeletedOn(deletedOn time.Time) Option {
	return real(func(doc *realDocument) error {
		doc.deletedOn = deletedOn
		return nil
	})
}

func real(fn func(*realDocument) error) Option {
	return func(doc Document) error {
		if d, ok := doc.(*realDocument); ok {
			return fn(d)
		}
		return errors.New("invalid document")
	}
}
