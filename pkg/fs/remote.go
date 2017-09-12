package fs

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3crypto"
	"github.com/pkg/errors"
)

const (
	defaultContentType = "application/octet-stream"
)

// RemoteAccessType defines what type of access to the remote S3 is going to be
// utilized.
type RemoteAccessType int

const (
	// KMSRemoteAccessType uses the KMS access
	KMSRemoteAccessType RemoteAccessType = iota

	// KeySecretRemoteAccessType uses the traditional Key Secret access
	KeySecretRemoteAccessType
)

// RemoteConfig creates a configuration to create a RemoteFilesystem.
type RemoteConfig struct {
	Type                         RemoteAccessType
	KMSKey, ServerSideEncryption string
	ID, Secret, Token            string
	Region, Bucket               string
}

func (c RemoteConfig) String() string {
	switch c.Type {
	case KMSRemoteAccessType:
		return fmt.Sprintf("RemoteConfig (encryption: true, Region: %q, Bucket: %q, SSE: %q)", c.Region, c.Bucket, c.ServerSideEncryption)
	default:
		return fmt.Sprintf("RemoteConfig (encryption: false, Region: %q, Bucket: %q, ID: %q)", c.Region, c.Bucket, c.ID)
	}
}

type remoteFilesystem struct {
	client remoteClient
	bucket *string
}

// NewRemoteFilesystem creates a new remote file system that abstracts over a S3
// bucket.
func NewRemoteFilesystem(config *RemoteConfig) (Filesystem, error) {
	creds := credentials.NewStaticCredentials(
		config.ID,
		config.Secret,
		config.Token,
	)
	if _, err := creds.Get(); err != nil {
		return nil, errors.Wrap(err, "invalid credentials")
	}

	var (
		client remoteClient
		cfg    = aws.NewConfig().
			WithRegion(config.Region).
			WithCredentials(creds).
			WithCredentialsChainVerboseErrors(true)
		sess = session.New(cfg)
	)

	switch config.Type {
	case KMSRemoteAccessType:
		if len(config.KMSKey) == 0 {
			return nil, errors.Errorf("expected valid KMSKey")
		}

		var (
			kms       = kms.New(sess)
			generator = s3crypto.NewKMSKeyGenerator(kms, config.KMSKey)
			crypto    = s3crypto.AESGCMContentCipherBuilder(generator)
		)
		client = newCryptoS3Client(
			s3.New(sess),
			s3crypto.NewEncryptionClient(sess, crypto),
			s3crypto.NewDecryptionClient(sess),
			config.KMSKey,
			config.ServerSideEncryption,
		)

	case KeySecretRemoteAccessType:
		client = newS3Client(s3.New(sess))
	default:
		return nil, errors.Errorf("invalid remote config type %v", config.Type)
	}

	return &remoteFilesystem{
		client: client,
		bucket: aws.String(config.Bucket),
	}, nil
}

func (fs *remoteFilesystem) Create(path string) (File, error) {
	file := newRemoteFile(fs,
		ioutil.NopCloser(bytes.NewReader(make([]byte, 0))),
		path,
		0,
	)
	return file, file.Sync()
}

