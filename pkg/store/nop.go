package store

import (
	"github.com/trussle/snowy/pkg/uuid"
)

type nop struct{}

// NewNopStore creates a Store that has methods that always succeed,
// but do nothing.
func NewNopStore() Store {
	return nop{}
}

func (nop) Get(resourceID uuid.UUID) (Entity, error) { return Entity{}, nil }
func (nop) Put(Entity) error                         { return nil }
func (nop) Run() error                               { return nil }
func (nop) Stop()                                    {}
