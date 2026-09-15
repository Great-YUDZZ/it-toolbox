package shortcut

import (
	"log"
	"os"
	"path/filepath"
)

// EnsureDesktopShortcut checks if a desktop shortcut exists for the current executable,
// and automatically creates it if missing (for both Windows and Linux).
func EnsureDesktopShortcut() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Shortcut auto-creation recovered from panic: %v\n", r)
			}
		}()

		exePath, err := os.Executable()
		if err != nil {
			return
		}
		exePath, err = filepath.EvalSymlinks(exePath)
		if err != nil {
			return
		}

		createDesktopShortcut(exePath)
	}()
}
