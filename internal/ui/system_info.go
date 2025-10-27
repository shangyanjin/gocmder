package ui

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
	"github.com/shirou/gopsutil/v4/process"
)

// SystemInfoCache holds cached system information with thread safety
type SystemInfoCache struct {
	mu               sync.RWMutex
	cachedInfo       string
	lastUpdate       time.Time
	updateInterval   time.Duration
	updateInProgress bool
}

// Global cache instance - default no refresh
var systemInfoCache = &SystemInfoCache{
	updateInterval: 0, // Default: no automatic refresh
}

// SystemInfo holds system information
type SystemInfo struct {
	Hostname     string
	OS           string
	Platform     string
	Architecture string
	Kernel       string
	Uptime       uint64

	CPUCount int
	CPUModel string
	CPUUsage float64

	MemoryTotal uint64
	MemoryUsed  uint64
	MemoryFree  uint64
	MemoryUsage float64

	DiskTotal uint64
	DiskUsed  uint64
	DiskFree  uint64
	DiskUsage float64

	LoadAvg1  float64
	LoadAvg5  float64
	LoadAvg15 float64

	NetworkRx uint64
	NetworkTx uint64

	ProcessCount int
	GoVersion    string
	LastUpdate   time.Time
}

// GetSystemInfo collects comprehensive system information
func GetSystemInfo() (*SystemInfo, error) {
	info := &SystemInfo{
		LastUpdate: time.Now(),
		GoVersion:  runtime.Version(),
	}

	// Host information
	if hostInfo, err := host.Info(); err == nil {
		info.Hostname = hostInfo.Hostname
		info.OS = hostInfo.Platform
		info.Platform = hostInfo.Platform
		info.Architecture = hostInfo.KernelArch
		info.Kernel = hostInfo.KernelVersion
		info.Uptime = hostInfo.Uptime
	}

	// CPU information
	if cpuInfo, err := cpu.Info(); err == nil && len(cpuInfo) > 0 {
		info.CPUCount = len(cpuInfo)
		info.CPUModel = cpuInfo[0].ModelName
	}

	if cpuUsage, err := cpu.Percent(time.Second, false); err == nil && len(cpuUsage) > 0 {
		info.CPUUsage = cpuUsage[0]
	}

	// Memory information
	if memInfo, err := mem.VirtualMemory(); err == nil {
		info.MemoryTotal = memInfo.Total
		info.MemoryUsed = memInfo.Used
		info.MemoryFree = memInfo.Free
		info.MemoryUsage = memInfo.UsedPercent
	}

	// Disk information
	if diskInfo, err := disk.Usage("/"); err == nil {
		info.DiskTotal = diskInfo.Total
		info.DiskUsed = diskInfo.Used
		info.DiskFree = diskInfo.Free
		info.DiskUsage = diskInfo.UsedPercent
	}

	// Load average (Unix-like systems)
	if loadAvg, err := load.Avg(); err == nil {
		info.LoadAvg1 = loadAvg.Load1
		info.LoadAvg5 = loadAvg.Load5
		info.LoadAvg15 = loadAvg.Load15
	}

	// Network information
	if netStats, err := net.IOCounters(false); err == nil && len(netStats) > 0 {
		info.NetworkRx = netStats[0].BytesRecv
		info.NetworkTx = netStats[0].BytesSent
	}

	// Process count
	if processes, err := process.Processes(); err == nil {
		info.ProcessCount = len(processes)
	}

	return info, nil
}

// FormatBytes formats bytes to human readable format
func FormatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// FormatUptime formats uptime to human readable format
func FormatUptime(seconds uint64) string {
	days := seconds / 86400
	hours := (seconds % 86400) / 3600
	minutes := (seconds % 3600) / 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	} else if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	} else {
		return fmt.Sprintf("%dm", minutes)
	}
}