func (fs *remoteFilesystem) Open(path string) (File, error) {
	object, err := fs.client.GetObject(&s3.GetObjectInput{
		Bucket: fs.bucket,
		Key:    aws.String(path),
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == s3.ErrCodeNoSuchKey {
			return nil, errNotFound{err}
		}
		return nil, err
	}

	rf := newRemoteFile(fs,
		object.Body,
		path,
		*object.ContentLength,
	)
	rf.contentType = *object.ContentType
	return rf, nil
}

func (fs *remoteFilesystem) Rename(oldname, newname string) error {
	source := fmt.Sprintf("%s/%s", *fs.bucket, oldname)

	if _, err := fs.client.CopyObject(&s3.CopyObjectInput{
		Bucket:     fs.bucket,
		Key:        aws.String(newname),
		CopySource: aws.String(source),
	}); err != nil {
		return err
	}
	return fs.Remove(oldname)
}

func (fs *remoteFilesystem) Exists(path string) bool {
	_, err := fs.client.GetObject(&s3.GetObjectInput{
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
	_, err := fs.client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: fs.bucket,
		Key:    aws.String(path),
	})
	return err
}

func (fs *remoteFilesystem) Walk(root string, walkFn filepath.WalkFunc) error {
	objects, err := fs.client.ListObjects(&s3.ListObjectsInput{
		Bucket: fs.bucket,
		Prefix: aws.String(root),
	})
	if err != nil {
		return err
	}

	var loopErr error
	for _, v := range objects.Contents {
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
	sys               *remoteFilesystem
	writeContent      []byte
	readContent       io.ReadCloser
	path, contentType string
	size              int64
}

func newRemoteFile(sys *remoteFilesystem, body io.ReadCloser, path string, size int64) *remoteFile {
	return &remoteFile{
		sys,
		make([]byte, 0),
		body,
		path,
		defaultContentType,
		size,
	}
}

func (f *remoteFile) Write(p []byte) (int, error) {
	f.writeContent = append(f.writeContent, p...)
	return len(p), nil
}

func (f *remoteFile) Read(p []byte) (int, error) {
	if f.readContent != nil {
		return f.readContent.Read(p)
	}
	return 0, nil
}

func (f *remoteFile) Close() error {
	if f.readContent != nil {
		return f.readContent.Close()
	}
	return nil
}

func (f *remoteFile) Name() string { return f.path }
func (f *remoteFile) Size() int64  { return f.size }

func (f *remoteFile) Sync() error {
	_, err := f.sys.client.PutObject(&s3.PutObjectInput{
		Bucket:      f.sys.bucket,
		Key:         aws.String(f.path),
		Body:        bytes.NewReader(f.writeContent),
		ContentType: aws.String(f.contentType),
	})

	return err
}

func (f *remoteFile) WriteContentType(t string) error {
	f.contentType = t
	return nil
}

func (f *remoteFile) ContentType() string {
	return f.contentType
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

// WithEncryption adds an ID option to the configuration
func WithEncryption(encryption bool) ConfigOption {
	return func(config *RemoteConfig) error {
		if encryption {
			config.Type = KMSRemoteAccessType
		} else {
			config.Type = KeySecretRemoteAccessType
		}
		return nil
	}
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

// WithKMSKey adds a KMSKey option to the configuration
func WithKMSKey(kmsKey string) ConfigOption {
	return func(config *RemoteConfig) error {
		config.KMSKey = kmsKey
		return nil
	}
}

// WithServerSideEncryption adds a ServerSideEncryption option to the configuration
func WithServerSideEncryption(serverSideEncryption string) ConfigOption {
	return func(config *RemoteConfig) error {
		config.ServerSideEncryption = serverSideEncryption
		return nil
	}
}

// Various clients to work with

type remoteClient interface {
	GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error)
	CopyObject(input *s3.CopyObjectInput) (*s3.CopyObjectOutput, error)
	DeleteObject(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error)
	ListObjects(input *s3.ListObjectsInput) (*s3.ListObjectsOutput, error)
	PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error)
}

type cryptoS3Client struct {
	service     *s3.S3
	encrypt     *s3crypto.EncryptionClient
	decrypt     *s3crypto.DecryptionClient
	kmsKey, sse string
}

func newCryptoS3Client(service *s3.S3,
	encrypt *s3crypto.EncryptionClient,
	decrypt *s3crypto.DecryptionClient,
	kmsKey, sse string,
) remoteClient {
	return &cryptoS3Client{
		service,
		encrypt,
		decrypt,
		kmsKey,
		sse,
	}
}

func (c *cryptoS3Client) GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	return c.service.GetObject(input)
}

func (c *cryptoS3Client) CopyObject(input *s3.CopyObjectInput) (*s3.CopyObjectOutput, error) {
	return c.service.CopyObject(input)
}

func (c *cryptoS3Client) DeleteObject(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
	return c.service.DeleteObject(input)
}

func (c *cryptoS3Client) ListObjects(input *s3.ListObjectsInput) (*s3.ListObjectsOutput, error) {
	return c.service.ListObjects(input)
}

func (c *cryptoS3Client) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	input.SetSSEKMSKeyId(c.kmsKey)
	input.SetServerSideEncryption(c.sse)

	return c.service.PutObject(input)
}

type s3Client struct {
	service *s3.S3
}

func newS3Client(service *s3.S3) remoteClient {
	return &s3Client{service}
}

func (c *s3Client) GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	return c.service.GetObject(input)
}

func (c *s3Client) CopyObject(input *s3.CopyObjectInput) (*s3.CopyObjectOutput, error) {
	return c.service.CopyObject(input)
}

func (c *s3Client) DeleteObject(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
	return c.service.DeleteObject(input)
}

func (c *s3Client) ListObjects(input *s3.ListObjectsInput) (*s3.ListObjectsOutput, error) {
	return c.service.ListObjects(input)
}

func (c *s3Client) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	return c.service.PutObject(input)
}
