package repository

import (
	"reflect"
	"sort"
	"testing"
	"testing/quick"

	"github.com/pkg/errors"
	"github.com/trussle/harness/generators"
)

func TestBuildingQuery(t *testing.T) {
	t.Parallel()

	t.Run("build", func(t *testing.T) {

		fn := func(tags generators.ASCIISlice, authorID string) bool {
			query, err := BuildQuery(
				WithQueryTags(tags.Slice()),
				WithQueryAuthorID(authorID),
			)
			if err != nil {
				t.Fatal(err)
			}

			if expected, actual := tags.Slice(), query.Tags; !reflect.DeepEqual(sortTags(expected), sortTags(actual)) {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}

			if expected, actual := authorID, *query.AuthorID; expected != actual {
				t.Errorf("expected: %q, actual: %q", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("invalid build", func(t *testing.T) {
		_, err := BuildQuery(
			func(query *Query) error {
				return errors.Errorf("bad")
			},
		)

		if expected, actual := false, err == nil; expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}
	})

	t.Run("empty query", func(t *testing.T) {
		query := BuildEmptyQuery()
		if expected, actual := emptyAuthID, *query.AuthorID; expected != actual {
			t.Errorf("expected: %q, actual: %q", expected, actual)
		}
	})
}

func TestNotFound(t *testing.T) {
	t.Parallel()

	t.Run("source", func(t *testing.T) {
		fn := func(source string) bool {
			err := errNotFound{errors.New(source)}

			if expected, actual := source, err.Error(); expected != actual {
				t.Errorf("expected: %q, actual: %q", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		fn := func(source string) bool {
			err := errNotFound{errors.New(source)}

			if expected, actual := true, err.NotFound(); expected != actual {
				t.Errorf("expected: %t, actual: %t", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("valid", func(t *testing.T) {
		fn := func(source string) bool {
			err := errNotFound{errors.New(source)}

			if expected, actual := true, ErrNotFound(err); expected != actual {
				t.Errorf("expected: %t, actual: %t", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		fn := func(source string) bool {
			err := errors.New(source)

			if expected, actual := false, ErrNotFound(err); expected != actual {
				t.Errorf("expected: %t, actual: %t", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})
}

func sortTags(tags []string) []string {
	res := make([]string, len(tags))
	copy(res, tags)
	sort.Strings(res)
	return res
}
