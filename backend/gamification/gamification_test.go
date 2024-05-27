package gamification_test

import (
	"strconv"
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/gamification"
	"github.com/stretchr/testify/require"
)

// TestPointCalculation tests the PointCalculation function for certain states
//
// Parameters:
//   - t (*testing.T): A pointer to an instance of the testing framework, used for reporting test results.
//
// No return values.
func TestPointCalculation(t *testing.T) {
	securityChecks := []checks.Check{}

	tests := []struct {
		name string
		gs   gamification.GameState
	}{
		{name: "GameState with no points and no point history", gs: gamification.GameState{Points: 0, PointsHistory: []int{}, LighthouseState: 0}},
		{name: "GameState with positive points and no point history", gs: gamification.GameState{Points: 29, PointsHistory: []int{}, LighthouseState: 3}},
		{name: "GameState with positive points and point history", gs: gamification.GameState{Points: 37, PointsHistory: []int{50, 28, 34}, LighthouseState: 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := gamification.PointCalculation(tt.gs, securityChecks)
			require.NoError(t, err)
		})
	}
}

// TestLighthouseStateTransition tests the LighthouseStateTransition function for various points inputs.
//
// Parameters:
//   - t (*testing.T): A pointer to an instance of the testing framework, used for reporting test results.
//
// No return values.
func TestLighthouseStateTransition(t *testing.T) {
	tests := []struct {
		points                  int
		expectedLighthouseState int
	}{
		{points: 4, expectedLighthouseState: 5},
		{points: 13, expectedLighthouseState: 4},
		{points: 26, expectedLighthouseState: 3},
		{points: 35, expectedLighthouseState: 2},
		{points: 44, expectedLighthouseState: 1},
		{points: 70, expectedLighthouseState: 0},
	}
	for i, tt := range tests {
		t.Run("Test "+strconv.Itoa(i), func(t *testing.T) {
			gs := gamification.GameState{Points: tt.points, PointsHistory: []int{}, LighthouseState: 9}
			got := gamification.LighthouseStateTransition(gs)
			require.Equal(t, tt.expectedLighthouseState, got.LighthouseState)
		})
	}
}
