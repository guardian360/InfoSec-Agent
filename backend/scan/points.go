// Package points ...
//
// Exported function(s):
package scan

import (
	"database/sql"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
)

type GameState struct {
	Points         int
	PointsHistory  []int
	LigthouseState int
}

func PointCalculation(gs GameState) GameState {
	gs.Points = 0

	for _, check := range securityChecks {
		result := check()
		result.ResultID = 0
		if result.Error != nil {
			logger.Log.ErrorWithErr("Error performing security checks", result.Error)
		} else {
			gs.Points++
		}
		db, err := sql.Open("sqlite", "./database.db")

		if err != nil {
			logger.Log.ErrorWithErr("Error opening database:", err)
		}
		sev, err := GetSeverity(db, result.IssueID, result.ResultID)

		gs.Points = gs.Points + sev
		gs.PointsHistory = append(gs.PointsHistory, gs.Points)

		_, err = db.Exec("SELECT severity FROM issues ")

	}
	return gs
}
