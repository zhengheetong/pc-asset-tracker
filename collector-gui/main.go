package main

import (
	"collector-gui/internal/api"
	"collector-gui/internal/config"
	"collector-gui/internal/database"
	"collector-gui/internal/scanner"
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
	if len(os.Args) > 1 && os.Args[1] == "--silent" {
		runSilentMode()
		return
	}

	// --- 2. NORMAL GUI LAUNCH ---
	app := NewApp()

	err := wails.Run(&options.App{
		Title:  "PC Tracker Agent",
		Width:  850,
		Height: 750,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 15, G: 23, B: 42, A: 255},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		log.Fatal("Error:", err)
	}
}

// --- 3. SILENT MODE LOGIC ---
func runSilentMode() {
	logPath := filepath.Join("./", "tracker-debug.log")

	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(logFile)
		defer logFile.Close()
	}

	log.Println("==================================================")
	log.Println("Starting Silent Background Scan...")

	// 1. Scan Hardware
	specs, err := scanner.ScanHardware()
	if err != nil {
		log.Printf("ERROR: Hardware scan encountered an issue: %v", err)
	}
	if specs.OS == "" {
		specs.OS = scanner.GetOSInfo()
	}

	// 2. Load Tags
	cfg := config.LoadConfig()
	specs.Tag1 = cfg.Tag1
	specs.Tag2 = cfg.Tag2
	specs.Tag3 = cfg.Tag3

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
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("CRITICAL: Failed to open local database: %v", err)
	}
	defer db.Close()

	// 4. Check for Changes
	currentHash := database.GenerateHash(specs)
	log.Printf("Generated hardware hash: %s", currentHash)

	if database.HasHardwareChanged(db, specs.Serial, currentHash) {
		log.Println("Changes or new PC detected. Attempting upload to Google Sheets...")

		// 5. Upload to Cloud
		err = api.UploadToGoogleSheets(specs)
		if err != nil {
			log.Printf("ERROR: Google Sheets upload failed: %v", err)
			return
		}

		log.Println("SUCCESS: Data uploaded to Google Sheets!")

		// 6. Save the new hash locally
		err = database.UpdateLocalHash(db, specs.Serial, currentHash)
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
