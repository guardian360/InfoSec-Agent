// Package (reporting page) main contains the entry point of the reporting page application
//
// Exported function(s): NewApp, NewTray
//
// Exported struct(s): App
package main

import (
	"context"
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
