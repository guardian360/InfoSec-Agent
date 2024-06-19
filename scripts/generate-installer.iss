; MyAppVersion should be defined as a command line parameter, e.g.: iscc generate-installer.iss /DMyAppVersion=0.2.0
; This parameter should be set to the most recent git version tag
#ifndef MyAppVersion
  #define MyAppVersion "0.0.0"
#endif

#define MyAppName "InfoSec-Agent"
#define MyAppPublisher "Little Brother"
#define MyAppURL "github.com/InfoSec-Agent/InfoSec-Agent/"
#define MyAppExeName "InfoSec-Agent.exe"

[Setup]
; NOTE: The value of AppId uniquely identifies this application. Do not use the same AppId value in installers for other applications.
AppId={{3AA5C750-3A11-4FDC-9405-F9D85FFF977A}
AppName={#MyAppName}
AppVersion={#MyAppVersion}
AppVerName={#MyAppName} {#MyAppVersion}
AppPublisher={#MyAppPublisher}
AppPublisherURL={#MyAppURL}
AppSupportURL={#MyAppURL}
AppUpdatesURL={#MyAppURL}
AppMutex=InfoSec-Agent,InfoSec-Agent-Reporting-Page
ArchitecturesAllowed=x64
ArchitecturesInstallIn64BitMode=x64
CloseApplications=yes
DefaultDirName={autopf}\{#MyAppName}
DisableProgramGroupPage=yes
; Uncomment the following line to run in non administrative install mode (install for current user only.)
; PrivilegesRequired=lowest
OutputBaseFilename=InfoSec-Agent-{#MyAppVersion}-Setup
OutputDir=.
Compression=lzma
SetupIconFile=..\reporting-page\build\windows\icon.ico
SetupMutex=InfoSec-Agent-Setup
SolidCompression=yes
UninstallDisplayIcon={userappdata}\{#MyAppName}\icon\icon.ico
WizardStyle=modern

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"
Name: "dutch"; MessagesFile: "compiler:Languages\Dutch.isl"
Name: "french"; MessagesFile: "compiler:Languages\French.isl"
Name: "german"; MessagesFile: "compiler:Languages\German.isl"
Name: "portuguese"; MessagesFile: "compiler:Languages\Portuguese.isl"
Name: "spanish"; MessagesFile: "compiler:Languages\Spanish.isl"

[Tasks]
Name: "desktopicon"; Description: "{cm:CreateDesktopIcon}"; GroupDescription: "{cm:AdditionalIcons}"; Flags: unchecked

[Files]
Source: "{#SourcePath}\..\{#MyAppExeName}"; DestDir: "{app}"; Flags: ignoreversion
Source: "{#SourcePath}\..\reporting-page\build\bin\InfoSec-Agent-Reporting-Page.exe"; DestDir: "{app}\reporting-page\"; Flags: ignoreversion
Source: "{#SourcePath}\..\backend\localization\localizations_src\*"; DestDir: "{app}\localization\"; Flags: ignoreversion recursesubdirs createallsubdirs
Source: "{#SourcePath}\..\reporting-page\frontend\src\assets\images\*"; DestDir: "{app}\reporting-page\images\"; Flags: ignoreversion recursesubdirs createallsubdirs
Source: "{#SourcePath}\..\reporting-page\frontend\src\databases\database.en-GB.json"; DestDir: "{app}\localization\en-GB"; Flags: ignoreversion
Source: "{#SourcePath}\..\reporting-page\build\windows\icon.ico"; DestDir: "{userappdata}\{#MyAppName}\icon"; Flags: ignoreversion
Source: "{#SourcePath}\..\reporting-page\build\icon128.ico"; DestDir: "{userappdata}\{#MyAppName}\icon"; Flags: ignoreversion
; NOTE: Don't use "Flags: ignoreversion" on any shared system files

[Icons]
Name: "{autoprograms}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"; WorkingDir: "{app}"; IconFilename: "{#SourcePath}\..\reporting-page\build\windows\icon.ico"
Name: "{autodesktop}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"; WorkingDir: "{app}"; IconFilename: "{#SourcePath}\..\reporting-page\build\windows\icon.ico"; Tasks: desktopicon
Name: "{userstartup}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"; WorkingDir: "{app}"; IconFilename: "{#SourcePath}\..\reporting-page\build\windows\icon.ico"

[Registry]
Root: HKLM; Subkey: "Software\{#MyAppName}"; Flags: uninsdeletevalue uninsdeletekeyifempty
Root: HKLM; Subkey: "Software\{#MyAppName}\Reporting Page"; ValueType: string; ValueName: "exe"; ValueData: "{app}\reporting-page\InfoSec-Agent-Reporting-Page.exe"; Flags: uninsdeletevalue uninsdeletekeyifempty
Root: HKCR; Subkey: "infosecagent"; ValueType: string; ValueName: ""; ValueData: URL:infosecagent Protocol; Flags: uninsdeletekey
Root: HKCR; Subkey: "infosecagent"; ValueType: string; ValueName: URL Protocol; Flags: uninsdeletekey
Root: HKCR; Subkey: "infosecagent\shell\open\command"; ValueType: expandsz; ValueName: ""; ValueData: "{app}\reporting-page\InfoSec-Agent-Reporting-Page.exe"; Flags: uninsdeletekey
Root: HKCU; Subkey: "Software\Microsoft\Windows\CurrentVersion\App Paths\InfoSec-Agent-Reporting-Page.exe"; ValueType: string; ValueName: ""; ValueData: "{app}\reporting-page\InfoSec-Agent-Reporting-Page.exe"; Flags: uninsdeletekey
Root: HKCU; Subkey: "Software\Microsoft\Windows\CurrentVersion\App Paths\InfoSec-Agent-Reporting-Page.exe"; ValueType: string; ValueName: "Path"; ValueData: "{app}\reporting-page"; Flags: uninsdeletekey

[UninstallDelete]
Type: filesandordirs; Name: "{userappdata}\{#MyAppName}"

[Run]
Filename: "{app}\{#MyAppExeName}"; Description: "{cm:LaunchProgram,{#StringChange(MyAppName, '&', '&&')}}"; Flags: nowait postinstall skipifsilent

