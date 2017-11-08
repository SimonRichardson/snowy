package journals

import (
	"fmt"
	"net/http"
	"net/textproto"
	"net/url"
	"strconv"
	"testing"
	"testing/quick"

	"github.com/trussle/harness/generators"
	"github.com/trussle/snowy/pkg/uuid"
)

func TestInsertQueryParams(t *testing.T) {
	t.Parallel()

	t.Run("DecodeFrom with invalid content-type", func(t *testing.T) {
		fn := func(uid uuid.UUID, contentType generators.ASCII) bool {
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
		fn := func(uid uuid.UUID, contentType generators.ASCII) bool {
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
