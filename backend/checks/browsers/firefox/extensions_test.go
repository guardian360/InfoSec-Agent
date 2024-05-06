package firefox_test

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/browsers"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/browsers/firefox"
	"github.com/stretchr/testify/require"
)

// MockProfileFinder is a mock implementation of the FirefoxProfileFinder interface
type MockProfileFinder struct {
	FirefoxFolderFunc func() ([]string, error)
}

func (m *MockProfileFinder) FirefoxFolder() ([]string, error) {
	return m.FirefoxFolderFunc()
}

func TestExtensionFirefox(t *testing.T) {
	t.Run("returns error when no firefox directory found", func(t *testing.T) {
		profileFinder := &MockProfileFinder{
			FirefoxFolderFunc: func() ([]string, error) {
				return nil, errors.New("directory not found")
			},
		}

		check1, check2 := firefox.ExtensionFirefox(profileFinder)
		require.Nil(t, check1.Result)
		require.Error(t, check1.Error)
		require.Nil(t, check2.Result)
		require.Error(t, check2.Error)
	})

	t.Run("returns error when unable to open extensions.json", func(t *testing.T) {
		Profilefinder = browsers.MockProfileFinder{
			MockFirefoxFolder: func() ([]string, error) {
				return []string{"\\valid\\directory"}, nil
			},
		}
		// setup a temporary directory for the test
		tmpDir, err := os.MkdirTemp("", "firefox")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tmpDir) // clean up

		profileFinder := &MockProfileFinder{
			FirefoxFolderFunc: func() ([]string, error) {
				return []string{tmpDir}, nil
			},
		}

		check1, check2 := firefox.ExtensionFirefox(profileFinder)
		require.Nil(t, check1.Result)
		require.Error(t, check1.Error)
		require.Nil(t, check2.Result)
		require.Error(t, check2.Error)
	})

	t.Run("returns correct results when extensions.json is valid", func(t *testing.T) {
		// setup a temporary directory for the test
		tmpDir, err := os.MkdirTemp("", "firefox")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tmpDir) // clean up

		// create a valid extensions.json file
		extensionsFile := filepath.Join(tmpDir, "extensions.json")
		err = os.WriteFile(extensionsFile, []byte(`{
            "addons": [
                {
                    "defaultLocale": {
                        "name": "Test Addon",
                        "creator": "Test Creator"
                    },
                    "type": "extension",
                    "active": true
                }
            ]
        }`), 0644)
		if err != nil {
			t.Fatal(err)
		}

		profileFinder := &MockProfileFinder{
			FirefoxFolderFunc: func() ([]string, error) {
				return []string{tmpDir}, nil
			},
		}

		check1, check2 := firefox.ExtensionFirefox(profileFinder)
		require.Nil(t, check1.Error)
		require.Nil(t, check2.Error)
		expected1 := checks.NewCheckResult(checks.ExtensionFirefoxID, 0, "Test Addon,extension,Test Creator,true")
		expected2 := checks.NewCheckResult(checks.AdblockFirefoxID, 1, strconv.FormatBool(false))
		require.Equal(t, expected1, check1)
		require.Equal(t, expected2, check2)
	})
}

func TestExtensionFirefox_DecodeError(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "firefox")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir) // clean up

	// Create an invalid extensions.json file
	extensionsFile := tmpDir + "\\extensions.json"
	err = os.WriteFile(extensionsFile, []byte(`invalid json`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Mock the FirefoxProfileFinder to return the temporary directory
	profileFinder := &browsers.MockProfileFinder{
		MockFirefoxFolder: func() ([]string, error) {
			return []string{tmpDir}, nil
		},
	}

	// Call the function
	check1, check2 := firefox.ExtensionFirefox(profileFinder)

	require.Nil(t, check1.Result)
	require.Error(t, check1.Error)
	require.Nil(t, check2.Result)
	require.Error(t, check2.Error)
}

func TestAdblockerFirefox(t *testing.T) {
	t.Run("returns true when extension is a known adblocker", func(t *testing.T) {
		// Arrange
		extensionName := "adblocker ultimate"

		// Act
		isAdblocker := firefox.AdblockerFirefox(extensionName)

		// Assert
		require.True(t, isAdblocker)
	})

	t.Run("returns false when extension is not a known adblocker", func(t *testing.T) {
		// Arrange
		extensionName := "unknown extension"

		// Act
		isAdblocker := firefox.AdblockerFirefox(extensionName)

		// Assert
		require.False(t, isAdblocker)
	})
}

func TestExtensionFirefoxAdblocker(t *testing.T) {
	// setup a temporary directory for the test
	tmpDir, err := os.MkdirTemp("", "firefox")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir) // clean up

	// create a valid extensions.json file with a known adblocker
	extensionsFile := filepath.Join(tmpDir, "extensions.json")
	err = os.WriteFile(extensionsFile, []byte(`{
		"addons": [
			{
				"defaultLocale": {
					"name": "adblocker ultimate",
					"creator": "Test Creator"
				},
				"type": "extension",
				"active": true
			}
		]
	}`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	profileFinder := &MockProfileFinder{
		FirefoxFolderFunc: func() ([]string, error) {
			return []string{tmpDir}, nil
		},
	}

	check1, check2 := firefox.ExtensionFirefox(profileFinder)
	require.Nil(t, check1.Error)
	require.Nil(t, check2.Error)
	expected1 := checks.NewCheckResult(checks.ExtensionFirefoxID, 0, "adblocker ultimate,extension,Test Creator,true")
	expected2 := checks.NewCheckResult(checks.AdblockFirefoxID, 0, strconv.FormatBool(true))
	require.Equal(t, expected1, check1)
	require.Equal(t, expected2, check2)
}
