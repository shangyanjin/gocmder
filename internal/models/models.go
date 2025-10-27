package models

// Tool represents a development tool that can be installed
type Tool struct {
	Name      string
	Version   string
	Size      string
	Selected  bool
	Installed bool
}

// Setting represents a system configuration setting
type Setting struct {
	Name     string
	Selected bool
}

// Scheme represents a predefined installation scheme
type Scheme struct {
	Name           string
	Description    string
	ToolIndices    []int // Indices of tools to select
	SettingIndices []int // Indices of settings to select
}

// InstallConfig holds the installation configuration
type InstallConfig struct {
	Tools         []Tool
	Settings      []Setting
	Schemes       []Scheme
	CurrentScheme int // Index of currently selected scheme (-1 for none)
}

