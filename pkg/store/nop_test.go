package store

import (
	"testing"
	"testing/quick"

	"github.com/trussle/uuid"
)

func TestNopStore(t *testing.T) {
	t.Parallel()

	t.Run("get", func(t *testing.T) {
		store := NewNopStore()
		_, err := store.Select(uuid.MustNew(), Query{})

		if expected, actual := true, err == nil; expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}
	})

	t.Run("insert", func(t *testing.T) {
		store := NewNopStore()

		fn := func(res uuid.UUID) bool {
			err := store.Insert(Entity{ResourceID: res})
			return err == nil
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("insert and get", func(t *testing.T) {
		store := NewNopStore()

		fn := func(res uuid.UUID) bool {
			if err := store.Insert(Entity{ResourceID: res}); err != nil {
				return false
			}

			_, err := store.Select(res, Query{})
			return err == nil
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("insert then query with no tags should select none", func(t *testing.T) {
		store := NewNopStore()

		fn := func(res uuid.UUID) bool {
			if err := store.Insert(Entity{ResourceID: res}); err != nil {
				return false
			}

			entities, err := store.SelectRevisions(res, Query{Tags: make([]string, 0)})
			if err != nil {
				return false
			}

			return len(entities) == 0
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("insert then query any match should return none", func(t *testing.T) {
		store := NewNopStore()

		fn := func(res uuid.UUID, tags []string) bool {
			entity := Entity{
				ResourceID: res,
				Tags:       tags,
			}
			if err := store.Insert(entity); err != nil {
				return false
			}

			entities, err := store.SelectRevisions(res, Query{
				Tags: tags,
			})
			if err != nil {
				return false
			}

			return len(entities) == 0
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("drop", func(t *testing.T) {
		store := NewNopStore()

		if err := store.Drop(); err != nil {
			t.Error(err)
		}
	})

	t.Run("run and stop", func(t *testing.T) {
		store := NewNopStore()

		if err := store.Run(); err != nil {
			t.Error(err)
		}

		store.Stop()
	})
}
