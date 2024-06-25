package gamification_test

import (
	"errors"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/usersettings"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/gamification"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
	"github.com/stretchr/testify/require"
)

// TestMain sets up the necessary environment for the gamification tests and runs them.
//
// This function sets up the logger for the tests and runs the tests.
//
// Parameters:
//   - m *testing.M: The testing framework that manages and runs the tests.
//
// Returns: None. The function calls os.Exit with the exit code returned by m.Run().
func TestMain(m *testing.M) {
	logger.SetupTests()

	// Run tests
	exitCode := m.Run()

	os.Exit(exitCode)
}

// TestUpdateGameState tests the general workings of the UpdateGameState function.
//
// Parameters:
//   - t (*testing.T): A pointer to an instance of the testing framework, used for reporting test results.
//
// No return values.
func TestUpdateGameState(t *testing.T) {
	// Mock the following functions and variables
	mockDatabasePath := "../../reporting-page/frontend/src/databases/database.en-GB.json"
	mockScanResults := []checks.Check{
		{
			IssueID:  29,
			ResultID: 1, // severity 2
		},
		{
			IssueID:  5,
			ResultID: 1, // severity 1
		},
	}
	t.Run("Test", func(t *testing.T) {
		_, err := gamification.UpdateGameState(mockScanResults, mockDatabasePath, gamification.RealPointCalculationGetter{}, usersettings.RealSaveUserSettingsGetter{})
		require.NoError(t, err)
	})
}

// TestPointCalculation tests the PointCalculation function for certain states
//
// Parameters:
//   - t (*testing.T): A pointer to an instance of the testing framework, used for reporting test results.
//
// No return values.
func TestPointCalculation(t *testing.T) {
	securityChecks := []checks.Check{
		{
			IssueID:  29,
			ResultID: 0,
		},
		{
			IssueID:  5,
			ResultID: 1,
		},
	}
	tests := []struct {
		name string
		gs   gamification.GameState
	}{
		{name: "GameState with no points and no point history", gs: gamification.GameState{Points: 0, PointsHistory: nil, TimeStamps: nil, LighthouseState: 0}},
		{name: "GameState with positive points and no point history", gs: gamification.GameState{Points: 29, PointsHistory: nil, TimeStamps: nil, LighthouseState: 3}},
		{name: "GameState with positive points and point history", gs: gamification.GameState{Points: 37, PointsHistory: []int{50, 28, 34}, TimeStamps: []time.Time{time.Now(), time.Now(), time.Now()}, LighthouseState: 2}},
	}
	getter := gamification.RealPointCalculationGetter{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := getter.PointCalculation(tt.gs, securityChecks, "../../reporting-page/frontend/src/databases/database.en-GB.json")
			require.NoError(t, err)
		})
	}
}

// TestLighthouseStateTransition tests the LighthouseStateTransition function for various time stamps and points input.
//
// Parameters:
//   - t (*testing.T): A pointer to an instance of the testing framework, used for reporting test results.
//
// No return values.
func TestLighthouseStateTransition(t *testing.T) {
	tests := []struct {
		points                  int
		timestamps              []time.Time
		expectedLighthouseState int
	}{
		{points: 4, timestamps: []time.Time{time.Date(2023, 6, 18, 12, 0, 0, 0, time.Now().Local().Location())}, expectedLighthouseState: 4},
		{points: 13, timestamps: []time.Time{time.Date(2023, 5, 14, 12, 0, 0, 0, time.Now().Local().Location())}, expectedLighthouseState: 3},
		{points: 26, timestamps: []time.Time{time.Date(2023, 2, 23, 12, 0, 0, 0, time.Now().Local().Location())}, expectedLighthouseState: 2},
		{points: 35, timestamps: []time.Time{time.Date(2023, 6, 7, 12, 0, 0, 0, time.Now().Local().Location())}, expectedLighthouseState: 1},
		{points: 44, timestamps: []time.Time{time.Date(2023, 1, 3, 12, 0, 0, 0, time.Now().Local().Location())}, expectedLighthouseState: 0},
		{points: 70, timestamps: []time.Time{time.Date(2023, 11, 24, 12, 0, 0, 0, time.Now().Local().Location())}, expectedLighthouseState: 0},
		{points: 13, timestamps: []time.Time{time.Now()}, expectedLighthouseState: 2},
	}
	for i, tt := range tests {
		t.Run("Test "+strconv.Itoa(i), func(t *testing.T) {
			gs := gamification.GameState{Points: tt.points, PointsHistory: nil, TimeStamps: tt.timestamps, LighthouseState: 9}
			got := gamification.LighthouseStateTransition(gs)
			require.Equal(t, tt.expectedLighthouseState, got.LighthouseState)
		})
	}
}

func TestSufficienActivityFail(t *testing.T) {
	testGamestate := gamification.GameState{Points: 0, PointsHistory: nil, TimeStamps: nil, LighthouseState: 0}
	got := gamification.SufficientActivity(testGamestate)
	require.False(t, got)
}

type MockPointCalculationGetter struct{}

func (m MockPointCalculationGetter) PointCalculation(gs gamification.GameState, _ []checks.Check, _ string) (gamification.GameState, error) {
	return gs, errors.New("mock error")
}

func TestUpdateGameState_Error(t *testing.T) {
	mockDatabasePath := "../../reporting-page/frontend/src/databases/database.en-GB.json"
	mockScanResults := []checks.Check{
		{
			IssueID:  29,
			ResultID: 1, // severity 2
		},
		{
			IssueID:  5,
			ResultID: 1, // severity 1
		},
	}
	getter := MockPointCalculationGetter{}
	_, err := gamification.UpdateGameState(mockScanResults, mockDatabasePath, getter, usersettings.RealSaveUserSettingsGetter{})
	require.Error(t, err)
	require.Equal(t, "mock error", err.Error())
}

type MockSaveUserSettingsGetter struct{}

func (m MockSaveUserSettingsGetter) SaveUserSettings(_ usersettings.UserSettings) error {
	return errors.New("mock error")
}

func TestUpdateGameState_SaveSettingsError(t *testing.T) {
	mockDatabasePath := "../../reporting-page/frontend/src/databases/database.en-GB.json"
	mockScanResults := []checks.Check{
		{
			IssueID:  29,
			ResultID: 1, // severity 2
		},
		{
			IssueID:  5,
			ResultID: 1, // severity 1
		},
	}
	getter := gamification.RealPointCalculationGetter{}
	userGetter := MockSaveUserSettingsGetter{}
	got, _ := gamification.UpdateGameState(mockScanResults, mockDatabasePath, getter, userGetter)
	require.NotNil(t, got)
}

func TestPointCalculation_UnmarshalError(t *testing.T) {
	getter := gamification.RealPointCalculationGetter{}
	gs := gamification.GameState{Points: 0, PointsHistory: nil, TimeStamps: nil, LighthouseState: 0}

	scanResults := []checks.Check{}
	invalidJSONFilePath := "../../invalid.json"

	_, err := getter.PointCalculation(gs, scanResults, invalidJSONFilePath)

	require.Error(t, err)
}
