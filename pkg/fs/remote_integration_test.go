// +build integration

package fs

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"reflect"
	"syscall"
	"testing"
)

const (
	defaultAWSID     = ""
	defaultAWSSecret = ""
	defaultAWSToken  = ""
	defaultAWSRegion = "eu-west-1"
	defaultAWSBucket = ""
	defaultAWSKMSKey = ""
	defaultAWSSSE    = "aws:kms"
)

type ByteSize uint64

const (
	B  ByteSize = 1
	KB          = B << 10
)

func TestRemoteFilesystemWithoutEncryption_Integration(t *testing.T) {
	t.Parallel()

	config, err := BuildConfig(
		WithRegion(GetEnv("AWS_REGION", defaultAWSRegion)),
		WithBucket(GetEnv("AWS_BUCKET", defaultAWSBucket)),

		WithEncryption(false),
		WithServerSideEncryption(GetEnv("AWS_SSE", defaultAWSSSE)),
		WithKMSKey(GetEnv("AWS_KMSKEY", defaultAWSKMSKey)),

		WithID(GetEnv("AWS_ID", defaultAWSID)),
		WithSecret(GetEnv("AWS_SECRET", defaultAWSSecret)),
		WithToken(GetEnv("AWS_TOKEN", defaultAWSToken)),
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("new", func(t *testing.T) {
		fsys, err := NewRemoteFilesystem(config)
		if err != nil {
			t.Fatal(err)
		}

		if expected, actual := false, fsys == nil; expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}
	})

	t.Run("create then remove empty file", func(t *testing.T) {
		fsys, err := NewRemoteFilesystem(config)
		if err != nil {
			t.Fatal(err)
		}

		b := make([]byte, 10*KB)
		if _, err := rand.Read(b); err != nil {
			t.Fatal(err)
		}

		path, err := ContentAddress(b)
		if err != nil {
			t.Fatal(err)
		}

		file, err := fsys.Create(path)
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()

		if expected, actual := true, fsys.Exists(path); expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}

		if fsys.Remove(path); err != nil {
			t.Fatal(err)
		}

		if expected, actual := false, fsys.Exists(path); expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}
	})

	t.Run("create then remove", func(t *testing.T) {
		fsys, err := NewRemoteFilesystem(config)
		if err != nil {
			t.Fatal(err)
		}

		b := make([]byte, 10*KB)
		if _, err := rand.Read(b); err != nil {
			t.Fatal(err)
		}

		path, err := ContentAddress(b)
		if err != nil {
			t.Fatal(err)
		}

		file, err := fsys.Create(path)
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()

		if expected, actual := true, fsys.Exists(path); expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}

		if _, err := file.Write(b); err != nil {
			t.Fatal(err)
		}

		if err := file.Sync(); err != nil {
			t.Fatal(err)
		}

		if fsys.Remove(path); err != nil {
			t.Fatal(err)
		}

		if expected, actual := false, fsys.Exists(path); expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}
	})

	t.Run("create then open", func(t *testing.T) {
		fsys, err := NewRemoteFilesystem(config)
		if err != nil {
			t.Fatal(err)
		}

		input := make([]byte, 10*KB)
		if _, err := rand.Read(input); err != nil {
			t.Fatal(err)
		}

		file, close := createFile(t, fsys, input)
		defer close()

		f, err := fsys.Open(file.Name())
		if err != nil {
			t.Fatal(err)
		}

		output, err := ioutil.ReadAll(f)
		if err != nil {
			t.Fatal(err)
		}

		if expected, actual := input, output; !reflect.DeepEqual(expected, actual) {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("open file has correct size", func(t *testing.T) {
		fsys, err := NewRemoteFilesystem(config)
		if err != nil {
			t.Fatal(err)
		}

		input := make([]byte, ByteSize(rand.Intn(10)+1)*KB)
		if _, err := rand.Read(input); err != nil {
			t.Fatal(err)
		}

		file, close := createFile(t, fsys, input)
		defer close()

		f, err := fsys.Open(file.Name())
		if err != nil {
			t.Fatal(err)
		}

		output, err := ioutil.ReadAll(f)
		if err != nil {
			t.Fatal(err)
		}

		if expected, actual := len(input), len(output); expected != actual {
			t.Errorf("expected: %d, actual: %d", expected, actual)
		}
		if expected, actual := int64(len(input)), f.Size(); expected != actual {
			t.Errorf("expected: %d, actual: %d", expected, actual)
		}
	})

	t.Run("open file has correct content-type", func(t *testing.T) {
		fsys, err := NewRemoteFilesystem(config)
		if err != nil {
			t.Fatal(err)
		}

		input := make([]byte, 1*KB)
		if _, err := rand.Read(input); err != nil {
			t.Fatal(err)
		}

		file, close := createFile(t, fsys, input)
		defer close()

		f, err := fsys.Open(file.Name())
		if err != nil {
			t.Fatal(err)
		}

		if expected, actual := defaultContentType, f.ContentType(); expected != actual {
			t.Errorf("expected: %s, actual: %s", expected, actual)
		}
	})

	t.Run("renaming", func(t *testing.T) {
		fsys, err := NewRemoteFilesystem(config)
		if err != nil {
			t.Fatal(err)
		}

		input := make([]byte, 1*KB)
		if _, err := rand.Read(input); err != nil {
			t.Fatal(err)
		}

		// Don't close as it's moved
		file, _ := createFile(t, fsys, input)

		newPath := fmt.Sprintf("%s-new", file.Name())
		if err := fsys.Rename(file.Name(), newPath); err != nil {
			t.Fatal(err)
		}

		if expected, actual := false, fsys.Exists(file.Name()); expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}
		if expected, actual := true, fsys.Exists(newPath); expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}

		if err := fsys.Remove(newPath); err != nil {
			t.Fatal(err)
		}
		if expected, actual := false, fsys.Exists(newPath); expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}
	})

	t.Run("walk", func(t *testing.T) {
		fsys, err := NewRemoteFilesystem(config)
		if err != nil {
			t.Fatal(err)
		}

		input := make([]byte, 1*KB)
		if _, err := rand.Read(input); err != nil {
			t.Fatal(err)
		}

		file, close := createFile(t, fsys, input)
		defer close()

		var (
			infos  []os.FileInfo
			called = false
		)
		err = fsys.Walk("", func(path string, info os.FileInfo, err error) error {
			called = true
			infos = append(infos, info)
			return err
		})

		if expected, actual := true, err == nil; expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}
		if expected, actual := true, called; expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}

		contains := func(haystack []os.FileInfo, needle File) bool {
			for _, v := range haystack {
				if v.Name() == needle.Name() {
					return true
				}
			}
			return false
		}

		if expected, actual := true, contains(infos, file); expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}
	})
}

