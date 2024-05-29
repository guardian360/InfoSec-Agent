:: This script is intended to be run from the root of the project
:: It will build the tray and reporting-page executables
@echo off
go install github.com/tc-hib/go-winres@latest
go generate
go build -buildmode=exe -ldflags="-H=windowsgui -s -w"
rm *.syso
cd reporting-page
wails build