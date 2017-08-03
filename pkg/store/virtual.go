package store

import (
	"github.com/trussle/snowy/pkg/uuid"
	"github.com/pkg/errors"
)

// virtualStore keeps track of a entity objects.
type virtualStore struct {
	entities map[string]Entity
}

// NewVirtualStore creates a new Store with the correct dependencies
func NewVirtualStore() Store {
	return &virtualStore{
		entities: make(map[string]Entity),
	}
}

func (r *virtualStore) Get(resourceID uuid.UUID) (Entity, error) {
	if entity, ok := r.entities[resourceID.String()]; ok {
		return entity, nil
	}
	return Entity{}, errNotFound{errors.New("not found")}
}

func (r *virtualStore) Put(entity Entity) error {
	r.entities[entity.ResourceID.String()] = entity
	return nil
}

// Run manages the store, keeping the store reliable.
func (r *virtualStore) Run() error { return nil }

// Stop closes the store and prevents any new actions running on the
// underlying datastore.
func (r *virtualStore) Stop() {}
