package main

import (
	"context"
	"fmt"
	"os"
)

// App struct
type App struct {
	ctx context.Context
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

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
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
	// 1. Prepare the config struct
	newConfig := AppConfig{
		Tag1: tag1,
		Tag2: tag2,
		Tag3: tag3,
	}

	// 2. Call your existing SaveConfig function from config.go
	err := SaveConfig(newConfig)

	// 3. Return a status message
	if err != nil {
		return "Error: " + err.Error()
	}
	return "Configuration saved successfully!"
}

func (a *App) CheckCredentials() bool {
	_, err := os.Stat("service-account.json")
	return err == nil // Returns true if file exists, false otherwise
}
