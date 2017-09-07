package contents

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strconv"
	"testing"
	"testing/quick"

	"github.com/go-kit/kit/log"
	"github.com/trussle/snowy/pkg/models"
	"github.com/trussle/snowy/pkg/uuid"
)

func TestSelectQueryParams(t *testing.T) {
	t.Parallel()

	t.Run("DecodeFrom with required empty url", func(t *testing.T) {
		var (
			qp SelectQueryParams

			u, err = url.Parse("")
		)
		if err != nil {
			t.Fatal(err)
		}

		err = qp.DecodeFrom(u, queryRequired)

		if expected, actual := false, err == nil; expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("DecodeFrom with optional empty url", func(t *testing.T) {
		var (
			qp SelectQueryParams

			u, err = url.Parse("")
		)
		if err != nil {
			t.Fatal(err)
		}

		err = qp.DecodeFrom(u, queryOptional)

		if expected, actual := true, err == nil; expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("DecodeFrom with invalid resource_id", func(t *testing.T) {
		var (
			qp SelectQueryParams

			u, err = url.Parse("/?resource_id=123asd")
		)
		if err != nil {
			t.Fatal(err)
		}

		err = qp.DecodeFrom(u, queryRequired)

		if expected, actual := false, err == nil; expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("DecodeFrom with valid resource_id", func(t *testing.T) {
		fn := func(uid uuid.UUID) bool {
			var (
				qp SelectQueryParams

				u, err = url.Parse(fmt.Sprintf("/?resource_id=%s", uid.String()))
			)
			if err != nil {
				t.Fatal(err)
			}

			err = qp.DecodeFrom(u, queryRequired)

			if expected, actual := true, err == nil; expected != actual {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}
			return uid.Equals(qp.ResourceID)
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})
}

func TestSelectQueryResult(t *testing.T) {
	t.Parallel()

	t.Run("EncodeTo includes the correct headers", func(t *testing.T) {
		fn := func(uid uuid.UUID) bool {
			var (
				qp SelectQueryParams

				u, err = url.Parse(fmt.Sprintf("/?resource_id=%s", uid.String()))
			)
			if err != nil {
				t.Fatal(err)
			}

			err = qp.DecodeFrom(u, queryRequired)
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()

			res := SelectQueryResult{Logger: log.NewNopLogger(), Params: qp}
			res.EncodeTo(recorder)

			headers := recorder.Header()
			return headers.Get(httpHeaderResourceID) == uid.String()
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("EncodeTo with no content has correct status code", func(t *testing.T) {
		fn := func(uid uuid.UUID) bool {
			var (
				qp SelectQueryParams

				u, err = url.Parse(fmt.Sprintf("/?resource_id=%s", uid.String()))
			)
			if err != nil {
				t.Fatal(err)
			}

			err = qp.DecodeFrom(u, queryRequired)
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()

			res := SelectQueryResult{Logger: log.NewNopLogger(), Params: qp}
			res.EncodeTo(recorder)

			return recorder.Code == 200
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("EncodeTo with no content has correct body", func(t *testing.T) {
		fn := func(uid uuid.UUID) bool {
			var (
				qp SelectQueryParams

				u, err = url.Parse(fmt.Sprintf("/?resource_id=%s", uid.String()))
			)
			if err != nil {
				t.Fatal(err)
			}

			err = qp.DecodeFrom(u, queryRequired)
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()

			res := SelectQueryResult{Logger: log.NewNopLogger(), Params: qp}
			res.EncodeTo(recorder)

			return string(recorder.Body.Bytes()) == ""
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("EncodeTo with content has correct body", func(t *testing.T) {
		fn := func(uid uuid.UUID, body []byte) bool {
			var (
				qp SelectQueryParams

				u, err = url.Parse(fmt.Sprintf("/?resource_id=%s", uid.String()))
			)
			if err != nil {
				t.Fatal(err)
			}

			err = qp.DecodeFrom(u, queryRequired)
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()

			res := SelectQueryResult{Logger: log.NewNopLogger(), Params: qp}
			res.Content, err = models.BuildContent(
				models.WithBytes(body),
			)
			if err != nil {
				t.Fatal(err)
			}

			res.EncodeTo(recorder)

			return bytesEqual(recorder.Body.Bytes(), body)
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})
}

func TestInsertQueryParams(t *testing.T) {
	t.Parallel()

	t.Run("DecodeFrom with required empty url", func(t *testing.T) {
		var (
			qp InsertQueryParams

			u, err = url.Parse("")
			h      = make(http.Header)
		)
		if err != nil {
			t.Fatal(err)
		}

		err = qp.DecodeFrom(u, h, queryRequired)

		if expected, actual := false, err == nil; expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("DecodeFrom with optional empty url", func(t *testing.T) {
		var (
			qp InsertQueryParams

			u, err = url.Parse("")
			h      = make(http.Header)
		)
		if err != nil {
			t.Fatal(err)
		}

		err = qp.DecodeFrom(u, h, queryOptional)

		if expected, actual := true, err == nil; expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("DecodeFrom with no content-length", func(t *testing.T) {
		fn := func(contentType string) bool {
			var (
				qp InsertQueryParams

				u, err = url.Parse("")
				h      = make(http.Header)
			)
			if err != nil {
				t.Fatal(err)
			}

			h.Set("content-type", contentType)

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

	t.Run("DecodeFrom with no content-type", func(t *testing.T) {
		fn := func(contentLength int64) bool {
			var (
				qp InsertQueryParams

				u, err = url.Parse("")
				h      = make(http.Header)
			)
			if err != nil {
				t.Fatal(err)
			}

			h.Set("content-length", strconv.FormatInt(contentLength, 10))

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

	t.Run("DecodeFrom with invalid content-length", func(t *testing.T) {
		fn := func(contentType, contentLength string) bool {
			var (
				qp InsertQueryParams

				u, err = url.Parse("")
				h      = make(http.Header)
			)
			if err != nil {
				t.Fatal(err)
			}

			h.Set("content-type", contentType)
			h.Set("content-length", contentLength)

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

	t.Run("DecodeFrom with content-length to large", func(t *testing.T) {
		fn := func(contentType string) bool {
			var (
				qp InsertQueryParams

				u, err = url.Parse("")
				h      = make(http.Header)
			)
			if err != nil {
				t.Fatal(err)
			}

			h.Set("content-type", contentType)
			h.Set("content-length", strconv.FormatInt(defaultMaxContentLength+1, 10))

			err = qp.DecodeFrom(u, h, queryRequired)

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
				qp InsertQueryParams

				u, err = url.Parse("")
				h      = make(http.Header)
			)
			if err != nil {
				t.Fatal(err)
			}

			h.Set("content-type", contentType)
			h.Set("content-length", strconv.FormatInt(-int64(contentLength), 10))

			err = qp.DecodeFrom(u, h, queryRequired)

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
				qp InsertQueryParams

				u, err = url.Parse("")
				h      = make(http.Header)
			)
			if err != nil {
				t.Fatal(err)
			}

			var (
				size    = int64(contentLength % (defaultMaxContentLength - 1))
				fmtSize = strconv.FormatInt(size, 10)
			)

			h.Set("content-type", contentType)
			h.Set("content-length", fmtSize)

			err = qp.DecodeFrom(u, h, queryRequired)

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

func TestInsertQueryResult(t *testing.T) {
	t.Parallel()

	t.Run("EncodeTo includes the correct headers", func(t *testing.T) {
		fn := func(contentType string, contentLength uint) bool {
			if contentType == "" {
				return true
			}

			var (
				qp InsertQueryParams

				u, err = url.Parse("")
				h      = make(http.Header)
			)
			if err != nil {
				t.Fatal(err)
			}

			var (
				size    = int64(contentLength % (defaultMaxContentLength - 1))
				fmtSize = strconv.FormatInt(size, 10)
			)

			h.Set("content-type", contentType)
			h.Set("content-length", fmtSize)

			err = qp.DecodeFrom(u, h, queryRequired)
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()

			res := InsertQueryResult{Logger: log.NewNopLogger(), Params: qp}
			res.EncodeTo(recorder)

			headers := recorder.Header()
			return headers.Get(httpHeaderContentType) == defaultContentType
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("EncodeTo renders the correct content", func(t *testing.T) {
		fn := func(address, contentType string, contentLength uint) bool {
			if contentType == "" {
				return true
			}

			var (
				qp InsertQueryParams

				u, err = url.Parse("")
				h      = make(http.Header)
			)
			if err != nil {
				t.Fatal(err)
			}

			var (
				size    = int64(contentLength % (defaultMaxContentLength - 1))
				fmtSize = strconv.FormatInt(size, 10)
			)

			h.Set("content-type", contentType)
			h.Set("content-length", fmtSize)

			err = qp.DecodeFrom(u, h, queryRequired)
			if err != nil {
				t.Fatal(err)
			}

			content, err := models.BuildContent(
				models.WithAddress(address),
				models.WithSize(size),
				models.WithContentType(contentType),
			)
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()

			res := InsertQueryResult{Logger: log.NewNopLogger(), Params: qp}
			res.Content = content
			res.EncodeTo(recorder)

			var output models.Content
			if err = json.Unmarshal(recorder.Body.Bytes(), &output); err != nil {
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
}

func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	for k, v := range a {
		if v != b[k] {
			return false
		}
	}

	return true
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
