package store

import (
	"time"

	"github.com/trussle/snowy/pkg/uuid"
)

// Entity represents a value with in the persistent store, that allows us to
// formally understand the underlying model. Entity in this case represents a
// document.Document without the file content.
type Entity struct {
	ID, Name             string
	ResourceID, AuthorID uuid.UUID
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

// BuildEntityWithID adds a type of id to the entity.
func BuildEntityWithID(id string) EntityOption {
	return func(entity *Entity) error {
		entity.ID = id
		return nil
	}
}

// BuildEntityWithResourceID adds a type of resourceID to the entity.
func BuildEntityWithResourceID(resourceID uuid.UUID) EntityOption {
	return func(entity *Entity) error {
		entity.ResourceID = resourceID
		return nil
	}
}
