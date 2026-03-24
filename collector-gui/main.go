package main

import (
	"embed"
	"log"
	"os"
	"path/filepath"

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
	// Setup File Logging
	//exePath, _ := os.Executable()
	logPath := filepath.Join("./", "tracker-debug.log")

	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(logFile)
		defer logFile.Close()
	}

	log.Println("==================================================")
	log.Println("Starting Silent Background Scan...")

	// 1. Scan Hardware
	specs, err := ScanHardware()
	if err != nil {
		log.Printf("ERROR: Hardware scan encountered an issue: %v", err)
	}
	if specs.OS == "" {
		specs.OS = getOSInfo()
	}

	// 2. Load Tags
	config := LoadConfig()
	specs.Tag1 = config.Tag1
	specs.Tag2 = config.Tag2
	specs.Tag3 = config.Tag3

	// --- NEW: PRINT ALL FOUND INFORMATION TO LOG ---
	log.Println("--- Hardware Data Found ---")
	log.Printf("Serial:     %s", specs.Serial)
	log.Printf("OS:         %s", specs.OS)
	log.Printf("CPU:        %s", specs.CPU)
	log.Printf("RAM Total:  %s", specs.RAMTotal)
	log.Printf("RAM Detail: %s", specs.RAMModules)
	log.Printf("Disks:      %s", specs.Disks)
	log.Printf("Dept Tag:   %s", specs.Tag1)
	log.Printf("Loc Tag:    %s", specs.Tag2)
	log.Printf("Type Tag:   %s", specs.Tag3)
	log.Println("---------------------------")

	// 3. Open Local Database
	db, err := InitDB()
	if err != nil {
		log.Fatalf("CRITICAL: Failed to open local database: %v", err)
	}
	defer db.Close()

	// 4. Check for Changes
	currentHash := GenerateHash(specs)
	log.Printf("Generated hardware hash: %s", currentHash)

	if HasHardwareChanged(db, specs.Serial, currentHash) {
		log.Println("Changes or new PC detected. Attempting upload to Google Sheets...")

		// 5. Upload to Cloud
		err = UploadToGoogleSheets(specs)
		if err != nil {
			log.Printf("ERROR: Google Sheets upload failed: %v", err)
			return
		}

		log.Println("SUCCESS: Data uploaded to Google Sheets!")

		// 6. Save the new hash locally
		err = UpdateLocalHash(db, specs.Serial, currentHash)
		if err != nil {
			log.Printf("WARNING: Failed to save hash to local DB: %v", err)
		} else {
			log.Println("Local database updated with new hash.")
		}

	} else {
		log.Println("No hardware changes detected. Skipping upload.")
	}

	log.Println("Scan complete. Shutting down.")
}
