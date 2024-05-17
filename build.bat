:: This script is intended to be run from the root of the project
:: It will build the tray and reporting-page executables
@echo off
go build -buildmode=exe -ldflags -H=windowsgui
cd reporting-page
wails build