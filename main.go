package main

import (
	"github.com/shangyanjin/gocmder/internal/bootstrap"
)

func main() {
	// Create new application instance (includes logger initialization)
	app := bootstrap.New()
	defer app.Logger.Close()

	app.Logger.Info("Application startup")

	// Setup application components
	if err := app.Setup(); err != nil {
		app.Logger.Error("Failed to setup application: %v", err)
		return
	}

	// Run the application
	if err := app.Run(); err != nil {
		app.Logger.Error("Application exited with error: %v", err)
		return
	}

	app.Logger.Info("Application shutdown normally")
}