func TestRemoteFilesystemWithEncryption_Integration(t *testing.T) {
	t.Parallel()

	config, err := BuildConfig(
		WithRegion(GetEnv("AWS_REGION", defaultAWSRegion)),
		WithBucket(GetEnv("AWS_BUCKET", defaultAWSBucket)),

		WithEncryption(true),
		WithServerSideEncryption(GetEnv("AWS_SSE", defaultAWSSSE)),
		WithKMSKey(GetEnv("AWS_KMSKEY", defaultAWSKMSKey)),

		WithID(GetEnv("AWS_ID", defaultAWSID)),
		WithSecret(GetEnv("AWS_SECRET", defaultAWSSecret)),
		WithToken(GetEnv("AWS_TOKEN", defaultAWSToken)),
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("new", func(t *testing.T) {
		fsys, err := NewRemoteFilesystem(config)
		if err != nil {
			t.Fatal(err)
		}

		if expected, actual := false, fsys == nil; expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}
	})

	t.Run("create then remove empty file", func(t *testing.T) {
		fsys, err := NewRemoteFilesystem(config)
		if err != nil {
			t.Fatal(err)
		}

		b := make([]byte, 10*KB)
		if _, err := rand.Read(b); err != nil {
			t.Fatal(err)
		}

		path, err := ContentAddress(b)
		if err != nil {
			t.Fatal(err)
		}

		file, err := fsys.Create(path)
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()

		if expected, actual := true, fsys.Exists(path); expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}

		if err := fsys.Remove(path); err != nil {
			t.Fatal(err)
		}

		if expected, actual := false, fsys.Exists(path); expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}
	})

	t.Run("create then open", func(t *testing.T) {
		fsys, err := NewRemoteFilesystem(config)
		if err != nil {
			t.Fatal(err)
		}

		input := make([]byte, 10*KB)
		if _, err := rand.Read(input); err != nil {
			t.Fatal(err)
		}

		file, close := createFile(t, fsys, input)
		defer close()

		f, err := fsys.Open(file.Name())
		if err != nil {
			t.Fatal(err)
		}

		output, err := ioutil.ReadAll(f)
		if err != nil {
			t.Fatal(err)
		}

		if expected, actual := input, output; !reflect.DeepEqual(expected, actual) {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("open file has correct size", func(t *testing.T) {
		fsys, err := NewRemoteFilesystem(config)
		if err != nil {
			t.Fatal(err)
		}

		input := make([]byte, ByteSize(rand.Intn(10)+1)*KB)
		if _, err := rand.Read(input); err != nil {
			t.Fatal(err)
		}

		file, close := createFile(t, fsys, input)
		defer close()

		f, err := fsys.Open(file.Name())
		if err != nil {
			t.Fatal(err)
		}

		output, err := ioutil.ReadAll(f)
		if err != nil {
			t.Fatal(err)
		}

		if expected, actual := len(input), len(output); expected != actual {
			t.Errorf("expected: %d, actual: %d", expected, actual)
		}
	})

	t.Run("open file has correct content-type", func(t *testing.T) {
		fsys, err := NewRemoteFilesystem(config)
		if err != nil {
			t.Fatal(err)
		}

		input := make([]byte, 1*KB)
		if _, err := rand.Read(input); err != nil {
			t.Fatal(err)
		}

		file, close := createFile(t, fsys, input)
		defer close()

		f, err := fsys.Open(file.Name())
		if err != nil {
			t.Fatal(err)
		}

		if expected, actual := defaultContentType, f.ContentType(); expected != actual {
			t.Errorf("expected: %s, actual: %s", expected, actual)
		}
	})

	t.Run("renaming", func(t *testing.T) {
		fsys, err := NewRemoteFilesystem(config)
		if err != nil {
			t.Fatal(err)
		}

		input := make([]byte, 1*KB)
		if _, err := rand.Read(input); err != nil {
			t.Fatal(err)
		}

		// Don't close as it's moved
		file, _ := createFile(t, fsys, input)

		newPath := fmt.Sprintf("%s-new", file.Name())
		if err := fsys.Rename(file.Name(), newPath); err != nil {
			t.Fatal(err)
		}

		if expected, actual := false, fsys.Exists(file.Name()); expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}
		if expected, actual := true, fsys.Exists(newPath); expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}

		if err := fsys.Remove(newPath); err != nil {
			t.Fatal(err)
		}
		if expected, actual := false, fsys.Exists(newPath); expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}
	})

	t.Run("walk", func(t *testing.T) {
		fsys, err := NewRemoteFilesystem(config)
		if err != nil {
			t.Fatal(err)
		}

		input := make([]byte, 1*KB)
		if _, err := rand.Read(input); err != nil {
			t.Fatal(err)
		}

		file, close := createFile(t, fsys, input)
		defer close()

		var (
			infos  []os.FileInfo
			called = false
		)
		err = fsys.Walk("", func(path string, info os.FileInfo, err error) error {
			called = true
			infos = append(infos, info)
			return err
		})

		if expected, actual := true, err == nil; expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}
		if expected, actual := true, called; expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}

		contains := func(haystack []os.FileInfo, needle File) bool {
			for _, v := range haystack {
				if v.Name() == needle.Name() {
					return true
				}
			}
			return false
		}

		if expected, actual := true, contains(infos, file); expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}
	})
}

func GetEnv(key string, defaultValue string) (value string) {
	var ok bool
	if value, ok = syscall.Getenv(key); ok {
		return
	}
	return defaultValue
}

func ContentAddress(bytes []byte) (string, error) {
	hash := sha256.New()
	if _, err := hash.Write(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(hash.Sum(nil)), nil
}

func createFile(t *testing.T, fsys Filesystem, b []byte) (File, func()) {
	path, err := ContentAddress(b)
	if err != nil {
		t.Fatal(err)
	}

	file, err := fsys.Create(path)
	if err != nil {
		t.Fatal(err)
	}

	if expected, actual := true, fsys.Exists(path); expected != actual {
		t.Errorf("expected: %t, actual: %t", expected, actual)
	}

	if _, err := file.Write(b); err != nil {
		t.Fatal(err)
	}

	if err := file.Sync(); err != nil {
		t.Fatal(err)
	}

	return file, func() {
		file.Close()

		if err := fsys.Remove(path); err != nil {
			t.Fatal(err)
		}

		if expected, actual := false, fsys.Exists(path); expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}
	}
}
