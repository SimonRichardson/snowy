package repository

import (
	"github.com/trussle/snowy/pkg/models"
	"github.com/trussle/uuid"
)

var (
	emptyAuthID = ""
)

// Query allows you to specify different qualifiers when querying the
// repository
type Query struct {
	Tags     []string
	AuthorID *string
}

// Repository is an abstraction over the underlying persistence storage, that
// provides a highlevel interface for simple interaction.
type Repository interface {

	// SelectLedger returns a Ledger corresponding to resourceID. If no ledger
	// exists it will return an error.
	SelectLedger(resourceID uuid.UUID, options Query) (models.Ledger, error)

	// InsertLedger inserts ledgers into the repository. If there is an
	// error inserting ledgers into the repository then it will return an
	// error.
	InsertLedger(doc models.Ledger) (models.Ledger, error)

	// AppendLedger adds a new ledger as a revision. If there is no head
	// ledger, it will return an error. If there is an error appending
	// ledgers into the repository then it will return an error.
	AppendLedger(resourceID uuid.UUID, doc models.Ledger) (models.Ledger, error)

	// ForkLedger adds a new ledger as a branch revision. If there is no head
	// ledger, it will return an error. If there is an error appending
	// ledgers into the repository then it will return an error.
	ForkLedger(resourceID uuid.UUID, doc models.Ledger) (models.Ledger, error)

	// SelectLedgers returns a set of Ledgers corresponding to a resourceID,
	// with some additional qualifiers. If no ledgers are found it will return
	// an empty slice. If there is an error parsing the ledgers then it will
	// return an error.
	SelectLedgers(resourceID uuid.UUID, options Query) ([]models.Ledger, error)

	// SelectForkLedgers adds a new ledger as a branch revision. If there is no
	// head ledger, it will return an error. If there is an error appending
	// ledgers into the repository then it will return an error.
	SelectForkLedgers(resourceID uuid.UUID) ([]models.Ledger, error)

	// LedgerStatistics returns some statistics about the ledgers
	LedgerStatistics() (models.LedgerStatistics, error)

	// SelectContent returns a content corresponding to the resourceID. If no
	// ledger or content exists, it will return an error.
	SelectContent(resourceID uuid.UUID, options Query) (models.Content, error)

	// PutContent inserts content into the repository. If there is an error
	// putting content into the repository then it will return an error.
	PutContent(content models.Content) (models.Content, error)

	// SelectContents returns a set of content corresponding to the resourceID. If no
	// ledger or content exists, it will return an error.
	SelectContents(resourceID uuid.UUID, options Query) ([]models.Content, error)

	// Close the underlying ledger store and returns an error if it fails.
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

// WithQueryAuthorID adds authorID to the Query to use for the configuration.
func WithQueryAuthorID(authorID string) QueryOption {
	return func(query *Query) error {
		if authorID == "" {
			query.AuthorID = &emptyAuthID
		} else {
			query.AuthorID = &authorID
		}
		return nil
	}
}

// BuildEmptyQuery creates a Query with empty values.
func BuildEmptyQuery() Query {
	return Query{
		AuthorID: &emptyAuthID,
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
