package integration

import (
	"testing"
	"time"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/gamification"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/usersettings"
	"github.com/stretchr/testify/require"
)

func TestIntegrationUpdateGameState(t *testing.T) {
	// Mock scan results and database path
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
	mockDatabasePath := "../../reporting-page/database.db"

	// Run UpdateGameState
	gs, err := gamification.UpdateGameState(mockScanResults, mockDatabasePath)
	require.NoError(t, err)

	// Verify that the points, points history, and lighthouse state are correctly updated
	expectedPoints := 3
	expectedPointsHistory := []int{3, 3}
	expectedLighthouseState := 1

	require.Equal(t, expectedPoints, gs.Points)
	require.Equal(t, expectedPointsHistory, gs.PointsHistory)
	require.Equal(t, expectedLighthouseState, gs.LighthouseState)

	// Verify that the updated game state is correctly saved to the user settings
	userSettings := usersettings.LoadUserSettings()
	require.Equal(t, gs.Points, userSettings.Points)
	require.Equal(t, gs.PointsHistory, userSettings.PointsHistory)
	require.Equal(t, gs.LighthouseState, userSettings.LighthouseState)
}

func TestIntegrationPointCalculation(t *testing.T) {
	// Mock scan results
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

	mockDatabasePath := "../../reporting-page/database.db"

	gs := gamification.GameState{Points: 0, PointsHistory: nil, TimeStamps: nil, LighthouseState: 0}

	// Run PointCalculation
	var err error
	gs, err = gamification.PointCalculation(gs, mockScanResults, mockDatabasePath)
	require.NoError(t, err)

	// Verify that the points are correctly calculated
	expectedPoints := 3
	require.Equal(t, expectedPoints, gs.Points)

	// Verify that the points history and timestamps are correctly updated in the database
	userSettings := usersettings.LoadUserSettings()
	require.Contains(t, userSettings.PointsHistory, gs.Points)
	require.Len(t, userSettings.TimeStamps, 1)
}

func TestIntegrationLighthouseStateTransition(t *testing.T) {
	// Mock points and activity level
	mockPoints := 55
	mockPointsHistory := []int{50, 28, 34}
	date1 := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	date2 := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	date3 := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	mockTimeStamps := []time.Time{date1, date2, date3}
	gs := gamification.GameState{Points: mockPoints, PointsHistory: mockPointsHistory, TimeStamps: mockTimeStamps, LighthouseState: 0}

	// Run LighthouseStateTransition
	lighthouseState := gamification.LighthouseStateTransition(gs)

	// Verify that the lighthouse state is correctly updated
	expectedLighthouseState := 1
	require.Equal(t, expectedLighthouseState, lighthouseState)

	// Verify that the updated lighthouse state is correctly saved to the user settings
	userSettings := usersettings.LoadUserSettings()
	require.Equal(t, lighthouseState, userSettings.LighthouseState)
}
