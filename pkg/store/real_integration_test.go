// +build integration

package store

import (
	"testing"
	"testing/quick"
	"time"

	"github.com/trussle/snowy/pkg/uuid"
	"github.com/go-kit/kit/log"
)

func TestRealStore_Integration(t *testing.T) {
	t.Parallel()

	config, err := BuildConfig(
		WithHostPort("store", 5432),
		WithUsername("postgres"),
		WithPassword("postgres"),
		WithSSLMode("disable"),
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("get", func(t *testing.T) {
		store := NewRealStore(config, log.NewNopLogger())
		go func() {
			defer store.Stop()
			store.Run()
		}()
		time.Sleep(time.Millisecond * 100)

		_, err := store.Get(uuid.New())
		if expected, actual := true, err != nil; expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}

		if expected, actual := true, ErrNotFound(err); expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}
	})

	t.Run("put", func(t *testing.T) {
		store := NewRealStore(config, log.NewNopLogger())
		go func() {
			defer store.Stop()
			store.Run()
		}()
		time.Sleep(time.Millisecond * 100)

		fn := func(resourceID, authorID uuid.UUID, name string, tags []string) bool {
			err := store.Put(Entity{
				ResourceID: resourceID,
				AuthorID:   authorID,
				Name:       name,
				Tags:       tags,
				CreatedOn:  time.Now(),
				DeletedOn:  time.Time{},
			})
			return err == nil
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("put then get", func(t *testing.T) {
		store := NewRealStore(config, log.NewNopLogger())
		go func() {
			defer store.Stop()
			store.Run()
		}()
		time.Sleep(time.Millisecond * 100)

		fn := func(resourceID, authorID uuid.UUID, name string, tags []string) bool {
			if err := store.Put(Entity{
				ResourceID: resourceID,
				AuthorID:   authorID,
				Name:       name,
				Tags:       tags,
				CreatedOn:  time.Now(),
				DeletedOn:  time.Time{},
			}); err != nil {
				t.Fatal(err)
			}

			entity, err := store.Get(resourceID)
			if err != nil {
				return false
			}
			return entity.ResourceID.Equals(resourceID)
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})
}
