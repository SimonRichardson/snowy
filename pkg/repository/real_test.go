package repository

import (
	"errors"
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
				doc, _ = document.Build(
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
				doc, _ = document.Build(
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
