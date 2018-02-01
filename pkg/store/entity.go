package store

import (
	"time"

	"github.com/trussle/uuid"
)

// Entity represents a value with in the persistent store, that allows us to
// formally understand the underlying model. Entity in this case represents a
// models.Ledger without the file content.
type Entity struct {
	ID, ParentID         uuid.UUID
	Name                 string
	ResourceID           uuid.UUID
	ResourceAddress      string
	ResourceSize         int64
	ResourceContentType  string
	AuthorID             string
	Tags                 []string
	CreatedOn, DeletedOn time.Time
}

// EntityOption defines a option for generating a entity
type EntityOption func(*Entity) error

// BuildEntity ingests configuration options to then yield a Config and return
// an error if it fails during setup.
func BuildEntity(opts ...EntityOption) (Entity, error) {
	var entity Entity
	for _, opt := range opts {
		err := opt(&entity)
		if err != nil {
			return Entity{}, err
		}
	}
	return entity, nil
}

// WithID adds a type of id to the entity.
func WithID(id uuid.UUID) EntityOption {
	return func(entity *Entity) error {
		entity.ID = id
		return nil
	}
}

// WithParentID adds a type of parent id to the entity.
func WithParentID(parentID uuid.UUID) EntityOption {
	return func(entity *Entity) error {
		entity.ParentID = parentID
		return nil
	}
}

// WithName adds a type of name to the entity.
func WithName(name string) EntityOption {
	return func(entity *Entity) error {
		entity.Name = name
		return nil
	}
}

// WithResourceID adds a type of resourceID to the entity.
func WithResourceID(resourceID uuid.UUID) EntityOption {
	return func(entity *Entity) error {
		entity.ResourceID = resourceID
		return nil
	}
}

// WithResourceAddress adds a type of resourceAddress to the entity.
func WithResourceAddress(resourceAddress string) EntityOption {
	return func(entity *Entity) error {
		entity.ResourceAddress = resourceAddress
		return nil
	}
}

// WithResourceSize adds a type of resourceSize to the entity.
func WithResourceSize(resourceSize int64) EntityOption {
	return func(entity *Entity) error {
		entity.ResourceSize = resourceSize
		return nil
	}
}

// WithResourceContentType adds a type of resourceContentType to the entity.
func WithResourceContentType(resourceContentType string) EntityOption {
	return func(entity *Entity) error {
		entity.ResourceContentType = resourceContentType
		return nil
	}
}

// WithAuthorID adds a type of authorID to the entity.
func WithAuthorID(authorID string) EntityOption {
	return func(entity *Entity) error {
		entity.AuthorID = authorID
		return nil
	}
}

// WithTags adds a type of tags to the entity.
func WithTags(tags []string) EntityOption {
	return func(entity *Entity) error {
		entity.Tags = tags
		return nil
	}
}

// WithCreatedOn adds a type of createdOn to the entity.
func WithCreatedOn(createdOn time.Time) EntityOption {
	return func(entity *Entity) error {
		entity.CreatedOn = createdOn
		return nil
	}
}

// WithDeletedOn adds a type of deletedOn to the entity.
func WithDeletedOn(deletedOn time.Time) EntityOption {
	return func(entity *Entity) error {
		entity.DeletedOn = deletedOn
		return nil
	}
}
