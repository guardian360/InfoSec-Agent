package main

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/tray"
)

type Tray struct {
}

func NewTray() *Tray {
	return &Tray{}
}

func (t *Tray) ScanNow() {
	tray.ScanNow()
}

func (t *Tray) ChangeLanguage() {
	tray.ChangeLanguage()
}

func (t *Tray) ChangeScanInterval() {
	tray.ChangeScanInterval()
}
