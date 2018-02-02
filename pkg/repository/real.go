package repository

import (
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	"github.com/trussle/fsys"
	"github.com/trussle/snowy/pkg/models"
	"github.com/trussle/snowy/pkg/store"
	"github.com/trussle/uuid"
)

const (
	defaultRootParentID = "00000000-0000-0000-0000-000000000000"
)

type realRepository struct {
	fs     fsys.Filesystem
	store  store.Store
	logger log.Logger
}

// NewRealRepository creates a store that backs on to a real filesystem, with the
// correct dependencies.
func NewRealRepository(fs fsys.Filesystem, store store.Store, logger log.Logger) Repository {
	return &realRepository{
		fs:     fs,
		store:  store,
		logger: logger,
	}
}

// SelectLedger returns a Ledger corresponding to the resource ID. If no
// ledger exists it will return an error.
func (r *realRepository) SelectLedger(resourceID uuid.UUID, options Query) (models.Ledger, error) {
	query, err := store.BuildQuery(
		store.WithQueryTags(options.Tags),
		store.WithQueryAuthorID(options.AuthorID),
	)
	if err != nil {
		return models.Ledger{}, err
	}

	entity, err := r.store.Select(resourceID, query)
	if err != nil {
		if store.ErrNotFound(err) {
			return models.Ledger{}, errNotFound{err}
		}
		return models.Ledger{}, err
	}

	return models.BuildLedger(
		models.WithID(entity.ID),
		models.WithParentID(entity.ParentID),
		models.WithName(entity.Name),
		models.WithResourceID(entity.ResourceID),
		models.WithResourceAddress(entity.ResourceAddress),
		models.WithResourceSize(entity.ResourceSize),
		models.WithResourceContentType(entity.ResourceContentType),
		models.WithAuthorID(entity.AuthorID),
		models.WithTags(entity.Tags),
		models.WithCreatedOn(entity.CreatedOn),
		models.WithDeletedOn(entity.DeletedOn),
	)
}

// InsertLedger inserts ledger into the repository. If there is an error
// putting ledgers into the repository then it will return an error.
func (r *realRepository) InsertLedger(doc models.Ledger) (models.Ledger, error) {
	parentID, err := uuid.Parse(defaultRootParentID)
	if err != nil {
		return models.Ledger{}, err
	}

	return r.insertLedgerWithParentID(doc, parentID)
}

// AppendLedger adds a new ledger as a revision. If there is no head
// ledger, it will return an error. If there is an error appending
// ledgers into the repository then it will return an error.
func (r *realRepository) AppendLedger(resourceID uuid.UUID, doc models.Ledger) (models.Ledger, error) {
	// We don't care what we get back, just that it exists.
	entity, err := r.SelectLedger(resourceID, Query{})
	if err != nil {
		return models.Ledger{}, err
	}

	return r.insertLedgerWithParentID(doc, entity.ID())
}

// ForkLedger adds a new ledger as a revision. If there is no head
// ledger, it will return an error. If there is an error appending
// ledgers into the repository then it will return an error.
func (r *realRepository) ForkLedger(resourceID uuid.UUID, doc models.Ledger) (models.Ledger, error) {
	// We don't care what we get back, just that it exists.
	entity, err := r.SelectLedger(resourceID, Query{})
	if err != nil {
		return models.Ledger{}, err
	}

	return r.insertLedgerWithParentID(doc, entity.ID())
}

func (r *realRepository) insertLedgerWithParentID(doc models.Ledger, parentID uuid.UUID) (models.Ledger, error) {
	entity, err := store.BuildEntity(
		store.WithParentID(parentID),
		store.WithName(doc.Name()),
		store.WithResourceID(doc.ResourceID()),
		store.WithResourceAddress(doc.ResourceAddress()),
		store.WithResourceSize(doc.ResourceSize()),
		store.WithResourceContentType(doc.ResourceContentType()),
		store.WithAuthorID(doc.AuthorID()),
		store.WithTags(doc.Tags()),
		store.WithCreatedOn(doc.CreatedOn()),
		store.WithDeletedOn(time.Time{}),
	)
	if err != nil {
		return models.Ledger{}, err
	}

	if err = r.store.Insert(entity); err != nil {
		return models.Ledger{}, err
	}

	// Reconstruct the models.
	return models.BuildLedger(
		models.WithID(entity.ID),
		models.WithParentID(entity.ParentID),
		models.WithName(entity.Name),
		models.WithResourceID(entity.ResourceID),
		models.WithResourceAddress(entity.ResourceAddress),
		models.WithResourceSize(entity.ResourceSize),
		models.WithResourceContentType(entity.ResourceContentType),
		models.WithAuthorID(entity.AuthorID),
		models.WithTags(entity.Tags),
		models.WithCreatedOn(entity.CreatedOn),
		models.WithDeletedOn(entity.DeletedOn),
	)
}

// SelectLedgers returns a set of Ledgers corresponding to a resourceID,
// with some additional qualifiers. If no ledgers are found it will return
// an empty slice. If there is an error parsing the ledgers then it will
// return an error.
func (r *realRepository) SelectLedgers(resourceID uuid.UUID, options Query) ([]models.Ledger, error) {
	query, err := store.BuildQuery(
		store.WithQueryTags(options.Tags),
		store.WithQueryAuthorID(options.AuthorID),
	)
	if err != nil {
		return nil, err
	}

	entities, err := r.store.SelectRevisions(resourceID, query)
	if err != nil {
		return nil, err
	}

	res := make([]models.Ledger, len(entities))
	for k, entity := range entities {
		doc, err := models.BuildLedger(
			models.WithID(entity.ID),
			models.WithParentID(entity.ParentID),
			models.WithName(entity.Name),
			models.WithResourceID(entity.ResourceID),
			models.WithResourceAddress(entity.ResourceAddress),
			models.WithResourceSize(entity.ResourceSize),
			models.WithResourceContentType(entity.ResourceContentType),
			models.WithAuthorID(entity.AuthorID),
			models.WithTags(entity.Tags),
			models.WithCreatedOn(entity.CreatedOn),
			models.WithDeletedOn(entity.DeletedOn),
		)
		if err != nil {
			return nil, err
		}

		res[k] = doc
	}

	return res, nil
}

