package models

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"
)

// Content is abstraction over a potential underlying file, which has additional
// meta data information that could be useful.
type Content struct {
	address     string
	size        int64
	contentType string
	reader      io.ReadCloser
}

// Address returns the content addressable value
func (c Content) Address() string {
	return c.address
}

// Size returns the size of the content
func (c Content) Size() int64 {
	return c.size
}

// ContentType returns the MIME type of the content body.
func (c Content) ContentType() string {
	return c.contentType
}

// Bytes returns the body of the content as a slice of bytes.
func (c Content) Bytes() ([]byte, error) {
	if c.reader == nil {
		return make([]byte, 0), nil
	}

	defer c.reader.Close()
	return ioutil.ReadAll(c.reader)
}

// Reader returns the body of the content as a unread reader.
func (c Content) Reader() io.ReadCloser {
	return c.reader
}

// MarshalJSON converts a UUID into a serialisable json format
func (c Content) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Address     string `json:"address"`
		Size        int64  `json:"size"`
		ContentType string `json:"content_type"`
	}{
		Address:     c.address,
		Size:        c.size,
		ContentType: c.contentType,
	})
}

// UnmarshalJSON unserialises the json format
func (c *Content) UnmarshalJSON(b []byte) error {
	var res struct {
		Address     string `json:"address"`
		Size        int64  `json:"size"`
		ContentType string `json:"content_type"`
	}
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}

	c.address = res.Address
	c.size = res.Size
	c.contentType = res.ContentType

	return nil
}

// ContentAddress gets the addressable value of the content from the body of the
// content.
func ContentAddress(bytes []byte) (string, error) {
	hash := sha256.New()
	if _, err := hash.Write(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// ContentOption defines a option for generating a content
type ContentOption func(*Content) error

// BuildContent ingests configuration options to then yield a Content and returns a
// error if it fails during setup.
func BuildContent(opts ...ContentOption) (Content, error) {
	var content Content
	for _, opt := range opts {
		err := opt(&content)
		if err != nil {
			return Content{}, err
		}
	}
	return content, nil
}

// WithAddress adds a Address to the content
func WithAddress(address string) ContentOption {
	return func(content *Content) error {
		content.address = address
		return nil
	}
}

// WithSize adds a size to the content
func WithSize(size int64) ContentOption {
	return func(content *Content) error {
		content.size = size
		return nil
	}
}

// WithContentType adds a ContentType to the content
func WithContentType(contentType string) ContentOption {
	return func(content *Content) error {
		content.contentType = contentType
		return nil
	}
}

// WithBytes adds a body to the content
func WithBytes(b []byte) ContentOption {
	return func(content *Content) error {
		content.reader = ioutil.NopCloser(bytes.NewBuffer(b))
		return nil
	}
}

// WithReader adds a reader for the content
func WithReader(r io.ReadCloser) ContentOption {
	return func(content *Content) error {
		content.reader = r
		return nil
	}
}

// WithContentBytes adds a content and address in one step
func WithContentBytes(b []byte) ContentOption {
	return func(content *Content) error {
		address, err := ContentAddress(b)
		if err != nil {
			return err
		}

		content.address = address
		content.reader = ioutil.NopCloser(bytes.NewBuffer(b))
		content.size = int64(len(b))
		return nil
	}
}
