# scripts/build.ps1 — build ghosttools with garble obfuscation

$ErrorActionPreference = "Continue"

Write-Host "[*] ensuring garble is installed..." -ForegroundColor Cyan
go install mvdan.cc/garble@latest

$garbleExe = "$env:USERPROFILE\go\bin\garble.exe"
if (-not (Test-Path $garbleExe)) {
    $garbleExe = "garble"
}

Write-Host "[*] go mod tidy..." -ForegroundColor Cyan
go mod tidy

Write-Host "[*] obfuscating + building..." -ForegroundColor Cyan
& $garbleExe -tiny -literals -seed=random build -ldflags="-s -w" -o ghosttools.exe .

if (Test-Path .\ghosttools.exe) {
    $size = (Get-Item .\ghosttools.exe).Length / 1KB
    Write-Host "[+] built ghosttools.exe ($([math]::Round($size,1)) KB)" -ForegroundColor Green
} else {
    Write-Host "[!] build failed" -ForegroundColor Red
}