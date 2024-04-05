package main

import (
	"context"

	"github.com/InfoSec-Agent/InfoSec-Agent/localization"
	"github.com/InfoSec-Agent/InfoSec-Agent/tray"
)

// App serves as the core structure for the Wails runtime.
//
// It encapsulates the application's context and provides methods for its management. This struct is pivotal for the lifecycle of the application, enabling the invocation of runtime methods and facilitating interactions with the frontend.
type App struct {
	ctx context.Context
}

// NewApp is a factory method that generates an instance of the App struct.
//
// This method is typically invoked at the start of the application's lifecycle. The created App instance serves as the primary interface for managing the application's context.
//
// Parameters: None.
//
// Returns: *App: A pointer to a newly created App instance.
func NewApp() *App {
	return &App{}
}

// startup is the initial setup function for the App instance.
//
// This method is invoked at the beginning of the application's lifecycle. It stores the application's context for future use, enabling the invocation of runtime methods.
//
// Parameters: ctx (context.Context) - The application's context.
//
// Returns: None. The method performs an action and does not return any value.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Localize retrieves the localized version of a given message ID.
//
// This method is used to fetch the localized text corresponding to a specific message ID, based on the currently set language. It is typically invoked when the frontend needs to display text to the user.
//
// Parameters: messageID (string) - The identifier of the message to be localized.
//
// Returns: string: The localized text corresponding to the provided message ID.
func (a *App) Localize(messageID string) string {
	return localization.Localize(tray.Language(), messageID)
}

// PrintFromFrontend logs the provided message to the console.
//
// This method is primarily used for debugging purposes, allowing messages from the frontend to be logged and viewed in the console.
//
// Parameters: message (string) - The text to be logged.
//
// Returns: None. The method performs an action and does not return any value.
func (a *App) PrintFromFrontend(message string) {
	// log.Println(message)
}
