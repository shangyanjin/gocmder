package models

// NewInstallConfig creates a new installation configuration with default tools and settings
func NewInstallConfig() *InstallConfig {
	config := &InstallConfig{
		Tools: []Tool{
			{Name: "Git", Version: "2.43.0", Size: "190 MB", Selected: false, Installed: false},
			{Name: "VSCode", Version: "1.84.2", Size: "280 MB", Selected: false, Installed: false},
			{Name: "Go", Version: "1.21.3", Size: "240 MB", Selected: false, Installed: false},
			{Name: "Node.js", Version: "20.10.0", Size: "180 MB", Selected: false, Installed: false},
			{Name: "PostgreSQL", Version: "16.0", Size: "320 MB", Selected: false, Installed: false},
			{Name: "MySQL", Version: "8.1.0", Size: "350 MB", Selected: false, Installed: false},
			{Name: "Redis", Version: "3.0.504", Size: "80 MB", Selected: false, Installed: false},
		},
		Settings: []Setting{
			{Name: "Add PATH", Selected: false},
			{Name: "PowerConfig", Selected: false},
			{Name: "SetUserDirs", Selected: false},
		},
		CurrentScheme: -1,
	}

	// Initialize predefined schemes
	config.Schemes = []Scheme{
		{
			Name:           "Help",
			Description:    "Show help and keyboard shortcuts",
			ToolIndices:    []int{},
			SettingIndices: []int{},
		},
		{
			Name:           "Minimal",
			Description:    "Basic tools: Git, VSCode",
			ToolIndices:    []int{ToolGit, ToolVSCode},
			SettingIndices: []int{SettingAddPath},
		},
		{
			Name:           "Go Developer",
			Description:    "Go: Git, VSCode, Go",
			ToolIndices:    []int{ToolGit, ToolVSCode, ToolGo},
			SettingIndices: []int{SettingAddPath},
		},
		{
			Name:           "Node Developer",
			Description:    "Node: Git, VSCode, Node.js",
			ToolIndices:    []int{ToolGit, ToolVSCode, ToolNode},
			SettingIndices: []int{SettingAddPath},
		},
		{
			Name:           "Backend",
			Description:    "Backend: Go, PostgreSQL, MySQL, Redis",
			ToolIndices:    []int{ToolGit, ToolVSCode, ToolGo, ToolPostgres, ToolMySQL, ToolRedis},
			SettingIndices: []int{SettingAddPath},
		},
		{
			Name:           "Full Stack",
			Description:    "Full Stack: All tools",
			ToolIndices:    []int{ToolGit, ToolVSCode, ToolGo, ToolNode, ToolPostgres, ToolMySQL, ToolRedis},
			SettingIndices: []int{SettingAddPath, SettingSetUserDirs},
		},
		{
			Name:           "Custom",
			Description:    "Manual selection",
			ToolIndices:    []int{},
			SettingIndices: []int{},
		},
		{
			Name:           "Personal Settings",
			Description:    "Configure personal settings (Add PATH, PowerConfig, SetUserDirs)",
			ToolIndices:    []int{},
			SettingIndices: []int{SettingAddPath, SettingPowerConfig, SettingSetUserDirs},
		},
		{
			Name:           "Exit",
			Description:    "Exit application",
			ToolIndices:    []int{},
			SettingIndices: []int{},
		},
	}

	return config
}

// ApplyScheme applies a predefined scheme to the configuration
func (ic *InstallConfig) ApplyScheme(schemeIndex int) {
	if schemeIndex < 0 || schemeIndex >= len(ic.Schemes) {
		return
	}

	ic.CurrentScheme = schemeIndex
	scheme := ic.Schemes[schemeIndex]

	// Clear all selections
	ic.DeselectAllTools()
	ic.DeselectAllSettings()

	// Apply tool selections
	for _, toolIdx := range scheme.ToolIndices {
		if toolIdx >= 0 && toolIdx < len(ic.Tools) {
			ic.Tools[toolIdx].Selected = true
		}
	}

	// Apply setting selections
	for _, settingIdx := range scheme.SettingIndices {
		if settingIdx >= 0 && settingIdx < len(ic.Settings) {
			ic.Settings[settingIdx].Selected = true
		}
	}
}

// GetSelectedToolsCount returns the number of selected tools
func (ic *InstallConfig) GetSelectedToolsCount() int {
	count := 0
	for _, tool := range ic.Tools {
		if tool.Selected {
			count++
		}
	}
	return count
}

// GetSelectedSettingsCount returns the number of selected settings
func (ic *InstallConfig) GetSelectedSettingsCount() int {
	count := 0
	for _, setting := range ic.Settings {
		if setting.Selected {
			count++
		}
	}
	return count
}

// SelectAllTools selects all tools
func (ic *InstallConfig) SelectAllTools() {
	for i := range ic.Tools {
		ic.Tools[i].Selected = true
	}
}

// DeselectAllTools deselects all tools
func (ic *InstallConfig) DeselectAllTools() {
	for i := range ic.Tools {
		ic.Tools[i].Selected = false
	}
}

// SelectAllSettings selects all settings
func (ic *InstallConfig) SelectAllSettings() {
	for i := range ic.Settings {
		ic.Settings[i].Selected = true
	}
}

// DeselectAllSettings deselects all settings
func (ic *InstallConfig) DeselectAllSettings() {
	for i := range ic.Settings {
		ic.Settings[i].Selected = false
	}
}
