#define MyAppName "InfoSec-Agent"
#define MyAppVersion "0.2.0"
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
ArchitecturesAllowed=x64
ArchitecturesInstallIn64BitMode=x64
DefaultDirName={autopf}\{#MyAppName}
DisableProgramGroupPage=yes
; Uncomment the following line to run in non administrative install mode (install for current user only.)
;PrivilegesRequired=lowest
OutputBaseFilename=InfoSec-Agent-{#MyAppVersion}-Setup
OutputDir=.
Compression=lzma
SolidCompression=yes
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
Source: "{#SourcePath}\{#MyAppExeName}"; DestDir: "{app}"; Flags: ignoreversion
Source: "{#SourcePath}\reporting-page\build\bin\InfoSec-Agent-Reporting-Page.exe"; DestDir: "{app}\reporting-page\build\bin"; Flags: ignoreversion
Source: "{#SourcePath}\backend\localization\localizations_src\*"; DestDir: "{app}\backend\localization\localizations_src"; Flags: ignoreversion recursesubdirs createallsubdirs
Source: "{#SourcePath}\reporting-page\database.db"; DestDir: "{app}\reporting-page"; Flags: ignoreversion
; NOTE: Don't use "Flags: ignoreversion" on any shared system files

[Icons]
Name: "{autoprograms}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"
Name: "{autodesktop}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"; Tasks: desktopicon

[Registry]
Root: HKLM; Subkey: "Software\{#MyAppName}"; Flags: uninsdeletevalue uninsdeletekeyifempty; Check: IsWin64
Root: HKLM; Subkey: "Software\{#MyAppName}\Reporting Page"; ValueType: string; ValueName: "exe"; ValueData: "{app}\reporting-page\InfoSec-Agent-Reporting-Page.exe"; Flags: uninsdeletevalue uninsdeletekeyifempty; Check: IsWin64

[Run]
Filename: "{app}\{#MyAppExeName}"; Description: "{cm:LaunchProgram,{#StringChange(MyAppName, '&', '&&')}}"; Flags: nowait postinstall skipifsilent

