@echo off
setlocal
cd /d "%~dp0"

echo ===================================================
echo   IT Toolbox - Buat Shortcut Desktop Otomatis
echo ===================================================
echo.

set "TARGET_EXE=%~dp0it-toolbox.exe"
set "ICON_FILE=%~dp0icon.ico"
set "DESKTOP_DIR=%USERPROFILE%\Desktop"
set "LNK_PATH=%DESKTOP_DIR%\IT Toolbox.lnk"

if not exist "%TARGET_EXE%" (
    echo [ERROR] it-toolbox.exe tidak ditemukan di folder ini!
    echo Pastikan Anda mengekstrak semua isi ZIP ke dalam satu folder.
    pause
    exit /b 1
)

set "VBS_SCRIPT=%TEMP%\create_it_toolbox_shortcut_%RANDOM%.vbs"

(
    echo Set WshShell = WScript.CreateObject("WScript.Shell"^)
    echo Set oShellLink = WshShell.CreateShortcut("%LNK_PATH%"^)
    echo oShellLink.TargetPath = "%TARGET_EXE%"
    echo oShellLink.WorkingDirectory = "%~dp0"
    echo oShellLink.WindowStyle = 1
    echo oShellLink.Description = "IT Toolbox - All-in-one Developer & Network Suite"
    if exist "%ICON_FILE%" (
        echo oShellLink.IconLocation = "%ICON_FILE%, 0"
    ) else (
        echo oShellLink.IconLocation = "%TARGET_EXE%, 0"
    )
    echo oShellLink.Save
) > "%VBS_SCRIPT%"

cscript //nologo "%VBS_SCRIPT%"
del "%VBS_SCRIPT%" 2>nul

echo [SUKSES] Shortcut 'IT Toolbox' telah berhasil dibuat di Desktop Anda!
echo.
timeout /t 3 >nul
exit /b 0
