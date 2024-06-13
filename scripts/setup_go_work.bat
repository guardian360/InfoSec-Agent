:: Simple script to setup go work environment
:: This script is intended to be run from the scripts directory
:: It will initialize the go work environment, and then use the current directory and the reporting-page directory
@echo off
if %cd:~-7%%==scripts cd ..
go work init
go work use .
go work use .\reporting-page
go work sync