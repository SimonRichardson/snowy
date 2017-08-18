package models

import (
	"encoding/json"
	"testing"
	"testing/quick"

	"github.com/pkg/errors"
)

func TestContent(t *testing.T) {
	t.Parallel()

	t.Run("json marshal", func(t *testing.T) {
		fn := func(address string, size int64, contentType string) bool {
			input := Content{
				address:     address,
				size:        size,
				contentType: contentType,
			}

			bytes, err := json.Marshal(input)
			if err != nil {
				t.Fatal(err)
			}

			var output Content
			if err = json.Unmarshal(bytes, &output); err != nil {
				t.Fatal(err)
			}

			return output.Address() == address &&
				output.Size() == size &&
				output.ContentType() == contentType
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("json unmarshal with malformed body", func(t *testing.T) {
		fn := func() bool {
			bytes := []byte("{!}")

			var output Content
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
}

func TestContentAddress(t *testing.T) {
	t.Parallel()

	t.Run("calculate", func(t *testing.T) {
		fn := func(bytes []byte) bool {
			res, err := ContentAddress(bytes)
			if err != nil {
				t.Fatal(err)
			}

			return len(res) > 0
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("known values", func(t *testing.T) {
		res, err := ContentAddress([]byte("hello"))
		if err != nil {
			t.Fatal(err)
		}

		want := "LPJNul-wow4m6DsqxbninhsWHlwfp0JecwQzYpOLmCQ="
		if expected, actual := want, res; expected != actual {
			t.Errorf("expected: %q, actual: %q", expected, actual)
		}
	})
}

func TestContentBuild(t *testing.T) {
	t.Parallel()

	t.Run("build", func(t *testing.T) {
		fn := func(address string, size int64, contentType string, body string) bool {
			content, err := BuildContent(
				WithAddress(address),
				WithSize(size),
				WithContentType(contentType),
				WithBytes([]byte(body)),
			)
			if err != nil {
				t.Fatal(err)
			}

			bytes, err := content.Bytes()
			if err != nil {
				t.Fatal(err)
			}

			return content.Address() == address &&
				content.Size() == size &&
				content.ContentType() == contentType &&
				string(bytes) == body
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("invalid build", func(t *testing.T) {
		_, err := BuildContent(
			func(content *Content) error {
				return errors.Errorf("bad")
			},
		)

		if expected, actual := false, err == nil; expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}
	})
}
