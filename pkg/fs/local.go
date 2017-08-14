package fs

import (
	"io"
	"os"
	"path/filepath"
)

type localFilesystem struct{}

// NewLocalFilesystem yields a local disk filesystem.
func NewLocalFilesystem() Filesystem {
	return localFilesystem{}
}

func (localFilesystem) Create(path string) (File, error) {
	f, err := os.Create(path)
	return localFile{
		File:   f,
		Reader: f,
		Closer: f,
	}, err
}

func (fs localFilesystem) Open(path string) (File, error) {
	f, err := os.Open(path)
	if err != nil {
		if err == os.ErrNotExist {
			return nil, errNotFound{err}
		}
		return nil, err
	}

	return localFile{
		File:   f,
		Reader: f,
		Closer: f,
	}, nil
}

func (localFilesystem) Rename(oldname, newname string) error {
	return os.Rename(oldname, newname)
}

func (localFilesystem) Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func (localFilesystem) Remove(path string) error {
	return os.Remove(path)
}

func (localFilesystem) Walk(root string, walkFn filepath.WalkFunc) error {
	return filepath.Walk(root, walkFn)
}

type localFile struct {
	*os.File
	io.Reader
	io.Closer
}

func (f localFile) Read(p []byte) (int, error) {
	return f.Reader.Read(p)
}

func (f localFile) Close() error {
	return f.Closer.Close()
}

func (f localFile) Size() int64 {
	fi, err := f.File.Stat()
	if err != nil {
		panic(err)
	}
	return fi.Size()
}

func (f localFile) WriteContentType(t string) error { return nil }
func (f localFile) ContentType() string             { return defaultContentType }
