package mocking

import (
	"errors"

	"golang.org/x/sys/windows/registry"
)

var (
	CurrentUser  = RegistryKey(NewRegistryKeyWrapper(registry.CURRENT_USER))
	LocalMachine = RegistryKey(NewRegistryKeyWrapper(registry.LOCAL_MACHINE))
)

// RegistryKey is an interface for reading values from the Windows registry
type RegistryKey interface {
	GetStringValue(name string) (string, uint32, error)
	GetBinaryValue(name string) ([]byte, uint32, error)
	GetIntegerValue(name string) (uint64, uint32, error)
	OpenKey(path string, access uint32) (RegistryKey, error)
	ReadValueNames(count int) ([]string, error)
	ReadSubKeyNames(count int) ([]string, error)
	Close() error
	Stat() (*registry.KeyInfo, error)
}

type RegistryKeyWrapper struct {
	key registry.Key
}

func NewRegistryKeyWrapper(key registry.Key) *RegistryKeyWrapper {
	return &RegistryKeyWrapper{key: key}
}

// GetStringValue retrieves the string value of a specified name from the Windows registry key.
//
// Parameters:
//   - name string: The name of the value to retrieve.
//
// Returns:
//   - string: The string value associated with the specified name.
//   - uint32: The data type of the value. This is always REG_SZ for string values.
//   - error: An error object if an error occurred while retrieving the value. Otherwise, nil.
//
// This function is a method of the RegistryKeyWrapper struct and uses the GetStringvalue method of the registry.Key object it wraps.
func (r *RegistryKeyWrapper) GetStringValue(name string) (string, uint32, error) {
	return r.key.GetStringValue(name)
}

// GetBinaryValue retrieves the binary value of a specified name from the Windows registry key.
//
// Parameters:
//   - name string: The name of the value to retrieve.
//
// Returns:
//   - []byte: The binary value associated with the specified name.
//   - uint32: The data type of the value. This is always REG_BINARY for binary values.
//   - error: An error object if an error occurred while retrieving the value. Otherwise, nil.
//
// This function is a method of the RegistryKeyWrapper struct and uses the GetBinaryValue method of the registry.Key object it wraps.
func (r *RegistryKeyWrapper) GetBinaryValue(name string) ([]byte, uint32, error) {
	return r.key.GetBinaryValue(name)
}

// GetIntegerValue retrieves the integer value of a specified name from the Windows registry key.
//
// Parameters:
//   - name string: The name of the value to retrieve.
//
// Returns:
//   - uint64: The integer value associated with the specified name.
//   - uint32: The data type of the value. This is always REG_DWORD for integer values.
//   - error: An error object if an error occurred while retrieving the value. Otherwise, nil.
//
// This function is a method of the RegistryKeyWrapper struct and uses the GetIntegerValue method of the registry.Key object it wraps.
func (r *RegistryKeyWrapper) GetIntegerValue(name string) (uint64, uint32, error) {
	return r.key.GetIntegerValue(name)
}

// OpenKey opens a registry key with a path relative to the current key and the specified access rights.
//
// Parameters:
//   - path string: The path of the registry key to open, relative to the current key.
//   - access uint32: The access rights to use when opening the key.
//
// Returns:
//   - RegistryKey: A RegistryKey object representing the opened key.
//   - error: An error object if an error occurred while opening the key. Otherwise, nil.
//
// This function is a method of the RegistryKeyWrapper struct and uses the OpenKey method of the registry.Key object it wraps.
func (r *RegistryKeyWrapper) OpenKey(path string, access uint32) (RegistryKey, error) {
	newKey, err := registry.OpenKey(r.key, path, access)
	return &RegistryKeyWrapper{key: newKey}, err
}

// ReadValueNames retrieves the names of the values in the Windows registry key.
//
// Parameters:
//   - count int: The maximum number of value names to retrieve. If count is less than or equal to zero, all value names are retrieved.
//
// Returns:
//   - []string: A slice of strings containing the names of the values in the registry key. The length of the slice is the lesser of the number of values in the key and the count parameter.
//   - error: An error object if an error occurred while retrieving the value names. Otherwise, nil.
//
// This function is a method of the RegistryKeyWrapper struct and uses the ReadValueNames method of the registry.Key object it wraps.
func (r *RegistryKeyWrapper) ReadValueNames(count int) ([]string, error) {
	return r.key.ReadValueNames(count)
}

// Close closes the Windows registry key.
//
// Returns:
//   - error: An error object if an error occurred while closing the key. Otherwise, nil.
//
// This function is a method of the RegistryKeyWrapper struct and uses the Close method of the registry.Key object it wraps.
func (r *RegistryKeyWrapper) Close() error {
	return r.key.Close()
}

