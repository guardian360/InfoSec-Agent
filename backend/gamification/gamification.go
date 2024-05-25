// Package gamification handles the gamification within the application, to reward users for performing security checks and staying secure.
//
// Exported function(s): PointCalculation, LighthouseStateTransition
package gamification

import (
	"database/sql"
	"strconv"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/scan"
)

// GameState is a struct that represents the state of the gamification.
// This consists of the user's points, a history of all previous points, and a lighthouse state.
type GameState struct {
	Points          int
	PointsHistory   []int
	LighthouseState int
}

// PointCalculation calculates the number of points for the user based on the check results.
//
// Parameters:
//   - gs (GameState): The current game state, which includes the user's points and lighthouse state.
//   - scanResults ([]checks.Check): The results of the scans.
//
// TO DO: points based on more factors than only the checks.
//
// Returns:
//   - GameState: The updated game state with the new points amount.
func PointCalculation(gs GameState, scanResults []checks.Check) (GameState, error) {
	var err error
	gs.Points = 0
	db, err := sql.Open("sqlite", "../../reporting-page/database.db")
	if err != nil {
		logger.Log.ErrorWithErr("Error opening database:", err)
		return gs, err
	}

	for _, result := range scanResults {
		if result.Error != nil {
			logger.Log.ErrorWithErr("Error reading scan result", result.Error)
			return gs, result.Error
		}
		sev, err := scan.GetSeverity(db, result.IssueID, result.ResultID)
		if err != nil {
			logger.Log.ErrorWithErr("Error getting severity:", err)
			return gs, result.Error
		}
		logger.Log.Info("Issue ID: " + strconv.Itoa(result.IssueID) + " Severity: " + strconv.Itoa(sev))

		// When severity is of the Informative level , we do not want to adjust the points
		if sev != 4 {
			gs.Points += sev
		}
		gs.PointsHistory = append(gs.PointsHistory, gs.Points)
	}
	return gs, nil
}

// LighthouseStateTransition determines the lighthouse state based on the user's points (the less points, the better)
//
// Parameters:
//   - gs (GameState): The current game state, which includes the user's points and lighthouse state.
//
// TO DO: transition based on more factors than only points (time, points history, etc.)
//
// Returns:
//   - GameState: The updated game state with the new lighthouse state.
func LighthouseStateTransition(gs GameState) GameState {
	switch {
	case (gs.Points < 10):
		gs.LighthouseState = 5
	case (gs.Points < 20):
		gs.LighthouseState = 4
	case (gs.Points < 30):
		gs.LighthouseState = 3
	case (gs.Points < 40):
		gs.LighthouseState = 2
	case (gs.Points < 50):
		gs.LighthouseState = 1
	default:
		gs.LighthouseState = 0
	}
	return gs
}
