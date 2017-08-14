package store

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"

	"github.com/go-kit/kit/log"
	"github.com/pkg/errors"
)

func TestBuildingStore(t *testing.T) {
	t.Parallel()

	t.Run("build", func(t *testing.T) {
		fn := func(name string) bool {
			config, err := Build(
				With(name),
				WithConfig(nil),
			)
			if err != nil {
				t.Fatal(err)
			}

			if expected, actual := name, config.name; expected != actual {
				t.Errorf("expected: %s, actual: %s", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("invalid build", func(t *testing.T) {
		_, err := Build(
			func(config *Config) error {
				return errors.Errorf("bad")
			},
		)

		if expected, actual := false, err == nil; expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}
	})
}

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("real", func(t *testing.T) {
		config, err := Build(
			With("real"),
		)
		if err != nil {
			t.Fatal(err)
		}

		_, err = New(config, log.NewNopLogger())
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("virtual", func(t *testing.T) {
		config, err := Build(
			With("virtual"),
		)
		if err != nil {
			t.Fatal(err)
		}

		_, err = New(config, log.NewNopLogger())
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("nop", func(t *testing.T) {
		config, err := Build(
			With("nop"),
		)
		if err != nil {
			t.Fatal(err)
		}

		_, err = New(config, log.NewNopLogger())
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		config, err := Build(
			With("invalid"),
		)
		if err != nil {
			t.Fatal(err)
		}

		_, err = New(config, log.NewNopLogger())
		if expected, actual := false, err == nil; expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}
	})
}

func TestBuildingQuery(t *testing.T) {
	t.Parallel()

	t.Run("build", func(t *testing.T) {

		fn := func(tags Tags, authorID string) bool {
			query, err := BuildQuery(
				WithQueryTags(tags.Slice()),
				WithQueryAuthorID(&authorID),
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

func equals(a, b []Entity) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if !entityEquals(v, b[k]) {
			return false
		}
	}
	return true
}

func entityEquals(a, b Entity) bool {
	if expected, actual := a.Name, b.Name; expected != actual {
		fmt.Printf("name - expected: %q, actual: %q\n", expected, actual)
		return false
	}
	if expected, actual := a.ResourceID, b.ResourceID; !expected.Equals(actual) {
		fmt.Printf("resource_id - expected: %v, actual: %v\n", expected, actual)
		return false
	}
	if expected, actual := a.AuthorID, b.AuthorID; expected != actual {
		fmt.Printf("author_id - expected: %q, actual: %q\n", expected, actual)
		//	return false
	}
	if expected, actual := sortTags(a.Tags), sortTags(b.Tags); !reflect.DeepEqual(expected, actual) {
		fmt.Printf("tags - expected: %v, actual: %v\n", expected, actual)
		return false
	}
	if expected, actual := a.CreatedOn, b.CreatedOn; !expected.Equal(actual) {
		fmt.Printf("created_on - expected: %v, actual: %v\n", expected, actual)
		return false
	}
	if expected, actual := a.DeletedOn, b.DeletedOn; !expected.Equal(actual) {
		fmt.Printf("deleted_on - expected: %v, actual: %v\n", expected, actual)
		return false
	}
	return true
}

func randomizeTags(a []string) []string {
	res := make([]string, len(a))
	copy(res, a)

	if len(res) <= 1 {
		return res
	}

	for i := range res {
		j := rand.Intn(i + 1)
		res[i], res[j] = res[j], res[i]
	}

	return res
}

func splitTags(a []string) []string {
	if len(a) < 2 {
		return a
	}
	return a[:len(a)/2]
}
