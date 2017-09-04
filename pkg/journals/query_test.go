package journals

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/textproto"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"testing/quick"

	"github.com/trussle/snowy/pkg/uuid"
)

func TestInsertQueryParams(t *testing.T) {
	t.Parallel()

	t.Run("DecodeFrom with invalid content-type", func(t *testing.T) {
		fn := func(uid uuid.UUID, contentType ASCII) bool {
			var (
				qp InsertQueryParams

				u, err = url.Parse(fmt.Sprintf("/?resource_id=%s", uid.String()))
				h      = make(http.Header, 0)
			)
			if err != nil {
				t.Fatal(err)
			}

			h.Set("Content-Type", contentType.String())

			err = qp.DecodeFrom(u, h, queryRequired)

			if expected, actual := false, err == nil; expected != actual {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("DecodeFrom with invalid content-type", func(t *testing.T) {
		fn := func(uid uuid.UUID) bool {
			var (
				qp InsertQueryParams

				u, err = url.Parse(fmt.Sprintf("/?resource_id=%s", uid.String()))
				h      = make(http.Header, 0)
			)
			if err != nil {
				t.Fatal(err)
			}

			h.Set("Content-Type", "multipart/form-data")

			err = qp.DecodeFrom(u, h, queryRequired)

			if expected, actual := true, err == nil; expected != actual {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})
}

func TestInsertFileQueryParams(t *testing.T) {
	t.Parallel()

	t.Run("DecodeFrom with required empty url", func(t *testing.T) {
		var (
			qp InsertFileQueryParams

			h = make(textproto.MIMEHeader)
		)

		err := qp.DecodeFrom(h, queryRequired)
		if expected, actual := false, err == nil; expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("DecodeFrom with no content-length", func(t *testing.T) {
		fn := func(contentType string) bool {
			var (
				qp InsertFileQueryParams

				h = make(textproto.MIMEHeader)
			)

			h.Set("content-type", contentType)

			err := qp.DecodeFrom(h, queryRequired)
			if expected, actual := false, err == nil; expected != actual {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("DecodeFrom with no content-type", func(t *testing.T) {
		fn := func(contentLength int64) bool {
			var (
				qp InsertFileQueryParams

				h = make(textproto.MIMEHeader)
			)

			h.Set("content-length", strconv.FormatInt(contentLength, 10))

			err := qp.DecodeFrom(h, queryRequired)
			if expected, actual := false, err == nil; expected != actual {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("DecodeFrom with invalid content-length", func(t *testing.T) {
		fn := func(contentType, contentLength string) bool {
			var (
				qp InsertFileQueryParams

				h = make(textproto.MIMEHeader)
			)

			h.Set("content-type", contentType)
			h.Set("content-length", contentLength)

			err := qp.DecodeFrom(h, queryRequired)
			if expected, actual := false, err == nil; expected != actual {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("DecodeFrom with content-length to large", func(t *testing.T) {
		fn := func(contentType string) bool {
			var (
				qp InsertFileQueryParams

				h = make(textproto.MIMEHeader)
			)

			h.Set("content-type", contentType)
			h.Set("content-length", strconv.FormatInt(defaultMaxContentLength+1, 10))

			err := qp.DecodeFrom(h, queryRequired)
			if expected, actual := false, err == nil; expected != actual {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}

			return qp.ContentType() == contentType
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("DecodeFrom with content-length to small", func(t *testing.T) {
		fn := func(contentType string, contentLength uint) bool {
			var (
				qp InsertFileQueryParams

				h = make(textproto.MIMEHeader)
			)

			h.Set("content-type", contentType)
			h.Set("content-length", strconv.FormatInt(-int64(contentLength), 10))

			err := qp.DecodeFrom(h, queryRequired)
			if expected, actual := false, err == nil; expected != actual {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}

			return qp.ContentType() == contentType
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("DecodeFrom", func(t *testing.T) {
		fn := func(contentType string, contentLength uint) bool {
			if contentType == "" {
				return true
			}

			var (
				qp InsertFileQueryParams

				h       = make(textproto.MIMEHeader)
				size    = int64(contentLength % (defaultMaxContentLength - 1))
				fmtSize = strconv.FormatInt(size, 10)
			)

			h.Set("content-type", contentType)
			h.Set("content-length", fmtSize)

			err := qp.DecodeFrom(h, queryRequired)
			if expected, actual := true, err == nil; expected != actual {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}

			return qp.ContentType() == contentType &&
				qp.ContentLength() == size
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})
}

func TestAppendQueryParams(t *testing.T) {
	t.Parallel()

	t.Run("DecodeFrom with required empty url", func(t *testing.T) {
		var (
			qp AppendQueryParams

			u, err = url.Parse("")
			h      = make(http.Header, 0)
		)
		if err != nil {
			t.Fatal(err)
		}

		h.Set("Content-Type", "multipart/form-data")

		err = qp.DecodeFrom(u, h, queryRequired)

		if expected, actual := false, err == nil; expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("DecodeFrom with optional empty url", func(t *testing.T) {
		var (
			qp AppendQueryParams

			u, err = url.Parse("")
			h      = make(http.Header, 0)
		)
		if err != nil {
			t.Fatal(err)
		}

		h.Set("Content-Type", "multipart/form-data")

		err = qp.DecodeFrom(u, h, queryOptional)

		if expected, actual := true, err == nil; expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("DecodeFrom with invalid resource_id", func(t *testing.T) {
		var (
			qp AppendQueryParams

			u, err = url.Parse("/?resource_id=123asd")
			h      = make(http.Header, 0)
		)
		if err != nil {
			t.Fatal(err)
		}

		h.Set("Content-Type", "multipart/form-data")

		err = qp.DecodeFrom(u, h, queryRequired)

		if expected, actual := false, err == nil; expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("DecodeFrom with invalid content-type", func(t *testing.T) {
		fn := func(uid uuid.UUID, contentType ASCII) bool {
			var (
				qp AppendQueryParams

				u, err = url.Parse(fmt.Sprintf("/?resource_id=%s", uid.String()))
				h      = make(http.Header, 0)
			)
			if err != nil {
				t.Fatal(err)
			}

			h.Set("Content-Type", contentType.String())

			err = qp.DecodeFrom(u, h, queryRequired)

			if expected, actual := false, err == nil; expected != actual {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("DecodeFrom with valid resource_id", func(t *testing.T) {
		fn := func(uid uuid.UUID) bool {
			var (
				qp AppendQueryParams

				u, err = url.Parse(fmt.Sprintf("/?resource_id=%s", uid.String()))
				h      = make(http.Header, 0)
			)
			if err != nil {
				t.Fatal(err)
			}

			h.Set("Content-Type", "multipart/form-data")

			err = qp.DecodeFrom(u, h, queryRequired)

			if expected, actual := true, err == nil; expected != actual {
				t.Errorf("expected: %t, actual: %t", expected, actual)
			}
			return uid.Equals(qp.ResourceID)
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

func (a Tags) String() string {
	return strings.Join(a.Slice(), ",")
}

// ASCII creates a series of tags that are ascii compliant.
type ASCII []byte

// Generate allows ASCII to be used within quickcheck scenarios.
func (ASCII) Generate(r *rand.Rand, size int) reflect.Value {
	var (
		chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		res   = make([]byte, size)
	)

	for k := range res {
		res[k] = byte(chars[r.Intn(len(chars)-1)])
	}

	return reflect.ValueOf(res)
}

func (a ASCII) Slice() []byte {
	return a
}

func (a ASCII) String() string {
	return string(a)
}
