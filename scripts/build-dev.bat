:: This script is intended to be run from the scripts directory or the root directory
:: It will build the tray and reporting-page executables for dev mode
@echo off

if %cd:~-7%%==scripts cd ..

go build
cd reporting-page
wails build -clean
