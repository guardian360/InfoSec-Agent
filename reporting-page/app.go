package main

import (
	"context"

	"github.com/InfoSec-Agent/InfoSec-Agent/localization"
	"github.com/InfoSec-Agent/InfoSec-Agent/tray"
)

// App is the main application struct, necessary for the Wails runtime
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
//
// Parameters: _
//
// Returns: a pointer to a new App struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved, so runtime methods can be called
//
// Parameters: ctx (context.Context) - the context of the application
//
// Returns: _
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Localize calls the Localize function from the localization package and passes the given language and ID.
// Wails binds this function to the frontend.
func (a *App) Localize(MessageID string) string {
	return localization.Localize(tray.Language(), MessageID)
}
