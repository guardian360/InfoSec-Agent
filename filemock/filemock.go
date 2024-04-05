package filemock

import (
	"os"
)

// File is an interface that represents a file. It provides a method for closing the file.
// The Close method should return an error if the file cannot be closed, and nil if the file is closed successfully.
type File interface {
	Close() error
}

// FileWrapper is a struct that wraps the os.File type. It provides a method for closing the file and encapsulates the os.File within it.
// This allows for additional functionality or modifications to be added without changing the behavior of the os.File type itself.
// The Close method should return an error if the file cannot be closed, and nil if the file is closed successfully.
type FileWrapper struct {
	file *os.File
}

// NewFileWrapper is a constructor function that creates and returns a new instance of the FileWrapper struct.
// It takes an os.File pointer as an argument and wraps it in a FileWrapper to provide additional functionality.
//
// Parameter:
//   - file: A pointer to an os.File instance that needs to be wrapped.
//
// Returns:
//   - A pointer to a newly created FileWrapper instance that encapsulates the provided os.File.
func NewFileWrapper(file *os.File) *FileWrapper {
	return &FileWrapper{file: file}
}

// Close is a method of the FileWrapper struct that attempts to close the encapsulated os.File instance.
// It delegates the operation to the Close method of the os.File instance.
//
// Parameters: None.
//
// Returns:
//   - An error if the os.File instance cannot be closed, typically due to an I/O problem.
//   - Nil if the os.File instance is closed successfully.
//
// Note: After the Close method is called, the FileWrapper and its encapsulated os.File instance should not be used.
func (f *FileWrapper) Close() error {
	return f.file.Close()
}

// FileMock is a struct that simulates a file for testing purposes. It contains fields that represent the state of the file.
//
// Fields:
//   - FileName: A string representing the name of the file.
//   - Open: A boolean indicating whether the file is open or closed.
//   - Err: An error that will be returned when attempting to close the file if it is set.
//
// This struct is typically used in tests to simulate various file-related scenarios without having to interact with the actual file system.
type FileMock struct {
	FileName string
	Open     bool
	Err      error
}

// Close attempts to close the FileMock instance, simulating the behavior of a real file.
//
// If the FileMock instance is already closed or nil, it returns an appropriate error.
// If the FileMock instance has an error set (Err field), it returns that error.
// If the FileMock instance is open and no error is set, it closes the instance and returns nil.
//
// Parameters: None.
//
// Returns:
//   - os.ErrInvalid if the FileMock instance is nil.
//   - os.ErrClosed if the FileMock instance is already closed.
//   - The error set in the Err field of the FileMock instance, if any.
//   - Nil if the FileMock instance is open and no error is set.
//
// Note: After the Close method is called, the FileMock instance should not be used.
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
