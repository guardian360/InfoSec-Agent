// Package gamification handles the gamification within the application, to reward users for performing security checks and staying secure.
//
// Exported function(s): UpdateGameState, PointCalculation, LighthouseStateTransition
package gamification

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/database"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/usersettings"
)

// GameState is a struct that represents the state of the gamification.
// This consists of the user's points, a history of all previous points, and a lighthouse state.
type GameState struct {
	Points          int
	PointsHistory   []int
	TimeStamps      []time.Time
	LighthouseState int
}

// UpdateGameState updates the game state based on the scan results and the current game state.
//
// Parameters:
//   - scanResults ([]checks.Check): The results of the scans.
//   - databasePath (string): The path to the database file.
//
// Returns: The updated game state with the new points amount and new lighthouse state.
func UpdateGameState(scanResults []checks.Check, databasePath string) (GameState, error) {
	gs := GameState{Points: 0, PointsHistory: nil, TimeStamps: nil, LighthouseState: 0}

	// Loading the game state from the user settings and putting it in the game state struct
	points := usersettings.LoadUserSettings().Points
	pointsHistory := usersettings.LoadUserSettings().PointsHistory
	timeStamps := usersettings.LoadUserSettings().TimeStamps
	lighthouseState := usersettings.LoadUserSettings().LighthouseState
	gs.Points = points
	gs.PointsHistory = pointsHistory
	gs.TimeStamps = timeStamps
	gs.LighthouseState = lighthouseState

	gs, err := PointCalculation(gs, scanResults, databasePath)
	if err != nil {
		logger.Log.ErrorWithErr("Error calculating points:", err)
		return gs, err
	}
	gs = LighthouseStateTransition(gs)

	// Saving the game state in the user settings
	err = usersettings.SaveUserSettings(usersettings.UserSettings{
		Points:          gs.Points,
		PointsHistory:   gs.PointsHistory,
		TimeStamps:      gs.TimeStamps,
		LighthouseState: gs.LighthouseState,
	})
	if err != nil {
		logger.Log.Warning("Language setting not saved to file")
	}
	return gs, nil
}

// PointCalculation calculatese the number of points for the user based on the check results.
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
	}
	gs.PointsHistory = append(gs.PointsHistory, gs.Points)
	gs.TimeStamps = append(gs.TimeStamps, time.Now())

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
	case gs.Points < 10 && sufficientActivity(gs):
		gs.LighthouseState = 5 // The best state
	case gs.Points < 20 && sufficientActivity(gs):
		gs.LighthouseState = 4
	case gs.Points < 30 && sufficientActivity(gs):
		gs.LighthouseState = 3
	case gs.Points < 40 && sufficientActivity(gs):
		gs.LighthouseState = 2
	case gs.Points < 50 && sufficientActivity(gs):
		gs.LighthouseState = 1
	default:
		gs.LighthouseState = 1
	}
	return gs
}

// sufficientActivity checks if the user has been active enough to transition to another lighthouse state
//
// Parameters:
//   - gs (GameState): The game state of the user.
//   - duration (time.Duration): The duration of time that the user needs to be active.
//
// Returns:
//   - bool: whether the user has been active enough.
func sufficientActivity(gs GameState) bool {
	// The duration threshold of which the user has been active
	// Note that we define active as having performed a security check more than [1 week ago]
	requiredDuration := 7 * 24 * time.Hour

	if len(gs.TimeStamps) == 0 {
		return false
	}
	oldestRecord := gs.TimeStamps[0] // The oldest record is the first timestamp made

	return time.Since(oldestRecord) > requiredDuration
}
