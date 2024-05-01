package mocking

import (
	"io"
	"math"
	"os"
)

type File interface {
	Close() error
	io.Closer
	Read(p []byte) (int, error)
	Write(p []byte) (int, error)
	Seek(offset int64, whence int) (int64, error)
	Copy(source File, destination File) (int64, error)
}

// FileWrapper is a wrapper for the os.File type
type FileWrapper struct {
	file   *os.File
	Buffer []byte
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
	fileInfo, err := file.Stat()
	if err != nil {
		return nil
	}
	fileSize := fileInfo.Size()
	return &FileWrapper{
		file:   file,
		Buffer: make([]byte, int64(math.Round(float64(fileSize)*1.1))),
		Writer: file,
		Reader: file,
	}
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
// Seek sets the offset for the next Read or Write on file to offset, interpreted according to whence: 0 means relative to the origin of the file, 1 means relative to the current offset, and 2 means relative to the end. It returns the new offset and an error, if any.
func (f *FileWrapper) Seek(offset int64, whence int) (int64, error) {
	return f.file.Seek(offset, whence)
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
	// Read from the source file
	bytesRead, err := source.Read(f.Buffer)
	if err != nil {
		return 0, err
	}

	// Write to the destination file
	bytesWritten, err := destination.Write(f.Buffer[:bytesRead])
	if err != nil {
		return 0, err
	}

	// Return the number of bytes written
	return int64(bytesWritten), nil
}

//func (f *FileWrapper) Size() (int64, error) {
//	fileInfo, err := f.file.Stat()
//	if err != nil {
//		return 0, err
//	}
//	return fileInfo.Size(), nil
//}

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

func (f *FileMock) Seek(offset int64, whence int) (int64, error) {
	return offset, f.Err
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
