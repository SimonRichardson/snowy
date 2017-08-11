// +build integration

package store

import (
	"fmt"
	"sync"
	"testing"
	"testing/quick"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/trussle/snowy/pkg/uuid"
)

const (
	defaultWaitTime = time.Millisecond * 100
)

func TestRealStore_Integration(t *testing.T) {
	// Note: do not run this with Parallel, otherwise it introduces flakey test
	// results.

	config, err := BuildConfig(
		WithHostPort("store", 5432),
		WithUsername("postgres"),
		WithPassword("postgres"),
		WithSSLMode("disable"),
	)
	if err != nil {
		t.Fatal(err)
	}

	runStore := func() Store {
		var wg sync.WaitGroup
		wg.Add(1)

		store := NewRealStore(config, log.NewNopLogger())

		go func() {
			wg.Done()
			if err := store.Run(); err != nil {
				t.Fatal(err)
			}
		}()

		wg.Wait()

		return store
	}

	t.Run("run and stop", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(1)

		store := NewRealStore(config, log.NewNopLogger())
		go func() {
			wg.Done()
			if err := store.Run(); err != nil {
				t.Fatal(err)
			}
		}()

		wg.Wait()

		store.Stop()
	})

	t.Run("get", func(t *testing.T) {
		store := runStore()
		defer store.Stop()

		_, err := store.Get(uuid.New(), Query{})
		if expected, actual := true, err != nil; expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}

		if expected, actual := true, ErrNotFound(err); expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}
	})

	t.Run("insert", func(t *testing.T) {
		store := runStore()
		defer store.Stop()

		fn := func(resourceID uuid.UUID,
			resourceAddress string,
			resourceSize int64,
			resourceContentType, authorID, name string,
			tags Tags,
		) bool {
			defer store.Drop()

			err := store.Insert(Entity{
				ResourceID:          resourceID,
				ResourceAddress:     resourceAddress,
				ResourceSize:        resourceSize,
				ResourceContentType: resourceContentType,
				AuthorID:            authorID,
				Name:                name,
				Tags:                tags.Slice(),
				CreatedOn:           time.Now(),
				DeletedOn:           time.Time{},
			})
			return err == nil
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("insert then get", func(t *testing.T) {
		store := runStore()
		defer store.Stop()

		fn := func(resourceID uuid.UUID,
			resourceAddress string,
			resourceSize int64,
			resourceContentType, authorID, name string,
			tags Tags,
		) bool {
			defer store.Drop()

			if err := store.Insert(Entity{
				ResourceID:          resourceID,
				ResourceAddress:     resourceAddress,
				ResourceSize:        resourceSize,
				ResourceContentType: resourceContentType,
				AuthorID:            authorID,
				Name:                name,
				Tags:                tags.Slice(),
				CreatedOn:           time.Now(),
				DeletedOn:           time.Time{},
			}); err != nil {
				t.Fatal(err)
			}

			entity, err := store.Get(resourceID, Query{})
			if err != nil {
				return false
			}
			return entity.ResourceID.Equals(resourceID) &&
				entity.ResourceAddress == resourceAddress &&
				entity.ResourceSize == resourceSize &&
				entity.ResourceContentType == resourceContentType &&
				entity.AuthorID == authorID
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("insert then query with no tags should select all", func(t *testing.T) {
		store := runStore()
		defer store.Stop()

		fn := func(resourceID uuid.UUID,
			resourceAddress string,
			resourceSize int64,
			resourceContentType, authorID, name string,
			tags Tags,
		) bool {
			defer store.Drop()

			if err := store.Insert(Entity{
				ResourceID:          resourceID,
				ResourceAddress:     resourceAddress,
				ResourceSize:        resourceSize,
				ResourceContentType: resourceContentType,
				AuthorID:            authorID,
				Name:                name,
				Tags:                tags.Slice(),
				CreatedOn:           time.Now(),
				DeletedOn:           time.Time{},
			}); err != nil {
				t.Fatal(err)
			}

			entities, err := store.GetMultiple(resourceID, Query{Tags: make([]string, 0)})
			if err != nil {
				t.Fatal(err)
			}

			return len(entities) == 1
		}

		if err := quick.Check(fn, &quick.Config{MaxCount: 1}); err != nil {
			t.Error(err)
		}
	})

	t.Run("insert then query exact match", func(t *testing.T) {
		store := runStore()
		defer store.Stop()

		fn := func(resourceID uuid.UUID,
			resourceAddress string,
			resourceSize int64,
			resourceContentType, authorID, name string,
			tags Tags,
		) bool {
			defer store.Drop()

			entity := Entity{
				ResourceID:          resourceID,
				ResourceAddress:     resourceAddress,
				ResourceSize:        resourceSize,
				ResourceContentType: resourceContentType,
				AuthorID:            authorID,
				Name:                name,
				Tags:                tags.Slice(),
				CreatedOn:           time.Now().Round(time.Millisecond),
				DeletedOn:           time.Time{},
			}
			if err := store.Insert(entity); err != nil {
				t.Fatal(err)
			}

			got, err := store.GetMultiple(resourceID, Query{
				Tags: tags.Slice(),
			})
			if err != nil {
				t.Fatal(err)
			}

			want := []Entity{entity}
			return equals(want, got)
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("multiple puts then query exact match", func(t *testing.T) {
		store := runStore()
		defer store.Stop()

		fn := func(resourceID uuid.UUID,
			resourceAddress string,
			resourceSize int64,
			resourceContentType, authorID, name string,
			tags Tags,
		) bool {
			defer store.Drop()

			want := make([]Entity, 10)
			for k := range want {
				entity := Entity{
					ResourceID:          resourceID,
					ResourceAddress:     resourceAddress,
					ResourceSize:        resourceSize,
					ResourceContentType: resourceContentType,
					AuthorID:            fmt.Sprintf("%s%d", authorID, k),
					Name:                name,
					Tags:                tags.Slice(),
					CreatedOn:           time.Now().Round(time.Millisecond),
					DeletedOn:           time.Time{},
				}
				if err := store.Insert(entity); err != nil {
					t.Fatal(err)
				}
				want[(len(want)-1)-k] = entity
			}

			got, err := store.GetMultiple(resourceID, Query{
				Tags: tags.Slice(),
			})
			if err != nil {
				t.Fatal(err)
			}

			return equals(want, got)
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("inserts then query partial match", func(t *testing.T) {
		store := runStore()
		defer store.Stop()

		fn := func(resourceID uuid.UUID,
			resourceAddress string,
			resourceSize int64,
			resourceContentType, authorID, name string,
			tags Tags,
		) bool {
			defer store.Drop()

			want := make([]Entity, 10)
			for k := range want {
				entity := Entity{
					ResourceID:          resourceID,
					ResourceAddress:     resourceAddress,
					ResourceSize:        resourceSize,
					ResourceContentType: resourceContentType,
					AuthorID:            fmt.Sprintf("%s%d", authorID, k),
					Name:                name,
					Tags:                tags.Slice(),
					CreatedOn:           time.Now().Round(time.Millisecond),
					DeletedOn:           time.Time{},
				}
				if err := store.Insert(entity); err != nil {
					t.Fatal(err)
				}
				want[(len(want)-1)-k] = entity
			}

			got, err := store.GetMultiple(resourceID, Query{
				Tags: splitTags(tags.Slice()),
			})
			if err != nil {
				t.Fatal(err)
			}

			return equals(want, got)
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})
}
