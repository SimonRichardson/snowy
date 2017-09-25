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

	t.Run("run and tick", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(1)

		store := NewRealStore(config, log.NewNopLogger())
		real := store.(*realStore)
		real.ticker = time.NewTicker(time.Millisecond)
		go func() {
			go func() {
				time.Sleep(time.Millisecond * 4)
				wg.Done()
			}()

			if err := store.Run(); err != nil {
				t.Fatal(err)
			}
		}()

		wg.Wait()

		store.Stop()

	})

	t.Run("get", func(t *testing.T) {
		store := runStore(config)
		defer store.Stop()

		_, err := store.Select(uuid.New(), Query{})
		if expected, actual := true, err != nil; expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}

		if expected, actual := true, ErrNotFound(err); expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}
	})

	t.Run("insert", func(t *testing.T) {
		store := runStore(config)
		defer store.Stop()

		fn := func(parentID, resourceID uuid.UUID,
			resourceAddress string,
			resourceSize int64,
			resourceContentType, authorID, name string,
			tags Tags,
		) bool {
			defer store.Drop()

			err := store.Insert(Entity{
				ParentID:            parentID,
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
		store := runStore(config)
		defer store.Stop()

		fn := func(parentID, resourceID uuid.UUID,
			resourceAddress string,
			resourceSize int64,
			resourceContentType, authorID, name string,
			tags Tags,
		) bool {
			defer store.Drop()

			if err := store.Insert(Entity{
				ParentID:            parentID,
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

			entity, err := store.Select(resourceID, Query{})
			if err != nil {
				return false
			}
			return entity.ParentID.Equals(parentID) &&
				entity.ResourceID.Equals(resourceID) &&
				entity.ResourceAddress == resourceAddress &&
				entity.ResourceSize == resourceSize &&
				entity.ResourceContentType == resourceContentType &&
				entity.AuthorID == authorID
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("select fork revisions not found failure", func(t *testing.T) {
		store := runStore(config)
		defer store.Stop()

		entities, err := store.SelectForkRevisions(uuid.New())
		if err != nil {
			t.Fatal(err)
		}

		if expected, actual := 0, len(entities); expected != actual {
			t.Errorf("expected: %d, actual: %d", expected, actual)
		}
	})

	t.Run("select fork revisions", func(t *testing.T) {
		store := runStore(config)
		defer store.Stop()

		fn := func() bool {
			var (
				resourceID     = uuid.New()
				firstAuthorID  = uuid.New().String()
				secondAuthorID = uuid.New().String()
			)

			if err := store.Insert(Entity{
				ParentID:            uuid.Empty,
				ResourceID:          uuid.New(),
				ResourceAddress:     "address",
				ResourceSize:        0,
				ResourceContentType: "application/octet-stream",
				AuthorID:            uuid.New().String(),
				Name:                "name",
				Tags:                []string{},
				CreatedOn:           time.Now().Add(-time.Minute),
				DeletedOn:           time.Time{},
			}); err != nil {
				t.Fatal(err)
			}

			if err := store.Insert(Entity{
				ParentID:            uuid.Empty,
				ResourceID:          resourceID,
				ResourceAddress:     "address",
				ResourceSize:        0,
				ResourceContentType: "application/octet-stream",
				AuthorID:            firstAuthorID,
				Name:                "name",
				Tags:                []string{},
				CreatedOn:           time.Now().Add(-time.Minute),
				DeletedOn:           time.Time{},
			}); err != nil {
				t.Fatal(err)
			}

			entity, err := store.Select(resourceID, Query{AuthorID: &firstAuthorID})
			if err != nil {
				t.Fatal(err)
			}

			// Fork
			if err := store.Insert(Entity{
				ParentID:            entity.ID,
				ResourceID:          resourceID,
				ResourceAddress:     "address",
				ResourceSize:        0,
				ResourceContentType: "application/octet-stream",
				AuthorID:            secondAuthorID,
				Name:                "name",
				Tags:                []string{},
				CreatedOn:           time.Now(),
				DeletedOn:           time.Time{},
			}); err != nil {
				t.Fatal(err)
			}

			entities, err := store.SelectForkRevisions(resourceID)
			if err != nil {
				t.Fatal(err)
			}

			if expected, actual := 2, len(entities); expected != actual {
				t.Errorf("expected: %d, actual: %d", expected, actual)
			}
			if expected, actual := firstAuthorID, entities[0].AuthorID; expected != actual {
				t.Errorf("expected: %q, actual: %q", expected, actual)
			}
			if expected, actual := secondAuthorID, entities[1].AuthorID; expected != actual {
				t.Errorf("expected: %q, actual: %q", expected, actual)
			}

			return true
		}
		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("transaction db failure", func(t *testing.T) {
		store := runStore(config)
		defer store.Stop()

		real := store.(*realStore)
		db := real.db
		real.db = nil

		err := store.Insert(Entity{})

		real.db = db

		if expected, actual := true, err != nil; expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}
	})
}

func TestRealStore_IntegrationQuery(t *testing.T) {
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

	t.Run("insert then query with no tags should select all", func(t *testing.T) {
		store := runStore(config)
		defer store.Stop()

		fn := func(parentID, resourceID uuid.UUID,
			resourceAddress string,
			resourceSize int64,
			resourceContentType, authorID, name ASCII,
			tags Tags,
		) bool {
			defer store.Drop()

			if err := store.Insert(Entity{
				ParentID:            parentID,
				ResourceID:          resourceID,
				ResourceAddress:     resourceAddress,
				ResourceSize:        resourceSize,
				ResourceContentType: resourceContentType.String(),
				AuthorID:            authorID.String(),
				Name:                name.String(),
				Tags:                tags.Slice(),
				CreatedOn:           time.Now(),
				DeletedOn:           time.Time{},
			}); err != nil {
				t.Fatal(err)
			}

			entities, err := store.SelectRevisions(resourceID, Query{Tags: make([]string, 0)})
			if err != nil {
				t.Fatal(err)
			}

			return len(entities) == 1
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("insert then query exact match", func(t *testing.T) {
		store := runStore(config)
		defer store.Stop()

		fn := func(parentID, resourceID uuid.UUID,
			resourceAddress string,
			resourceSize int64,
			resourceContentType, authorID, name ASCII,
			tags Tags,
		) bool {
			defer store.Drop()

			entity := Entity{
				ParentID:            parentID,
				ResourceID:          resourceID,
				ResourceAddress:     resourceAddress,
				ResourceSize:        resourceSize,
				ResourceContentType: resourceContentType.String(),
				AuthorID:            authorID.String(),
				Name:                name.String(),
				Tags:                tags.Slice(),
				CreatedOn:           time.Now().Round(time.Millisecond),
				DeletedOn:           time.Time{},
			}
			if err := store.Insert(entity); err != nil {
				t.Fatal(err)
			}

			got, err := store.SelectRevisions(resourceID, Query{
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

	t.Run("revisions puts then query exact match", func(t *testing.T) {
		store := runStore(config)
		defer store.Stop()

		fn := func(parentID, resourceID uuid.UUID,
			resourceAddress string,
			resourceSize int64,
			resourceContentType, authorID, name ASCII,
			tags Tags,
		) bool {
			defer store.Drop()

			want := make([]Entity, 10)
			for k := range want {
				entity := Entity{
					ParentID:            parentID,
					ResourceID:          resourceID,
					ResourceAddress:     resourceAddress,
					ResourceSize:        resourceSize,
					ResourceContentType: resourceContentType.String(),
					AuthorID:            fmt.Sprintf("%s%d", authorID.String(), k),
					Name:                name.String(),
					Tags:                tags.Slice(),
					CreatedOn:           time.Now().Add(time.Duration(k) * time.Second).Round(time.Millisecond),
					DeletedOn:           time.Time{},
				}
				if err := store.Insert(entity); err != nil {
					t.Fatal(err)
				}
				want[(len(want)-1)-k] = entity
			}

			got, err := store.SelectRevisions(resourceID, Query{
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
		store := runStore(config)
		defer store.Stop()

		fn := func(parentID, resourceID uuid.UUID,
			resourceAddress string,
			resourceSize int64,
			resourceContentType, authorID, name ASCII,
			tags Tags,
		) bool {
			defer store.Drop()

			want := make([]Entity, 10)
			for k := range want {
				entity := Entity{
					ParentID:            parentID,
					ResourceID:          resourceID,
					ResourceAddress:     resourceAddress,
					ResourceSize:        resourceSize,
					ResourceContentType: resourceContentType.String(),
					AuthorID:            fmt.Sprintf("%d%s", k, authorID.String()),
					Name:                name.String(),
					Tags:                tags.Slice(),
					CreatedOn:           time.Now().Add(time.Duration(k) * time.Second).Round(time.Millisecond),
					DeletedOn:           time.Time{},
				}
				if err := store.Insert(entity); err != nil {
					t.Fatal(err)
				}
				want[(len(want)-1)-k] = entity
			}

			got, err := store.SelectRevisions(resourceID, Query{
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

	t.Run("insert then query exact match with AuthorID", func(t *testing.T) {
		store := runStore(config)
		defer store.Stop()

		fn := func(parentID, resourceID uuid.UUID,
			resourceAddress string,
			resourceSize int64,
			resourceContentType, authorID, name ASCII,
			tags Tags,
		) bool {
			defer store.Drop()

			entity := Entity{
				ParentID:            parentID,
				ResourceID:          resourceID,
				ResourceAddress:     resourceAddress,
				ResourceSize:        resourceSize,
				ResourceContentType: resourceContentType.String(),
				AuthorID:            authorID.String(),
				Name:                name.String(),
				Tags:                tags.Slice(),
				CreatedOn:           time.Now().Round(time.Millisecond),
				DeletedOn:           time.Time{},
			}
			if err := store.Insert(entity); err != nil {
				t.Fatal(err)
			}

			authID := authorID.String()
			got, err := store.SelectRevisions(resourceID, Query{
				Tags:     tags.Slice(),
				AuthorID: &authID,
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

	t.Run("revisions puts then query exact match with AuthorID", func(t *testing.T) {
		store := runStore(config)
		defer store.Stop()

		fn := func(parentID, resourceID uuid.UUID,
			resourceAddress string,
			resourceSize int64,
			resourceContentType, authorID, name ASCII,
			tags Tags,
		) bool {
			defer store.Drop()

			want := make([]Entity, 10)
			for k := range want {
				entity := Entity{
					ParentID:            parentID,
					ResourceID:          resourceID,
					ResourceAddress:     resourceAddress,
					ResourceSize:        resourceSize,
					ResourceContentType: resourceContentType.String(),
					AuthorID:            fmt.Sprintf("%d%s", k, authorID.String()),
					Name:                name.String(),
					Tags:                tags.Slice(),
					CreatedOn:           time.Now().Round(time.Millisecond),
					DeletedOn:           time.Time{},
				}
				if err := store.Insert(entity); err != nil {
					t.Fatal(err)
				}
				want[k] = entity
			}

			authID := fmt.Sprintf("0%s", authorID.String())
			got, err := store.SelectRevisions(resourceID, Query{
				Tags:     tags.Slice(),
				AuthorID: &authID,
			})
			if err != nil {
				t.Fatal(err)
			}

			if expected, actual := 1, len(got); expected != actual {
				t.Fatalf("expected: %d, actual: %d", expected, actual)
			}

			return entityEquals(want[0], got[0])
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})
}

func runStore(config *RealConfig) Store {
	var wg sync.WaitGroup
	wg.Add(1)

	store := NewRealStore(config, log.NewNopLogger())

	go func() {
		go func() {
			// Make sure we breathe before closing the latch
			time.Sleep(time.Millisecond * 50)

			wg.Done()
		}()
		if err := store.Run(); err != nil {
			panic(err)
		}
	}()

	wg.Wait()

	return store
}
