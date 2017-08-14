package store

import (
	"sort"

	"github.com/pkg/errors"
	"github.com/trussle/snowy/pkg/uuid"
)

// virtualStore keeps track of a entity objects.
type virtualStore struct {
	entities map[string][]Entity
}

// NewVirtualStore creates a new Store with the correct dependencies
func NewVirtualStore() Store {
	return &virtualStore{
		entities: make(map[string][]Entity),
	}
}

func (r *virtualStore) Get(resourceID uuid.UUID, query Query) (Entity, error) {
	entities, err := r.GetMultiple(resourceID, query)
	if err != nil {
		return Entity{}, err
	}
	if len(entities) == 0 {
		return Entity{}, errNotFound{errors.New("not found")}
	}
	return entities[0], nil
}

func (r *virtualStore) Insert(entity Entity) error {
	// Normalize the tags of the entity
	entity.Tags = sortTags(entity.Tags)

	id := entity.ResourceID.String()
	r.entities[id] = append(r.entities[id], entity)
	return nil
}

func (r *virtualStore) GetMultiple(resourceID uuid.UUID, query Query) ([]Entity, error) {
	if entities, ok := r.entities[resourceID.String()]; ok {

		// Filter by authorID before filtering by tags
		if query.AuthorID != nil {
			var (
				filtered []Entity
				authorID = *query.AuthorID
			)

			for _, v := range entities {
				if v.AuthorID == authorID {
					filtered = append(filtered, v)
				}
			}

			entities = filtered
		}

		// Filter by tags
		if len(query.Tags) == 0 {
			return entities, nil
		}

		var res []Entity
		for _, v := range entities {
			if intersection(v.Tags, query.Tags) {
				res = append(res, v)
			}
		}

		sort.Slice(res, func(a, b int) bool {
			return res[b].CreatedOn.Before(res[a].CreatedOn)
		})

		return res, nil
	}
	return make([]Entity, 0), nil
}

// Run manages the store, keeping the store reliable.
func (r *virtualStore) Run() error { return nil }

// Stop closes the store and prevents any new actions running on the
// underlying datastore.
func (r *virtualStore) Stop() {}

// Drop removes all of the stored documents
func (r *virtualStore) Drop() error {
	r.entities = make(map[string][]Entity)
	return nil
}

// intersection checks if there are any values that overlap between slices. It
// returns true if there was.
func intersection(a, b []string) bool {
	m := map[string]struct{}{}
	for _, v := range a {
		m[v] = struct{}{}
	}
	for _, v := range b {
		if _, ok := m[v]; ok {
			return true
		}
	}
	return false
}
