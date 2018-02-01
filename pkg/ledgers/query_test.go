package ledgers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"sort"
	"testing"
	"testing/quick"

	"github.com/go-kit/kit/log"
	"github.com/trussle/harness/generators"
	errs "github.com/trussle/snowy/pkg/http"
	"github.com/trussle/snowy/pkg/models"
	"github.com/trussle/uuid"
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

	t.Run("DecodeFrom with empty tags", func(t *testing.T) {
		fn := func(uid uuid.UUID) bool {
			var (
				qp SelectQueryParams

				u, err = url.Parse(fmt.Sprintf("/?resource_id=%s&query.tags=", uid.String()))
			)
			if err != nil {
				t.Fatal(err)
			}

			err = qp.DecodeFrom(u, queryRequired)

			if expected, actual := true, err == nil; expected != actual {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}
			return len(qp.Tags) == 0
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("DecodeFrom with non-empty tags", func(t *testing.T) {
		fn := func(uid uuid.UUID, tags generators.ASCIISlice) bool {
			var (
				qp SelectQueryParams

				u, err = url.Parse(fmt.Sprintf("/?resource_id=%s&query.tags=%s", uid.String(), tags.String()))
			)
			if err != nil {
				t.Fatal(err)
			}

			err = qp.DecodeFrom(u, queryRequired)

			if expected, actual := true, err == nil; expected != actual {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}

			a := tags.Slice()

			sort.Strings(a)
			sort.Strings(qp.Tags)

			return reflect.DeepEqual(a, qp.Tags)
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("DecodeFrom with authorID", func(t *testing.T) {
		fn := func(uid uuid.UUID, authorID generators.ASCII) bool {
			var (
				qp SelectQueryParams

				u, err = url.Parse(fmt.Sprintf("/?resource_id=%s&query.author_id=%s", uid.String(), authorID))
			)
			if err != nil {
				t.Fatal(err)
			}

			err = qp.DecodeFrom(u, queryRequired)

			if expected, actual := true, err == nil; expected != actual {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}

			return qp.AuthorID == authorID.String()
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})
}

func TestSelectQueryResult(t *testing.T) {
	t.Parallel()

	emptyDoc, err := json.Marshal(models.Ledger{})
	if err != nil {
		t.Fatal(err)
	}

	t.Run("EncodeTo includes the correct headers", func(t *testing.T) {
		fn := func(uid uuid.UUID, tags generators.ASCIISlice) bool {
			var (
				qp SelectQueryParams

				u, err = url.Parse(fmt.Sprintf("/?resource_id=%s&query.tags=%s", uid.String(), tags.String()))
			)
			if err != nil {
				t.Fatal(err)
			}

			err = qp.DecodeFrom(u, queryRequired)
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()

			res := SelectQueryResult{Errors: errs.NewError(log.NewNopLogger()), Params: qp}
			res.EncodeTo(recorder)

			headers := recorder.Header()
			return headers.Get(httpHeaderResourceID) == uid.String() &&
				headers.Get(httpHeaderQueryTags) == tags.String()
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("EncodeTo with no ledger has correct status code", func(t *testing.T) {
		fn := func(uid uuid.UUID, tags generators.ASCIISlice) bool {
			var (
				qp SelectQueryParams

				u, err = url.Parse(fmt.Sprintf("/?resource_id=%s&query.tags=%s", uid.String(), tags.String()))
			)
			if err != nil {
				t.Fatal(err)
			}

			err = qp.DecodeFrom(u, queryRequired)
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()

			res := SelectQueryResult{Errors: errs.NewError(log.NewNopLogger()), Params: qp}
			res.EncodeTo(recorder)

			return recorder.Code == 200
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("EncodeTo with no ledger has correct body", func(t *testing.T) {
		fn := func(uid uuid.UUID, tags generators.ASCIISlice) bool {
			var (
				qp SelectQueryParams

				u, err = url.Parse(fmt.Sprintf("/?resource_id=%s&query.tags=%s", uid.String(), tags.String()))
			)
			if err != nil {
				t.Fatal(err)
			}

			err = qp.DecodeFrom(u, queryRequired)
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()

			res := SelectQueryResult{Errors: errs.NewError(log.NewNopLogger()), Params: qp}
			res.EncodeTo(recorder)

			return string(recorder.Body.Bytes()) == string(emptyDoc)+"\n"
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("EncodeTo with a ledger has correct body", func(t *testing.T) {
		fn := func(uid uuid.UUID, tags generators.ASCIISlice) bool {
			var (
				qp SelectQueryParams

				u, err = url.Parse(fmt.Sprintf("/?resource_id=%s&query.tags=%s", uid.String(), tags.String()))
			)
			if err != nil {
				t.Fatal(err)
			}

			err = qp.DecodeFrom(u, queryRequired)
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()

			res := SelectQueryResult{Errors: errs.NewError(log.NewNopLogger()), Params: qp}
			res.Ledger, err = models.BuildLedger(
				models.WithResourceID(uid),
			)
			if err != nil {
				t.Fatal(err)
			}
			res.EncodeTo(recorder)

			var doc models.Ledger
			if err := json.Unmarshal(recorder.Body.Bytes(), &doc); err != nil {
				t.Fatal(err)
			}

			return doc.ResourceID().Equals(uid)
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})
}

func TestSelectRevisionsQueryResult(t *testing.T) {
	t.Parallel()

	emptyDocs, err := json.Marshal(make([]models.Ledger, 0))
	if err != nil {
		t.Fatal(err)
	}

	t.Run("EncodeTo includes the correct headers", func(t *testing.T) {
		fn := func(uid uuid.UUID, tags generators.ASCIISlice) bool {
			var (
				qp SelectQueryParams

				u, err = url.Parse(fmt.Sprintf("/?resource_id=%s&query.tags=%s", uid.String(), tags.String()))
			)
			if err != nil {
				t.Fatal(err)
			}

			err = qp.DecodeFrom(u, queryRequired)
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()

			res := SelectRevisionsQueryResult{Errors: errs.NewError(log.NewNopLogger()), Params: qp}
			res.EncodeTo(recorder)

			headers := recorder.Header()
			return headers.Get(httpHeaderResourceID) == uid.String() &&
				headers.Get(httpHeaderQueryTags) == tags.String()
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("EncodeTo with no ledger has correct status code", func(t *testing.T) {
		fn := func(uid uuid.UUID, tags generators.ASCIISlice) bool {
			var (
				qp SelectQueryParams

				u, err = url.Parse(fmt.Sprintf("/?resource_id=%s&query.tags=%s", uid.String(), tags.String()))
			)
			if err != nil {
				t.Fatal(err)
			}

			err = qp.DecodeFrom(u, queryRequired)
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()

			res := SelectRevisionsQueryResult{Errors: errs.NewError(log.NewNopLogger()), Params: qp}
			res.EncodeTo(recorder)

			return recorder.Code == 200
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("EncodeTo with no ledger has correct body", func(t *testing.T) {
		fn := func(uid uuid.UUID, tags generators.ASCIISlice) bool {
			var (
				qp SelectQueryParams

				u, err = url.Parse(fmt.Sprintf("/?resource_id=%s&query.tags=%s", uid.String(), tags.String()))
			)
			if err != nil {
				t.Fatal(err)
			}

			err = qp.DecodeFrom(u, queryRequired)
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()

			res := SelectRevisionsQueryResult{Errors: errs.NewError(log.NewNopLogger()), Params: qp}
			res.EncodeTo(recorder)

			return string(recorder.Body.Bytes()) == string(emptyDocs)+"\n"
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("EncodeTo with a ledger has correct body", func(t *testing.T) {
		fn := func(uid uuid.UUID, tags generators.ASCIISlice) bool {
			var (
				qp SelectQueryParams

				u, err = url.Parse(fmt.Sprintf("/?resource_id=%s&query.tags=%s", uid.String(), tags.String()))
			)
			if err != nil {
				t.Fatal(err)
			}

			err = qp.DecodeFrom(u, queryRequired)
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()

			docs := make([]models.Ledger, 1)
			docs[0], err = models.BuildLedger(
				models.WithResourceID(uid),
			)
			if err != nil {
				t.Fatal(err)
			}

			res := SelectRevisionsQueryResult{Errors: errs.NewError(log.NewNopLogger()), Params: qp}
			res.Ledgers = docs
			res.EncodeTo(recorder)

			var resDocs []models.Ledger
			if err := json.Unmarshal(recorder.Body.Bytes(), &resDocs); err != nil {
				t.Fatal(err)
			}

			if expected, actual := len(resDocs), 1; expected != actual {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}

			return resDocs[0].ResourceID().Equals(uid)
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})
}

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

			h.Set("Content-Type", "application/json")

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

		h.Set("Content-Type", "application/json")

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

		h.Set("Content-Type", "application/json")

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

		h.Set("Content-Type", "application/json")

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

			h.Set("Content-Type", "application/json")

			err = qp.DecodeFrom(u, h, queryRequired)

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

func TestForkQueryParams(t *testing.T) {
	t.Parallel()

	t.Run("DecodeFrom with required empty url", func(t *testing.T) {
		var (
			qp ForkQueryParams

			u, err = url.Parse("")
			h      = make(http.Header, 0)
		)
		if err != nil {
			t.Fatal(err)
		}

		h.Set("Content-Type", "application/json")

		err = qp.DecodeFrom(u, h, queryRequired)

		if expected, actual := false, err == nil; expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("DecodeFrom with optional empty url", func(t *testing.T) {
		var (
			qp ForkQueryParams

			u, err = url.Parse("")
			h      = make(http.Header, 0)
		)
		if err != nil {
			t.Fatal(err)
		}

		h.Set("Content-Type", "application/json")

		err = qp.DecodeFrom(u, h, queryOptional)

		if expected, actual := true, err == nil; expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("DecodeFrom with invalid resource_id", func(t *testing.T) {
		var (
			qp ForkQueryParams

			u, err = url.Parse("/?resource_id=123asd")
			h      = make(http.Header, 0)
		)
		if err != nil {
			t.Fatal(err)
		}

		h.Set("Content-Type", "application/json")

		err = qp.DecodeFrom(u, h, queryRequired)

		if expected, actual := false, err == nil; expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("DecodeFrom with invalid content-type", func(t *testing.T) {
		fn := func(uid uuid.UUID, contentType generators.ASCII) bool {
			var (
				qp ForkQueryParams

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
				qp ForkQueryParams

				u, err = url.Parse(fmt.Sprintf("/?resource_id=%s", uid.String()))
				h      = make(http.Header, 0)
			)
			if err != nil {
				t.Fatal(err)
			}

			h.Set("Content-Type", "application/json")

			err = qp.DecodeFrom(u, h, queryRequired)

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
