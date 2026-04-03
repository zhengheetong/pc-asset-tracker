package main

import (
	"collector-gui/internal/config"
	"collector-gui/internal/scanner"
	"context"
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

func (a *App) GetSpecs() scanner.PCSpecs {
	// 1. Get the hardware info
	specs, _ := scanner.ScanHardware()

	// 2. Load the current tags from config.json
	cfg := config.LoadConfig()

	// 3. Attach the tags to the specs struct
	specs.Tag1 = cfg.Tag1
	specs.Tag2 = cfg.Tag2
	specs.Tag3 = cfg.Tag3

	return specs
}

func (a *App) SaveConfig(tag1, tag2, tag3 string) string {
	cfg := config.AppConfig{
		Tag1: tag1,
		Tag2: tag2,
		Tag3: tag3,
	}

	err := config.SaveConfig(cfg)
	if err != nil {
		return "Failed to save config.json!"
	}

	return "Tags successfully saved to config.json!"
}

// CheckCredentials checks if the credentials file sits next to the executable
func (a *App) CheckCredentials() bool {
	_, err := os.Stat("service-account.json")
	return err == nil
}

func (a *App) InstallToPC() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("could not find executable: %v", err)
	}

	appData, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("could not find AppData folder: %v", err)
	}

	installDir := filepath.Join(appData, "PCTracker")
	os.MkdirAll(installDir, os.ModePerm)

	destExe := filepath.Join(installDir, filepath.Base(exePath))
	destJson := filepath.Join(installDir, "service-account.json")
	destConfig := filepath.Join(installDir, "config.json")

	// 1. Copy the Executable
	if err := copyFile(exePath, destExe); err != nil {
		return "", fmt.Errorf("failed to copy program: %v", err)
	}

	// 2. Copy the Google Credentials
	if err := copyFile("service-account.json", destJson); err != nil {
		return "", fmt.Errorf("failed to copy credentials: %v", err)
	}

	// 3. Copy the Saved Tags
	if _, err := os.Stat("config.json"); err == nil {
		copyFile("config.json", destConfig)
	}

	// 4. Set Registry for Startup
	runCommand := fmt.Sprintf("\"%s\" --silent", destExe)
	cmd := exec.Command("reg", "add", "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Run", "/v", "PCTracker", "/t", "REG_SZ", "/d", runCommand, "/f")

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to set startup key: %v", err)
	}

	return "Successfully installed! The tracker will now run silently in the background on startup.", nil
}

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
