package models

import (
	"encoding/json"
	"reflect"
	"testing"
	"testing/quick"
	"time"

	"github.com/pkg/errors"
	"github.com/trussle/harness/generators"
	"github.com/trussle/uuid"
)

func TestLedger(t *testing.T) {
	t.Parallel()

	t.Run("fields", func(t *testing.T) {
		fn := func(id, parentID uuid.UUID,
			name string,
			resourceID uuid.UUID,
			resourceAddress string,
			resourceSize int64,
			resourceContentType, authorID string,
			tags generators.ASCIISlice,
		) bool {
			now := time.Now()
			output := Ledger{
				id:                  id,
				parentID:            parentID,
				name:                name,
				resourceID:          resourceID,
				resourceAddress:     resourceAddress,
				resourceSize:        resourceSize,
				resourceContentType: resourceContentType,
				authorID:            authorID,
				tags:                tags,
				createdOn:           now,
				deletedOn:           now,
			}

			return output.ID().Equals(id) &&
				output.ParentID().Equals(parentID) &&
				output.Name() == name &&
				output.ResourceID().Equals(resourceID) &&
				output.ResourceAddress() == resourceAddress &&
				output.ResourceSize() == resourceSize &&
				output.ResourceContentType() == resourceContentType &&
				output.AuthorID() == authorID &&
				reflect.DeepEqual(output.Tags(), tags.Slice()) &&
				output.CreatedOn().Equal(now) &&
				output.DeletedOn().Equal(now)
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("json marshal", func(t *testing.T) {
		fn := func(id, parentID uuid.UUID,
			name string,
			resourceID uuid.UUID,
			resourceAddress string,
			resourceSize int64,
			resourceContentType, authorID string,
			tags generators.ASCIISlice,
		) bool {
			now := time.Now().Round(time.Second)
			input := Ledger{
				id:                  id,
				parentID:            parentID,
				name:                name,
				resourceID:          resourceID,
				resourceAddress:     resourceAddress,
				resourceSize:        resourceSize,
				resourceContentType: resourceContentType,
				authorID:            authorID,
				tags:                tags,
				createdOn:           now,
				deletedOn:           now,
			}

			bytes, err := json.Marshal(input)
			if err != nil {
				t.Fatal(err)
			}

			var output Ledger
			if err = json.Unmarshal(bytes, &output); err != nil {
				t.Fatal(err)
			}

			return output.Name() == name &&
				output.ResourceID().Equals(resourceID) &&
				output.ResourceAddress() == resourceAddress &&
				output.ResourceSize() == resourceSize &&
				output.ResourceContentType() == resourceContentType &&
				output.AuthorID() == authorID &&
				reflect.DeepEqual(output.Tags(), tags.Slice()) &&
				output.CreatedOn().Equal(now) &&
				output.DeletedOn().Equal(now)
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("json marshal with empty tags", func(t *testing.T) {
		fn := func(id, parentID uuid.UUID,
			name string,
			resourceID uuid.UUID,
			resourceAddress string,
			resourceSize int64,
			resourceContentType, authorID string,
		) bool {
			now := time.Now().Round(time.Second)
			input := Ledger{
				id:                  id,
				parentID:            parentID,
				name:                name,
				resourceID:          resourceID,
				resourceAddress:     resourceAddress,
				resourceSize:        resourceSize,
				resourceContentType: resourceContentType,
				authorID:            authorID,
				createdOn:           now,
				deletedOn:           now,
			}

			bytes, err := json.Marshal(input)
			if err != nil {
				t.Fatal(err)
			}

			var output Ledger
			if err = json.Unmarshal(bytes, &output); err != nil {
				t.Fatal(err)
			}

			return output.Name() == name &&
				output.ResourceID().Equals(resourceID) &&
				output.ResourceAddress() == resourceAddress &&
				output.ResourceSize() == resourceSize &&
				output.ResourceContentType() == resourceContentType &&
				output.AuthorID() == authorID &&
				reflect.DeepEqual(output.Tags(), make([]string, 0)) &&
				output.CreatedOn().Equal(now) &&
				output.DeletedOn().Equal(now)
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("json unmarshal with malformed body", func(t *testing.T) {
		fn := func() bool {
			bytes := []byte("{!}")

			var output Ledger
			err := output.UnmarshalJSON(bytes)
			if expected, actual := false, err == nil; expected != actual {
				t.Errorf("expected: %t, actual: %t", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("json unmarshal with malformed created_on", func(t *testing.T) {
		fn := func(name string, resourceID uuid.UUID, authorID string, tags generators.ASCIISlice) bool {
			now := time.Now().Round(time.Second)
			input := struct {
				Name       string    `json:"name"`
				ResourceID uuid.UUID `json:"resource_id"`
				AuthorID   string    `json:"author_id"`
				Tags       []string  `json:"tags"`
				CreatedOn  string    `json:"created_on"`
				DeletedOn  string    `json:"deleted_on"`
			}{
				Name:       name,
				ResourceID: resourceID,
				AuthorID:   authorID,
				Tags:       tags,
				CreatedOn:  "bad",
				DeletedOn:  now.Format(time.RFC3339),
			}

			bytes, err := json.Marshal(input)
			if err != nil {
				t.Fatal(err)
			}

			var output Ledger
			err = output.UnmarshalJSON(bytes)
			if expected, actual := false, err == nil; expected != actual {
				t.Errorf("expected: %t, actual: %t", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("json unmarshal with malformed deleted_on", func(t *testing.T) {
		fn := func(name string, resourceID uuid.UUID, authorID string, tags generators.ASCIISlice) bool {
			now := time.Now().Round(time.Second)
			input := struct {
				Name       string    `json:"name"`
				ResourceID uuid.UUID `json:"resource_id"`
				AuthorID   string    `json:"author_id"`
				Tags       []string  `json:"tags"`
				CreatedOn  string    `json:"created_on"`
				DeletedOn  string    `json:"deleted_on"`
			}{
				Name:       name,
				ResourceID: resourceID,
				AuthorID:   authorID,
				Tags:       tags,
				CreatedOn:  now.Format(time.RFC3339),
				DeletedOn:  "bad",
			}

			bytes, err := json.Marshal(input)
			if err != nil {
				t.Fatal(err)
			}

			var output Ledger
			err = output.UnmarshalJSON(bytes)
			if expected, actual := false, err == nil; expected != actual {
				t.Errorf("expected: %t, actual: %t", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})
}

func TestLedgerBuild(t *testing.T) {
	t.Parallel()

	t.Run("build", func(t *testing.T) {
		fn := func(id, parentID uuid.UUID,
			name string,
			resourceID uuid.UUID,
			resourceAddress string,
			resourceSize int64,
			resourceContentType, authorID string,
			tags generators.ASCIISlice,
		) bool {
			now := time.Now()
			doc, err := BuildLedger(
				WithID(id),
				WithParentID(parentID),
				WithName(name),
				WithResourceID(resourceID),
				WithResourceAddress(resourceAddress),
				WithResourceSize(resourceSize),
				WithResourceContentType(resourceContentType),
				WithAuthorID(authorID),
				WithTags(tags),
				WithCreatedOn(now),
				WithDeletedOn(now),
			)
			if err != nil {
				t.Fatal(err)
			}
			return doc.ID().Equals(id) &&
				doc.ParentID().Equals(parentID) &&
				doc.Name() == name &&
				doc.ResourceID().Equals(resourceID) &&
				doc.ResourceAddress() == resourceAddress &&
				doc.ResourceSize() == resourceSize &&
				doc.ResourceContentType() == resourceContentType &&
				doc.AuthorID() == authorID &&
				reflect.DeepEqual(doc.Tags(), tags.Slice()) &&
				doc.CreatedOn().Equal(now) &&
				doc.DeletedOn().Equal(now)
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("build with new resource_id", func(t *testing.T) {
		fn := func() bool {
			doc, err := BuildLedger(
				WithNewResourceID(),
			)
			if err != nil {
				t.Fatal(err)
			}
			return !doc.ResourceID().Zero()
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("invalid build", func(t *testing.T) {
		_, err := BuildLedger(
			func(doc *Ledger) error {
				return errors.Errorf("bad")
			},
		)

		if expected, actual := false, err == nil; expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}
	})
}
