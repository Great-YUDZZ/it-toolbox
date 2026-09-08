@echo off
REM =====================================================================
REM IT Toolbox - Windows Build & Inno Setup Packaging Script
REM Jalankan script ini di Windows untuk mengompilasi it-toolbox.exe
REM dan menghasilkan IT-Toolbox-Setup-x64.exe secara otomatis.
REM =====================================================================

setlocal enabledelayedexpansion

echo ========================================================
echo  Mempersiapkan Build IT Toolbox untuk Windows (x64)
echo ========================================================

cd /d "%~dp0..\.."

REM 1. Cek apakah Go terpasang
where go >nul 2>nul
if %errorlevel% neq 0 (
    echo [ERROR] Golang tidak ditemukan di PATH! Silakan install Go terlebih dahulu.
    pause
    exit /b 1
)

REM 2. Embed Icon Windows (jika rsrc/go-winres tersedia atau hasilkan syso)
echo [1/3] Memeriksa Windows Resource Icon...
if not exist "rsrc_windows_amd64.syso" (
    echo Menghasilkan resource syso untuk icon exe...
    go run github.com/tc-hib/go-winres@latest simply --icon assets/icon.png
)

REM 3. Kompilasi binary Windows tanpa console cmd (-H=windowsgui)
echo [2/3] Mengompilasi it-toolbox.exe (GUI mode, no CMD window)...
set GOOS=windows
set GOARCH=amd64
go build -ldflags "-H=windowsgui -s -w" -o it-toolbox.exe .
if %errorlevel% neq 0 (
    echo [ERROR] Kompilasi Go gagal!
    pause
    exit /b 1
)

echo [OK] it-toolbox.exe berhasil dikompilasi!

REM 4. Kompilasi Installer Inno Setup jika ISCC tersedia
echo [3/3] Memeriksa Inno Setup Compiler (ISCC)...
set "ISCC_PATH="
if exist "%ProgramFiles(x86)%\Inno Setup 6\ISCC.exe" (
    set "ISCC_PATH=%ProgramFiles(x86)%\Inno Setup 6\ISCC.exe"
) else if exist "%ProgramFiles%\Inno Setup 6\ISCC.exe" (
    set "ISCC_PATH=%ProgramFiles%\Inno Setup 6\ISCC.exe"
) else (
    where ISCC.exe >nul 2>nul
    if !errorlevel! equ 0 set "ISCC_PATH=ISCC.exe"
)

if defined ISCC_PATH (
    echo Mengompilasi installer Inno Setup...
    "!ISCC_PATH!" "build\windows\installer.iss"
    if !errorlevel! equ 0 (
        echo.
        echo ========================================================
        echo  SUKSES! Installer Windows telah dibuat di:
        echo  dist\IT-Toolbox-Setup-x64.exe
        echo ========================================================
    ) else (
        echo [WARNING] Pembuatan installer via ISCC gagal.
    )
) else (
    echo [INFO] Inno Setup 6 belum terpasang di komputer ini.
    echo Anda dapat mengunduh Inno Setup gratis di https://jrsoftware.org/isdl.php
    echo File it-toolbox.exe sudah siap untuk digunakan langsung atau dikemas.
)

echo Selesai.
pause
