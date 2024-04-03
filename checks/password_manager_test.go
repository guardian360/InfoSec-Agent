package checks_test

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/stretchr/testify/mock"
)

// MockProgramLister is a mock type for the ProgramLister interface
type MockProgramLister struct {
	mock.Mock
}

// ListInstalledPrograms mocks the ProgramLister interface ListInstalledPrograms method
func (m *MockProgramLister) ListInstalledPrograms(directory string) ([]string, error) {
	args := m.Called(directory)
	return args.Get(0).([]string), args.Error(1)
}

// TestPasswordManager tests the PasswordManager function
//
// Parameters: t *testing.T - The testing framework
//
// Returns: _
func TestPasswordManager(t *testing.T) {
	tests := []struct {
		name         string
		mockPrograms []string
		want         checks.Check
	}{
		{
			name:         "With Known Password Manager",
			mockPrograms: []string{"1Password"},
			want:         checks.NewCheckResult("PasswordManager", "1Password"),
		},
		{
			name:         "No Password Manager",
			mockPrograms: []string{"RandomSoftware"},
			want:         checks.NewCheckResult("PasswordManager", "No password manager found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLister := new(MockProgramLister)
			mockLister.On("ListInstalledPrograms", mock.Anything).Return(tt.mockPrograms, nil)

			result := checks.PasswordManager(mockLister)
			require.Equal(t, tt.want, result)
		})
	}
}

// TestListInstalledPrograms tests the ListInstalledPrograms function
//
// Parameters: t *testing.T - The testing framework
//
// Returns: _
func TestListInstalledPrograms(t *testing.T) {
	tests := []struct {
		name      string
		directory string
		want      []string
	}{
		{
			name:      "With Programs",
			directory: "C:\\Program Files",
			want:      []string{"Program1", "Program2"},
		},
		{
			name:      "No Programs",
			directory: "C:\\Program Files",
			want:      []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLister := new(MockProgramLister)
			mockLister.On("ListInstalledPrograms", mock.Anything).Return(tt.want, nil)

			result, err := mockLister.ListInstalledPrograms(tt.directory)
			require.Equal(t, tt.want, result)
			if err != nil {
				t.Errorf("Test %s failed. Expected no error, got %v", tt.name, err)
			}
		})
	}
}
