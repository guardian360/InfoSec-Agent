package main

import (
	"context"

	"github.com/InfoSec-Agent/InfoSec-Agent/localization"
)

var language = 0

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Localize calls the Localize function from the localization package and passes the given language and ID.
// Wails binds this function to the frontend.
func (a *App) Localize(MessageID string) string {
	return localization.Localize(language, MessageID)
}
