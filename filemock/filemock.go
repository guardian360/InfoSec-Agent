package filemock

import (
	"io"
	"os"
)

type File interface {
	Close() error
	io.Closer
	Read(p []byte) (int, error)
	Write(p []byte) (int, error)
	Copy(destination File, source File) (int64, error)
}

// FileWrapper is a wrapper for the os.File type
type FileWrapper struct {
	file   *os.File
	Err    error
	Writer *os.File
	Reader *os.File
}

// Wrap creates a new FileWrapper struct
//
// Parameters: file - the file to wrap
//
// Returns: a pointer to a new FileWrapper struct
func Wrap(file *os.File) *FileWrapper {
	return &FileWrapper{file: file, Writer: file, Reader: file}
}

// Close closes the file
//
// Parameters: _
//
// Returns: an error if the file cannot be closed, nil if the file is closed successfully
func (f *FileWrapper) Close() error {
	return f.file.Close()
}

// TODO: fix this docstring according to new documentation standard

// Read reads from the file
//
// Parameters: p - the buffer to read into
//
// Returns: the number of bytes read and an error if the file cannot be read from
func (f *FileWrapper) Read(p []byte) (int, error) {
	return f.file.Read(p)
}

// TODO: fix this docstring according to new documentation standard

// Write writes to the file
//
// Parameters: p - the buffer to write from
//
// Returns: the number of bytes written and an error if the file cannot be written to
func (f *FileWrapper) Write(p []byte) (int, error) {
	return f.file.Write(p)
}

func (f *FileWrapper) Copy(source File, destination File) (int64, error) {
	return io.Copy(destination, source)
}

type FileMock struct {
	FileName string
	IsOpen   bool
	Buffer   []byte
	Bytes    int
	Err      error
}

// Close closes the mock file
//
// Parameters: _
//
// Returns: an error if the file cannot be closed, nil if the file is closed successfully
func (f *FileMock) Close() error {
	if f == nil {
		return os.ErrInvalid
	}
	if f.IsOpen {
		f.IsOpen = false
	} else {
		return os.ErrClosed
	}
	return f.Err
}

// TODO: fix this docstring according to new documentation standard

// Read reads from the mock file
//
// Parameters: p - the buffer to read into
//
// Returns: the number of bytes read and an error if the file cannot be read from
func (f *FileMock) Read(_ []byte) (int, error) {
	return f.Bytes, f.Err
}

// TODO: fix this docstring according to new documentation standard

// Write writes to the mock file
//
// Parameters: p - the buffer to write from
//
// Returns: the number of bytes written and an error if the file cannot be written to
func (f *FileMock) Write(_ []byte) (int, error) {
	return f.Bytes, f.Err
}

func (f *FileMock) Copy(source File, destination File) (int64, error) {
	var err error
	_, err = source.Read([]byte{})
	if err != nil {
		return 0, err
	}
	_, err = destination.Read([]byte{})
	if err != nil {
		return 0, err
	}
	return int64(f.Bytes), f.Err
}
