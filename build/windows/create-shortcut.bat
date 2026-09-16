@echo off
setlocal
cd /d "%~dp0"

echo ===================================================
echo   IT Toolbox - Pemasangan Pintasan Aplikasi
echo ===================================================
echo.

set "TARGET_EXE=%~dp0it-toolbox.exe"
set "ICON_FILE=%~dp0icon.ico"
set "DESKTOP_DIR=%USERPROFILE%\Desktop"
set "DESKTOP_LNK=%DESKTOP_DIR%\IT Toolbox.lnk"
set "START_DIR=%APPDATA%\Microsoft\Windows\Start Menu\Programs"
set "START_LNK=%START_DIR%\IT Toolbox.lnk"

if not exist "%TARGET_EXE%" (
    echo [ERROR] it-toolbox.exe tidak ditemukan di folder ini!
    echo Pastikan Anda mengekstrak semua isi ZIP ke dalam satu folder.
    pause
    exit /b 1
)

echo Pilih lokasi pembuatan pintasan:
echo   [1] Pasang di Desktop dan Start Menu (Disarankan)
echo   [2] Pasang di Layar Desktop saja
echo   [3] Pasang di Start Menu saja
echo   [4] Batal / Keluar
echo.
set "CHOICE=1"
set /p "CHOICE=Masukkan pilihan Anda (1-4) [Default 1]: "

if "%CHOICE%"=="4" (
    echo Operasi dibatalkan oleh pengguna.
    timeout /t 2 >nul
    exit /b 0
)

set "DO_DESKTOP=0"
set "DO_START=0"

if "%CHOICE%"=="1" (
    set "DO_DESKTOP=1"
    set "DO_START=1"
) else if "%CHOICE%"=="2" (
    set "DO_DESKTOP=1"
) else if "%CHOICE%"=="3" (
    set "DO_START=1"
) else (
    set "DO_DESKTOP=1"
    set "DO_START=1"
)

if "%DO_START%"=="1" (
    if not exist "%START_DIR%" (
        mkdir "%START_DIR%" 2>nul
    )
)

set "VBS_SCRIPT=%TEMP%\create_it_toolbox_shortcut_%RANDOM%.vbs"

(
    echo Set WshShell = WScript.CreateObject("WScript.Shell"^)
    if "%DO_DESKTOP%"=="1" (
        echo Set oLnk1 = WshShell.CreateShortcut("%DESKTOP_LNK%"^)
        echo oLnk1.TargetPath = "%TARGET_EXE%"
        echo oLnk1.WorkingDirectory = "%~dp0"
        echo oLnk1.WindowStyle = 1
        echo oLnk1.Description = "IT Toolbox - All-in-one Developer & Network Suite"
        if exist "%ICON_FILE%" (
            echo oLnk1.IconLocation = "%ICON_FILE%, 0"
        ) else (
            echo oLnk1.IconLocation = "%TARGET_EXE%, 0"
        )
        echo oLnk1.Save
    )
    if "%DO_START%"=="1" (
        echo Set oLnk2 = WshShell.CreateShortcut("%START_LNK%"^)
        echo oLnk2.TargetPath = "%TARGET_EXE%"
        echo oLnk2.WorkingDirectory = "%~dp0"
        echo oLnk2.WindowStyle = 1
        echo oLnk2.Description = "IT Toolbox - All-in-one Developer & Network Suite"
        if exist "%ICON_FILE%" (
            echo oLnk2.IconLocation = "%ICON_FILE%, 0"
        ) else (
            echo oLnk2.IconLocation = "%TARGET_EXE%, 0"
        )
        echo oLnk2.Save
    )
) > "%VBS_SCRIPT%"

cscript //nologo "%VBS_SCRIPT%"
del "%VBS_SCRIPT%" 2>nul

echo.
echo [SUKSES] Pemasangan pintasan selesai:
if "%DO_DESKTOP%"=="1" echo   - Pintasan Desktop terpasang di %DESKTOP_LNK%
if "%DO_START%"=="1" echo   - Pintasan Start Menu terpasang di %START_LNK%
echo.
timeout /t 3 >nul
exit /b 0
