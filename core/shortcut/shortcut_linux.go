//go:build !windows

package shortcut

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

//go:embed icon.png
var embeddedIconPng []byte

func createDesktopShortcut(exePath string) {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return
	}

	// 1. Ensure local icon exists
	iconDir := filepath.Join(home, ".local", "share", "icons", "hicolor", "512x512", "apps")
	_ = os.MkdirAll(iconDir, 0755)
	iconPath := filepath.Join(iconDir, "it-toolbox.png")
	if _, err := os.Stat(iconPath); os.IsNotExist(err) && len(embeddedIconPng) > 0 {
		_ = os.WriteFile(iconPath, embeddedIconPng, 0644)
	}

	// 2. Build .desktop content
	desktopContent := fmt.Sprintf(`[Desktop Entry]
Version=1.0
Type=Application
Name=IT Toolbox
GenericName=Network & Developer Suite
Comment=All-in-one Developer and Network Utility Suite
Exec="%s"
Icon=it-toolbox
Terminal=false
Categories=Network;Development;Utility;
Keywords=network;subnet;calculator;converter;cisco;
StartupNotify=true
StartupWMClass=it-toolbox
`, exePath)

	// 3. Place in ~/.local/share/applications/ (App Menu / Launcher)
	appMenuDir := filepath.Join(home, ".local", "share", "applications")
	_ = os.MkdirAll(appMenuDir, 0755)
	appMenuFile := filepath.Join(appMenuDir, "it-toolbox.desktop")
	if _, err := os.Stat(appMenuFile); os.IsNotExist(err) {
		_ = os.WriteFile(appMenuFile, []byte(desktopContent), 0755)
	}

	// 4. Place in ~/Desktop
	desktopDirs := []string{
		filepath.Join(home, "Desktop"),
		filepath.Join(home, "desktop"),
	}

	for _, d := range desktopDirs {
		if fi, err := os.Stat(d); err == nil && fi.IsDir() {
			desktopFile := filepath.Join(d, "it-toolbox.desktop")
			if _, err := os.Stat(desktopFile); os.IsNotExist(err) {
				if err := os.WriteFile(desktopFile, []byte(desktopContent), 0755); err != nil {
					log.Printf("Failed to write desktop shortcut: %v\n", err)
				} else {
					// Mark trusted on GNOME if gio is available
					_ = exec.Command("gio", "set", desktopFile, "metadata::trusted", "true").Run()
				}
			}
		}
	}
}
