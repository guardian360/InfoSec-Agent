package mocking

import (
	"errors"
	"io"
	"math"
	"os"
	"time"
)

type File interface {
	Close() error
	io.Closer
	Read(p []byte) (int, error)
	Write(p []byte) (int, error)
	Seek(offset int64, _ int) (int64, error)
	Stat() (os.FileInfo, error)
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

func (f *FileWrapper) Stat() (os.FileInfo, error) {
	return f.file.Stat()
}

type FileInfoMock struct {
	file *FileMock
}

func (f *FileInfoMock) Size() int64 {
	return int64(len(f.file.Buffer))
}

func (f *FileInfoMock) Name() string       { return "" }
func (f *FileInfoMock) Mode() os.FileMode  { return 0 }
func (f *FileInfoMock) ModTime() time.Time { return time.Time{} }
func (f *FileInfoMock) IsDir() bool        { return false }
func (f *FileInfoMock) Sys() interface{}   { return nil }

// func (f *FileWrapper) Size() (int64, error) {
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
	FileInfo *FileInfoMock
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
func (f *FileMock) Read(p []byte) (int, error) {
	if f.Err != nil {
		return 0, f.Err
	}

	n := copy(p, f.Buffer)
	f.Buffer = f.Buffer[n:]
	return n, nil
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

// Seek sets the offset for the next Read or Write on file to offset, interpreted according to whence:
// 0 means relative to the origin of the file,
// 1 means relative to the current offset,
// 2 means relative to the end.
// It returns the new offset and an error, if any.
func (f *FileMock) Seek(offset int64, whence int) (int64, error) {
	if f.Err != nil {
		return 0, f.Err
	}

	switch whence {
	case 0: // relative to the origin of the file
		if offset < 0 || offset > int64(len(f.Buffer)) {
			return 0, io.EOF
		}
		f.Buffer = f.Buffer[offset:]
	case 1: // relative to the current offset
		if offset < 0 || offset > int64(len(f.Buffer)) {
			return 0, io.EOF
		}
		f.Buffer = f.Buffer[offset:]
	case 2: // relative to the end
		if offset > 0 || -offset > int64(len(f.Buffer)) {
			return 0, io.EOF
		}
		f.Buffer = f.Buffer[:len(f.Buffer)+int(offset)]
	default:
		return 0, errors.New("invalid whence")
	}

	return int64(len(f.Buffer)), nil
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

func (f *FileMock) Stat() (os.FileInfo, error) {
	return f.FileInfo, f.Err
}
