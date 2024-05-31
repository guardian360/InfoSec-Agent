# Get version from latest git tag
$VERSION=(git describe --tags $(git rev-list --tags --max-count=1)).Substring(1)

# Build executables for production
$currentDir = Get-Location
if ($currentDir -like "*scripts") {
    Set-Location ..
}
#.\scripts\build-prod.bat

# Generate installer with Inno Setup
iscc .\scripts\generate-installer.iss /DMyAppVersion=$VERSION