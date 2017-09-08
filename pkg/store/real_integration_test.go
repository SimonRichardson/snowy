// +build integration

package store

import (
	"fmt"
	"math/rand"
	"reflect"
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

		_, err := store.Get(uuid.New(), Query{})
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

			entity, err := store.Get(resourceID, Query{})
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
			got, err := store.GetMultiple(resourceID, Query{
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

	t.Run("multiple puts then query exact match with AuthorID", func(t *testing.T) {
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
			got, err := store.GetMultiple(resourceID, Query{
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

// ASCII creates a series of tags that are ascii compliant.
type ASCII []byte

// Generate allows ASCII to be used within quickcheck scenarios.
func (ASCII) Generate(r *rand.Rand, size int) reflect.Value {
	var (
		chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		res   = make([]byte, size)
	)

	for k := range res {
		res[k] = byte(chars[r.Intn(len(chars)-1)])
	}

	return reflect.ValueOf(res)
}

func (a ASCII) Slice() []byte {
	return a
}

func (a ASCII) String() string {
	return string(a)
}
