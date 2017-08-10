package fs

import (
	"bytes"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
)

// RemoteConfig creates a configuration to create a RemoteFilesystem.
type RemoteConfig struct {
	ID, Secret, Token, Region, Bucket string
}

type remoteFilesystem struct {
	service *s3.S3
	bucket  *string
}

// NewRemoteFilesystem creates a new remote file system that abstracts over a S3
// bucket.
func NewRemoteFilesystem(config *RemoteConfig) (Filesystem, error) {
	creds := credentials.NewStaticCredentials(config.ID, config.Secret, config.Token)
	if _, err := creds.Get(); err != nil {
		return nil, errors.Wrap(err, "invalid credentials")
	}

	cfg := aws.NewConfig().WithRegion(config.Region).WithCredentials(creds)
	return &remoteFilesystem{
		service: s3.New(session.New(), cfg),
		bucket:  aws.String(config.Bucket),
	}, nil
}

func (fs *remoteFilesystem) Create(path string) (File, error) {
	return newRemoteFile(fs,
		ioutil.NopCloser(bytes.NewReader(make([]byte, 0))),
		path,
		0,
	), nil
}

func (fs *remoteFilesystem) Open(path string) (File, error) {
	object, err := fs.service.GetObject(&s3.GetObjectInput{
		Bucket: fs.bucket,
		Key:    aws.String(path),
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == s3.ErrCodeNoSuchKey {
			return nil, errNotFound{err}
		}
		return nil, err
	}

	return newRemoteFile(fs,
		object.Body,
		path,
		*object.ContentLength,
	), nil
}

func (fs *remoteFilesystem) Rename(oldname, newname string) error {
	if _, err := fs.service.CopyObject(&s3.CopyObjectInput{
		Bucket:     fs.bucket,
		Key:        aws.String(newname),
		CopySource: aws.String(oldname),
	}); err != nil {
		return err
	}
	return fs.Remove(oldname)
}

func (fs *remoteFilesystem) Exists(path string) bool {
	_, err := fs.service.GetObject(&s3.GetObjectInput{
		Bucket: fs.bucket,
		Key:    aws.String(path),
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == s3.ErrCodeNoSuchKey {
			return false
		}
	}
	return true
}

func (fs *remoteFilesystem) Remove(path string) error {
	_, err := fs.service.DeleteObject(&s3.DeleteObjectInput{
		Bucket: fs.bucket,
		Key:    aws.String(path),
	})
	return err
}

func (fs *remoteFilesystem) Walk(root string, walkFn filepath.WalkFunc) error {
	objects, err := fs.service.ListObjects(&s3.ListObjectsInput{
		Bucket: fs.bucket,
		Prefix: aws.String(root),
	})
	if err != nil {
		return err
	}

	var loopErr error
	for _, v := range objects.Contents {
		if strings.HasPrefix(*v.Key, root) {
			continue
		}

		info := virtualFileInfo{
			name:  *v.Key,
			mtime: *v.LastModified,
			size:  *v.Size,
		}

		if err := walkFn(*v.Key, info, loopErr); err != nil {
			loopErr = err
		}
	}

	return loopErr
}

type remoteFile struct {
	sys          *remoteFilesystem
	writeContent []byte
	readContent  io.ReadCloser
	path         string
	size         int64
}

func newRemoteFile(sys *remoteFilesystem, body io.ReadCloser, path string, size int64) *remoteFile {
	return &remoteFile{
		sys,
		make([]byte, 0),
		body,
		path,
		size,
	}
}

func (f *remoteFile) Write(p []byte) (int, error) {
	f.writeContent = append(f.writeContent, p...)
	return len(p), nil
}

func (f *remoteFile) Read(p []byte) (int, error) {
	return f.readContent.Read(p)
}

func (f *remoteFile) Close() error {
	return f.readContent.Close()
}

func (f *remoteFile) Name() string { return f.path }
func (f *remoteFile) Size() int64  { return f.size }

const defaultFileType = "text/plain; charset=utf-8"

func (f *remoteFile) Sync() error {
	_, err := f.sys.service.PutObject(&s3.PutObjectInput{
		Bucket:        f.sys.bucket,
		Key:           aws.String(f.path),
		Body:          bytes.NewReader(f.writeContent),
		ContentLength: aws.Int64(int64(len(f.writeContent))),
		ContentType:   aws.String(defaultFileType),
	})
	return err
}

// ConfigOption defines a option for generating a RemoteConfig
type ConfigOption func(*RemoteConfig) error

// BuildConfig ingests configuration options to then yield a
// RemoteConfig, and return an error if it fails during configuring.
func BuildConfig(opts ...ConfigOption) (*RemoteConfig, error) {
	var config RemoteConfig
	for _, opt := range opts {
		err := opt(&config)
		if err != nil {
			return nil, err
		}
	}
	return &config, nil
}

// WithID adds an ID option to the configuration
func WithID(id string) ConfigOption {
	return func(config *RemoteConfig) error {
		config.ID = id
		return nil
	}
}

// WithSecret adds an Secret option to the configuration
func WithSecret(secret string) ConfigOption {
	return func(config *RemoteConfig) error {
		config.Secret = secret
		return nil
	}
}

// WithToken adds an Token option to the configuration
func WithToken(token string) ConfigOption {
	return func(config *RemoteConfig) error {
		config.Token = token
		return nil
	}
}

// WithRegion adds an Region option to the configuration
func WithRegion(region string) ConfigOption {
	return func(config *RemoteConfig) error {
		config.Region = region
		return nil
	}
}

// WithBucket adds an Bucket option to the configuration
func WithBucket(bucket string) ConfigOption {
	return func(config *RemoteConfig) error {
		config.Bucket = bucket
		return nil
	}
}
