package repository

import (
	"fmt"

	"github.com/trussle/snowy/pkg/document"
	"github.com/trussle/snowy/pkg/fs"
	"github.com/trussle/snowy/pkg/store"
	"github.com/trussle/snowy/pkg/uuid"
	"github.com/go-kit/kit/log"
)

type realRepository struct {
	fs     fs.Filesystem
	store  store.Store
	logger log.Logger
}

// NewRealRepository creates a store that backs on to a real filesystem, with the
// correct dependencies.
func NewRealRepository(fs fs.Filesystem, store store.Store, logger log.Logger) Repository {
	return &realRepository{
		fs:     fs,
		store:  store,
		logger: logger,
	}
}

// GetDocument returns a Document corresponding to the resource ID. If no
// document exists it will return an error.
func (r *realRepository) GetDocument(resourceID uuid.UUID) (document.Document, error) {
	entity, err := r.store.Get(resourceID)
	if err != nil {
		if store.ErrNotFound(err) {
			return nil, errNotFound{err}
		}
		return nil, err
	}

	return document.Build(
		document.WithID(entity.ID),
	)
}

// PutDocument inserts document into the repository. If there is an error
// putting documents into the repository then it will return an error.
func (r *realRepository) PutDocument(doc document.Document) (uuid.UUID, error) {
	entity, err := store.BuildEntity(
		store.BuildEntityWithID(doc.ID()),
		store.BuildEntityWithResourceID(doc.ResourceID()),
	)
	if err != nil {
		return uuid.Empty, err
	}

	if err = r.store.Put(entity); err != nil {
		return uuid.Empty, err
	}
	return doc.ResourceID(), nil
}

// PutContent inserts content into the repository, this will make sure that
// links to the content is managed by the document storage. If there is an error
// during the saving of the content to the underlying storage it will then
// return an error.
func (r *realRepository) GetContent(resourceID uuid.UUID) (document.Content, error) {
	document, err := r.GetDocument(resourceID)
	if err != nil {
		return nil, err
	}

	fmt.Println(document)

	return nil, nil
}

// PutContent inserts content into the repository. If there is an error
// putting content into the repository then it will return an error.
func (r *realRepository) PutContent(content document.Content) (uuid.UUID, error) {
	return uuid.Empty, nil
}

// Close the underlying document store and returns an error if it fails.
func (r *realRepository) Close() error {
	return nil
}
