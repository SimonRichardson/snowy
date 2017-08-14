package repository

import (
	"errors"
	"reflect"
	"testing"
	"testing/quick"
	"time"

	"github.com/go-kit/kit/log"
	gomock "github.com/golang/mock/gomock"
	"github.com/trussle/snowy/pkg/document"
	"github.com/trussle/snowy/pkg/fs"
	"github.com/trussle/snowy/pkg/store"
	storeMocks "github.com/trussle/snowy/pkg/store/mocks"
	"github.com/trussle/snowy/pkg/uuid"
)

func TestGetDocument(t *testing.T) {
	t.Parallel()

	t.Run("get document with invalid id", func(t *testing.T) {
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

			_, err := repo.GetDocument(uid, Query{})
			if expected, actual := true, ErrNotFound(err); expected != actual {
				t.Errorf("expected: %t, actual: %t", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("get document with store error", func(t *testing.T) {
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

			_, err := repo.GetDocument(uid, Query{})
			if expected, actual := false, ErrNotFound(err); expected != actual {
				t.Errorf("expected: %t, actual: %t", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("get document", func(t *testing.T) {
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

			doc, err := repo.GetDocument(uid, Query{})
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

func TestInsertDocument(t *testing.T) {
	t.Parallel()

	t.Run("insert document with store failure", func(t *testing.T) {
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
				doc, _ = document.BuildDocument(
					document.WithName(name),
					document.WithResourceID(resourceID),
					document.WithAuthorID(authorID),
					document.WithTags(tags),
					document.WithCreatedOn(createdOn),
				)

				fsys = fs.NewVirtualFilesystem()
				mock = storeMocks.NewMockStore(ctrl)
				repo = NewRealRepository(fsys, mock, log.NewNopLogger())
			)

			mock.EXPECT().
				Insert(gomock.Eq(entity)).
				Return(errNotFound{errors.New("not found")})

			_, err := repo.InsertDocument(doc)
			if expected, actual := true, ErrNotFound(err); expected != actual {
				t.Errorf("expected: %t, actual: %t", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("insert document", func(t *testing.T) {
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
				doc, _ = document.BuildDocument(
					document.WithResourceID(resourceID),
					document.WithName(name),
					document.WithAuthorID(authorID),
					document.WithTags(tags),
					document.WithCreatedOn(createdOn),
				)

				fsys = fs.NewVirtualFilesystem()
				mock = storeMocks.NewMockStore(ctrl)
				repo = NewRealRepository(fsys, mock, log.NewNopLogger())
			)

			mock.EXPECT().
				Insert(gomock.Eq(entity)).
				Return(nil)

			res, err := repo.InsertDocument(doc)
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

func TestAppendDocument(t *testing.T) {
	t.Parallel()

	t.Run("append document with exists store failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fn := func(resourceID uuid.UUID, authorID, name string, tags []string) bool {
			var (
				createdOn = time.Now()
				doc, _    = document.BuildDocument(
					document.WithName(name),
					document.WithResourceID(resourceID),
					document.WithAuthorID(authorID),
					document.WithTags(tags),
					document.WithCreatedOn(createdOn),
				)

				fsys = fs.NewVirtualFilesystem()
				mock = storeMocks.NewMockStore(ctrl)
				repo = NewRealRepository(fsys, mock, log.NewNopLogger())
			)

			mock.EXPECT().
				Get(resourceID, store.Query{}).
				Return(store.Entity{}, errNotFound{errors.New("not found")})

			_, err := repo.AppendDocument(resourceID, doc)
			if expected, actual := true, ErrNotFound(err); expected != actual {
				t.Errorf("expected: %t, actual: %t", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("append document with insert store failure", func(t *testing.T) {
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
				doc, _ = document.BuildDocument(
					document.WithName(name),
					document.WithResourceID(resourceID),
					document.WithAuthorID(authorID),
					document.WithTags(tags),
					document.WithCreatedOn(createdOn),
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

			_, err := repo.AppendDocument(resourceID, doc)
			if expected, actual := true, ErrNotFound(err); expected != actual {
				t.Errorf("expected: %t, actual: %t", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("append document", func(t *testing.T) {
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
				doc, _ = document.BuildDocument(
					document.WithResourceID(resourceID),
					document.WithName(name),
					document.WithAuthorID(authorID),
					document.WithTags(tags),
					document.WithCreatedOn(createdOn),
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

			res, err := repo.AppendDocument(resourceID, doc)
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

func TestGetDocuments(t *testing.T) {
	t.Parallel()

	t.Run("get documents with invalid id", func(t *testing.T) {
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

			_, err := repo.GetDocuments(uid, Query{})
			if expected, actual := false, err == nil; expected != actual {
				t.Errorf("expected: %t, actual: %t", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("get documents with store error", func(t *testing.T) {
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

			_, err := repo.GetDocuments(uid, Query{})
			if expected, actual := false, err == nil; expected != actual {
				t.Errorf("expected: %t, actual: %t", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("get documents", func(t *testing.T) {
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

			doc, err := repo.GetDocuments(uid, Query{})
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

	t.Run("get content with no document", func(t *testing.T) {
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

			_, err := repo.GetContent(uid)

			if expected, actual := true, ErrNotFound(err); expected != actual {
				t.Errorf("expected: %t, actual: %t", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("get content with store failure for document", func(t *testing.T) {
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

			_, err := repo.GetContent(uid)

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

			_, err := repo.GetContent(uid)

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

			content, err := repo.GetContent(uid)
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

			content, err := document.BuildContent(
				document.WithAddress(address),
				document.WithContentType("application/octet-stream"),
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

			content, err := document.BuildContent(
				document.WithContentBytes(body),
				document.WithSize(int64(len(body))),
				document.WithContentType(contentType),
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
