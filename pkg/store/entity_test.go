package store

import (
	"testing"
	"testing/quick"
	"time"

	"github.com/pkg/errors"
	"github.com/trussle/snowy/pkg/uuid"
)

func TestBuildingEntity(t *testing.T) {
	t.Parallel()

	t.Run("build", func(t *testing.T) {

		fn := func(id, parentID uuid.UUID,
			name string,
			resourceID uuid.UUID,
			resourceAddress string,
			resourceSize int64,
			resourceContentType, authorID string,
			tags Tags,
		) bool {
			now := time.Now()
			entity, err := BuildEntity(
				BuildEntityWithID(id),
				BuildEntityWithParentID(parentID),
				BuildEntityWithName(name),
				BuildEntityWithResourceID(resourceID),
				BuildEntityWithResourceAddress(resourceAddress),
				BuildEntityWithResourceSize(resourceSize),
				BuildEntityWithResourceContentType(resourceContentType),
				BuildEntityWithAuthorID(authorID),
				BuildEntityWithTags(tags),
				BuildEntityWithCreatedOn(now),
				BuildEntityWithDeletedOn(time.Time{}),
			)
			if err != nil {
				t.Fatal(err)
			}

			want := Entity{
				ID:                  id,
				ParentID:            parentID,
				Name:                name,
				ResourceID:          resourceID,
				ResourceAddress:     resourceAddress,
				ResourceSize:        resourceSize,
				ResourceContentType: resourceContentType,
				AuthorID:            authorID,
				Tags:                tags,
				CreatedOn:           now,
				DeletedOn:           time.Time{},
			}

			if expected, actual := want, entity; !entityEquals(expected, actual) {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("invalid build", func(t *testing.T) {
		_, err := BuildEntity(
			func(entity *Entity) error {
				return errors.Errorf("bad")
			},
		)

		if expected, actual := false, err == nil; expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}
	})
}
