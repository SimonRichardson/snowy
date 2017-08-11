package repository

import (
	"time"

	"github.com/go-kit/kit/log"
	"github.com/trussle/snowy/pkg/document"
	"github.com/trussle/snowy/pkg/fs"
	"github.com/trussle/snowy/pkg/store"
	"github.com/trussle/snowy/pkg/uuid"
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
func (r *realRepository) GetDocument(resourceID uuid.UUID, options Query) (document.Document, error) {
	query, err := store.BuildQuery(
		store.WithQueryTags(options.Tags),
	)
	if err != nil {
		return document.Document{}, err
	}

	entity, err := r.store.Get(resourceID, query)
	if err != nil {
		if store.ErrNotFound(err) {
			return document.Document{}, errNotFound{err}
		}
		return document.Document{}, err
	}

	return document.BuildDocument(
		document.WithID(entity.ID),
		document.WithName(entity.Name),
		document.WithResourceID(entity.ResourceID),
		document.WithResourceAddress(entity.ResourceAddress),
		document.WithResourceSize(entity.ResourceSize),
		document.WithResourceContentType(entity.ResourceContentType),
		document.WithAuthorID(entity.AuthorID),
		document.WithTags(entity.Tags),
		document.WithCreatedOn(entity.CreatedOn),
		document.WithDeletedOn(entity.DeletedOn),
	)
}

// InsertDocument inserts document into the repository. If there is an error
// putting documents into the repository then it will return an error.
func (r *realRepository) InsertDocument(doc document.Document) (document.Document, error) {
	entity, err := store.BuildEntity(
		store.BuildEntityWithName(doc.Name()),
		store.BuildEntityWithResourceID(doc.ResourceID()),
		store.BuildEntityWithResourceAddress(doc.ResourceAddress()),
		store.BuildEntityWithResourceSize(doc.ResourceSize()),
		store.BuildEntityWithResourceContentType(doc.ResourceContentType()),
		store.BuildEntityWithAuthorID(doc.AuthorID()),
		store.BuildEntityWithTags(doc.Tags()),
		store.BuildEntityWithCreatedOn(doc.CreatedOn()),
		store.BuildEntityWithDeletedOn(time.Time{}),
	)
	if err != nil {
		return document.Document{}, err
	}

	if err = r.store.Insert(entity); err != nil {
		return document.Document{}, err
	}

	// Reconstruct the document.
	return document.BuildDocument(
		document.WithID(entity.ID),
		document.WithName(entity.Name),
		document.WithResourceID(entity.ResourceID),
		document.WithResourceAddress(entity.ResourceAddress),
		document.WithResourceSize(entity.ResourceSize),
		document.WithResourceContentType(entity.ResourceContentType),
		document.WithAuthorID(entity.AuthorID),
		document.WithTags(entity.Tags),
		document.WithCreatedOn(entity.CreatedOn),
		document.WithDeletedOn(entity.DeletedOn),
	)
}

// AppendDocument adds a new document as a revision. If there is no head
// document, it will return an error. If there is an error appending
// documents into the repository then it will return an error.
func (r *realRepository) AppendDocument(resourceID uuid.UUID, doc document.Document) (document.Document, error) {
	// We don't care what we get back, just that it exists.
	_, err := r.GetDocument(resourceID, Query{})
	if err != nil {
		return document.Document{}, err
	}

	return r.InsertDocument(doc)
}

// GetDocuments returns a set of Documents corresponding to a resourceID,
// with some additional qualifiers. If no documents are found it will return
// an empty slice. If there is an error parsing the documents then it will
// return an error.
func (r *realRepository) GetDocuments(resourceID uuid.UUID, options Query) ([]document.Document, error) {
	query, err := store.BuildQuery(
		store.WithQueryTags(options.Tags),
	)
	if err != nil {
		return nil, err
	}

	entities, err := r.store.GetMultiple(resourceID, query)
	if err != nil {
		return nil, err
	}

	res := make([]document.Document, len(entities))
	for k, entity := range entities {
		doc, err := document.BuildDocument(
			document.WithID(entity.ID),
			document.WithName(entity.Name),
			document.WithResourceID(entity.ResourceID),
			document.WithResourceAddress(entity.ResourceAddress),
			document.WithResourceSize(entity.ResourceSize),
			document.WithResourceContentType(entity.ResourceContentType),
			document.WithAuthorID(entity.AuthorID),
			document.WithTags(entity.Tags),
			document.WithCreatedOn(entity.CreatedOn),
			document.WithDeletedOn(entity.DeletedOn),
		)
		if err != nil {
			return nil, err
		}

		res[k] = doc
	}

	return res, nil
}

// PutContent inserts content into the repository, this will make sure that
// links to the content is managed by the document storage. If there is an error
// during the saving of the content to the underlying storage it will then
// return an error.
func (r *realRepository) GetContent(resourceID uuid.UUID) (content document.Content, err error) {
	var doc document.Document
	doc, err = r.GetDocument(resourceID, Query{})
	if err != nil {
		return
	}

	var file fs.File
	file, err = r.fs.Open(doc.ResourceAddress())
	if err != nil {
		if fs.ErrNotFound(err) {
			err = errNotFound{err}
			return
		}
		return
	}

	return document.BuildContent(
		document.WithAddress(doc.ResourceAddress()),
		document.WithSize(file.Size()),
		document.WithContentType(doc.ResourceContentType()),
		document.WithReader(file),
	)
}

// PutContent inserts content into the repository. If there is an error
// putting content into the repository then it will return an error.
func (r *realRepository) PutContent(content document.Content) (document.Content, error) {
	return document.Content{}, nil
}

// Close the underlying document store and returns an error if it fails.
func (r *realRepository) Close() error {
	return nil
}
