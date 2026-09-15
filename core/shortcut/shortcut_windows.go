//go:build windows

package shortcut

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//go:embed icon.ico
var embeddedIconIco []byte

func createDesktopShortcut(exePath string) {
	desktopDir := getWindowsDesktopDir()
	if desktopDir == "" {
		return
	}

	lnkPath := filepath.Join(desktopDir, "IT Toolbox.lnk")
	if _, err := os.Stat(lnkPath); err == nil {
		// Shortcut already exists
		return
	}

	exeDir := filepath.Dir(exePath)
	icoPath := filepath.Join(exeDir, "icon.ico")

	// If icon.ico doesn't exist in exeDir, extract embedded icon
	if _, err := os.Stat(icoPath); os.IsNotExist(err) && len(embeddedIconIco) > 0 {
		_ = os.WriteFile(icoPath, embeddedIconIco, 0644)
	}

	// Create temporary VBScript to safely create the .lnk shortcut
	vbsContent := fmt.Sprintf(`Set WshShell = WScript.CreateObject("WScript.Shell")
Set oShellLink = WshShell.CreateShortcut("%s")
oShellLink.TargetPath = "%s"
oShellLink.WorkingDirectory = "%s"
oShellLink.WindowStyle = 1
oShellLink.Description = "IT Toolbox - All-in-one Developer & Network Suite"
oShellLink.IconLocation = "%s, 0"
oShellLink.Save
`, strings.ReplaceAll(lnkPath, `"`, `""`),
		strings.ReplaceAll(exePath, `"`, `""`),
		strings.ReplaceAll(exeDir, `"`, `""`),
		strings.ReplaceAll(icoPath, `"`, `""`))

	tempVBS := filepath.Join(os.TempDir(), "create_it_toolbox_shortcut.vbs")
	if err := os.WriteFile(tempVBS, []byte(vbsContent), 0600); err != nil {
		log.Printf("Failed to write temporary VBScript for shortcut: %v\n", err)
		return
	}
	defer os.Remove(tempVBS)

	cmd := exec.Command("wscript.exe", tempVBS)
	if err := cmd.Run(); err != nil {
		log.Printf("Failed to execute VBScript for shortcut: %v\n", err)
	}
}

func getWindowsDesktopDir() string {
	userProfile := os.Getenv("USERPROFILE")
	if userProfile != "" {
		d := filepath.Join(userProfile, "Desktop")
		if _, err := os.Stat(d); err == nil {
			return d
		}
	}
	home, err := os.UserHomeDir()
	if err == nil {
		d := filepath.Join(home, "Desktop")
		if _, err := os.Stat(d); err == nil {
			return d
		}
	}
	return ""
}
