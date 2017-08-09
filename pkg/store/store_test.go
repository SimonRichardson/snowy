package store

import (
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

		fn := func(tags Tags) bool {
			query, err := BuildQuery(
				WithQueryTags(tags.Slice()),
			)
			if err != nil {
				t.Fatal(err)
			}

			if expected, actual := tags.Slice(), query.Tags; !reflect.DeepEqual(sortTags(expected), sortTags(actual)) {
				t.Errorf("expected: %v, actual: %v", expected, actual)
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
	return a.Name == b.Name &&
		a.ResourceID.Equals(b.ResourceID) &&
		reflect.DeepEqual([]byte(a.AuthorID), []byte(b.AuthorID)) &&
		reflect.DeepEqual(sortTags(a.Tags), sortTags(b.Tags)) &&
		a.CreatedOn.Equal(b.CreatedOn) &&
		a.DeletedOn.Equal(b.DeletedOn)
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
	if len(a) <= 1 {
		return a
	}
	return a[:len(a)/2]
}
