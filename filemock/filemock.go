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

// NewFileWrapper creates a new FileWrapper struct
//
// Parameters: file - the file to wrap
//
// Returns: a pointer to a new FileWrapper struct
func NewFileWrapper(file *os.File) *FileWrapper {
	return &FileWrapper{file: file}
}

// Close closes the file
//
// Parameters: _
//
// Returns: an error if the file cannot be closed, nil if the file is closed successfully
func (f *FileWrapper) Close() error {
	return f.file.Close()
}

type FileMock struct {
	FileName string
	Open     bool
	Err      error
}

// Close closes the mock file
//
// Parameters: _
//
// Returns: an error if the file cannot be closed, nil if the file is closed successfully
func (file *FileMock) Close() error {
	if file == nil {
		return os.ErrInvalid
	}
	if file.Open {
		file.Open = false
	} else {
		return os.ErrClosed
	}
	return file.Err
}
