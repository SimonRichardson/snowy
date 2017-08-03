package store

import (
	"fmt"
	"reflect"
	"testing"
	"testing/quick"
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
}
