package store

import (
	"testing"
	"testing/quick"

	"github.com/trussle/snowy/pkg/uuid"
)

func TestNopStore(t *testing.T) {
	t.Parallel()

	t.Run("get", func(t *testing.T) {
		store := NewNopStore()
		_, err := store.Get(uuid.New())

		if expected, actual := true, err == nil; expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}
	})

	t.Run("put", func(t *testing.T) {
		store := NewNopStore()

		fn := func(res uuid.UUID) bool {
			err := store.Put(Entity{ResourceID: res})
			return err == nil
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("put and get", func(t *testing.T) {
		store := NewNopStore()

		fn := func(res uuid.UUID) bool {
			if err := store.Put(Entity{ResourceID: res}); err != nil {
				return false
			}

			_, err := store.Get(res)
			return err == nil
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})
}
