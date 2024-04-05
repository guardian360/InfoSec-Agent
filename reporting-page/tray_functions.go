package main

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/tray"
)

// TODO: fix this comment once copilot decides to cooperate

// Tray represents the system tray, used for calling system tray functions from the front-end
type Tray struct {
}

// TODO: fix this comment once copilot decides to cooperate

// NewTray creates a new Tray struct
//
// Parameters: _
//
// Returns: a pointer to a new Tray struct
func NewTray() *Tray {
	return &Tray{}
}

// TODO: fix this comment once copilot decides to cooperate

// ScanNow calls the ScanNow function from the tray package
//
// Parameters: t (*Tray) - a pointer to the Tray struct
//
// Returns: _
func (t *Tray) ScanNow() ([]checks.Check, error) {
	return tray.ScanNow()
}

// TODO: fix this comment once copilot decides to cooperate

// ChangeLanguage calls the ChangeLanguage function from the tray package
//
// Parameters: t (*Tray) - a pointer to the Tray struct
//
// Returns: _
func (t *Tray) ChangeLanguage() {
	tray.ChangeLanguage()
	tray.RefreshMenu()
}

// TODO: fix this comment once copilot decides to cooperate

// ChangeScanInterval calls the ChangeScanInterval function from the tray package
//
// Parameters: t (*Tray) - a pointer to the Tray struct
//
// Returns: _
func (t *Tray) ChangeScanInterval() {
	tray.ChangeScanInterval()
}
