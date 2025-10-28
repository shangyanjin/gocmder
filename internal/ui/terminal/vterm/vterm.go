package vterm

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shangyanjin/gocmder/internal/ui/style"
	"github.com/shangyanjin/gocmder/internal/ui/utils"
)

const (
	maxOutputLines = 1000
)

// VtermDialog is a virtual terminal dialog primitive
type VtermDialog struct {
	*tview.Box

	layout      *tview.Flex
	outputView  *tview.TextView
	inputField  *tview.InputField
	display     bool
	output      []string
	outputMutex sync.Mutex
	cmd         *exec.Cmd
	cancelFunc  func()
	workingDir  string
	shellPath   string
}

// NewVtermDialog returns a new virtual terminal dialog
func NewVtermDialog() *VtermDialog {
	bgColor := style.DialogBgColor

	vterm := &VtermDialog{
		Box:        tview.NewBox(),
		output:     make([]string, 0),
		workingDir: getDefaultWorkingDir(),
		shellPath:  getDefaultShell(),
	}

	// Create output view
	vterm.outputView = tview.NewTextView()
	vterm.outputView.SetBackgroundColor(bgColor)
	vterm.outputView.SetTextColor(style.DialogFgColor)
	vterm.outputView.SetDynamicColors(true)
	vterm.outputView.SetScrollable(true)
	vterm.outputView.SetBorder(true)
	vterm.outputView.SetBorderColor(style.DialogBorderColor)
	vterm.outputView.SetTitle(" Terminal Output ")

	// Create input field
	vterm.inputField = tview.NewInputField()
	vterm.inputField.SetBackgroundColor(bgColor)
	vterm.inputField.SetFieldBackgroundColor(style.BgColor)
	vterm.inputField.SetLabel("> ")
	vterm.inputField.SetLabelColor(style.StatusInstalledColor)
	vterm.inputField.SetFieldTextColor(style.FgColor)

	// Set input handler
	vterm.inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			command := vterm.inputField.GetText()
			if command != "" {
				vterm.executeCommand(command)
				vterm.inputField.SetText("")
			}
		}
	})

	// Create layout
	vterm.layout = tview.NewFlex().SetDirection(tview.FlexRow)
	vterm.layout.AddItem(vterm.outputView, 0, 1, false)
	vterm.layout.AddItem(vterm.inputField, 1, 0, true)
	vterm.layout.SetBackgroundColor(bgColor)

	vterm.addOutput(fmt.Sprintf("GoCmder Terminal - %s", runtime.GOOS))
	vterm.addOutput(fmt.Sprintf("Working Directory: %s", vterm.workingDir))
	vterm.addOutput(fmt.Sprintf("Shell: %s", vterm.shellPath))
	vterm.addOutput("Type 'help' for available commands, 'clear' to clear output")
	vterm.addOutput("")

	return vterm
}

// getDefaultWorkingDir returns the default working directory
func getDefaultWorkingDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return os.TempDir()
	}
	return dir
}

// getDefaultShell returns the default shell path
func getDefaultShell() string {
	if runtime.GOOS == "windows" {
		// Try PowerShell first, fallback to cmd
		if path, err := exec.LookPath("powershell.exe"); err == nil {
			return path
		}
		return "cmd.exe"
	}
	// Unix-like systems
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/sh"
	}
	return shell
}

// Display displays this primitive
func (v *VtermDialog) Display() {
	v.display = true
}

// IsDisplay returns true if primitive is shown
func (v *VtermDialog) IsDisplay() bool {
	return v.display
}

// Hide stops displaying this primitive
func (v *VtermDialog) Hide() {
	v.display = false
	if v.cmd != nil && v.cmd.Process != nil {
		v.cmd.Process.Kill()
	}
}

// HasFocus returns whether or not this primitive has focus
func (v *VtermDialog) HasFocus() bool {
	return v.inputField.HasFocus() || v.outputView.HasFocus()
}

// Focus is called when this primitive receives focus
func (v *VtermDialog) Focus(delegate func(p tview.Primitive)) {
	delegate(v.inputField)
}

// SetCancelFunc sets the cancel function
func (v *VtermDialog) SetCancelFunc(handler func()) {
	v.cancelFunc = handler
}

// addOutput adds a line to the output
func (v *VtermDialog) addOutput(line string) {
	v.outputMutex.Lock()
	defer v.outputMutex.Unlock()

	v.output = append(v.output, line)

	// Trim output if too long
	if len(v.output) > maxOutputLines {
		v.output = v.output[len(v.output)-maxOutputLines:]
	}

	v.updateOutputView()
}

