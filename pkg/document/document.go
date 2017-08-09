package document

import (
	"encoding/json"
	"time"

	"github.com/trussle/snowy/pkg/uuid"
)

// Document encapsulates all values that are required to represent a document of
// the system.
type Document struct {
	id                   uuid.UUID
	name                 string
	resourceID           uuid.UUID
	authorID             string
	tags                 []string
	createdOn, deletedOn time.Time
}

// ID returns the id of the document resource, this is the unique identifier
// for each document.
func (d Document) ID() uuid.UUID {
	return d.id
}

// ResourceID returns the id associated with the document resource.
func (d Document) ResourceID() uuid.UUID {
	return d.resourceID
}

// AuthorID returns the id associated with the document resource.
func (d Document) AuthorID() string {
	return d.authorID
}

// Name returns the name of the document
func (d Document) Name() string {
	return d.name
}

// Tags returns the associated tags that categorize the document.
func (d Document) Tags() []string {
	return d.tags
}

// CreatedOn returns the time of creation for the document
func (d Document) CreatedOn() time.Time {
	return d.createdOn
}

// DeletedOn returns the time of deletion for the document
func (d Document) DeletedOn() time.Time {
	return d.deletedOn
}

// MarshalJSON converts a UUID into a serialisable json format
func (d Document) MarshalJSON() ([]byte, error) {
	tags := d.tags
	if d.tags == nil {
		tags = make([]string, 0)
	}

	return json.Marshal(struct {
		Name       string    `json:"name"`
		ResourceID uuid.UUID `json:"resource_id"`
		AuthorID   string    `json:"author_id"`
		Tags       []string  `json:"tags"`
		CreatedOn  string    `json:"created_on"`
		DeletedOn  string    `json:"deleted_on"`
	}{
		Name:       d.name,
		ResourceID: d.resourceID,
		AuthorID:   d.authorID,
		Tags:       tags,
		CreatedOn:  d.createdOn.Format(time.RFC3339),
		DeletedOn:  d.deletedOn.Format(time.RFC3339),
	})
}

// UnmarshalJSON unserialises the json format and converts it into a UUID
func (d *Document) UnmarshalJSON(b []byte) error {
	var res struct {
		Name       string    `json:"name"`
		ResourceID uuid.UUID `json:"resource_id"`
		AuthorID   string    `json:"author_id"`
		Tags       []string  `json:"tags"`
		CreatedOn  string    `json:"created_on"`
		DeletedOn  string    `json:"deleted_on"`
	}
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}

	var err error

	d.name = res.Name
	d.resourceID = res.ResourceID
	d.authorID = res.AuthorID
	d.tags = res.Tags

	d.createdOn, err = time.Parse(time.RFC3339, res.CreatedOn)
	if err != nil {
		return err
	}

	d.deletedOn, err = time.Parse(time.RFC3339, res.DeletedOn)
	if err != nil {
		return err
	}

	return nil
}

// Option defines a option for generating a document
type Option func(*Document) error

// Build ingests configuration options to then yield a Document and returns a
// error if it fails during setup.
func Build(opts ...Option) (Document, error) {
	var doc Document
	for _, opt := range opts {
		err := opt(&doc)
		if err != nil {
			return Document{}, err
		}
	}
	return doc, nil
}

// WithID adds a ID to the document
func WithID(id uuid.UUID) Option {
	return func(doc *Document) error {
		doc.id = id
		return nil
	}
}

// WithName adds a Name to the document
func WithName(name string) Option {
	return func(doc *Document) error {
		doc.name = name
		return nil
	}
}

// WithNewResourceID adds a new ResourceID to the document
func WithNewResourceID() Option {
	return func(doc *Document) error {
		doc.resourceID = uuid.New()
		return nil
	}
}

// WithResourceID adds a ResourceID to the document
func WithResourceID(resourceID uuid.UUID) Option {
	return func(doc *Document) error {
		doc.resourceID = resourceID
		return nil
	}
}

// WithAuthorID adds a AuthorID to the document
func WithAuthorID(authorID string) Option {
	return func(doc *Document) error {
		doc.authorID = authorID
		return nil
	}
}

// WithTags adds a Tags to the document
func WithTags(tags []string) Option {
	return func(doc *Document) error {
		doc.tags = tags
		return nil
	}
}

// WithCreatedOn adds a CreatedOn to the document
func WithCreatedOn(createdOn time.Time) Option {
	return func(doc *Document) error {
		doc.createdOn = createdOn
		return nil
	}
}

// WithDeletedOn adds a DeletedOn to the document
func WithDeletedOn(deletedOn time.Time) Option {
	return func(doc *Document) error {
		doc.deletedOn = deletedOn
		return nil
	}
}
