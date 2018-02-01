package store

import (
	"github.com/trussle/uuid"
)

type nop struct{}

// NewNopStore creates a Store that has methods that always succeed,
// but do nothing.
func NewNopStore() Store {
	return nop{}
}

func (nop) Select(resourceID uuid.UUID, query Query) (Entity, error) { return Entity{}, nil }
func (nop) Insert(entity Entity) error                               { return nil }
func (nop) SelectRevisions(resourceID uuid.UUID, query Query) ([]Entity, error) {
	return make([]Entity, 0), nil
}
func (nop) SelectForkRevisions(resourceID uuid.UUID) ([]Entity, error) {
	return make([]Entity, 0), nil
}
func (nop) Run() error  { return nil }
func (nop) Stop()       {}
func (nop) Drop() error { return nil }
