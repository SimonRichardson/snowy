package repository

import (
	"errors"
	"reflect"
	"testing"
	"testing/quick"
	"time"

	"github.com/go-kit/kit/log"
	gomock "github.com/golang/mock/gomock"
	"github.com/trussle/fsys"
	"github.com/trussle/snowy/pkg/models"
	"github.com/trussle/snowy/pkg/store"
	storeMocks "github.com/trussle/snowy/pkg/store/mocks"
	"github.com/trussle/uuid"
)

func TestSelectLedger(t *testing.T) {
	t.Parallel()

	t.Run("get ledger with invalid id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fn := func(uid uuid.UUID) bool {
			var (
				fsys = fsys.NewVirtualFilesystem()
				mock = storeMocks.NewMockStore(ctrl)
				repo = NewRealRepository(fsys, mock, log.NewNopLogger())
			)

			mock.EXPECT().
				Select(uid, store.Query{}).
				Return(store.Entity{}, errNotFound{errors.New("not found")})

			_, err := repo.SelectLedger(uid, Query{})
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
				fsys = fsys.NewVirtualFilesystem()
				mock = storeMocks.NewMockStore(ctrl)
				repo = NewRealRepository(fsys, mock, log.NewNopLogger())
			)

			mock.EXPECT().
				Select(uid, store.Query{}).
				Return(store.Entity{}, errors.New("not found"))

			_, err := repo.SelectLedger(uid, Query{})
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
				fsys = fsys.NewVirtualFilesystem()
				mock = storeMocks.NewMockStore(ctrl)
				repo = NewRealRepository(fsys, mock, log.NewNopLogger())
			)

			mock.EXPECT().
				Select(uid, store.Query{}).
				Return(store.Entity{ID: id}, nil)

			doc, err := repo.SelectLedger(uid, Query{})
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
					ParentID:   uuid.MustParse(defaultRootParentID),
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

				fsys = fsys.NewVirtualFilesystem()
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
					ParentID:   uuid.MustParse(defaultRootParentID),
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

				fsys = fsys.NewVirtualFilesystem()
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

				fsys = fsys.NewVirtualFilesystem()
				mock = storeMocks.NewMockStore(ctrl)
				repo = NewRealRepository(fsys, mock, log.NewNopLogger())
			)

			mock.EXPECT().
				Select(resourceID, store.Query{}).
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

				fsys = fsys.NewVirtualFilesystem()
				mock = storeMocks.NewMockStore(ctrl)
				repo = NewRealRepository(fsys, mock, log.NewNopLogger())
			)

			mock.EXPECT().
				Select(resourceID, store.Query{}).
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

				fsys = fsys.NewVirtualFilesystem()
				mock = storeMocks.NewMockStore(ctrl)
				repo = NewRealRepository(fsys, mock, log.NewNopLogger())
			)

			mock.EXPECT().
				Select(resourceID, store.Query{}).
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

