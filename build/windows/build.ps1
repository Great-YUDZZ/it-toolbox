<#
 =====================================================================
 IT Toolbox - PowerShell Build & Inno Setup Script
 =====================================================================
#>

$ErrorActionPreference = "Stop"
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Definition
$RootDir = Resolve-Path "$ScriptDir\..\.."
Set-Location $RootDir

Write-Host "========================================================" -ForegroundColor Cyan
Write-Host " Mempersiapkan Build IT Toolbox untuk Windows (x64)" -ForegroundColor Cyan
Write-Host "========================================================" -ForegroundColor Cyan

# 1. Cek Go
if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
    Write-Error "Golang tidak ditemukan di PATH! Pastikan Go telah terinstall."
}

# 2. Embed Icon Windows
Write-Host "[1/3] Memeriksa Windows Resource Icon..." -ForegroundColor Yellow
if (-not (Test-Path "rsrc_windows_amd64.syso")) {
    Write-Host "Menghasilkan syso resource icon..." -ForegroundColor Gray
    go run github.com/tc-hib/go-winres@latest simply --icon assets/icon.png
}

# 3. Kompilasi it-toolbox.exe
Write-Host "[2/3] Mengompilasi it-toolbox.exe (GUI mode)..." -ForegroundColor Yellow
$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -ldflags "-H=windowsgui -s -w" -o it-toolbox.exe .
Write-Host "[OK] it-toolbox.exe berhasil dikompilasi!" -ForegroundColor Green

# 4. Inno Setup Compiler
Write-Host "[3/3] Memeriksa Inno Setup Compiler (ISCC)..." -ForegroundColor Yellow
$isccCandidates = @(
    "${env:ProgramFiles(x86)}\Inno Setup 6\ISCC.exe",
    "${env:ProgramFiles}\Inno Setup 6\ISCC.exe",
    (Get-Command ISCC.exe -ErrorAction SilentlyContinue).Source
)

$isccPath = $isccCandidates | Where-Object { $_ -and (Test-Path $_) } | Select-Object -First 1

if ($isccPath) {
    Write-Host "Mengompilasi installer Inno Setup via $isccPath..." -ForegroundColor Green
    & "$isccPath" "build\windows\installer.iss"
    Write-Host ""
    Write-Host "========================================================" -ForegroundColor Green
    Write-Host " SUKSES! Installer Windows telah dibuat di:" -ForegroundColor Green
    Write-Host " dist\IT-Toolbox-Setup-x64.exe" -ForegroundColor Green
    Write-Host "========================================================" -ForegroundColor Green
} else {
    Write-Host "[INFO] Inno Setup 6 belum terpasang di sistem." -ForegroundColor Cyan
    Write-Host "Download gratis di https://jrsoftware.org/isdl.php" -ForegroundColor Gray
    Write-Host "Binary 'it-toolbox.exe' sudah siap dipakai." -ForegroundColor Green
}
