package main

import (
	_ "embed"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"github.com/yudz/it-toolbox/database"
	"github.com/yudz/it-toolbox/ui"
)

//go:embed assets/icon.png
var appIconBytes []byte

func main() {
	// Initialize local SQLite database
	if _, err := database.InitDB(); err != nil {
		log.Printf("Warning: failed to initialize database: %v\n", err)
	}

	myApp := app.NewWithID("com.yudz.it-toolbox")
	myApp.Settings().SetTheme(ui.NewCustomCyanTheme())

	// Set app icon from embedded resource
	if len(appIconBytes) > 0 {
		myApp.SetIcon(fyne.NewStaticResource("icon.png", appIconBytes))
	}

	mainWindow := ui.NewMainWindow(myApp)
	mainWindow.ShowAndRun()
}
