# Running/Using the program
From the terminal, within the SystemTray folder, run the command `go run .`
This will start the program and a new icon should appear in your system tray. (Icon currently only works for Windows, not sure about the rest)

Clicking on the icon will show a menu containing actions that the program can execute.
Any confirmations/messages that the application sends will be send to the console.

The program can be exited by selecting 'Quit' in the menu or by manually interrupting the command line (Ctrl + c)

# Running the tests
From the terminal, within the SystemTray folder, run the command `go test ./...` which will run all tests in the current folder and all subfolders
If you wish to run a specific test, you can run the command `go test -run regexp`. This will only run the tests that match the regexp.