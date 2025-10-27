package bootstrap

import (
	"github.com/rivo/tview"
	"github.com/shangyanjin/gocmder/internal/config"
	"github.com/shangyanjin/gocmder/internal/logger"
	"github.com/shangyanjin/gocmder/internal/ui"
)

// Application represents the main application structure
type Application struct {
	Tview     *tview.Application
	Config    *config.Config
	dashboard *ui.Dashboard
	Logger    *logger.Logger
}

// New creates and initializes a new Application instance
func New() *Application {
	// Initialize logger with console output to logs directory
	lg, err := logger.NewLogger(true, "logs")
	if err != nil {
		panic(err)
	}

	// Initialize configuration
	cfg := config.NewConfig()

	app := &Application{
		Tview:  tview.NewApplication(),
		Config: cfg,
		Logger: lg,
	}

	app.Logger.Info("Application instance created")
	app.Logger.Info("Configuration initialized")
	return app
}

// Setup initializes the application components
func (a *Application) Setup() error {
	a.logInfo("Application setup started")

	// Create the main dashboard
	a.logInfo("Creating dashboard")
	a.dashboard = ui.NewDashboard(a.Tview)
	a.logInfo("Dashboard created successfully")

	// Set the root primitive
	a.Tview.SetRoot(a.dashboard.GetRoot(), true)

	a.logInfo("Application setup completed")
	return nil
}

// Run starts the application
func (a *Application) Run() error {
	a.logInfo("Application started")
	if err := a.Tview.Run(); err != nil {
		a.logError("Application error: %v", err)
		return err
	}
	return nil
}

// Stop gracefully stops the application
func (a *Application) Stop() {
	a.logInfo("Application stopping")
	a.Tview.Stop()
	if a.Logger != nil {
		a.Logger.Sync()
	}
}

// logInfo logs an info message
func (a *Application) logInfo(msg string, args ...interface{}) {
	if a.Logger != nil {
		a.Logger.Info(msg, args...)
	}
}

// logError logs an error message
func (a *Application) logError(msg string, args ...interface{}) {
	if a.Logger != nil {
		a.Logger.Error(msg, args...)
	}
}

// GetDashboard returns the main dashboard
func (a *Application) GetDashboard() *ui.Dashboard {
	return a.dashboard
}
