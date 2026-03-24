package main

import (
	"embed"
	"log"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// --- 1. INTERCEPT SILENT MODE ---
	// If the user runs "collector-gui.exe --silent"
	if len(os.Args) > 1 && os.Args[1] == "--silent" {
		runSilentMode()
		return // Exit the program, skipping the GUI entirely
	}

	// --- 2. NORMAL GUI LAUNCH ---
	app := NewApp()

	err := wails.Run(&options.App{
		Title:         "PC Tracker Agent",
		Width:         850,  // Widen to perfectly fit the two columns
		Height:        750,  // Shrink height so there is no dead space
		DisableResize: true, // Lock the window size
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 15, G: 23, B: 42, A: 255}, // Matches our Dark Mode Slate background
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

// --- 3. SILENT MODE LOGIC ---
func runSilentMode() {
	log.Println("Running in background mode...")

	// 1. Scan Hardware & Load Tags
	specs, _ := ScanHardware()
	config := LoadConfig()
	specs.Tag1 = config.Tag1
	specs.Tag2 = config.Tag2
	specs.Tag3 = config.Tag3

	// 2. Open Local Database
	db, err := InitDB()
	if err != nil {
		log.Fatalf("Failed to open local database: %v", err)
	}
	defer db.Close()

	// 3. Check for Changes
	currentHash := GenerateHash(specs)
	if HasHardwareChanged(db, specs.Serial, currentHash) {
		log.Println("Changes detected. Uploading to Google Sheets...")

		// 4. Upload to Cloud
		err := UploadToGoogleSheets(specs)
		if err != nil {
			log.Fatalf("Failed to upload: %v", err)
		}

		// (Make sure to update your local DB here so it doesn't upload again next time!)
		// e.g., UpdateDatabase(db, specs.Serial, currentHash)

		log.Println("Upload successful.")
	} else {
		log.Println("No hardware changes detected. Exiting.")
	}
}
