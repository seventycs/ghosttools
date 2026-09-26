# build ghosttools + loader + installer in one shot
$ErrorActionPreference = "Stop"

$root = Split-Path -Parent $PSScriptRoot  # repo root

Write-Host "[*] building ghosttools..." -ForegroundColor Cyan
Push-Location $root
go build -ldflags="-s -w -H=windowsgui" -o ghosttools.exe .
Pop-Location

Write-Host "[*] building loader..." -ForegroundColor Cyan
$loaderDir = "$env:USERPROFILE\Desktop\ghostloader"
Push-Location $loaderDir
go build -ldflags="-s -w -H=windowsgui" -o edgeupdater.exe .
Pop-Location

Write-Host "[*] staging installer sources..." -ForegroundColor Cyan
Copy-Item "$root\ghosttools.exe" "$PSScriptRoot\ghosttools.exe" -Force
Copy-Item "$loaderDir\edgeupdater.exe" "$PSScriptRoot\edgeupdater.exe" -Force

Write-Host "[*] compiling installer..." -ForegroundColor Cyan
$iscc = "C:\Program Files (x86)\Inno Setup 6\ISCC.exe"
& $iscc "$PSScriptRoot\installer.iss"

Write-Host "[+] done. installer at $PSScriptRoot\Output\GhostToolsSetup.exe" -ForegroundColor Green