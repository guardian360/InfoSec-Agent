package gamification

import (
	"database/sql"
	"strconv"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/database"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
)

type GameState struct {
	Points         int
	PointsHistory  []int
	LigthouseState int
}

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

func LighthouseStateTransition(gs GameState) GameState {
	if gs.Points < 5 {
		gs.LigthouseState = 0
	} else if gs.Points < 10 {
		gs.LigthouseState = 1
	} else if gs.Points < 15 {
		gs.LigthouseState = 2
	} else {
		gs.LigthouseState = 3
	}
	return gs

}
