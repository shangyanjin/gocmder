package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Level represents log level
type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

var levelNames = map[Level]string{
	DebugLevel: "DEBUG",
	InfoLevel:  "INFO",
	WarnLevel:  "WARN",
	ErrorLevel: "ERROR",
}

// Logger handles application logging with file and console output
type Logger struct {
	mu       sync.Mutex
	level    Level
	file     *os.File
	filePath string
	console  bool
	fileDir  string
}

var globalLogger *Logger
var once sync.Once

// Init initializes the global logger instance
func Init(enableConsole bool, logDir ...string) error {
	var err error
	once.Do(func() {
		globalLogger, err = NewLogger(enableConsole, logDir...)
	})
	return err
}

// NewLogger creates a new Logger instance
func NewLogger(enableConsole bool, logDir ...string) (*Logger, error) {
	lg := &Logger{
		level:   InfoLevel,
		console: enableConsole,
	}

	// Set log directory
	dir := "logs"
	if len(logDir) > 0 && logDir[0] != "" {
		dir = logDir[0]
	}
	lg.fileDir = dir

	// Create log directory if it doesn't exist
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// Open log file
	filePath := filepath.Join(dir, fmt.Sprintf("gocmder_%s.log", time.Now().Format("2006-01-02")))
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	lg.file = file
	lg.filePath = filePath

	return lg, nil
}

// SetLevel sets the log level
func (l *Logger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// log writes a log entry
func (l *Logger) log(level Level, msg string, args ...interface{}) {
	if l == nil {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if level < l.level {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	levelName := levelNames[level]
	content := fmt.Sprintf("[%s] %s: %s\n", timestamp, levelName, fmt.Sprintf(msg, args...))

	// Write to console
	if l.console {
		fmt.Fprint(os.Stdout, content)
	}

	// Write to file
	if l.file != nil {
		_, _ = l.file.WriteString(content)
	}
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, args ...interface{}) {
	l.log(DebugLevel, msg, args...)
}

// Info logs an info message
func (l *Logger) Info(msg string, args ...interface{}) {
	l.log(InfoLevel, msg, args...)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, args ...interface{}) {
	l.log(WarnLevel, msg, args...)
}

// Error logs an error message
func (l *Logger) Error(msg string, args ...interface{}) {
	l.log(ErrorLevel, msg, args...)
}

// Close closes the logger and releases resources
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// Sync flushes buffered data to disk
func (l *Logger) Sync() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file != nil {
		return l.file.Sync()
	}
	return nil
}

// GetLogFile returns the current log file path
func (l *Logger) GetLogFile() string {
	return l.filePath
}

// GetLogDir returns the log directory
func (l *Logger) GetLogDir() string {
	return l.fileDir
}

// Global logger functions

// SetLevel sets the global logger level
func SetLevel(level Level) {
	if globalLogger != nil {
		globalLogger.SetLevel(level)
	}
}

// Debug logs a debug message using global logger
func Debug(msg string, args ...interface{}) {
	if globalLogger != nil {
		globalLogger.Debug(msg, args...)
	}
}

// Info logs an info message using global logger
func Info(msg string, args ...interface{}) {
	if globalLogger != nil {
		globalLogger.Info(msg, args...)
	}
}

// Warn logs a warning message using global logger
func Warn(msg string, args ...interface{}) {
	if globalLogger != nil {
		globalLogger.Warn(msg, args...)
	}
}

// Error logs an error message using global logger
func Error(msg string, args ...interface{}) {
	if globalLogger != nil {
		globalLogger.Error(msg, args...)
	}
}

// Close closes the global logger
func Close() error {
	if globalLogger != nil {
		return globalLogger.Close()
	}
	return nil
}

// Sync flushes the global logger
func Sync() error {
	if globalLogger != nil {
		return globalLogger.Sync()
	}
	return nil
}

// GetLogFile returns the current log file path
func GetLogFile() string {
	if globalLogger != nil {
		return globalLogger.GetLogFile()
	}
	return ""
}

// GetLogDir returns the log directory
func GetLogDir() string {
	if globalLogger != nil {
		return globalLogger.GetLogDir()
	}
	return ""
}
