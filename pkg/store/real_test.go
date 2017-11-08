package store

import (
	"fmt"
	"reflect"
	"testing"
	"testing/quick"

	"github.com/lib/pq"

	"github.com/trussle/harness/generators"
	"github.com/trussle/snowy/pkg/uuid"
)

func TestConfig(t *testing.T) {
	t.Parallel()

	t.Run("empty config with connection string", func(t *testing.T) {
		got := ConnectionString(&RealConfig{})

		if expected, actual := "port=0", got; expected != actual {
			t.Errorf("expected: %q, actual: %q", expected, actual)
		}
	})

	t.Run("config connection string with host and port", func(t *testing.T) {
		fn := func(host string, port int) bool {
			config, err := BuildConfig(
				WithHostPort(host, port),
			)
			if err != nil {
				t.Fatal(err)
			}

			var want string
			if host != "" {
				want = fmt.Sprintf("host=%s port=%d", host, port)
			} else {
				want = fmt.Sprintf("port=%d", port)
			}

			got := ConnectionString(config)
			expected, actual := want, got

			return expected == actual
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("config connection string with username and password", func(t *testing.T) {
		fn := func(username, password string) bool {
			config, err := BuildConfig(
				WithUsername(username),
				WithPassword(password),
			)
			if err != nil {
				t.Fatal(err)
			}

			var want string
			if username != "" && password != "" {
				want = fmt.Sprintf("port=0 user=%s password=%s", username, password)
			} else if username == "" {
				want = fmt.Sprintf("port=0 password=%s", password)
			} else if password == "" {
				want = fmt.Sprintf("port=0 user=%s", username)
			}

			got := ConnectionString(config)
			expected, actual := want, got

			return reflect.DeepEqual([]byte(expected), []byte(actual))
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("config connection string with all settings", func(t *testing.T) {
		fn := func(username, password, dbName, host string, port int) bool {
			if username == "" || password == "" || dbName == "" || host == "" {
				return true
			}

			config, err := BuildConfig(
				WithUsername(username),
				WithPassword(password),
				WithDBName(dbName),
				WithSSLMode("disable"),
				WithHostPort(host, port),
			)
			if err != nil {
				t.Fatal(err)
			}

			var (
				want = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
					host, port,
					username, password,
					dbName,
					"disable",
				)
				got              = ConnectionString(config)
				expected, actual = want, got
			)

			return reflect.DeepEqual([]byte(expected), []byte(actual))
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("invalid config", func(t *testing.T) {
		_, err := BuildConfig(
			WithSSLMode("bad"),
		)
		if expected, actual := false, err == nil; expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}
	})
}

func TestSQLBuilder(t *testing.T) {
	t.Parallel()

	t.Run("select", func(t *testing.T) {
		fn := func(resourceID uuid.UUID) bool {
			statement, args := buildSQLFromQuery(resourceID, Query{})
			return statement == defaultSelectQuery &&
				reflect.DeepEqual(args, []interface{}{resourceID.String()})
		}
		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("select with empty authorID", func(t *testing.T) {
		fn := func(resourceID uuid.UUID) bool {
			s := ""
			statement, args := buildSQLFromQuery(resourceID, Query{
				AuthorID: &s,
			})
			return statement == defaultSelectQuery &&
				reflect.DeepEqual(args, []interface{}{resourceID.String()})
		}
		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("select with tags", func(t *testing.T) {
		fn := func(resourceID uuid.UUID, tags generators.ASCIISlice) bool {
			statement, args := buildSQLFromQuery(resourceID, Query{
				Tags: tags.Slice(),
			})
			return statement == defaultSelectQueryTags &&
				reflect.DeepEqual(args, []interface{}{
					resourceID.String(),
					pq.Array(tags.Slice()),
				})
		}
		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("select with tags and empty authorID", func(t *testing.T) {
		fn := func(resourceID uuid.UUID, tags generators.ASCIISlice) bool {
			s := ""
			statement, args := buildSQLFromQuery(resourceID, Query{
				Tags:     tags.Slice(),
				AuthorID: &s,
			})
			return statement == defaultSelectQueryTags &&
				reflect.DeepEqual(args, []interface{}{
					resourceID.String(),
					pq.Array(tags.Slice()),
				})
		}
		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("select with tags", func(t *testing.T) {
		fn := func(resourceID uuid.UUID, tags generators.ASCIISlice, authorID generators.ASCII) bool {
			s := authorID.String()
			statement, args := buildSQLFromQuery(resourceID, Query{
				Tags:     tags.Slice(),
				AuthorID: &s,
			})
			return statement == defaultSelectQueryTagsAuthorID &&
				reflect.DeepEqual(args, []interface{}{
					resourceID.String(),
					authorID.String(),
					pq.Array(tags.Slice()),
				})
		}
		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})
}

func TestSortTags(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Input, Output []string
	}{
		{[]string{"a", "b", "c"}, []string{"a", "b", "c"}},
		{[]string{"c", "a", "b"}, []string{"a", "b", "c"}},
		{[]string{"A", "a", "b"}, []string{"A", "a", "b"}},
	}

	for _, v := range testCases {
		if expected, actual := sortTags(v.Input), v.Output; !reflect.DeepEqual(expected, actual) {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	}
}
