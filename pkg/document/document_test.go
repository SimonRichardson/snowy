package document

import (
	"encoding/json"
	"math/rand"
	"reflect"
	"strings"
	"testing"
	"testing/quick"
	"time"

	"github.com/pkg/errors"
	"github.com/trussle/snowy/pkg/uuid"
)

func TestDocument(t *testing.T) {
	t.Parallel()

	t.Run("fields", func(t *testing.T) {
		fn := func(id uuid.UUID,
			name string,
			resourceID uuid.UUID,
			resourceAddress string,
			resourceSize int64,
			resourceContentType, authorID string,
			tags Tags,
		) bool {
			now := time.Now()
			output := Document{
				id:                  id,
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
		fn := func(id uuid.UUID,
			name string,
			resourceID uuid.UUID,
			resourceAddress string,
			resourceSize int64,
			resourceContentType, authorID string,
			tags Tags,
		) bool {
			now := time.Now().Round(time.Second)
			input := Document{
				id:                  id,
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

			var output Document
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
		fn := func(id uuid.UUID,
			name string,
			resourceID uuid.UUID,
			resourceAddress string,
			resourceSize int64,
			resourceContentType, authorID string,
		) bool {
			now := time.Now().Round(time.Second)
			input := Document{
				id:                  id,
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

			var output Document
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

			var output Document
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
		fn := func(name string, resourceID uuid.UUID, authorID string, tags Tags) bool {
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

			var output Document
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
		fn := func(name string, resourceID uuid.UUID, authorID string, tags Tags) bool {
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

			var output Document
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

func TestDocumentBuild(t *testing.T) {
	t.Parallel()

	t.Run("build", func(t *testing.T) {
		fn := func(id uuid.UUID,
			name string,
			resourceID uuid.UUID,
			resourceAddress string,
			resourceSize int64,
			resourceContentType, authorID string,
			tags Tags,
		) bool {
			now := time.Now()
			doc, err := BuildDocument(
				WithID(id),
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
			doc, err := BuildDocument(
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
		_, err := BuildDocument(
			func(doc *Document) error {
				return errors.Errorf("bad")
			},
		)

		if expected, actual := false, err == nil; expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}
	})
}

// Tags creates a series of tags that are ascii compliant.
type Tags []string

// Generate allows Tags to be used within quickcheck scenarios.
func (Tags) Generate(r *rand.Rand, size int) reflect.Value {
	var (
		chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		res   = make([]string, size)
	)

	for k := range res {
		str := make([]byte, r.Intn(50)+1)
		for k := range str {
			str[k] = chars[r.Intn(len(chars)-1)]
		}
		res[k] = string(str)
	}

	return reflect.ValueOf(res)
}

func (a Tags) Slice() []string {
	return a
}

func (a Tags) String() string {
	return strings.Join(a.Slice(), ",")
}
