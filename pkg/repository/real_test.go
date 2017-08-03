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
	"github.com/trussle/snowy/pkg/uuid"
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

		fn := func(id, resourceID, authorID uuid.UUID, name string, tags []string) bool {
			var (
				createdOn = time.Now()
				entity    = store.Entity{
					ID:         id.String(),
					Name:       name,
					ResourceID: resourceID,
					AuthorID:   authorID,
					Tags:       tags,
					CreatedOn:  createdOn,
					DeletedOn:  time.Time{},
				}
				doc, _ = document.Build(
					document.WithID(entity.ID),
					document.WithName(name),
					document.WithResourceID(resourceID),
					document.WithAuthorID(authorID),
					document.WithTags(tags),
					document.WithCreatedOn(createdOn),
					document.WithDeletedOn(time.Time{}),
				)

				fsys = fs.NewVirtualFilesystem()
				mock = store.NewMockStore(ctrl)
				repo = NewRealRepository(fsys, mock, log.NewNopLogger())
			)

			mock.EXPECT().
				Put(gomock.Eq(entity)).
				Return(errNotFound{errors.New("not found")})

			_, err := repo.PutDocument(doc)
			if expected, actual := true, ErrNotFound(err); expected != actual {
				t.Errorf("expected: %t, actual: %t", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("put document", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fn := func(id, resourceID, authorID uuid.UUID, name string, tags []string) bool {
			var (
				createdOn = time.Now()
				entity    = store.Entity{
					ID:         id.String(),
					Name:       name,
					ResourceID: resourceID,
					AuthorID:   authorID,
					Tags:       tags,
					CreatedOn:  createdOn,
					DeletedOn:  time.Time{},
				}
				doc, _ = document.Build(
					document.WithID(entity.ID),
					document.WithName(name),
					document.WithResourceID(resourceID),
					document.WithAuthorID(authorID),
					document.WithTags(tags),
					document.WithCreatedOn(createdOn),
					document.WithDeletedOn(time.Time{}),
				)

				fsys = fs.NewVirtualFilesystem()
				mock = store.NewMockStore(ctrl)
				repo = NewRealRepository(fsys, mock, log.NewNopLogger())
			)

			mock.EXPECT().
				Put(gomock.Eq(entity)).
				Return(nil)

			res, err := repo.PutDocument(doc)
			if err != nil {
				t.Fatal(err)
			}

			if expected, actual := resourceID.String(), res.String(); expected != actual {
				t.Errorf("expected: %q, actual: %q", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})
}
