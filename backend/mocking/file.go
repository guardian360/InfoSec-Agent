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

// TODO: Update documentation
// FileWrapper is a wrapper for the os.File type
type FileWrapper struct {
	file   *os.File
	Buffer []byte
	Err    error
	Writer *os.File
	Reader *os.File
}

// TODO: Update documentation
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

// TODO: Update documentation
// Close is a method of the FileWrapper struct that closes the underlying os.File.
//
// It calls the Close method of the os.File that the FileWrapper is wrapping.
//
// Returns:
//   - An error if the underlying os.File cannot be closed. If the file is closed successfully, it returns nil.
func (f *FileWrapper) Close() error {
	return f.file.Close()
}

// TODO: Update documentation
// Read is a method of the FileWrapper struct that reads from the underlying os.File.
//
// It calls the Read method of the os.File that the FileWrapper is wrapping.
//
// Parameters:
//
//	p - A byte slice that serves as the buffer into which the data is read.
//
// Returns:
//   - The number of bytes read. If the file is at the end, it returns 0.
//   - An error if the underlying os.File cannot be read. If the file is read successfully, it returns nil.
func (f *FileWrapper) Read(p []byte) (int, error) {
	return f.file.Read(p)
}

// TODO: Update documentation
// Seek is a method of the FileWrapper struct that sets the offset for the next Read or Write operation on the underlying os.File.
//
// It calls the Seek method of the os.File that the FileWrapper is wrapping.
//
// Parameters:
//
//	offset - The offset for the next Read or Write operation. The interpretation of the offset is determined by the whence parameter.
//	whence - The point relative to which the offset is interpreted. It can be 0 (relative to the start of the file), 1 (relative to the current offset), or 2 (relative to the end of the file).
//
// Returns:
//   - The new offset relative to the start of the file.
//   - An error if the underlying os.File cannot seek to the specified offset. If the seek operation is successful, it returns nil.func (f *FileWrapper) Seek(offset int64, whence int) (int64, error) {
func (f *FileWrapper) Seek(offset int64, whence int) (int64, error) {
	return f.file.Seek(offset, whence)
}

// TODO: Update documentation
// Write is a method of the FileWrapper struct that writes to the underlying os.File.
//
// It calls the Write method of the os.File that the FileWrapper is wrapping.
//
// Parameters:
//
//	p - A byte slice that serves as the buffer from which the data is written.
//
// Returns:
//   - The number of bytes written. If the file is at the end, it returns 0.
//   - An error if the underlying os.File cannot be written to. If the file is written successfully, it returns nil.
func (f *FileWrapper) Write(p []byte) (int, error) {
	return f.file.Write(p)
}

// TODO: Update documentation
// Copy is a method of the FileWrapper struct that copies data from a source File to a destination File.
//
// It reads from the source File into the FileWrapper's Buffer, then writes from the Buffer to the destination File.
//
// Parameters:
//
//	source - The source File to copy data from.
//	destination - The destination File to copy data to.
//
// Returns:
//   - The number of bytes written to the destination File.
//   - An error if the source File cannot be read or the destination File cannot be written to. If the copy operation is successful, it returns nil.
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

// TODO: Update documentation
// Stat is a method of the FileWrapper struct that retrieves the file descriptor's metadata.
//
// It calls the Stat method of the os.File that the FileWrapper is wrapping.
//
// Returns:
//   - An os.FileInfo object that describes the file. If the method is successful, it returns this object and a nil error.
//   - An error if the underlying os.File's metadata cannot be retrieved. If the method is unsuccessful, it returns a nil os.FileInfo and the error.
func (f *FileWrapper) Stat() (os.FileInfo, error) {
	return f.file.Stat()
}

// TODO: Update documentation
// FileInfoMock is a struct that mocks the os.FileInfo interface for testing purposes.
//
// It contains a single field, file, which is a pointer to a FileMock. This allows the FileInfoMock to return
// information about the mocked file when its methods are called.
type FileInfoMock struct {
	file *FileMock
}

// TODO: Update documentation
// Size returns the length of the buffer in the FileMock that FileInfoMock is associated with.
// Name returns an empty string as it's not relevant in this mock implementation.
// Mode returns 0 as it's not relevant in this mock implementation.
// ModTime returns the zero value for time.Time as it's not relevant in this mock implementation.
// IsDir returns false as it's not relevant in this mock implementation.
// Sys returns nil as it's not relevant in this mock implementation.
func (f *FileInfoMock) Size() int64 {
	return int64(len(f.file.Buffer))
}
func (f *FileInfoMock) Name() string       { return "" }
func (f *FileInfoMock) Mode() os.FileMode  { return 0 }
func (f *FileInfoMock) ModTime() time.Time { return time.Time{} }
func (f *FileInfoMock) IsDir() bool        { return false }
func (f *FileInfoMock) Sys() interface{}   { return nil }

// TODO: Update documentation
// FileMock is a struct that mocks the File interface for testing purposes.
//
// It contains the following fields:
//   - FileName: The name of the file. This is a string.
//   - IsOpen: A boolean indicating whether the file is open or not.
//   - Buffer: A byte slice that serves as the buffer for the file data.
//   - Bytes: The number of bytes in the buffer.
//   - Err: An error that can be set to simulate an error condition.
//   - FileInfo: A pointer to a FileInfoMock struct. This allows the FileMock to return
//     information about the mocked file when its methods are called.
type FileMock struct {
	FileName string
	IsOpen   bool
	Buffer   []byte
	Bytes    int
	Err      error
	FileInfo *FileInfoMock
}

