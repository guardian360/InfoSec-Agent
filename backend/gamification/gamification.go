// Package gamification handles the gamification within the application, to reward users for performing security checks and staying secure.
//
// Exported function(s): PointCalculation, LighthouseStateTransition
package gamification

import (
	"database/sql"
	"strconv"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/database"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
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
//   - databasePath (string): The path to the database file.
//
// Returns:
//   - GameState: The updated game state with the new points amount.
func PointCalculation(gs GameState, scanResults []checks.Check, databasePath string) (GameState, error) {
	gs.Points = 0
	db, err := sql.Open("sqlite", databasePath)
	if err != nil {
		logger.Log.ErrorWithErr("Error opening database:", err)
		return gs, err
	}

	for _, result := range scanResults {
		if result.Error != nil {
			logger.Log.ErrorWithErr("Error reading scan result", result.Error)
			return gs, result.Error
		}
		sev, err1 := database.GetSeverity(db, result.IssueID, result.ResultID)
		if err1 != nil {
			logger.Log.ErrorWithErr("Error getting severity:", err1)
			return gs, err1
		}
		logger.Log.Info("Issue ID: " + strconv.Itoa(result.IssueID) + " Severity: " + strconv.Itoa(sev))

		// When severity is of the Informative level , we do not want to adjust the points
		if sev != 4 {
			gs.Points += sev
		}
		gs.PointsHistory = append(gs.PointsHistory, gs.Points)
	}
	// Close the database
	logger.Log.Debug("Closing database")
	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			logger.Log.ErrorWithErr("Error closing database:", err)
		}
	}(db)

	return gs, nil
}

// LighthouseStateTransition determines the lighthouse state based on the user's points (the less points, the better)
//
// Parameters:
//   - gs (GameState): The current game state, which includes the user's points and lighthouse state.
//
// Returns:
//   - GameState: The updated game state with the new lighthouse state.
func LighthouseStateTransition(gs GameState) GameState {
	switch {
	case gs.Points < 10:
		gs.LighthouseState = 5
	case gs.Points < 20:
		gs.LighthouseState = 4
	case gs.Points < 30:
		gs.LighthouseState = 3
	case gs.Points < 40:
		gs.LighthouseState = 2
	case gs.Points < 50:
		gs.LighthouseState = 1
	default:
		gs.LighthouseState = 0
	}
	return gs
}
