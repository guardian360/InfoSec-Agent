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

type GameState struct {
	Points          int
	PointsHistory   []int
	LighthouseState int
}

// PointCalculation calculates the number of points for the user based on the check results.
//
// Parameters:
//   - gs (GameState): The current game state, which includes the user's points and lighthouse state.
//   - securityChecks ([] func() checks.Check): A slice of check functions, which will be executed to calculate the points amount.
//
// Returns:
//   - GameState: The updated game state with the new points amount.
func PointCalculation(gs GameState, securityChecks []func() checks.Check) GameState {
	gs.Points = 0

	for _, check := range securityChecks {
		result := check()
		result.ResultID = 0
		if result.Error != nil {
			logger.Log.ErrorWithErr("Error performing security checks", result.Error)
		} else {
			gs.Points++
		}
		db, err := sql.Open("sqlite", "backend/database.db")

		// Note that due to opening of non-existent database, it will create one, so there can not be an error.
		// This is a potential bug, as the database is created in the current directory, which is not the intended location.
		if err != nil {
			logger.Log.ErrorWithErr("Error opening database:", err)
		}
		sev, err := database.GetSeverity(db, result.IssueID, result.ResultID)

		if err != nil {
			logger.Log.ErrorWithErr("Error getting severity:", err)
		}
		logger.Log.Info("Issue ID: " + strconv.Itoa(result.IssueID) + " Severity: " + strconv.Itoa(sev))

		// When severity is of the Informative level , we do not want to adjust the points
		if sev != 4 {
			gs.Points = gs.Points + sev
		}

		gs.PointsHistory = append(gs.PointsHistory, gs.Points)

	}
	return gs
}

// LighthouseStateTransition determines the lighthouse state based on the user's points. (TO DO: it should be based also on point history / time of use)
//
// Parameters:
//   - gs (GameState): The current game state, which includes the user's points and lighthouse state.
//
// Returns:
//   - GameState: The updated game state with the new ligthouse state.
func LighthouseStateTransition(gs GameState) GameState {
	if gs.Points < 10 {
		gs.LighthouseState = 0
	} else if gs.Points < 20 {
		gs.LighthouseState = 1
	} else if gs.Points < 30 {
		gs.LighthouseState = 2
	} else if gs.Points < 40 {
		gs.LighthouseState = 3
	} else if gs.Points < 50 {
		gs.LighthouseState = 4
	} else {
		gs.LighthouseState = 5
	}
	return gs
}