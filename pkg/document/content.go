package document

// Content is abstraction over a potential underlying file, which has additional
// meta data information that could be useful.
type Content interface {

	// Size returns the size of the content
	Size() int64

	// Bytes return the content as a series of bytes.
	Bytes() []byte
}