// SelectForkLedgers returns a set of Ledgers corresponding to a resourceID,
// with some additional qualifiers. If no ledgers are found it will return
// an empty slice. If there is an error parsing the ledgers then it will
// return an error.
func (r *realRepository) SelectForkLedgers(resourceID uuid.UUID) ([]models.Ledger, error) {
	entities, err := r.store.SelectForkRevisions(resourceID)
	if err != nil {
		return nil, err
	}

	res := make([]models.Ledger, len(entities))
	for k, entity := range entities {
		doc, err := models.BuildLedger(
			models.WithID(entity.ID),
			models.WithParentID(entity.ParentID),
			models.WithName(entity.Name),
			models.WithResourceID(entity.ResourceID),
			models.WithResourceAddress(entity.ResourceAddress),
			models.WithResourceSize(entity.ResourceSize),
			models.WithResourceContentType(entity.ResourceContentType),
			models.WithAuthorID(entity.AuthorID),
			models.WithTags(entity.Tags),
			models.WithCreatedOn(entity.CreatedOn),
			models.WithDeletedOn(entity.DeletedOn),
		)
		if err != nil {
			return nil, err
		}

		res[k] = doc
	}

	return res, nil
}

func (r *realRepository) LedgerStatistics() (models.LedgerStatistics, error) {
	stats, err := r.store.Statistics()
	if err != nil {
		return models.LedgerStatistics{}, err
	}

	return models.LedgerStatistics{
		TotalLedgers: stats.Total,
	}, nil
}

// PutContent inserts content into the repository, this will make sure that
// links to the content is managed by the ledger storage. If there is an error
// during the saving of the content to the underlying storage it will then
// return an error.
func (r *realRepository) SelectContent(resourceID uuid.UUID, options Query) (content models.Content, err error) {
	var doc models.Ledger
	doc, err = r.SelectLedger(resourceID, options)
	if err != nil {
		level.Error(r.logger).Log("action", "content", "case", "get", "err", err.Error())
		return
	}

	var file fsys.File
	file, err = r.fs.Open(doc.ResourceAddress())
	if err != nil {
		level.Error(r.logger).Log("action", "content", "case", "open", "err", err.Error(), "resource", doc.ResourceAddress())
		if fsys.ErrNotFound(err) {
			err = errNotFound{err}
			return
		}
		return
	}

	return models.BuildContent(
		models.WithAddress(doc.ResourceAddress()),
		models.WithSize(file.Size()),
		models.WithContentType(doc.ResourceContentType()),
		models.WithReader(file),
	)
}

// PutContent inserts content into the repository. If there is an error
// putting content into the repository then it will return an error.
func (r *realRepository) PutContent(content models.Content) (res models.Content, err error) {
	var bytes []byte
	bytes, err = content.Bytes()
	if err != nil {
		return
	}

	if len(bytes) < 1 {
		err = errors.Errorf("no content")
		return
	}

	// Content already exists, return out quickly.
	if r.fs.Exists(content.Address()) {
		res = content
		return
	}

	var file fsys.File
	file, err = r.fs.Create(content.Address())
	if err != nil {
		return
	}

	if err = file.SetContentType(content.ContentType()); err != nil {
		return
	}

	if _, err = file.Write(bytes); err != nil {
		return
	}

	if err = file.Sync(); err != nil {
		return
	}

	return content, nil
}

// PutContent inserts content into the repository, this will make sure that
// links to the content is managed by the ledger storage. If there is an error
// during the saving of the content to the underlying storage it will then
// return an error.
func (r *realRepository) SelectContents(resourceID uuid.UUID, options Query) (contents []models.Content, err error) {
	var docs []models.Ledger
	docs, err = r.SelectLedgers(resourceID, options)
	if err != nil {
		return
	}
	if len(docs) == 0 {
		contents = make([]models.Content, 0)
		return
	}

	var (
		notFound      = make(chan struct{})
		internalError = make(chan error)
		result        = make(chan models.Content)
	)
	for _, k := range docs {
		go func(doc models.Ledger) {
			var file fsys.File
			file, err = r.fs.Open(doc.ResourceAddress())
			if err != nil {
				if fsys.ErrNotFound(err) {
					notFound <- struct{}{}
					return
				}
				internalError <- err
				return
			}

			content, err := models.BuildContent(
				models.WithAddress(doc.ResourceAddress()),
				models.WithSize(file.Size()),
				models.WithContentType(doc.ResourceContentType()),
				models.WithReader(file),
			)
			if err != nil {
				internalError <- err
				return
			}
			result <- content
		}(k)
	}

	var res []models.Content
	for i := 0; i < len(docs); i++ {
		select {
		case <-notFound:
			continue
		case err := <-internalError:
			return nil, err
		case content := <-result:
			res = append(res, content)
		}
	}
	return res, nil
}

// Close the underlying ledger store and returns an error if it fails.
func (r *realRepository) Close() error {
	return nil
}
