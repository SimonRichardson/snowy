package repository

import (
	"errors"
	"testing"

	"github.com/trussle/snowy/pkg/document"
	"github.com/trussle/snowy/pkg/fs"
	"github.com/trussle/snowy/pkg/store"
	"github.com/trussle/snowy/pkg/uuid"
	"github.com/go-kit/kit/log"
	gomock "github.com/golang/mock/gomock"
)

func TestRealRepository(t *testing.T) {
	t.Parallel()

	t.Run("get document with invalid id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			fsys = fs.NewVirtualFilesystem()
			mock = store.NewMockStore(ctrl)
			repo = NewRealRepository(fsys, mock, log.NewNopLogger())

			uid = uuid.New()
		)

		mock.EXPECT().
			Get(uid).
			Return(store.Entity{}, errNotFound{errors.New("not found")})

		_, err := repo.GetDocument(uid)
		if expected, actual := true, ErrNotFound(err); expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}
	})

	t.Run("get document", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			fsys = fs.NewVirtualFilesystem()
			mock = store.NewMockStore(ctrl)
			repo = NewRealRepository(fsys, mock, log.NewNopLogger())

			id  = uuid.New().String()
			uid = uuid.New()
		)

		mock.EXPECT().
			Get(uid).
			Return(store.Entity{ID: id}, nil)

		doc, err := repo.GetDocument(uid)
		if err != nil {
			t.Error(err)
		}

		if expected, actual := id, doc.ID(); expected != actual {
			t.Errorf("expected: %q, actual: %q", expected, actual)
		}
	})

	t.Run("put document with store failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			fsys = fs.NewVirtualFilesystem()
			mock = store.NewMockStore(ctrl)
			repo = NewRealRepository(fsys, mock, log.NewNopLogger())

			uid    = uuid.New()
			entity = store.Entity{
				ID:         uuid.New().String(),
				ResourceID: uid,
			}
			doc, _ = document.Build(
				document.WithID(entity.ID),
				document.WithResourceID(uid),
			)
		)

		mock.EXPECT().
			Put(entity).
			Return(errNotFound{errors.New("not found")})

		_, err := repo.PutDocument(doc)
		if expected, actual := true, ErrNotFound(err); expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}
	})

	t.Run("put document", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			fsys = fs.NewVirtualFilesystem()
			mock = store.NewMockStore(ctrl)
			repo = NewRealRepository(fsys, mock, log.NewNopLogger())

			uid    = uuid.New()
			entity = store.Entity{
				ID:         uuid.New().String(),
				ResourceID: uid,
			}
			doc, _ = document.Build(
				document.WithID(entity.ID),
				document.WithResourceID(uid),
			)
		)

		mock.EXPECT().
			Put(entity).
			Return(nil)

		id, err := repo.PutDocument(doc)
		if err != nil {
			t.Error(err)
		}

		if expected, actual := uid.String(), id.String(); expected != actual {
			t.Errorf("expected: %q, actual: %q", expected, actual)
		}
	})
}
