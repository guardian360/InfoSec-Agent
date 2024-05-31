:: This script is intended to be run from the root of the project
:: It will build the tray and reporting-page executables for production mode
@echo off

:: Use go-winres for windows resource file generation
go install github.com/tc-hib/go-winres@latest
go generate

:: Build the tray executable
go build -buildmode=exe -ldflags="-H=windowsgui -s -w" -tags prod

:: Prepare localization files for the reporting-page
mkdir localization
cp -r backend/localization/localizations_src/* localization/

:: Build the reporting-page executable
cd reporting-page
wails build -clean -tags prod

:: Cleanup resource and localization files
cd ..
rm *.syso
rm -rf localization