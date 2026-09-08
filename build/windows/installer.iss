; =====================================================================
; IT Toolbox - Inno Setup Script
; Menghasilkan installer Windows (.exe) dengan wizard instalasi ramah
; pengguna ("tinggal klik-klik aja") dan auto-shortcut di Desktop & Start Menu.
; =====================================================================

#define MyAppName "IT Toolbox"
#define MyAppVersion "1.0.0"
#define MyAppPublisher "Yudz"
#define MyAppURL "https://github.com/yudz/it-toolbox"
#define MyAppExeName "it-toolbox.exe"

[Setup]
; ID Unik Aplikasi (Jangan diubah agar update versi menimpa instalasi lama)
AppId={{D37E88A1-8012-429C-A2C3-5D7E815F8587}
AppName={#MyAppName}
AppVersion={#MyAppVersion}
AppPublisher={#MyAppPublisher}
AppPublisherURL={#MyAppURL}
AppSupportURL={#MyAppURL}
AppUpdatesURL={#MyAppURL}
DefaultDirName={autopf}\{#MyAppName}
DefaultGroupName={#MyAppName}
DisableProgramGroupPage=yes
; Tempat installer .exe disimpan setelah dikompilasi
OutputDir=..\..\dist
OutputBaseFilename=IT-Toolbox-Setup-x64
SetupIconFile=..\..\assets\icon.ico
UninstallDisplayIcon={app}\{#MyAppExeName}
Compression=lzma2/ultra64
SolidCompression=yes
WizardStyle=modern
ArchitecturesInstallIn64BitMode=x64compatible
PrivilegesRequired=lowest
PrivilegesRequiredOverridesAllowed=dialog

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"
Name: "indonesian"; MessagesFile: "compiler:Languages\Indonesian.isl"

[Tasks]
; Dicentang secara default sehingga shortcut Desktop langsung dibuat otomatis!
Name: "desktopicon"; Description: "{cm:CreateDesktopIcon}"; GroupDescription: "{cm:AdditionalIcons}"

[Files]
Source: "..\..\it-toolbox.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: "..\..\assets\icon.ico"; DestDir: "{app}"; Flags: ignoreversion
; CATATAN: Jangan gunakan "Flags: ignoreversion" pada file bersama sistem

[Icons]
; Shortcut di Start Menu / Program Files
Name: "{autoprograms}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"; IconFilename: "{app}\icon.ico"
; Shortcut di Desktop (Langsung muncul setelah instalasi)
Name: "{autodesktop}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"; Tasks: desktopicon; IconFilename: "{app}\icon.ico"

[Run]
; Opsi "Jalankan IT Toolbox" setelah wizard selesai (otomatis tercentang)
Filename: "{app}\{#MyAppExeName}"; Description: "{cm:LaunchProgram,{#StringChange(MyAppName, '&', '&&')}}"; Flags: nowait postinstall skipifsilent
