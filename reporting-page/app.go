package main

import (
	"context"
	"log"

	"github.com/InfoSec-Agent/InfoSec-Agent/localization"
	"github.com/InfoSec-Agent/InfoSec-Agent/tray"
)

// TODO: fix this comment once copilot decides to cooperate

// App is the main application struct, necessary for the Wails runtime
type App struct {
	ctx context.Context
}

// TODO: fix this comment once copilot decides to cooperate

// NewApp creates a new App application struct
//
// Parameters: _
//
// Returns: a pointer to a new App struct
func NewApp() *App {
	return &App{}
}

// TODO: fix this comment once copilot decides to cooperate

// startup is called when the app starts. The context is saved, so runtime methods can be called
//
// Parameters: ctx (context.Context) - the context of the application
//
// Returns: _
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// TODO: fix this comment once copilot decides to cooperate

// Localize is called when the frontend loads text. It localizes the text based on the current language.
//
// Parameters: messageID (string) - the ID of the message to localize
//
// Returns: localized string (string)
func (a *App) Localize(messageID string) string {
	return localization.Localize(tray.Language(), messageID)
}

// TODO: fix this comment once copilot decides to cooperate

// PrintFromFrontend Print prints the given message to the console
//
// Parameters: message (string) - the message to print
//
// Returns: _
func (a *App) PrintFromFrontend(message string) {
	log.Println(message)
}
