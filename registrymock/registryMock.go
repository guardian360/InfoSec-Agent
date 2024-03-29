package registrymock

import (
	"errors"
	"golang.org/x/sys/windows/registry"
)

var (
	CLASSES_ROOT  = RegistryKey(NewRegistryKeyWrapper(registry.CLASSES_ROOT))
	CURRENT_USER  = RegistryKey(NewRegistryKeyWrapper(registry.CURRENT_USER))
	LOCAL_MACHINE = RegistryKey(NewRegistryKeyWrapper(registry.LOCAL_MACHINE))
)

// RegistryKey is an interface for reading values from the Windows registry
type RegistryKey interface {
	GetStringValue(name string) (val string, valtype uint32, err error)
	GetBinaryValue(name string) (val []byte, valtype uint32, err error)
	GetIntegerValue(name string) (val uint64, valtype uint32, err error)
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

func (r *RegistryKeyWrapper) GetStringValue(name string) (val string, valtype uint32, err error) {
	return r.key.GetStringValue(name)
}

func (r *RegistryKeyWrapper) GetBinaryValue(name string) (val []byte, valtype uint32, err error) {
	return r.key.GetBinaryValue(name)
}

func (r *RegistryKeyWrapper) GetIntegerValue(name string) (val uint64, valtype uint32, err error) {
	return r.key.GetIntegerValue(name)
}

func (r *RegistryKeyWrapper) OpenKey(path string, access uint32) (RegistryKey, error) {
	newKey, err := registry.OpenKey(r.key, path, access)
	return &RegistryKeyWrapper{key: newKey}, err
}

func (r *RegistryKeyWrapper) ReadValueNames(count int) ([]string, error) {
	return r.key.ReadValueNames(count)
}
func (r *RegistryKeyWrapper) Close() error {
	return r.key.Close()
}

func (r *RegistryKeyWrapper) Stat() (*registry.KeyInfo, error) {
	return r.key.Stat()
}

func (r *RegistryKeyWrapper) ReadSubKeyNames(count int) ([]string, error) {
	return r.key.ReadSubKeyNames(count)
}

// MockRegistryKey is a mock implementation of the RegistryKey interface
type MockRegistryKey struct {
	KeyName       string
	StringValues  map[string]string
	BinaryValues  map[string][]byte
	IntegerValues map[string]uint64
	SubKeys       []MockRegistryKey
	StatReturn    *registry.KeyInfo
	Err           error
}

// GetStringValue returns the string value of the key
func (m *MockRegistryKey) GetStringValue(name string) (string, uint32, error) {
	// TODO: add correct errors
	return m.StringValues[name], 0, nil
}

// GetBinaryValue returns the binary value of the key
func (m *MockRegistryKey) GetBinaryValue(name string) ([]byte, uint32, error) {
	return m.BinaryValues[name], 0, nil
}

// GetIntegerValue returns the integer value of the key
func (m *MockRegistryKey) GetIntegerValue(name string) (uint64, uint32, error) {
	return m.IntegerValues[name], 0, nil
}

// OpenKey opens a registry key with a path relative to the current key
func (m *MockRegistryKey) OpenKey(path string, access uint32) (RegistryKey, error) {
	for _, key := range m.SubKeys {
		if key.KeyName == path {
			return &key, nil
		}
	}
	return m, errors.New("key not found")
}

// ReadValueNames reads the value names of the key
func (m *MockRegistryKey) ReadValueNames(maxCount int) ([]string, error) {
	uniqueValueNames := make(map[string]string)
	valueNames := make([]string, 0, len(m.StringValues)+len(m.BinaryValues)+len(m.IntegerValues))
	count := 0
	if maxCount <= count {
		return []string{}, nil
	}
	for key := range m.StringValues {
		if uniqueValueNames[key] == "" {
			uniqueValueNames[key] = key
		}
	}
	for key := range m.BinaryValues {
		if uniqueValueNames[key] == "" {
			uniqueValueNames[key] = key
		}
	}
	for key := range m.IntegerValues {
		if uniqueValueNames[key] == "" {
			uniqueValueNames[key] = key
		}
	}
	for key := range uniqueValueNames {
		if maxCount == count {
			break
		}
		valueNames = append(valueNames, key)
		count++
	}
	return valueNames, nil
}

// Close closes the registry key
func (m *MockRegistryKey) Close() error {
	return nil
}

// Stat returns the key info of the registry key
func (m *MockRegistryKey) Stat() (*registry.KeyInfo, error) {
	return m.StatReturn, nil
}

// ReadSubKeyNames reads the subkey names of the registry key
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
