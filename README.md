<img src="https://github.com/InfoSec-Agent/InfoSec-Agent/raw/main/reporting-page/frontend/src/assets/images/logoTeamA-transformed.png" alt="InfoSec-Agent" height="192" />

[![Go Reference](https://pkg.go.dev/badge/github.com/InfoSec-Agent/InfoSec-Agent.svg)](https://pkg.go.dev/github.com/InfoSec-Agent/InfoSec-Agent)
[![Go Report Card](https://goreportcard.com/badge/github.com/InfoSec-Agent/InfoSec-Agent)](https://goreportcard.com/report/github.com/InfoSec-Agent/InfoSec-Agent)
[![Open Source Insights](https://img.shields.io/badge/Open%20Source%20Insights-2ea5b3)](https://deps.dev/go/github.com%2FInfoSec-Agent%2FInfoSec-Agent/)
[![GPL Licence](https://badges.frapsoft.com/os/gpl/gpl.svg?v=103)](LICENSE)

# InfoSec Agent
The InfoSec Agent is a security and privacy tool for Windows 10 and 11.

## Summary
The InfoSec Agent project aims to improve the security and privacy of Windows computer users. Currently, there are applications available that do this, but they are mainly targeted at large companies. The goal of this project is to make this accessible to everyone. An application is being developed that collects information about the user's system to discover any security or privacy related vulnerabilities. The results will be presented to the user in a special dashboard, showing the current status of the system, including recommended actions to improve it.

## Affiliations
This project is a collaborative effort involving nine students from Utrecht University in The Netherlands, in partnership with the Dutch IT company [Guardian360](https://www.guardian360.net/). It serves as the Software Project for the [Bachelor's Programme in Computing Sciences at the UU](https://www.uu.nl/en/organisation/department-of-information-and-computing-sciences/education/bachelors-programmes/computing-sciences).


This project is also supported by funding from the [SIDN Fund](https://www.sidnfonds.nl/projecten/infosec-agent) (Stichting Internet Domeinregistratie Nederland), the Dutch domain name registrar.

## Contributing
InfoSec-Agent is an Open-Source project licensed under the GPL-3.0 License. However, due to its origins as a Utrecht University assignment, public contributions to this repository will only be merged after the completion of this assignment, which is scheduled for June 24, 2024.

Feel free to report any bugs or issues you encounter. Your feedback is valuable and helps improve the InfoSec-Agent project.

# Program instructions
## Running/Using the program
From the terminal, within the InfoSec-Agent folder, run the command `go run .`
This will start the program and a new icon should appear in your system tray.

Clicking on the icon will show a menu containing actions that the program can execute.
Any confirmations/messages/errors that the application sends will be sent to the log file located in the %AppData%/InfoSec-Agent directory.
The reporting page has its own log file located in the same directory.

The program can be exited by selecting 'Quit' in the menu or by manually interrupting the command line (Ctrl + c)

# Running tests
## Frontend
From the terminal, within the ***InfoSec-Agent/reporting-page/frontend*** folder, run the command `npm test` which will run all test found in the current folder and all subfolders. The tests are located in the ***InfoSec-Agent/reporting-page/frontend/test*** folder and only tests defined with the `.test.js` extension will be run.

To receive a coverage report of the tests, run the command `npm test -- --coverage`, which will show a table containing coverage of files being tested by the tests.

If you wish to run a specific test, you can run the command `npm test -- --testPathPattern=test/specific-test` where you would change specific-test to the filename of the test you would like to run. To get coverage from this single test, add `--coverage` to the end of the command.

## Backend
From the terminal, within the InfoSec-Agent folder, run the command `go test ./...` which will run all tests in the current folder and all subfolders.
If you wish to run a specific test, you can run the command `go test -run regexp`. This will only run the tests that match the regular expression.

# Linters
## Frontend

### Installation

The ESLint (for JavaScript) and Stylelint (for css) are installed trough npm.

The configuration for these linters is already defined in reporting-page/package.json and reporting-page/frontend/package.json

To install the linters, open a terminal in the ***InfoSec-Agent/reporting-page/frontend*** directory and run:

```
npm ci
```

### Usage
To run ESlint on all JavaScript code, open a terminal in the ***InfoSec-Agent/reporting-page/frontend*** directory and run:

```
npx eslint **/*.js
```

To run Stylelint on all CSS code, open a terminal in the ***InfoSec-Agent/reporting-page/frontend*** directory and run:

```
npx stylelint **/*.css
```

Both linters accept a `--fix` flag, to let the linters fix all issues they are able to fix automatically.

Adding this flag will change the code files, make sure to inspect the changes before commiting/pushing them.
Not all found issues can be fixed automatically, you may need to fix some issues yourself.

### Configuration

Configuration for these two linters can be found in the reporting-page/frontend/package.json file.

## Backend

### Installation

To install the golangci-lint binary on Windows run the following commands in ***git bash***
(not in a cmd or PowerShell terminal):

```
# binary will be $(go env GOPATH)/bin/golangci-lint
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.57.1
```

See [golangci-ci install documentation](https://golangci-lint.run/welcome/install/) for more information about installation.

After installation of the golangci-lint binary, open an ***elevated*** PowerShell and run the following commands:

```
# This is the default installation path for git for Windows,
# you may need to change it to your own custom installation path 
$Env:Path += ";C:\Program Files\Git\usr\bin"
[Environment]::SetEnvironmentVariable("Path", $env:Path, [System.EnvironmentVariableTarget]::Machine)
```

The linter should now be installed, you can check the version of the linter by opening any terminal and running:

```
golangci-lint --version
```

### Usage

To run the linter, open a terminal in the ***root*** of the InfoSec-Agent repository and run:

```
golangci-lint run
```

The linter will output the found issues in the CLI and make a golangci-lint-report.json file.
You can format this json using any (online) software you like.

To let the linter fix the found issues it is able to fix, open a terminal in the ***root*** of the InfoSec-Agent repository and run:

```
golangci-lint run --fix
```

This command will change the code files, make sure to inspect the changes before commiting/pushing them.
Not all found issues can be fixed automatically, you may need to fix some issues yourself.

### Configuration

Configuration for golangci-lint can be found in the .golangci.yml file in the root of the repository.

# Front-end Wails information
 
All the following information applies to the reporting-page folder.

## Build Directory

The build directory is used to house all the build files and assets for your application.

The structure is:

* bin - Output directory
* darwin - macOS specific files
* windows - Windows specific files

## Windows

The `windows` directory contains the manifest and rc files used when building with `wails build`.
These may be customised for your application. To return these files to the default state, simply delete them and
build with `wails build`.

- `icon.ico` - The icon used for the application. This is used when building using `wails build`. If you wish to
  use a different icon, simply replace this file with your own. If it is missing, a new `icon.ico` file
  will be created using the `appicon.png` file in the build directory.
- `installer/*` - The files used to create the Windows installer. These are used when building using `wails build`.
- `info.json` - Application details used for Windows builds. The data here will be used by the Windows installer,
  as well as the application itself (right-click the exe -> properties -> details)
- `wails.exe.manifest` - The main application manifest file.

# Generating Installer

This project can generate an installer for the InfoSec-Agent application using Inno Setup.

## Requirements

To generate the installer the following software has to be installed:

- [Golang](https://go.dev/doc/install) - To generate the executable for the system tray application.
- [Wails](https://wails.io/docs/gettingstarted/installation) - To generate the executable for the reporting page.
  - Wails requires GoLang and [NPM](https://nodejs.org/en/download/)
- [Inno Setup](https://jrsoftware.org/isdl.php) - To generate the installer itself.
  - To use the installer script provided by this project, the Inno Setup directory has to be added to the Path, as the script calls Inno Setup's console-mode compiler, ISCC.exe

## Creating installer executable
How the installer executable is generated is defined in the generate-installer.iss Inno Setup Script file.
There is a generate-installer.ps1 PowerShell script to easily generate the installer.

The version of the software is set using a parameter given to the iscc command.
The PowerShell scripts obtains the version information from the most recent git tag, as the tags indicate version numbers.

To generate the installer yourself, without using the provided PowerShell script, you first need to build both the tray and reporting-page executables.
There is a build.bat script provided to easily do this.

Then, run the following command from the root of the repository:

```
iscc .\generate-installer.iss /DMyAppVersion=$VERSION
```

Here, the VERSION variable needs to be set beforehand to the desired version number.

This command will output a InfoSec-Agent-{Version number}-Setup.exe file in the root of the repository.