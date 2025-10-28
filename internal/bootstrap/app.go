package bootstrap

import (
	"github.com/rivo/tview"
	"github.com/shangyanjin/gocmder/internal/config"
	"github.com/shangyanjin/gocmder/internal/logger"
	"github.com/shangyanjin/gocmder/internal/models"
	"github.com/shangyanjin/gocmder/internal/ui"
)

// Application represents the main application structure
type Application struct {
	Tview        *tview.Application
	Config       *config.Config
	UI           *ui.App
	Logger       *logger.Logger
	toolsData    []models.Tool
	settingsData []models.Setting
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

	// Initialize default data
	app.initializeData()

	return app
}

// initializeData initializes default tools and settings data
func (a *Application) initializeData() {
	// Initialize tools
	a.toolsData = []models.Tool{
		{Name: "Git", Version: "2.43.0", Size: "~50 MB", Selected: false, Installed: false},
		{Name: "VSCode", Version: "1.84.2", Size: "~100 MB", Selected: false, Installed: false},
		{Name: "Go", Version: "1.21.3", Size: "~130 MB", Selected: false, Installed: false},
		{Name: "Node.js", Version: "20.10.0", Size: "~40 MB", Selected: false, Installed: false},
		{Name: "PostgreSQL", Version: "16.0", Size: "~200 MB", Selected: false, Installed: false},
		{Name: "MySQL", Version: "8.1.0", Size: "~350 MB", Selected: false, Installed: false},
		{Name: "Redis", Version: "3.0.504", Size: "~5 MB", Selected: false, Installed: false},
	}

	// Initialize settings
	a.settingsData = []models.Setting{
		{Name: "Add to PATH", Selected: false},
		{Name: "Configure Power Settings", Selected: false},
		{Name: "Set User Directories", Selected: false},
	}
}

// Setup initializes the application components
func (a *Application) Setup() error {
	a.logInfo("Application setup started")

	// Create the main UI
	a.logInfo("Creating UI")
	a.UI = ui.NewApp(a.Tview)
	a.logInfo("UI created successfully")

	// Set handlers
	a.UI.SetInstallHandler(a.handleInstallTool)
	a.UI.SetApplySettingsHandler(a.handleApplySettings)
	a.UI.SetRefreshHandler(a.handleRefresh)

	// Update initial data
	a.UI.UpdateToolsData(a.toolsData)
	a.UI.UpdateSettingsData(a.settingsData)
	a.UI.RefreshSystemInfo()

	// Set the root primitive
	a.Tview.SetRoot(a.UI.GetRoot(), true)

	a.logInfo("Application setup completed")
	return nil
}

// handleInstallTool handles tool installation
func (a *Application) handleInstallTool(toolName string) {
	a.logInfo("Install requested for tool: %s", toolName)
	// TODO: Implement actual installation logic
	// For now, just log the request
}

// handleApplySettings handles settings application
func (a *Application) handleApplySettings(settings []models.Setting) {
	a.logInfo("Apply settings requested for %d settings", len(settings))
	// TODO: Implement actual settings application logic
	// For now, just log the request
	for _, setting := range settings {
		a.logInfo("  - %s", setting.Name)
	}
}

// handleRefresh handles data refresh
func (a *Application) handleRefresh() {
	a.logInfo("Refresh requested")
	// TODO: Implement actual refresh logic
	// For now, just update with current data
	a.UI.UpdateToolsData(a.toolsData)
	a.UI.UpdateSettingsData(a.settingsData)
	a.UI.RefreshSystemInfo()
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

// GetUI returns the main UI
func (a *Application) GetUI() *ui.App {
	return a.UI
}
