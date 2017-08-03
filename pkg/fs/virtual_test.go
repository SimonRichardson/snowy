package fs

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestVirtualFilesystem(t *testing.T) {
	t.Parallel()

	t.Run("create", func(t *testing.T) {
		dir := fmt.Sprintf("tmpdir-%d", rand.Intn(1000))
		fsys := NewVirtualFilesystem()
		testFilesystemCreate(fsys, dir, t)
	})

	t.Run("open", func(t *testing.T) {
		dir := fmt.Sprintf("tmpdir-%d", rand.Intn(1000))
		fsys := NewVirtualFilesystem()
		testFilesystemOpen(fsys, dir, t)
	})

	t.Run("rename", func(t *testing.T) {
		dir := fmt.Sprintf("tmpdir-%d", rand.Intn(1000))
		fsys := NewVirtualFilesystem()
		testFilesystemRename(fsys, dir, t)
	})

	t.Run("exists", func(t *testing.T) {
		dir := fmt.Sprintf("tmpdir-%d", rand.Intn(1000))
		fsys := NewVirtualFilesystem()
		testFilesystemExists(fsys, dir, t)
	})

	t.Run("remove", func(t *testing.T) {
		dir := fmt.Sprintf("tmpdir-%d", rand.Intn(1000))
		fsys := NewVirtualFilesystem()
		testFilesystemRemove(fsys, dir, t)
	})

	t.Run("walk", func(t *testing.T) {
		dir := fmt.Sprintf("tmpdir-%d", rand.Intn(1000))
		fsys := NewVirtualFilesystem()
		testFilesystemWalk(fsys, dir, t)
	})
}

func TestVirtualFile(t *testing.T) {
	t.Parallel()

	t.Run("name", func(t *testing.T) {
		dir := fmt.Sprintf("tmpdir-%d", rand.Intn(1000))
		fsys := NewVirtualFilesystem()
		testFileName(fsys, dir, t)
	})

	t.Run("size", func(t *testing.T) {
		dir := fmt.Sprintf("tmpdir-%d", rand.Intn(1000))
		fsys := NewVirtualFilesystem()
		testFileSize(fsys, dir, t)
	})

	t.Run("read and write", func(t *testing.T) {
		dir := fmt.Sprintf("tmpdir-%d", rand.Intn(1000))
		fsys := NewVirtualFilesystem()
		testFileReadWrite(fsys, dir, t)
	})
}
