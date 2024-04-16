// Package (reporting page) main contains the entry point of the reporting page application
//
// Exported function(s): NewApp, NewTray
//
// Exported struct(s): App
package main

import (
	"embed"
	"net/http"
	"os"
	"strings"

	"github.com/InfoSec-Agent/InfoSec-Agent/localization"
	"github.com/InfoSec-Agent/InfoSec-Agent/logger"
	"github.com/InfoSec-Agent/InfoSec-Agent/tray"
	"github.com/InfoSec-Agent/InfoSec-Agent/usersettings"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

type FileLoader struct {
	http.Handler
}

func NewFileLoader() *FileLoader {
	return &FileLoader{}
}

func (h *FileLoader) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var err error
	requestedFilename := strings.TrimPrefix(req.URL.Path, "/")
	logger.Log.Info("Requesting file:" + requestedFilename)
	fileData, err := os.ReadFile(requestedFilename)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		logger.Log.Info("Could not load file" + requestedFilename)
	}
	if _, err = res.Write(fileData); err != nil {
		logger.Log.Info("Could not write file" + requestedFilename)
	}
}

// main is the entry point of the reporting page application. It initializes the localization settings, creates a new instance of the application, tray, and database, and starts the Wails application.
//
// This function first calls the Init function from the localization package to set up the localization settings for the application. It then creates new instances of the App, Tray, and DataBase structs.
// After setting up these instances, it creates a Wails application with specific options including the title, dimensions, startup behavior, asset server, background color, startup function, and bound interfaces.
// It also sets up the Windows-specific options for the Wails application, including the theme and custom theme settings.
// If an error occurs during the creation or startup of the Wails application, it is logged and the program terminates.
//
// Parameters: None.
//
// Returns: None. This function does not return a value as it is the entry point of the application.
func main() {
	logger.Setup(0, -1)
	logger.Log.Info("Reporting page starting")

	// Create a new instance of the app and tray struct
	app := NewApp()
	systemTray := NewTray(logger.Log)
	database := NewDataBase()
	customLogger := logger.Log
	localization.Init("../")
	lang := usersettings.LoadUserSettings("../usersettings").Language
	tray.Language = lang

	// Create a Wails application with the specified options
	err := wails.Run(&options.App{
		Title:       "reporting-page",
		Width:       1024,
		Height:      768,
		StartHidden: true,
		AssetServer: &assetserver.Options{
			Assets:  assets,
			Handler: NewFileLoader(),
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
			systemTray,
			database,
		},
		Logger: customLogger,
		Windows: &windows.Options{
			Theme: windows.SystemDefault,
			CustomTheme: &windows.ThemeSettings{
				DarkModeTitleBar:   windows.RGB(20, 20, 20),
				DarkModeTitleText:  windows.RGB(200, 200, 200),
				DarkModeBorder:     windows.RGB(20, 0, 20),
				LightModeTitleBar:  windows.RGB(200, 200, 200),
				LightModeTitleText: windows.RGB(20, 20, 20),
				LightModeBorder:    windows.RGB(200, 200, 200),
			},
		},
	})
	if err != nil {
		logger.Log.ErrorWithErr("Error creating Wails application:", err)
	}
}
