package checks

// TestLastPasswordChange tests if response is correct when password was not changed recently
//
// Parameters: _
//
// Returns: _
// func TestLastPasswordChange(t *testing.T) {
// 	// Save the original functions and restore them at the end
// 	originalCommand := executor
// 	originalUsername := fetcher
// 	defer func() {
// 		executor = originalCommand
// 		fetcher = originalUsername
// 	}()

// 	// Test when exec.Command returns an error
// 	utils.CurrentUsername = func() (string, error) {
// 		return "testuser", nil
// 	}
// 	exec.Command = func(command string, args ...string) *exec.Cmd {
// 		return exec.Command("nonexistent")
// 	}
// 	result = LastPasswordChange()
// 	assert.NotNil(t, result.Error)

// 	// Test when exec.Command returns valid output
// 	exec.Command = func(command string, args ...string) *exec.Cmd {
// 		output := []byte(`User name                    testuser
// 	Full Name                    Test User
// 	Comment
// 	User's comment
// 	Country/region code          000 (System Default)
// 	Account active               Yes
// 	Account expires              Never

// 	Password last set            01-01-2022 12:00:00 AM
// 	Password expires             Never
// 	Password changeable          01-01-2022 12:00:00 AM
// 	Password required            Yes
// 	User may change password     Yes

// 	Workstations allowed         All
// 	Logon script
// 	User profile
// 	Home directory
// 	Last logon                   Never

// 	Logon hours allowed          All`)
// 		cmd := exec.Command(command, args...)
// 		cmd.Stdout = bytes.NewBuffer(output)
// 		return cmd
// 	}
// 	expectedResult := NewCheckResult("LastPasswordChange", "You changed your password recently on 01-01-2022")
// 	result = LastPasswordChange()
// 	assert.Equal(t, expectedResult, result)
// }