// updateOutputView updates the output view with current output
func (v *VtermDialog) updateOutputView() {
	v.outputView.Clear()
	for _, line := range v.output {
		fmt.Fprintln(v.outputView, line)
	}
	v.outputView.ScrollToEnd()
}

// executeCommand executes a command
func (v *VtermDialog) executeCommand(command string) {
	v.addOutput(fmt.Sprintf("> %s", command))

	// Handle special commands
	switch strings.TrimSpace(strings.ToLower(command)) {
	case "clear":
		v.clearOutput()
		return
	case "help":
		v.showHelp()
		return
	case "exit", "quit":
		if v.cancelFunc != nil {
			v.cancelFunc()
		}
		return
	}

	// Handle cd command
	if strings.HasPrefix(command, "cd ") {
		newDir := strings.TrimSpace(strings.TrimPrefix(command, "cd "))
		if newDir == "" {
			newDir = os.Getenv("HOME")
			if newDir == "" {
				newDir = os.Getenv("USERPROFILE")
			}
		}
		if err := os.Chdir(newDir); err != nil {
			v.addOutput(fmt.Sprintf("Error: %v", err))
		} else {
			v.workingDir, _ = os.Getwd()
			v.addOutput(fmt.Sprintf("Changed directory to: %s", v.workingDir))
		}
		return
	}

	// Execute command
	go v.runCommand(command)
}

// runCommand runs a command and captures output
func (v *VtermDialog) runCommand(command string) {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		// Use cmd.exe /C for Windows
		cmd = exec.Command("cmd.exe", "/C", command)
	} else {
		// Use shell -c for Unix-like systems
		cmd = exec.Command("sh", "-c", command)
	}

	cmd.Dir = v.workingDir

	// Capture output
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		v.addOutput(fmt.Sprintf("Error creating stdout pipe: %v", err))
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		v.addOutput(fmt.Sprintf("Error creating stderr pipe: %v", err))
		return
	}

	// Start command
	if err := cmd.Start(); err != nil {
		v.addOutput(fmt.Sprintf("Error starting command: %v", err))
		return
	}

	v.cmd = cmd

	// Read output
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			v.addOutput(scanner.Text())
		}
	}()

	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			v.addOutput(fmt.Sprintf("[ERROR] %s", scanner.Text()))
		}
	}()

	wg.Wait()

	// Wait for command to finish
	if err := cmd.Wait(); err != nil {
		v.addOutput(fmt.Sprintf("Command finished with error: %v", err))
	} else {
		v.addOutput("Command completed successfully")
	}

	v.cmd = nil
}

// clearOutput clears the output
func (v *VtermDialog) clearOutput() {
	v.outputMutex.Lock()
	defer v.outputMutex.Unlock()

	v.output = make([]string, 0)
	v.updateOutputView()
}

// showHelp shows help information
func (v *VtermDialog) showHelp() {
	v.addOutput("")
	v.addOutput("Available Commands:")
	v.addOutput("  help      - Show this help message")
	v.addOutput("  clear     - Clear terminal output")
	v.addOutput("  exit/quit - Close terminal")
	v.addOutput("  cd <dir>  - Change working directory")
	v.addOutput("")
	v.addOutput("Any other command will be executed in the shell")
	v.addOutput(fmt.Sprintf("Current shell: %s", v.shellPath))
	v.addOutput("")
}

// InputHandler returns input handler function for this primitive
func (v *VtermDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return v.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		// ESC key to close
		if event.Key() == utils.CloseDialogKey.Key {
			if v.cancelFunc != nil {
				v.cancelFunc()
			}
			return
		}

		// Pass to input field
		if v.inputField.HasFocus() {
			if inputHandler := v.inputField.InputHandler(); inputHandler != nil {
				inputHandler(event, setFocus)
				return
			}
		}
	})
}

// SetRect sets rects for this primitive
func (v *VtermDialog) SetRect(x, y, width, height int) {
	v.Box.SetRect(x, y, width, height)

	x, y, width, height = v.GetInnerRect()
	v.layout.SetRect(x, y, width, height)
}

// Draw draws this primitive onto the screen
func (v *VtermDialog) Draw(screen tcell.Screen) {
	if !v.display {
		return
	}

	v.DrawForSubclass(screen, v)
	v.layout.Draw(screen)
}
