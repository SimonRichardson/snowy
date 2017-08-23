package repository

import (
	"errors"
	"reflect"
	"testing"
	"testing/quick"
	"time"

	"github.com/go-kit/kit/log"
	gomock "github.com/golang/mock/gomock"
	"github.com/trussle/snowy/pkg/fs"
	"github.com/trussle/snowy/pkg/models"
	"github.com/trussle/snowy/pkg/store"
	storeMocks "github.com/trussle/snowy/pkg/store/mocks"
	"github.com/trussle/snowy/pkg/uuid"
)

func TestGetLedger(t *testing.T) {
	t.Parallel()

	t.Run("get ledger with invalid id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fn := func(uid uuid.UUID) bool {
			var (
				fsys = fs.NewVirtualFilesystem()
				mock = storeMocks.NewMockStore(ctrl)
				repo = NewRealRepository(fsys, mock, log.NewNopLogger())
			)

			mock.EXPECT().
				Get(uid, store.Query{}).
				Return(store.Entity{}, errNotFound{errors.New("not found")})

			_, err := repo.GetLedger(uid, Query{})
			if expected, actual := true, ErrNotFound(err); expected != actual {
				t.Errorf("expected: %t, actual: %t", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("get ledger with store error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fn := func(uid uuid.UUID) bool {
			var (
				fsys = fs.NewVirtualFilesystem()
				mock = storeMocks.NewMockStore(ctrl)
				repo = NewRealRepository(fsys, mock, log.NewNopLogger())
			)

			mock.EXPECT().
				Get(uid, store.Query{}).
				Return(store.Entity{}, errors.New("not found"))

			_, err := repo.GetLedger(uid, Query{})
			if expected, actual := false, ErrNotFound(err); expected != actual {
				t.Errorf("expected: %t, actual: %t", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("get ledger", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fn := func(id, uid uuid.UUID) bool {
			var (
				fsys = fs.NewVirtualFilesystem()
				mock = storeMocks.NewMockStore(ctrl)
				repo = NewRealRepository(fsys, mock, log.NewNopLogger())
			)

			mock.EXPECT().
				Get(uid, store.Query{}).
				Return(store.Entity{ID: id}, nil)

			doc, err := repo.GetLedger(uid, Query{})
			if err != nil {
				t.Error(err)
			}

			if expected, actual := id, doc.ID(); expected != actual {
				t.Errorf("expected: %q, actual: %q", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})
}

func TestInsertLedger(t *testing.T) {
	t.Parallel()

	t.Run("insert ledger with store failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fn := func(resourceID uuid.UUID, authorID, name string, tags []string) bool {
			var (
				createdOn = time.Now()
				entity    = store.Entity{
					Name:       name,
					ResourceID: resourceID,
					AuthorID:   authorID,
					Tags:       tags,
					CreatedOn:  createdOn,
					DeletedOn:  time.Time{},
				}
				doc, _ = models.BuildLedger(
					models.WithName(name),
					models.WithResourceID(resourceID),
					models.WithAuthorID(authorID),
					models.WithTags(tags),
					models.WithCreatedOn(createdOn),
				)

				fsys = fs.NewVirtualFilesystem()
				mock = storeMocks.NewMockStore(ctrl)
				repo = NewRealRepository(fsys, mock, log.NewNopLogger())
			)

			mock.EXPECT().
				Insert(gomock.Eq(entity)).
				Return(errNotFound{errors.New("not found")})

			_, err := repo.InsertLedger(doc)
			if expected, actual := true, ErrNotFound(err); expected != actual {
				t.Errorf("expected: %t, actual: %t", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("insert ledger", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fn := func(resourceID uuid.UUID, authorID, name string, tags []string) bool {
			var (
				createdOn = time.Now()
				entity    = store.Entity{
					Name:       name,
					ResourceID: resourceID,
					AuthorID:   authorID,
					Tags:       tags,
					CreatedOn:  createdOn,
					DeletedOn:  time.Time{},
				}
				doc, _ = models.BuildLedger(
					models.WithResourceID(resourceID),
					models.WithName(name),
					models.WithAuthorID(authorID),
					models.WithTags(tags),
					models.WithCreatedOn(createdOn),
				)

				fsys = fs.NewVirtualFilesystem()
				mock = storeMocks.NewMockStore(ctrl)
				repo = NewRealRepository(fsys, mock, log.NewNopLogger())
			)

			mock.EXPECT().
				Insert(gomock.Eq(entity)).
				Return(nil)

			res, err := repo.InsertLedger(doc)
			if err != nil {
				t.Fatal(err)
			}

			if expected, actual := resourceID.String(), res.ResourceID().String(); expected != actual {
				t.Errorf("expected: %q, actual: %q", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})
}

func TestAppendLedger(t *testing.T) {
	t.Parallel()

	t.Run("append ledger with exists store failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fn := func(resourceID uuid.UUID, authorID, name string, tags []string) bool {
			var (
				createdOn = time.Now()
				doc, _    = models.BuildLedger(
					models.WithName(name),
					models.WithResourceID(resourceID),
					models.WithAuthorID(authorID),
					models.WithTags(tags),
					models.WithCreatedOn(createdOn),
				)

				fsys = fs.NewVirtualFilesystem()
				mock = storeMocks.NewMockStore(ctrl)
				repo = NewRealRepository(fsys, mock, log.NewNopLogger())
			)

			mock.EXPECT().
				Get(resourceID, store.Query{}).
				Return(store.Entity{}, errNotFound{errors.New("not found")})

			_, err := repo.AppendLedger(resourceID, doc)
			if expected, actual := true, ErrNotFound(err); expected != actual {
				t.Errorf("expected: %t, actual: %t", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("append ledger with insert store failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fn := func(resourceID uuid.UUID, authorID, name string, tags []string) bool {
			var (
				createdOn = time.Now()
				entity    = store.Entity{
					Name:       name,
					ResourceID: resourceID,
					AuthorID:   authorID,
					Tags:       tags,
					CreatedOn:  createdOn,
					DeletedOn:  time.Time{},
				}
				doc, _ = models.BuildLedger(
					models.WithName(name),
					models.WithResourceID(resourceID),
					models.WithAuthorID(authorID),
					models.WithTags(tags),
					models.WithCreatedOn(createdOn),
				)

				fsys = fs.NewVirtualFilesystem()
				mock = storeMocks.NewMockStore(ctrl)
				repo = NewRealRepository(fsys, mock, log.NewNopLogger())
			)

			mock.EXPECT().
				Get(resourceID, store.Query{}).
				Return(entity, nil)
			mock.EXPECT().
				Insert(entity).
				Return(errNotFound{errors.New("not found")})

			_, err := repo.AppendLedger(resourceID, doc)
			if expected, actual := true, ErrNotFound(err); expected != actual {
				t.Errorf("expected: %t, actual: %t", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("append ledger", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fn := func(resourceID uuid.UUID, authorID, name string, tags []string) bool {
			var (
				createdOn = time.Now()
				entity    = store.Entity{
					Name:       name,
					ResourceID: resourceID,
					AuthorID:   authorID,
					Tags:       tags,
					CreatedOn:  createdOn,
					DeletedOn:  time.Time{},
				}
				doc, _ = models.BuildLedger(
					models.WithResourceID(resourceID),
					models.WithName(name),
					models.WithAuthorID(authorID),
					models.WithTags(tags),
					models.WithCreatedOn(createdOn),
				)

				fsys = fs.NewVirtualFilesystem()
				mock = storeMocks.NewMockStore(ctrl)
				repo = NewRealRepository(fsys, mock, log.NewNopLogger())
			)

			mock.EXPECT().
				Get(resourceID, store.Query{}).
				Return(entity, nil)
			mock.EXPECT().
				Insert(entity).
				Return(nil)

			res, err := repo.AppendLedger(resourceID, doc)
			if err != nil {
				t.Fatal(err)
			}

			if expected, actual := resourceID.String(), res.ResourceID().String(); expected != actual {
				t.Errorf("expected: %q, actual: %q", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})
}

func TestGetLedgers(t *testing.T) {
	t.Parallel()

	t.Run("get ledgers with invalid id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fn := func(uid uuid.UUID) bool {
			var (
				fsys = fs.NewVirtualFilesystem()
				mock = storeMocks.NewMockStore(ctrl)
				repo = NewRealRepository(fsys, mock, log.NewNopLogger())
			)

			mock.EXPECT().
				GetMultiple(uid, store.Query{}).
				Return([]store.Entity{}, errNotFound{errors.New("not found")})

			_, err := repo.GetLedgers(uid, Query{})
			if expected, actual := false, err == nil; expected != actual {
				t.Errorf("expected: %t, actual: %t", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("get ledgers with store error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fn := func(uid uuid.UUID) bool {
			var (
				fsys = fs.NewVirtualFilesystem()
				mock = storeMocks.NewMockStore(ctrl)
				repo = NewRealRepository(fsys, mock, log.NewNopLogger())
			)

			mock.EXPECT().
				GetMultiple(uid, store.Query{}).
				Return([]store.Entity{}, errors.New("not found"))

			_, err := repo.GetLedgers(uid, Query{})
			if expected, actual := false, err == nil; expected != actual {
				t.Errorf("expected: %t, actual: %t", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("get ledgers", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fn := func(id, uid uuid.UUID) bool {
			var (
				fsys = fs.NewVirtualFilesystem()
				mock = storeMocks.NewMockStore(ctrl)
				repo = NewRealRepository(fsys, mock, log.NewNopLogger())
			)

			mock.EXPECT().
				GetMultiple(uid, store.Query{}).
				Return([]store.Entity{store.Entity{ID: id}}, nil)

			doc, err := repo.GetLedgers(uid, Query{})
			if err != nil {
				t.Error(err)
			}

			if expected, actual := 1, len(doc); expected != actual {
				t.Errorf("expected: %d, actual: %d", expected, actual)
			}

			if expected, actual := id, doc[0].ID(); expected != actual {
				t.Errorf("expected: %q, actual: %q", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})
}

func TestGetContent(t *testing.T) {
	t.Parallel()

	t.Run("get content with no ledger", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fn := func(id, uid uuid.UUID) bool {
			var (
				fsys = fs.NewVirtualFilesystem()
				mock = storeMocks.NewMockStore(ctrl)
				repo = NewRealRepository(fsys, mock, log.NewNopLogger())
			)

			mock.EXPECT().
				Get(uid, store.Query{}).
				Return(store.Entity{}, errNotFound{errors.New("not found")})

			_, err := repo.GetContent(uid, Query{})

			if expected, actual := true, ErrNotFound(err); expected != actual {
				t.Errorf("expected: %t, actual: %t", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("get content with store failure for ledger", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fn := func(id, uid uuid.UUID) bool {
			var (
				fsys = fs.NewVirtualFilesystem()
				mock = storeMocks.NewMockStore(ctrl)
				repo = NewRealRepository(fsys, mock, log.NewNopLogger())
			)

			mock.EXPECT().
				Get(uid, store.Query{}).
				Return(store.Entity{}, errors.New("not found"))

			_, err := repo.GetContent(uid, Query{})

			if expected, actual := false, err == nil; expected != actual {
				t.Errorf("expected: %t, actual: %t", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("get content with no file", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fn := func(id, uid uuid.UUID) bool {
			var (
				fsys = fs.NewVirtualFilesystem()
				mock = storeMocks.NewMockStore(ctrl)
				repo = NewRealRepository(fsys, mock, log.NewNopLogger())
			)

			mock.EXPECT().
				Get(uid, store.Query{}).
				Return(store.Entity{
					ResourceID: uid,
				}, nil)

			_, err := repo.GetContent(uid, Query{})

			if expected, actual := false, err == nil; expected != actual {
				t.Errorf("expected: %t, actual: %t", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("get content with file", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fn := func(id, uid uuid.UUID, body []byte) bool {
			var (
				fsys = fs.NewVirtualFilesystem()
				mock = storeMocks.NewMockStore(ctrl)
				repo = NewRealRepository(fsys, mock, log.NewNopLogger())
			)

			file, err := fsys.Create(uid.String())
			if err != nil {
				t.Fatal(err)
			}

			if _, err = file.Write(body); err != nil {
				t.Fatal(err)
			}

			mock.EXPECT().
				Get(uid, store.Query{}).
				Return(store.Entity{
					ResourceID:      uid,
					ResourceAddress: uid.String(),
				}, nil)

			content, err := repo.GetContent(uid, Query{})
			if expected, actual := true, err == nil; expected != actual {
				t.Errorf("expected: %t, actual: %t", expected, actual)
			}

			b, err := content.Bytes()
			if err != nil {
				t.Fatal(err)
			}

			return reflect.DeepEqual(b, body)
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})
}

func TestGetContents(t *testing.T) {
	t.Parallel()

	t.Run("get contents with no ledger", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fn := func(id, uid uuid.UUID) bool {
			var (
				fsys = fs.NewVirtualFilesystem()
				mock = storeMocks.NewMockStore(ctrl)
				repo = NewRealRepository(fsys, mock, log.NewNopLogger())
			)

			mock.EXPECT().
				GetMultiple(uid, store.Query{}).
				Return(nil, errNotFound{errors.New("not found")})

			_, err := repo.GetContents(uid, Query{})

			if expected, actual := true, ErrNotFound(err); expected != actual {
				t.Fatalf("expected: %t, actual: %t", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("get contents with store failure for ledger", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fn := func(id, uid uuid.UUID) bool {
			var (
				fsys = fs.NewVirtualFilesystem()
				mock = storeMocks.NewMockStore(ctrl)
				repo = NewRealRepository(fsys, mock, log.NewNopLogger())
			)

			mock.EXPECT().
				GetMultiple(uid, store.Query{}).
				Return(nil, errors.New("not found"))

			_, err := repo.GetContents(uid, Query{})

			if expected, actual := false, err == nil; expected != actual {
				t.Fatalf("expected: %t, actual: %t", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("get content with no file", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fn := func(id, uid uuid.UUID) bool {
			var (
				fsys = fs.NewVirtualFilesystem()
				mock = storeMocks.NewMockStore(ctrl)
				repo = NewRealRepository(fsys, mock, log.NewNopLogger())
			)

			mock.EXPECT().
				GetMultiple(uid, store.Query{}).
				Return([]store.Entity{
					store.Entity{
						ResourceID: uid,
					},
				}, nil)

			_, err := repo.GetContents(uid, Query{})

			if expected, actual := true, err == nil; expected != actual {
				t.Fatalf("expected: %t, actual: %t", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("get contents with file", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fn := func(id, uid uuid.UUID, body []byte) bool {
			var (
				fsys = fs.NewVirtualFilesystem()
				mock = storeMocks.NewMockStore(ctrl)
				repo = NewRealRepository(fsys, mock, log.NewNopLogger())
			)

			file, err := fsys.Create(uid.String())
			if err != nil {
				t.Fatal(err)
			}

			if _, err = file.Write(body); err != nil {
				t.Fatal(err)
			}

			mock.EXPECT().
				GetMultiple(uid, store.Query{}).
				Return([]store.Entity{
					store.Entity{
						ResourceID:      uid,
						ResourceAddress: uid.String(),
					},
				}, nil)

			contents, err := repo.GetContents(uid, Query{})
			if expected, actual := true, err == nil; expected != actual {
				t.Fatalf("expected: %t, actual: %t", expected, actual)
			}

			if expected, actual := 1, len(contents); expected != actual {
				t.Fatalf("expected: %d, actual: %d", expected, actual)
			}

			b, err := contents[0].Bytes()
			if err != nil {
				t.Fatal(err)
			}

			return reflect.DeepEqual(b, body)
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})
}

func TestPutContent(t *testing.T) {
	t.Parallel()

	t.Run("put without content", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fn := func(id, uid uuid.UUID, address string) bool {
			var (
				fsys = fs.NewVirtualFilesystem()
				mock = storeMocks.NewMockStore(ctrl)
				repo = NewRealRepository(fsys, mock, log.NewNopLogger())
			)

			content, err := models.BuildContent(
				models.WithAddress(address),
				models.WithContentType("application/octet-stream"),
			)
			if err != nil {
				t.Fatal(err)
			}

			_, err = repo.PutContent(content)

			if expected, actual := false, err == nil; expected != actual {
				t.Errorf("expected: %t, actual: %t", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("put content", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fn := func(id, uid uuid.UUID, body []byte, contentType string) bool {
			if len(body) < 1 {
				return true
			}

			var (
				fsys = fs.NewVirtualFilesystem()
				mock = storeMocks.NewMockStore(ctrl)
				repo = NewRealRepository(fsys, mock, log.NewNopLogger())
			)

			content, err := models.BuildContent(
				models.WithContentBytes(body),
				models.WithSize(int64(len(body))),
				models.WithContentType(contentType),
			)
			if err != nil {
				t.Fatal(err)
			}

			res, err := repo.PutContent(content)

			if expected, actual := true, err == nil; expected != actual {
				t.Errorf("expected: %t, actual: %t", expected, actual)
			}

			return content.Address() == res.Address() &&
				content.ContentType() == res.ContentType()

		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})
}
