package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

// App struct
type App struct {
	ctx context.Context
}

type Config struct {
	Tag1 string `json:"tag1"`
	Tag2 string `json:"tag2"`
	Tag3 string `json:"tag3"`
}

func getConfigPath() string {
	exePath, _ := os.Executable()
	return filepath.Join(filepath.Dir(exePath), "config.json")
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetSpecs() PCSpecs {
	// 1. Get the hardware info
	specs, _ := ScanHardware()

	// 2. Load the current tags from config.json
	config := LoadConfig()

	// 3. Attach the tags to the specs struct
	specs.Tag1 = config.Tag1
	specs.Tag2 = config.Tag2
	specs.Tag3 = config.Tag3

	return specs
}

func (a *App) SaveConfig(tag1, tag2, tag3 string) string {
	cfg := Config{
		Tag1: tag1,
		Tag2: tag2,
		Tag3: tag3,
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return "Failed to encode config: " + err.Error()
	}

	err = os.WriteFile(getConfigPath(), data, 0644)
	if err != nil {
		return "Failed to save config.json!"
	}

	return "Tags successfully saved to config.json!"
}

// 1. Check if the credentials file sits next to the executable
func (a *App) CheckCredentials() bool {
	_, err := os.Stat("service-account.json")
	return err == nil
}

// 2. The main Installation logic
func (a *App) InstallToPC() (string, error) {
	// Find where the user is running the current .exe from
	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("could not find executable: %v", err)
	}

	// Define the hidden install location: C:\Users\Username\AppData\Local\PCTracker
	appData, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("could not find AppData folder: %v", err)
	}

	installDir := filepath.Join(appData, "PCTracker")
	os.MkdirAll(installDir, os.ModePerm)

	// Define destination paths
	destExe := filepath.Join(installDir, filepath.Base(exePath))
	destJson := filepath.Join(installDir, "service-account.json")

	// Copy the Executable
	if err := copyFile(exePath, destExe); err != nil {
		return "", fmt.Errorf("failed to copy program: %v", err)
	}

	// Copy the JSON Credentials
	if err := copyFile("service-account.json", destJson); err != nil {
		return "", fmt.Errorf("failed to copy credentials: %v", err)
	}

	// Tell Windows to run this silently on startup via the Registry
	runCommand := fmt.Sprintf("\"%s\" --silent", destExe)
	cmd := exec.Command("reg", "add", "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Run", "/v", "PCTracker", "/t", "REG_SZ", "/d", runCommand, "/f")

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to set startup key: %v", err)
	}

	return "Successfully installed! The tracker will now run silently in the background on startup.", nil
}

// 3. A simple helper function to copy files
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