// GetSystemInfoText returns formatted system information as text
func GetSystemInfoText() string {
	info, err := GetSystemInfo()
	if err != nil {
		return fmt.Sprintf("Error getting system info: %v", err)
	}

	// Create a safe, truncated version
	text := fmt.Sprintf(`System Information:
Hostname: %s
OS: %s %s
Kernel: %s
Uptime: %s
Go Version: %s

CPU:
Model: %s
Cores: %d
Usage: %.1f%%

Memory:
Total: %s
Used: %s (%.1f%%)
Free: %s

Disk:
Total: %s
Used: %s (%.1f%%)
Free: %s

Load Average:
1m: %.2f, 5m: %.2f, 15m: %.2f

Network:
RX: %s, TX: %s

Processes: %d
Last Update: %s`,
		truncateString(info.Hostname, 20),
		truncateString(info.OS, 15),
		truncateString(info.Platform, 10),
		truncateString(info.Kernel, 30),
		FormatUptime(info.Uptime),
		info.GoVersion,

		truncateString(info.CPUModel, 30),
		info.CPUCount,
		info.CPUUsage,

		FormatBytes(info.MemoryTotal),
		FormatBytes(info.MemoryUsed),
		info.MemoryUsage,
		FormatBytes(info.MemoryFree),

		FormatBytes(info.DiskTotal),
		FormatBytes(info.DiskUsed),
		info.DiskUsage,
		FormatBytes(info.DiskFree),

		info.LoadAvg1, info.LoadAvg5, info.LoadAvg15,

		FormatBytes(info.NetworkRx),
		FormatBytes(info.NetworkTx),

		info.ProcessCount,
		info.LastUpdate.Format("15:04:05"),
	)

	return text
}

// truncateString safely truncates a string to max length
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

// InitializeSystemInfoCache initializes the cache with initial data
func InitializeSystemInfoCache() {
	// Initialize with basic info
	info := fetchSystemInfoCompact()
	systemInfoCache.mu.Lock()
	systemInfoCache.cachedInfo = info
	systemInfoCache.lastUpdate = time.Now()
	systemInfoCache.mu.Unlock()
}

// GetSystemInfoCompact returns cached system information (thread-safe)
func GetSystemInfoCompact() string {
	systemInfoCache.mu.RLock()
	defer systemInfoCache.mu.RUnlock()

	// Return cached data (default: no automatic refresh)
	return systemInfoCache.cachedInfo
}

// UpdateSystemInfoAsync triggers async update of system information
func UpdateSystemInfoAsync() {
	// Check if update is already in progress
	systemInfoCache.mu.Lock()
	if systemInfoCache.updateInProgress {
		systemInfoCache.mu.Unlock()
		return
	}
	systemInfoCache.updateInProgress = true
	systemInfoCache.mu.Unlock()

	// Perform update in background goroutine
	go func() {
		info := fetchSystemInfoCompact()

		// Update cache with new data
		systemInfoCache.mu.Lock()
		systemInfoCache.cachedInfo = info
		systemInfoCache.lastUpdate = time.Now()
		systemInfoCache.updateInProgress = false
		systemInfoCache.mu.Unlock()
	}()
}

// fetchSystemInfoCompact fetches system information (potentially slow operations)
func fetchSystemInfoCompact() string {
	var info string

	// Add title/header
	info = "[GoCmder - Auto install DevTools and setup environment]\n"

	// Get host information (fast operation)
	if hostInfo, err := host.Info(); err == nil {
		info += fmt.Sprintf("OS: %s %s | Hostname: %s\n", hostInfo.OS, hostInfo.PlatformVersion, hostInfo.Hostname)
	}

	// Get CPU count (fast operation)
	if cpuCount, err := cpu.Counts(true); err == nil {
		info += fmt.Sprintf("CPU: %d cores ", cpuCount)
	}

	// Get load average (fast operation)
	if loadAvg, err := load.Avg(); err == nil {
		info += fmt.Sprintf("| Load: %.2f,%.2f,%.2f\n", loadAvg.Load1, loadAvg.Load5, loadAvg.Load15)
	}

	// Get memory information (fast operation)
	if memInfo, err := mem.VirtualMemory(); err == nil {
		info += fmt.Sprintf("Memory: %.2fGB/%.2fGB (%.2f%%) ",
			float64(memInfo.Used)/1024/1024/1024,
			float64(memInfo.Total)/1024/1024/1024,
			memInfo.UsedPercent)
	}

	// Get disk information (potentially slow, but necessary)
	if diskInfo, err := disk.Usage("/"); err == nil {
		info += fmt.Sprintf("| Disk: %.2fGB/%.2fGB (%.2f%%)\n",
			float64(diskInfo.Used)/1024/1024/1024,
			float64(diskInfo.Total)/1024/1024/1024,
			diskInfo.UsedPercent)
	}

	// Get network interfaces (fast operation)
	if netInterfaces, err := net.Interfaces(); err == nil && len(netInterfaces) > 0 {
		// Find non-loopback interface with IP
		for _, iface := range netInterfaces {
			if len(iface.Addrs) > 0 && iface.Name != "lo" {
				info += fmt.Sprintf("IP: %s | Interface: %s\n", iface.Addrs[0], iface.Name)
				break
			}
		}
	}

	// Get shortcut hints
	info += "[F2] Files  [F5] Input  [F9] Output  [Esc] Quit"

	return info
}