// Stat retrieves the information about the Windows registry key.
//
// Returns:
//   - *registry.KeyInfo: A pointer to a KeyInfo object that contains information about the registry key.
//   - error: An error object if an error occurred while retrieving the information. Otherwise, nil.
//
// This function is a method of the RegistryKeyWrapper struct and uses the Stat method of the registry.Key object it wraps.
func (r *RegistryKeyWrapper) Stat() (*registry.KeyInfo, error) {
	return r.key.Stat()
}

// ReadSubKeyNames retrieves the names of the subkeys in the Windows registry key.
//
// Parameters:
//   - count int: The maximum number of subkey names to retrieve. If count is less than or equal to zero, all subkey names are retrieved.
//
// Returns:
//   - []string: A slice of strings containing the names of the subkeys in the registry key. The length of the slice is the lesser of the number of subkeys in the key and the count parameter.
//   - error: An error object if an error occurred while retrieving the subkey names. Otherwise, nil.
//
// This function is a method of the RegistryKeyWrapper struct and uses the ReadSubKeyNames method of the registry.Key object it wraps.
func (r *RegistryKeyWrapper) ReadSubKeyNames(count int) ([]string, error) {
	return r.key.ReadSubKeyNames(count)
}

// MockRegistryKey is a struct that provides a mock implementation of the RegistryKey interface.
// It is designed for testing purposes and allows for the simulation of registry key operations
// without interacting with the actual Windows registry.
//
// The struct contains fields that represent the possible values a registry key can hold,
// including string, binary, and integer values. It also includes a list of subkeys,
// allowing for the simulation of a registry key hierarchy.
//
// The methods of this struct mimic the behavior of their counterparts in the RegistryKey interface,
// returning the values stored in the struct's fields instead of interacting with the Windows registry.
type MockRegistryKey struct {
	KeyName       string
	StringValues  map[string]string
	BinaryValues  map[string][]byte
	IntegerValues map[string]uint64
	SubKeys       []MockRegistryKey
	StatReturn    *registry.KeyInfo
	Err           error
}

// GetStringValue retrieves the string value associated with a specified name from the MockRegistryKey.
//
// Parameters:
//   - name string: The name of the value to retrieve.
//
// Returns:
//   - string: The string value associated with the specified name. If the name does not exist in the StringValues map, an empty string is returned.
//   - uint32: The data type of the value. This is always 0 for string values in the MockRegistryKey.
//   - error: An error object if the specified name does not exist in the StringValues map. Otherwise, nil.
//
// This function is a method of the MockRegistryKey struct and simulates the GetStringvalue method of the RegistryKey interface by returning the values stored in the StringValues map of the MockRegistryKey.
func (m *MockRegistryKey) GetStringValue(name string) (string, uint32, error) {
	if m.StringValues[name] == "" {
		return "", 0, errors.New("key not found")
	}
	return m.StringValues[name], 0, nil
}

// GetBinaryValue retrieves the binary value associated with a specified name from the MockRegistryKey.
//
// Parameters:
//   - name string: The name of the value to retrieve.
//
// Returns:
//   - []byte: The binary value associated with the specified name. If the name does not exist in the BinaryValues map, nil is returned.
//   - uint32: The data type of the value. This is always 0 for binary values in the MockRegistryKey.
//   - error: An error object if the specified name does not exist in the BinaryValues map. Otherwise, nil.
//
// This function is a method of the MockRegistryKey struct and simulates the GetBinaryValue method of the RegistryKey interface by returning the values stored in the BinaryValues map of the MockRegistryKey.
func (m *MockRegistryKey) GetBinaryValue(name string) ([]byte, uint32, error) {
	if m.BinaryValues[name] == nil {
		return nil, 0, errors.New("key not found")
	}
	return m.BinaryValues[name], 0, nil
}

// GetIntegerValue retrieves the integer value associated with a specified name from the MockRegistryKey.
//
// Parameters:
//   - name string: The name of the value to retrieve.
//
// Returns:
//   - uint64: The integer value associated with the specified name. If the name does not exist in the IntegerValues map, zero is returned.
//   - uint32: The data type of the value. This is always 0 for integer values in the MockRegistryKey.
//   - error: An error object if the specified name does not exist in the IntegerValues map. Otherwise, nil.
//
// This function is a method of the MockRegistryKey struct and simulates the GetIntegerValue method of the RegistryKey interface by returning the values stored in the IntegerValues map of the MockRegistryKey.
func (m *MockRegistryKey) GetIntegerValue(name string) (uint64, uint32, error) {
	return m.IntegerValues[name], 0, nil
}

