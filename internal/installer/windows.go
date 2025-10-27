package installer

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// WindowsInstaller handles installation on Windows
type WindowsInstaller struct {
	LocalSourcePath string
	Tools           map[string]*ToolInfo
}

// ToolInfo contains information about a tool
type ToolInfo struct {
	Name        string
	FileName    string
	DownloadURL string
	InstallCmd  string
}

// NewWindowsInstaller creates a new Windows installer
func NewWindowsInstaller(localSourcePath string) *WindowsInstaller {
	return &WindowsInstaller{
		LocalSourcePath: localSourcePath,
		Tools: map[string]*ToolInfo{
			"Git": {
				Name:        "Git",
				FileName:    "Git-2.43.0-64-bit.exe",
				DownloadURL: "https://github.com/git-for-windows/git/releases/download/v2.43.0.windows.1/Git-2.43.0-64-bit.exe",
				InstallCmd:  "/VERYSILENT /NORESTART",
			},
			"VSCode": {
				Name:        "VSCode",
				FileName:    "VSCodeSetup-x64-1.84.2.exe",
				DownloadURL: "https://aka.ms/win32-x64-user-stable",
				InstallCmd:  "/VERYSILENT /NORESTART",
			},
			"Go": {
				Name:        "Go",
				FileName:    "go1.21.3.windows-amd64.msi",
				DownloadURL: "https://go.dev/dl/go1.21.3.windows-amd64.msi",
				InstallCmd:  "/quiet /norestart",
			},
			"Node.js": {
				Name:        "Node.js",
				FileName:    "node-v20.10.0-x64.msi",
				DownloadURL: "https://nodejs.org/dist/v20.10.0/node-v20.10.0-x64.msi",
				InstallCmd:  "/quiet /norestart",
			},
			"PostgreSQL": {
				Name:        "PostgreSQL",
				FileName:    "postgresql-16.0-1-windows-x64.exe",
				DownloadURL: "https://get.enterprisedb.com/postgresql/postgresql-16.0-1-windows-x64.exe",
				InstallCmd:  "--unattendedmodeui minimal --mode unattended --superpassword postgres",
			},
			"MySQL": {
				Name:        "MySQL",
				FileName:    "mysql-8.1.0-winx64.msi",
				DownloadURL: "https://dev.mysql.com/get/Downloads/MySQLInstaller/mysql-installer-community-8.1.0.0.msi",
				InstallCmd:  "/quiet /norestart",
			},
			"Redis": {
				Name:        "Redis",
				FileName:    "Redis-x64-3.0.504.msi",
				DownloadURL: "https://github.com/microsoftarchive/redis/releases/download/win-3.0.504/Redis-x64-3.0.504.msi",
				InstallCmd:  "/quiet /norestart",
			},
		},
	}
}

// GetInstallFile gets the installer file, downloading if necessary
func (wi *WindowsInstaller) GetInstallFile(tool *ToolInfo) (string, error) {
	localFile := filepath.Join(wi.LocalSourcePath, tool.FileName)

	if _, err := os.Stat(localFile); err == nil {
		return localFile, nil
	}

	if err := os.MkdirAll(wi.LocalSourcePath, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	if err := wi.downloadFile(tool.DownloadURL, localFile); err != nil {
		return "", fmt.Errorf("failed to download %s: %w", tool.Name, err)
	}

	return localFile, nil
}

// downloadFile downloads a file from URL
func (wi *WindowsInstaller) downloadFile(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// InstallTool installs a tool
func (wi *WindowsInstaller) InstallTool(toolName string) error {
	tool, exists := wi.Tools[toolName]
	if !exists {
		return fmt.Errorf("tool %s not found", toolName)
	}

	installerPath, err := wi.GetInstallFile(tool)
	if err != nil {
		return err
	}

	var cmd *exec.Cmd
	if strings.HasSuffix(installerPath, ".msi") {
		args := append([]string{"/i", installerPath}, strings.Fields(tool.InstallCmd)...)
		cmd = exec.Command("msiexec.exe", args...)
	} else {
		args := strings.Fields(tool.InstallCmd)
		cmd = exec.Command(installerPath, args...)
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("installation failed: %w", err)
	}

	return nil
}

// AddSystemPaths adds paths to system PATH
func (wi *WindowsInstaller) AddSystemPaths(paths []string) error {
	envPath := os.Getenv("PATH")

	for _, path := range paths {
		if !strings.Contains(envPath, path) {
			envPath += ";" + path
		}
	}

	cmd := exec.Command("powershell", "-Command",
		fmt.Sprintf(`[System.Environment]::SetEnvironmentVariable("Path", "%s", "Machine")`, envPath))

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to set PATH: %w", err)
	}

	return nil
}

// ConfigurePowerSettings configures power settings
func (wi *WindowsInstaller) ConfigurePowerSettings() error {
	settings := []struct {
		cmd  string
		args []string
	}{
		{"powercfg", []string{"/change", "monitor-timeout-dc", "180"}},
		{"powercfg", []string{"/change", "monitor-timeout-ac", "180"}},
		{"powercfg", []string{"/change", "standby-timeout-dc", "0"}},
		{"powercfg", []string{"/change", "standby-timeout-ac", "0"}},
	}

	for _, setting := range settings {
		cmd := exec.Command(setting.cmd, setting.args...)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to configure power settings: %w", err)
		}
	}

	return nil
}
