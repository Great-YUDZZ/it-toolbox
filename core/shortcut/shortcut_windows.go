//go:build windows

package shortcut

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

//go:embed icon.ico
var embeddedIconIco []byte

func createDesktopShortcut(exePath string) {
	desktopDirs := getWindowsDesktopDirs()
	if len(desktopDirs) == 0 {
		return
	}

	exeDir := filepath.Dir(exePath)
	icoPath := filepath.Join(exeDir, "icon.ico")

	// If icon.ico doesn't exist in exeDir, extract embedded icon
	if _, err := os.Stat(icoPath); os.IsNotExist(err) && len(embeddedIconIco) > 0 {
		_ = os.WriteFile(icoPath, embeddedIconIco, 0644)
	}

	for _, desktopDir := range desktopDirs {
		lnkPath := filepath.Join(desktopDir, "IT Toolbox.lnk")
		if _, err := os.Stat(lnkPath); err == nil {
			// Shortcut already exists in this desktop directory
			continue
		}

		// Try Method 1: PowerShell (Modern, direct, no temporary files)
		psScript := fmt.Sprintf(
			`$ws = New-Object -ComObject WScript.Shell; $s = $ws.CreateShortcut('%s'); $s.TargetPath = '%s'; $s.WorkingDirectory = '%s'; $s.Description = 'IT Toolbox - All-in-one Developer & Network Suite'; if (Test-Path '%s') { $s.IconLocation = '%s,0' }; $s.Save()`,
			strings.ReplaceAll(lnkPath, `'`, `''`),
			strings.ReplaceAll(exePath, `'`, `''`),
			strings.ReplaceAll(exeDir, `'`, `''`),
			strings.ReplaceAll(icoPath, `'`, `''`),
			strings.ReplaceAll(icoPath, `'`, `''`),
		)

		psCmd := exec.Command("powershell.exe", "-NoProfile", "-NonInteractive", "-ExecutionPolicy", "Bypass", "-Command", psScript)
		psCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		if err := psCmd.Run(); err == nil {
			if _, chk := os.Stat(lnkPath); chk == nil {
				continue // Successfully created via PowerShell
			}
		}

		// Try Method 2: Temporary VBScript (Fallback for legacy systems or if PS is restricted)
		vbsContent := fmt.Sprintf(`Set WshShell = WScript.CreateObject("WScript.Shell")
Set oShellLink = WshShell.CreateShortcut("%s")
oShellLink.TargetPath = "%s"
oShellLink.WorkingDirectory = "%s"
oShellLink.WindowStyle = 1
oShellLink.Description = "IT Toolbox - All-in-one Developer & Network Suite"
If CreateObject("Scripting.FileSystemObject").FileExists("%s") Then
    oShellLink.IconLocation = "%s, 0"
End If
oShellLink.Save
`, strings.ReplaceAll(lnkPath, `"`, `""`),
			strings.ReplaceAll(exePath, `"`, `""`),
			strings.ReplaceAll(exeDir, `"`, `""`),
			strings.ReplaceAll(icoPath, `"`, `""`),
			strings.ReplaceAll(icoPath, `"`, `""`))

		tempVBS := filepath.Join(os.TempDir(), "create_it_toolbox_shortcut.vbs")
		if err := os.WriteFile(tempVBS, []byte(vbsContent), 0600); err == nil {
			vbsCmd := exec.Command("wscript.exe", tempVBS)
			vbsCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			_ = vbsCmd.Run()
			_ = os.Remove(tempVBS)
		}
	}
}

// getWindowsDesktopDirs returns all candidate desktop directories,
// including user desktop, OneDrive Desktop, and Public Desktop.
func getWindowsDesktopDirs() []string {
	var dirs []string
	seen := make(map[string]bool)

	add := func(p string) {
		if p == "" {
			return
		}
		clean := filepath.Clean(p)
		if !seen[clean] {
			if info, err := os.Stat(clean); err == nil && info.IsDir() {
				seen[clean] = true
				dirs = append(dirs, clean)
			}
		}
	}

	// 1. User Profile Desktop
	if userProfile := os.Getenv("USERPROFILE"); userProfile != "" {
		add(filepath.Join(userProfile, "Desktop"))
		// OneDrive Desktop (common in Windows 10/11)
		add(filepath.Join(userProfile, "OneDrive", "Desktop"))
	}

	// 2. OneDrive commercial or localized Desktop
	if oneDrive := os.Getenv("OneDrive"); oneDrive != "" {
		add(filepath.Join(oneDrive, "Desktop"))
	}
	if oneDriveConsumer := os.Getenv("OneDriveConsumer"); oneDriveConsumer != "" {
		add(filepath.Join(oneDriveConsumer, "Desktop"))
	}

	// 3. UserHomeDir fallback
	if home, err := os.UserHomeDir(); err == nil {
		add(filepath.Join(home, "Desktop"))
		add(filepath.Join(home, "OneDrive", "Desktop"))
	}

	return dirs
}
