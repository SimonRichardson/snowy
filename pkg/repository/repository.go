package repository

import (
	"github.com/trussle/snowy/pkg/document"
	"github.com/trussle/snowy/pkg/uuid"
)

// Repository is an abstraction over the underlying persistence storage, that
// provides a highlevel interface for simple interaction.
type Repository interface {

	// GetDocument returns a Document corresponding to resourceID. If no document
	// exists it will return an error.
	GetDocument(resourceID uuid.UUID) (document.Document, error)

	// PutDocument inserts documents into the repository. If there is an error
	// putting documents into the repository then it will return an error.
	PutDocument(doc document.Document) (uuid.UUID, error)

	// GetContent returns a content corresponding to the resourceID. If no
	// document or content exists, it will return an error.
	GetContent(resourceID uuid.UUID) (document.Content, error)

	// PutContent inserts content into the repository. If there is an error
	// putting content into the repository then it will return an error.
	PutContent(content document.Content) (uuid.UUID, error)

	// Close the underlying document store and returns an error if it fails.
	Close() error
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