// TODO: Update documentation
// ReadDir is a method of the FileMock struct that simulates the behavior of reading a directory.
// It doesn't actually read a directory, but instead returns a predefined result that can be set for testing purposes.
//
// Parameters:
//   - _: This method ignores its input parameter. The underscore character is a convention in Go for discarding a variable.
//
// Returns:
//   - A slice of os.DirEntry: In this mock implementation, it always returns nil.
//   - An error: The error that was previously set in the FileMock struct. If no error was set, it returns nil.
func (f *FileMock) ReadDir(_ string) ([]os.DirEntry, error) {
	return nil, f.Err
}

// TODO: Update documentation
// Close is a method of the FileMock struct that simulates the behavior of closing a file.
//
// It checks if the FileMock is nil or if the file is already closed, and returns an appropriate error in each case.
// If the file is open, it sets the IsOpen field to false, simulating the closing of the file.
//
// Returns:
//   - os.ErrInvalid if the FileMock is nil.
//   - os.ErrClosed if the file is already closed.
//   - The error that was previously set in the FileMock struct. If no error was set and the file is closed successfully, it returns nil.
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

// TODO: Update documentation
// Read is a method of the FileMock struct that simulates the behavior of reading from a file.
//
// It checks if there is an error set in the FileMock. If there is, it returns 0 and the error.
// If there is no error, it copies data from the FileMock's Buffer into the provided byte slice and updates the Buffer.
// If the Buffer is empty, it returns 0 and io.EOF to simulate the end of the file.
//
// Parameters:
//
//	p - A byte slice that serves as the buffer into which the data is read.
//
// Returns:
//   - The number of bytes read. If the Buffer is empty, it returns 0.
//   - An error if one was set in the FileMock. If the Buffer is empty, it returns io.EOF. If the read operation is successful, it returns nil.
func (f *FileMock) Read(p []byte) (int, error) {
	if f.Err != nil {
		return 0, f.Err
	}

	if len(f.Buffer) == 0 {
		return 0, io.EOF
	}

	n := copy(p, f.Buffer)
	f.Buffer = f.Buffer[n:]
	return n, nil
}

// TODO: Update documentation
// Write is a method of the FileMock struct that simulates the behavior of writing to a file.
//
// It ignores its input parameter and instead returns the number of bytes that were previously set in the FileMock struct.
// This allows you to control the behavior of the Write method for testing purposes.
//
// Parameters:
//
//	_ - A byte slice that serves as the buffer from which the data is written. In this mock implementation, this parameter is ignored.
//
// Returns:
//   - The number of bytes that were previously set in the FileMock struct. If no number was set, it returns 0.
//   - An error if one was set in the FileMock. If no error was set, it returns nil.
func (f *FileMock) Write(_ []byte) (int, error) {
	return f.Bytes, f.Err
}

// TODO: Update documentation
// Seek is a method of the FileMock struct that simulates the behavior of setting the offset for the next Read or Write operation on a file.
//
// It checks if there is an error set in the FileMock. If there is, it returns 0 and the error.
// If there is no error, it adjusts the Buffer according to the offset and whence parameters.
// If the offset is out of range, it returns 0 and io.EOF to simulate the end of the file.
//
// Parameters:
//
//	offset - The offset for the next Read or Write operation. The interpretation of the offset is determined by the whence parameter.
//	whence - The point relative to which the offset is interpreted. It can be 0 (relative to the start of the file), 1 (relative to the current offset), or 2 (relative to the end of the file).
//
// Returns:
//   - The new offset relative to the start of the file.
//   - An error if one was set in the FileMock. If the offset is out of range, it returns io.EOF. If the seek operation is successful, it returns nil.
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

// TODO: Update documentation
// Copy is a method of the FileMock struct that simulates the behavior of copying data from a source File to a destination File.
//
// It reads from the source File and writes to the destination File. The actual data transfer is not simulated in this mock implementation.
// Instead, it returns the number of bytes that were previously set in the FileMock struct and the error that was previously set, if any.
// This allows you to control the behavior of the Copy method for testing purposes.
//
// Parameters:
//
//	source - The source File to copy data from. In this mock implementation, this parameter is used but no data is actually read from it.
//	destination - The destination File to copy data to. In this mock implementation, this parameter is used but no data is actually written to it.
//
// Returns:
//   - The number of bytes that were previously set in the FileMock struct. If no number was set, it returns 0.
//   - An error if one was set in the FileMock. If no error was set, it returns nil.
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

// TODO: Update documentation
// Stat is a method of the FileMock struct that simulates the behavior of retrieving the file descriptor's metadata.
//
// It returns the FileInfo and error that were previously set in the FileMock struct. This allows you to control the behavior of the Stat method for testing purposes.
//
// Returns:
//   - An os.FileInfo object that describes the file. If no error was set and the method is successful, it returns this object and a nil error.
//   - An error if one was set in the FileMock. If no error was set, it returns nil.
func (f *FileMock) Stat() (os.FileInfo, error) {
	return f.FileInfo, f.Err
}
