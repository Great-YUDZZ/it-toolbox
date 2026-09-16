package shortcut

import (
	"log"
	"os"
	"path/filepath"
)

// CreateShortcuts creates shortcuts selectively on Desktop, Start Menu, or both.
func CreateShortcuts(desktop, startMenu bool) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Shortcut creation recovered from panic: %v\n", r)
		}
	}()

	if !desktop && !startMenu {
		return
	}

	exePath, err := os.Executable()
	if err != nil {
		return
	}
	exePath, err = filepath.EvalSymlinks(exePath)
	if err != nil {
		return
	}

	createPlatformShortcuts(exePath, desktop, startMenu)
}

// EnsureDesktopShortcut creates shortcuts in both Desktop and Start Menu for backward compatibility.
func EnsureDesktopShortcut() {
	CreateShortcuts(true, true)
}
