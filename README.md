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
Any confirmations/messages/errors that the application sends will be sent to the log file located in the root of the project.
The reporting page has its own log file located in its root folder (folder with main.go).

The program can be exited by selecting 'Quit' in the menu or by manually interrupting the command line (Ctrl + c)

## Running the tests
From the terminal, within the InfoSec-Agent folder, run the command `go test ./...` which will run all tests in the current folder and all subfolders.
If you wish to run a specific test, you can run the command `go test -run regexp`. This will only run the tests that match the regular expression.

# Linters
## Frontend

### Installation

The ESLint (for JavaScript) and Stylelint (for css) are installed trough npm.

The configuration for these linters is already defined in reporting-page/package.json and reporting-page/frontend/package.json

To install the linters, open a terminal in the ***InfoSec-Agent/reporting-page*** directory and run:

```
npm ci
```

### Usage
To run ESlint on all JavaScript code, open a terminal in the ***InfoSec-Agent/reporting-page*** directory and run:

```
npx eslint **/*.js
```

To run Stylelint on all CSS code, open a terminal in the ***InfoSec-Agent/reporting-page*** directory and run:

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
