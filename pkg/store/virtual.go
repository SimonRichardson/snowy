package store

import (
	"sort"
	"sync"

	"github.com/pkg/errors"
	"github.com/trussle/uuid"
)

// virtualStore keeps track of a entity objects.
type virtualStore struct {
	mutex    sync.RWMutex
	entities map[string][]Entity
	links    map[string]Entity
	stop     chan chan struct{}
}

// NewVirtualStore creates a new Store with the correct dependencies
func NewVirtualStore() Store {
	return &virtualStore{
		mutex:    sync.RWMutex{},
		entities: make(map[string][]Entity),
		links:    make(map[string]Entity),
		stop:     make(chan chan struct{}),
	}
}

func (r *virtualStore) Select(resourceID uuid.UUID, query Query) (Entity, error) {
	entities, err := r.SelectRevisions(resourceID, query)
	if err != nil {
		return Entity{}, err
	}
	if len(entities) == 0 {
		return Entity{}, errNotFound{errors.New("not found")}
	}
	return entities[0], nil
}

func (r *virtualStore) Insert(entity Entity) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Normalize the tags of the entity
	entity.Tags = sortTags(entity.Tags)

	id := entity.ResourceID.String()
	r.entities[id] = append(r.entities[id], entity)
	r.links[entity.ID.String()] = entity
	return nil
}

func (r *virtualStore) SelectRevisions(resourceID uuid.UUID, query Query) ([]Entity, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if entities, ok := r.entities[resourceID.String()]; ok {

		// Sort all the entities first
		sort.Slice(entities, func(a, b int) bool {
			return entities[a].CreatedOn.Before(entities[b].CreatedOn)
		})

		// Filter by authorID before filtering by tags
		if query.AuthorID != nil && *query.AuthorID != "" {
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

		return res, nil
	}
	return make([]Entity, 0), nil
}

func (r *virtualStore) SelectForkRevisions(resourceID uuid.UUID) ([]Entity, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if entities, ok := r.entities[resourceID.String()]; ok {
		if len(entities) == 0 {
			return entities, nil
		}

		// Sort all the entities first
		sort.Slice(entities, func(a, b int) bool {
			return entities[a].CreatedOn.Before(entities[b].CreatedOn)
		})

		var (
			last = entities[len(entities)-1]
			res  = []Entity{last}
		)
		for {
			// We're at the root
			if last.ParentID.Zero() {
				break
			}

			if entity, ok := r.links[last.ParentID.String()]; ok {
				last = entity
				res = append(res, last)
				continue
			}

			// Dead link, break out
			break
		}

		sort.Slice(res, func(a, b int) bool {
			return res[a].CreatedOn.Before(res[b].CreatedOn)
		})

		return res, nil
	}
	return make([]Entity, 0), nil
}

func (r *virtualStore) Statistics() (Statistics, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var stats Statistics
	for range r.entities {
		stats.Total++
	}
	return stats, nil
}

// Run manages the store, keeping the store reliable.
func (r *virtualStore) Run() error {
	for {
		select {
		case c := <-r.stop:
			close(c)
			return nil
		}
	}
}

// Stop closes the store and prevents any new actions running on the
// underlying datastore.
func (r *virtualStore) Stop() {
	c := make(chan struct{})
	r.stop <- c
	<-c
}

// Drop removes all of the stored ledgers
func (r *virtualStore) Drop() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.entities = make(map[string][]Entity)
	r.links = make(map[string]Entity)
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
