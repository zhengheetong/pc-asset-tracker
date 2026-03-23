package main

import (
	"flag"
	"fmt"
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
	fmt.Println("=== Starting Background Asset Scan ===")

	// 1. Initialize Local Database
	db, err := InitDB()
	if err != nil {
		fmt.Printf("[!] Database Error: %v\n", err)
		return
	}
	defer db.Close()

	// 2. Scan Physical Hardware
	specs, err := ScanHardware()
	if err != nil {
		fmt.Printf("[!] Scan Error: %v\n", err)
		return
	}

	// 3. ASSIGN TAGS (Placeholder logic for now)
	// In the future, these could come from a local config.json file
	config := LoadConfig()
	specs.Tag1 = config.Tag1
	specs.Tag2 = config.Tag2
	specs.Tag3 = config.Tag3

	// 4. Debug Logs: Show what we found
	fmt.Println("--------------------------------------------------")
	fmt.Printf("SERIAL:  %s\n", specs.Serial)
	fmt.Printf("CPU:     %s\n", specs.CPU)
	fmt.Printf("RAM:     %s (%s)\n", specs.RAMTotal, specs.RAMModules)
	fmt.Printf("DISKS:   %s\n", specs.Disks)
	fmt.Printf("TAGS:    [%s] [%s] [%s]\n", specs.Tag1, specs.Tag2, specs.Tag3)
	fmt.Println("--------------------------------------------------")

	// 5. Generate System Fingerprint (Hash)
	currentHash := GenerateHash(specs)

	// 6. Check if anything (Hardware or Tags) has changed
	if HasHardwareChanged(db, specs.Serial, currentHash) {
		fmt.Println("[+] Change detected! Syncing with Google Sheets...")

		// Upload to Cloud
		err := UploadToGoogleSheets(specs)
		if err != nil {
			fmt.Printf("[!] Cloud Upload Failed: %v\n", err)
			return
		}

		// Update SQLite so we don't upload again until the next change
		err = UpdateLocalHash(db, specs.Serial, currentHash)
		if err != nil {
			fmt.Printf("[!] Local Cache Update Failed: %v\n", err)
		} else {
			fmt.Println("[*] SUCCESS: Cloud synced and local cache updated.")
		}
	} else {
		fmt.Println("[*] STATUS: System state unchanged. Skipping upload.")
	}

	fmt.Println("=== Scan Complete ===")
}

// runInteractiveMode opens the Management UI
func runInteractiveMode() {
	fmt.Println("Launching Configuration UI...")
	fmt.Println("Opening Management UI...")
	// For now, we just print a message since we are focusing on the scanner
	fmt.Println("Currently in development. Use --silent to run the scanner.")
	fmt.Println("Press Enter to exit...")
	fmt.Scanln()
}