// OpenKey opens a subkey with a specified path relative to the current key in the MockRegistryKey.
//
// Parameters:
//   - path string: The path of the subkey to open, relative to the current key.
//   - _ uint32: This parameter is ignored in the MockRegistryKey implementation as no access rights are required to open a mock key.
//
// Returns:
//   - RegistryKey: A RegistryKey object representing the opened subkey. If the specified path does not exist in the SubKeys slice, the method returns the current key.
//   - error: An error object if the specified path does not exist in the SubKeys slice. Otherwise, nil.
//
// This function is a method of the MockRegistryKey struct and simulates the OpenKey method of the RegistryKey interface by returning a subkey from the SubKeys slice of the MockRegistryKey.
func (m *MockRegistryKey) OpenKey(path string, _ uint32) (RegistryKey, error) {
	for _, key := range m.SubKeys {
		if key.KeyName == path {
			return &key, nil
		}
	}
	return m, errors.New("key not found")
}

// ReadValueNames retrieves the names of the values stored in the MockRegistryKey.
//
// Parameters:
//   - maxCount int: The maximum number of value names to retrieve. If maxCount is less than or equal to zero, all value names are retrieved.
//
// Returns:
//   - []string: A slice of strings containing the names of the values in the MockRegistryKey. The length of the slice is the lesser of the number of values in the key and the maxCount parameter.
//   - error: An error object if an error occurred while retrieving the value names. Otherwise, nil.
//
// This function is a method of the MockRegistryKey struct and simulates the ReadValueNames method of the RegistryKey interface by returning the names of the values stored in the StringValues, BinaryValues, and IntegerValues maps of the MockRegistryKey. It removes any duplicate names before returning the result.
func (m *MockRegistryKey) ReadValueNames(maxCount int) ([]string, error) {
	var valueNames []string
	for key := range m.StringValues {
		valueNames = append(valueNames, key)
	}
	for key := range m.BinaryValues {
		valueNames = append(valueNames, key)
	}
	for key := range m.IntegerValues {
		valueNames = append(valueNames, key)
	}
	// remove duplicate keys from valueNames
	keys := make(map[string]bool)
	var uniqueValueNames []string
	for _, entry := range valueNames {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			uniqueValueNames = append(uniqueValueNames, entry)
		}
	}
	if maxCount <= 0 || maxCount >= len(uniqueValueNames) {
		return uniqueValueNames, nil
	}
	return uniqueValueNames[:maxCount], nil
}

// Close terminates the connection to the MockRegistryKey.
//
// This function is a method of the MockRegistryKey struct and simulates the Close method of the RegistryKey interface.
// It is used to clean up any resources associated with the MockRegistryKey.
// In the context of the MockRegistryKey, this function does not perform any operations as there are no resources to release.
//
// Returns:
//   - error: Always returns nil as there are no resources to release in the MockRegistryKey.
func (m *MockRegistryKey) Close() error {
	return nil
}

// Stat retrieves the information about the MockRegistryKey.
//
// Returns:
//   - *registry.KeyInfo: A pointer to a KeyInfo object that contains information about the MockRegistryKey. This includes the number of subkeys and values, and the time of the last write operation.
//   - error: An error object if an error occurred while retrieving the information. Otherwise, nil.
//
// This function is a method of the MockRegistryKey struct and simulates the Stat method of the RegistryKey interface by returning the KeyInfo stored in the StatReturn field of the MockRegistryKey.
func (m *MockRegistryKey) Stat() (*registry.KeyInfo, error) {
	return m.StatReturn, nil
}

// ReadSubKeyNames retrieves the names of the subkeys in the MockRegistryKey.
//
// Parameters:
//   - count int: The maximum number of subkey names to retrieve. If count is less than or equal to zero, all subkey names are retrieved.
//
// Returns:
//   - []string: A slice of strings containing the names of the subkeys in the MockRegistryKey. The length of the slice is the lesser of the number of subkeys in the key and the count parameter.
//   - error: An error object if an error occurred while retrieving the subkey names. Otherwise, nil.
//
// This function is a method of the MockRegistryKey struct and simulates the ReadSubKeyNames method of the RegistryKey interface by returning the names of the subkeys stored in the SubKeys slice of the MockRegistryKey.
func (m *MockRegistryKey) ReadSubKeyNames(count int) ([]string, error) {
	var subKeyNames []string
	maxCount := 0
	for _, key := range m.SubKeys {
		if maxCount == count {
			break
		}
		subKeyNames = append(subKeyNames, key.KeyName)
		maxCount++
	}
	return subKeyNames, nil
}