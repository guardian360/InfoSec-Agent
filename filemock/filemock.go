package filemock

import (
	"os"
)

type File interface {
	Close() error
}

// FileWrapper is a wrapper for the os.File type
type FileWrapper struct {
	file *os.File
}

// TODO: Fix this one
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

// TODO: Fix this one
type FileMock struct {
	FileName string
	Open     bool
	Err      error
}

// TODO: Fix this one
// Close closes the mock file
//
// Parameters: _
//
// Returns: an error if the file cannot be closed, nil if the file is closed successfully
func (f *FileMock) Close() error {
	if f == nil {
		return os.ErrInvalid
	}
	if f.Open {
		f.Open = false
	} else {
		return os.ErrClosed
	}
	return f.Err
}
