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

	// 1. Ensure local icon exists in standard hicolor & pixmaps locations
	iconDir := filepath.Join(home, ".local", "share", "icons", "hicolor", "512x512", "apps")
	_ = os.MkdirAll(iconDir, 0755)
	iconPath := filepath.Join(iconDir, "it-toolbox.png")
	if len(embeddedIconPng) > 0 {
		_ = os.WriteFile(iconPath, embeddedIconPng, 0644)
	}

	pixmapsDir := filepath.Join(home, ".local", "share", "pixmaps")
	_ = os.MkdirAll(pixmapsDir, 0755)
	if len(embeddedIconPng) > 0 {
		_ = os.WriteFile(filepath.Join(pixmapsDir, "it-toolbox.png"), embeddedIconPng, 0644)
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

	// 3. Place in ~/.local/share/applications/ (Start Menu / Application Launcher)
	appMenuDir := filepath.Join(home, ".local", "share", "applications")
	_ = os.MkdirAll(appMenuDir, 0755)
	appMenuFile := filepath.Join(appMenuDir, "it-toolbox.desktop")
	_ = os.WriteFile(appMenuFile, []byte(desktopContent), 0755)

	// 4. Place in ~/Desktop & ~/desktop
	desktopDirs := []string{
		filepath.Join(home, "Desktop"),
		filepath.Join(home, "desktop"),
	}

	for _, d := range desktopDirs {
		if fi, err := os.Stat(d); err == nil && fi.IsDir() {
			desktopFile := filepath.Join(d, "it-toolbox.desktop")
			if err := os.WriteFile(desktopFile, []byte(desktopContent), 0755); err != nil {
				log.Printf("Failed to write desktop shortcut: %v\n", err)
			} else {
				// Mark trusted on GNOME if gio is available
				_ = exec.Command("gio", "set", desktopFile, "metadata::trusted", "true").Run()
			}
		}
	}

	// 5. Trigger desktop database update if command is available
	_ = exec.Command("update-desktop-database", appMenuDir).Run()
}
