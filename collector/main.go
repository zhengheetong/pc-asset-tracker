package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// 1. Define the command-line flags
	// Use --silent when adding to Windows Startup
	silentMode := flag.Bool("silent", false, "Run the collector in background mode")
	flag.Parse()

	if *silentMode {
		runBackgroundMode()
	} else {
		runInteractiveMode()
	}
}

// runBackgroundMode handles the "Ghost" execution
func runBackgroundMode() {
	fmt.Println("Running in background...")

	// TODO:
	// 1. Check for Internet
	// 2. Scan Hardware
	// 3. Compare Hash with SQLite
	// 4. Upload to Google Sheets if changed

	os.Exit(0)
}

// runInteractiveMode opens the Management UI
func runInteractiveMode() {
	fmt.Println("Opening Management UI...")

	// TODO:
	// 1. Initialize GUI (Fyne/Wails)
	// 2. Load current config (User Name/Dept)
	// 3. Show Hardware Specs

	fmt.Println("Press Enter to exit the Management UI...")
	fmt.Scanln()
}
