package gindump

import "bytes"

// LimitedBuffer is a wrapper of bytes.Buffer which limits its size.
// The buffer will store no more than the specified size, and discards the rest.
// This is useful for preventing gindump from storing large request/response body in memory.
// LimitedBuffer implements Writer and Reader interface.
type LimitedBuffer struct {
	buf bytes.Buffer
	size int
}

// Create a new LimitedBuffer with a limited size in bytes.
func NewLimitedBuffer(size int) *LimitedBuffer {
	return &LimitedBuffer{
		buf: bytes.Buffer{},
		size: size,
	}
}

// Read implements the io.Reader interface.
func (lb *LimitedBuffer) Read(p []byte) (n int, err error) {
	return lb.buf.Read(p)
}

// Write implements the io.Writer interface.
func (lb *LimitedBuffer) Write(p []byte) (n int, err error) {
	if lb.buf.Len() >= lb.size {
		return 0, nil
	}

	sizeToWrite := lb.size - lb.buf.Len()
	if sizeToWrite >= len(p) {
		sizeToWrite = len(p)
	}
	return lb.buf.Write(p[:sizeToWrite])
}

// Bytes returns the bytes of the underlying buffer.
func (lb *LimitedBuffer) Bytes() []byte {
	return lb.buf.Bytes()
}