func TestForkLedger(t *testing.T) {
	t.Parallel()

	t.Run("fork ledger with exists store failure", func(t *testing.T) {
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

				fsys = fsys.NewVirtualFilesystem()
				mock = storeMocks.NewMockStore(ctrl)
				repo = NewRealRepository(fsys, mock, log.NewNopLogger())
			)

			mock.EXPECT().
				Select(resourceID, store.Query{}).
				Return(store.Entity{}, errNotFound{errors.New("not found")})

			_, err := repo.ForkLedger(resourceID, doc)
			if expected, actual := true, ErrNotFound(err); expected != actual {
				t.Errorf("expected: %t, actual: %t", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("fork ledger with insert store failure", func(t *testing.T) {
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

				fsys = fsys.NewVirtualFilesystem()
				mock = storeMocks.NewMockStore(ctrl)
				repo = NewRealRepository(fsys, mock, log.NewNopLogger())
			)

			mock.EXPECT().
				Select(resourceID, store.Query{}).
				Return(entity, nil)
			mock.EXPECT().
				Insert(entity).
				Return(errNotFound{errors.New("not found")})

			_, err := repo.ForkLedger(resourceID, doc)
			if expected, actual := true, ErrNotFound(err); expected != actual {
				t.Errorf("expected: %t, actual: %t", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("fork ledger", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fn := func(id, resourceID, parentID, forkedID uuid.UUID, authorID, name string, tags []string) bool {
			var (
				createdOn    = time.Now()
				parentEntity = store.Entity{
					ID:         id,
					Name:       name,
					ResourceID: resourceID,
					ParentID:   parentID,
					AuthorID:   authorID,
					Tags:       tags,
					CreatedOn:  createdOn,
					DeletedOn:  time.Time{},
				}
				forkedEntity = store.Entity{
					ID:         uuid.MustNew(),
					Name:       name,
					ResourceID: forkedID,
					ParentID:   id,
					AuthorID:   authorID,
					Tags:       tags,
					CreatedOn:  createdOn,
					DeletedOn:  time.Time{},
				}
				doc, _ = models.BuildLedger(
					models.WithResourceID(forkedID),
					models.WithParentID(parentID),
					models.WithName(name),
					models.WithAuthorID(authorID),
					models.WithTags(tags),
					models.WithCreatedOn(createdOn),
					models.WithDeletedOn(time.Time{}),
				)

				fsys = fsys.NewVirtualFilesystem()
				mock = storeMocks.NewMockStore(ctrl)
				repo = NewRealRepository(fsys, mock, log.NewNopLogger())
			)

			mock.EXPECT().
				Select(resourceID, store.Query{}).
				Return(parentEntity, nil)
			mock.EXPECT().
				Insert(Entity(forkedEntity)).
				Return(nil)

			res, err := repo.ForkLedger(resourceID, doc)
			if err != nil {
				t.Fatal(err)
			}

			if expected, actual := forkedID.String(), res.ResourceID().String(); expected != actual {
				t.Errorf("expected: %q, actual: %q", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})
}

func TestSelectLedgers(t *testing.T) {
	t.Parallel()

	t.Run("get ledgers with invalid id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fn := func(uid uuid.UUID) bool {
			var (
				fsys = fsys.NewVirtualFilesystem()
				mock = storeMocks.NewMockStore(ctrl)
				repo = NewRealRepository(fsys, mock, log.NewNopLogger())
			)

			mock.EXPECT().
				SelectRevisions(uid, store.Query{}).
				Return([]store.Entity{}, errNotFound{errors.New("not found")})

			_, err := repo.SelectLedgers(uid, Query{})
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
				fsys = fsys.NewVirtualFilesystem()
				mock = storeMocks.NewMockStore(ctrl)
				repo = NewRealRepository(fsys, mock, log.NewNopLogger())
			)

			mock.EXPECT().
				SelectRevisions(uid, store.Query{}).
				Return([]store.Entity{}, errors.New("not found"))

			_, err := repo.SelectLedgers(uid, Query{})
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
				fsys = fsys.NewVirtualFilesystem()
				mock = storeMocks.NewMockStore(ctrl)
				repo = NewRealRepository(fsys, mock, log.NewNopLogger())
			)

			mock.EXPECT().
				SelectRevisions(uid, store.Query{}).
				Return([]store.Entity{store.Entity{ID: id}}, nil)

			doc, err := repo.SelectLedgers(uid, Query{})
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

func TestSelectForkLedgers(t *testing.T) {
	t.Parallel()

	t.Run("get ledgers with invalid id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fn := func(uid uuid.UUID) bool {
			var (
				fsys = fsys.NewVirtualFilesystem()
				mock = storeMocks.NewMockStore(ctrl)
				repo = NewRealRepository(fsys, mock, log.NewNopLogger())
			)

			mock.EXPECT().
				SelectForkRevisions(uid).
				Return([]store.Entity{}, errNotFound{errors.New("not found")})

			_, err := repo.SelectForkLedgers(uid)
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
				fsys = fsys.NewVirtualFilesystem()
				mock = storeMocks.NewMockStore(ctrl)
				repo = NewRealRepository(fsys, mock, log.NewNopLogger())
			)

			mock.EXPECT().
				SelectForkRevisions(uid).
				Return([]store.Entity{}, errors.New("not found"))

			_, err := repo.SelectForkLedgers(uid)
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
				fsys = fsys.NewVirtualFilesystem()
				mock = storeMocks.NewMockStore(ctrl)
				repo = NewRealRepository(fsys, mock, log.NewNopLogger())
			)

			mock.EXPECT().
				SelectForkRevisions(uid).
				Return([]store.Entity{store.Entity{ID: id}}, nil)

			doc, err := repo.SelectForkLedgers(uid)
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

func TestSelectContent(t *testing.T) {
	t.Parallel()

	t.Run("get content with no ledger", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fn := func(id, uid uuid.UUID) bool {
			var (
				fsys = fsys.NewVirtualFilesystem()
				mock = storeMocks.NewMockStore(ctrl)
				repo = NewRealRepository(fsys, mock, log.NewNopLogger())
			)

			mock.EXPECT().
				Select(uid, store.Query{}).
				Return(store.Entity{}, errNotFound{errors.New("not found")})

			_, err := repo.SelectContent(uid, Query{})

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
				fsys = fsys.NewVirtualFilesystem()
				mock = storeMocks.NewMockStore(ctrl)
				repo = NewRealRepository(fsys, mock, log.NewNopLogger())
			)

			mock.EXPECT().
				Select(uid, store.Query{}).
				Return(store.Entity{}, errors.New("not found"))

			_, err := repo.SelectContent(uid, Query{})

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
				fsys = fsys.NewVirtualFilesystem()
				mock = storeMocks.NewMockStore(ctrl)
				repo = NewRealRepository(fsys, mock, log.NewNopLogger())
			)

			mock.EXPECT().
				Select(uid, store.Query{}).
				Return(store.Entity{
					ResourceID: uid,
				}, nil)

			_, err := repo.SelectContent(uid, Query{})

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
				fsys = fsys.NewVirtualFilesystem()
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
				Select(uid, store.Query{}).
				Return(store.Entity{
					ResourceID:      uid,
					ResourceAddress: uid.String(),
				}, nil)

			content, err := repo.SelectContent(uid, Query{})
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

	t.Run("get same content multiple times", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fn := func(id, uid uuid.UUID, body []byte) bool {
			var (
				fsys = fsys.NewVirtualFilesystem()
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

			total := 3

			mock.EXPECT().
				Select(uid, store.Query{}).
				Return(store.Entity{
					ResourceID:      uid,
					ResourceAddress: uid.String(),
				}, nil).
				Times(total)

			for i := 0; i < total; i++ {
				content, err := repo.SelectContent(uid, Query{})
				if expected, actual := true, err == nil; expected != actual {
					t.Errorf("expected: %t, actual: %t, for iteration: %d", expected, actual, i)
				}

				b, err := content.Bytes()
				if err != nil {
					t.Fatal(err)
				}

				if expected, actual := true, reflect.DeepEqual(b, body); expected != actual {
					t.Errorf("expected: %t, actual: %t, for iteration: %d", expected, actual, i)
				}
			}
			return true
		}

		if err := quick.Check(fn, &quick.Config{MaxCount: 1}); err != nil {
			t.Error(err)
		}
	})
}

func TestSelectContents(t *testing.T) {
	t.Parallel()

	t.Run("get contents with no ledger", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fn := func(id, uid uuid.UUID) bool {
			var (
				fsys = fsys.NewVirtualFilesystem()
				mock = storeMocks.NewMockStore(ctrl)
				repo = NewRealRepository(fsys, mock, log.NewNopLogger())
			)

			mock.EXPECT().
				SelectRevisions(uid, store.Query{}).
				Return(nil, errNotFound{errors.New("not found")})

			_, err := repo.SelectContents(uid, Query{})

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
				fsys = fsys.NewVirtualFilesystem()
				mock = storeMocks.NewMockStore(ctrl)
				repo = NewRealRepository(fsys, mock, log.NewNopLogger())
			)

			mock.EXPECT().
				SelectRevisions(uid, store.Query{}).
				Return(nil, errors.New("not found"))

			_, err := repo.SelectContents(uid, Query{})

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
				fsys = fsys.NewVirtualFilesystem()
				mock = storeMocks.NewMockStore(ctrl)
				repo = NewRealRepository(fsys, mock, log.NewNopLogger())
			)

			mock.EXPECT().
				SelectRevisions(uid, store.Query{}).
				Return([]store.Entity{
					store.Entity{
						ResourceID: uid,
					},
				}, nil)

			_, err := repo.SelectContents(uid, Query{})

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
				fsys = fsys.NewVirtualFilesystem()
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
				SelectRevisions(uid, store.Query{}).
				Return([]store.Entity{
					store.Entity{
						ResourceID:      uid,
						ResourceAddress: uid.String(),
					},
				}, nil)

			contents, err := repo.SelectContents(uid, Query{})
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
				fsys = fsys.NewVirtualFilesystem()
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
				fsys = fsys.NewVirtualFilesystem()
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

type entityMatcher struct {
	doc store.Entity
}

func (m entityMatcher) Matches(x interface{}) bool {
	d, ok := x.(store.Entity)
	if !ok {
		return false
	}

	return m.doc.ParentID.Equals(d.ParentID) &&
		m.doc.Name == d.Name &&
		m.doc.AuthorID == d.AuthorID &&
		reflect.DeepEqual(m.doc.Tags, d.Tags)
}

func (entityMatcher) String() string {
	return "is entity"
}

func Entity(doc store.Entity) gomock.Matcher { return entityMatcher{doc} }
