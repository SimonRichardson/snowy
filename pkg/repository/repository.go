package repository

import (
	"github.com/trussle/snowy/pkg/document"
	"github.com/trussle/snowy/pkg/uuid"
)

// Query allows you to specify different qualifiers when querying the
// repository
type Query struct {
	Tags []string
}

// Repository is an abstraction over the underlying persistence storage, that
// provides a highlevel interface for simple interaction.
type Repository interface {

	// GetDocument returns a Document corresponding to resourceID. If no document
	// exists it will return an error.
	GetDocument(resourceID uuid.UUID, options Query) (document.Document, error)

	// InsertDocument inserts documents into the repository. If there is an
	// error inserting documents into the repository then it will return an
	// error.
	InsertDocument(doc document.Document) (document.Document, error)

	// AppendDocument adds a new document as a revision. If there is no head
	// document, it will return an error. If there is an error appending
	// documents into the repository then it will return an error.
	AppendDocument(resourceID uuid.UUID, doc document.Document) (document.Document, error)

	// GetDocuments returns a set of Documents corresponding to a resourceID,
	// with some additional qualifiers. If no documents are found it will return
	// an empty slice. If there is an error parsing the documents then it will
	// return an error.
	GetDocuments(resourceID uuid.UUID, options Query) ([]document.Document, error)

	// GetContent returns a content corresponding to the resourceID. If no
	// document or content exists, it will return an error.
	GetContent(resourceID uuid.UUID) (document.Content, error)

	// PutContent inserts content into the repository. If there is an error
	// putting content into the repository then it will return an error.
	PutContent(content document.Content) (document.Content, error)

	// Close the underlying document store and returns an error if it fails.
	Close() error
}

// QueryOption defines a option for generating a filesystem Query
type QueryOption func(*Query) error

// BuildQuery ingests configuration options to then yield a Query and return an
// error if it fails during setup.
func BuildQuery(opts ...QueryOption) (Query, error) {
	var config Query
	for _, opt := range opts {
		err := opt(&config)
		if err != nil {
			return Query{}, err
		}
	}
	return config, nil
}

// WithQueryTags adds tags to the Query to use for the configuration.
func WithQueryTags(tags []string) QueryOption {
	return func(query *Query) error {
		query.Tags = tags
		return nil
	}
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
