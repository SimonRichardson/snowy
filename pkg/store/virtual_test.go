package store

import (
	"testing"
	"testing/quick"

	"github.com/trussle/snowy/pkg/uuid"
)

func TestVirtualStore(t *testing.T) {
	t.Parallel()

	t.Run("get", func(t *testing.T) {
		store := NewVirtualStore()
		_, err := store.Get(uuid.New())

		if expected, actual := true, err != nil; expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}

		if expected, actual := true, ErrNotFound(err); expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}
	})

	t.Run("put", func(t *testing.T) {
		store := NewVirtualStore()

		fn := func(res uuid.UUID) bool {
			err := store.Put(Entity{ResourceID: res})
			return err == nil
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("put then get", func(t *testing.T) {
		store := NewVirtualStore()

		fn := func(res uuid.UUID) bool {
			if err := store.Put(Entity{ResourceID: res}); err != nil {
				return false
			}

			entity, err := store.Get(res)
			if err != nil {
				return false
			}

			return entity.ResourceID == res
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})
}
