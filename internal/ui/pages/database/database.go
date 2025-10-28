package database

import (
	"fmt"
	"strings"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shangyanjin/gocmder/internal/db"
	"github.com/shangyanjin/gocmder/internal/ui/components/dialogs"
	"github.com/shangyanjin/gocmder/internal/ui/style"
)

// Database implements the database management page primitive
type Database struct {
	*tview.Box

	title           string
	mainFlex        *tview.Flex
	leftPanel       *tview.TreeView
	rightPanel      *tview.Flex
	sqlEditor       *tview.TextArea
	resultTable     *tview.Table
	statusBar       *tview.TextView
	errorDialog     *dialogs.ErrorDialog
	messageDialog   *dialogs.MessageDialog
	connDialog      *ConnectionDialog
	driver          db.Driver
	connected       bool
	currentDatabase string
	mu              sync.Mutex
	focusedElement  int // 0=tree, 1=editor, 2=result
	appFocusHandler func()
}

const (
	focusTree = iota
	focusEditor
	focusResult
)

// NewDatabase returns database page view
func NewDatabase() *Database {
	database := &Database{
		Box:            tview.NewBox(),
		title:          "database",
		errorDialog:    dialogs.NewErrorDialog(),
		messageDialog:  dialogs.NewMessageDialog(""),
		connected:      false,
		focusedElement: focusTree,
	}

	// Create tree view for databases and tables
	database.leftPanel = tview.NewTreeView()
	database.leftPanel.SetBorder(true)
	database.leftPanel.SetTitle(" Database Tree ")
	database.leftPanel.SetTitleColor(style.FgColor)
	database.leftPanel.SetBorderColor(style.BorderColor)
	database.leftPanel.SetBackgroundColor(style.BgColor)
	database.leftPanel.SetGraphicsColor(style.StatusInstalledColor)

	rootNode := tview.NewTreeNode("Not Connected")
	rootNode.SetColor(style.StatusNotInstalledColor)
	database.leftPanel.SetRoot(rootNode)
	database.leftPanel.SetCurrentNode(rootNode)

	// Create SQL editor
	database.sqlEditor = tview.NewTextArea()
	database.sqlEditor.SetBorder(true)
	database.sqlEditor.SetTitle(" SQL Editor (Ctrl+R: Execute) ")
	database.sqlEditor.SetTitleColor(style.FgColor)
	database.sqlEditor.SetBorderColor(style.BorderColor)
	database.sqlEditor.SetBackgroundColor(style.DialogBgColor)
	database.sqlEditor.SetPlaceholder("Enter SQL query here...")
	database.sqlEditor.SetTextStyle(tcell.StyleDefault.
		Foreground(style.FgColor).
		Background(style.DialogBgColor))

	// Create result table
	database.resultTable = tview.NewTable()
	database.resultTable.SetBorder(true)
	database.resultTable.SetTitle(" Query Results ")
	database.resultTable.SetTitleColor(style.FgColor)
	database.resultTable.SetBorderColor(style.BorderColor)
	database.resultTable.SetBackgroundColor(style.BgColor)
	database.resultTable.SetSelectable(true, false)
	database.resultTable.SetFixed(1, 0)

	// Create status bar
	database.statusBar = tview.NewTextView()
	database.statusBar.SetBackgroundColor(style.InfoBarBgColor)
	database.statusBar.SetTextColor(style.InfoBarFgColor)
	database.statusBar.SetDynamicColors(true)
	database.updateStatusBar("Ready. Press Ctrl+N to connect")

	// Create right panel layout
	database.rightPanel = tview.NewFlex().SetDirection(tview.FlexRow)
	database.rightPanel.AddItem(database.sqlEditor, 0, 1, true)
	database.rightPanel.AddItem(database.resultTable, 0, 2, false)

	// Create main layout
	database.mainFlex = tview.NewFlex().SetDirection(tview.FlexRow)

	contentFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	contentFlex.AddItem(database.leftPanel, 0, 1, true)
	contentFlex.AddItem(database.rightPanel, 0, 3, false)

	database.mainFlex.AddItem(contentFlex, 0, 1, true)
	database.mainFlex.AddItem(database.statusBar, 1, 0, false)

	// Create connection dialog
	database.connDialog = NewConnectionDialog(database.handleConnect)

	// Set dialog handlers with focus restoration
	database.errorDialog.SetDoneFunc(func() {
		database.errorDialog.Hide()
		if database.appFocusHandler != nil {
			database.appFocusHandler()
		}
	})
	database.messageDialog.SetCancelFunc(func() {
		database.messageDialog.Hide()
		if database.appFocusHandler != nil {
			database.appFocusHandler()
		}
	})

	// Set connection dialog app focus handler to restore focus after closing
	database.connDialog.SetAppFocusHandler(func() {
		if database.appFocusHandler != nil {
			database.appFocusHandler()
		}
	})

	// Set tree selection handler
	database.leftPanel.SetSelectedFunc(database.handleTreeSelection)

	return database
}

// GetTitle returns primitive title
func (d *Database) GetTitle() string {
	return d.title
}

// HasFocus returns whether or not this primitive has focus
func (d *Database) HasFocus() bool {
	return d.mainFlex.HasFocus() || d.errorDialog.HasFocus() ||
		d.messageDialog.HasFocus() || d.connDialog.HasFocus() || d.Box.HasFocus()
}

// Focus is called when this primitive receives focus
func (d *Database) Focus(delegate func(p tview.Primitive)) {
	if d.errorDialog.IsDisplay() {
		delegate(d.errorDialog)
		return
	}
	if d.messageDialog.IsDisplay() {
		delegate(d.messageDialog)
		return
	}
	if d.connDialog.IsDisplay() {
		delegate(d.connDialog)
		return
	}

	// Focus based on current element
	switch d.focusedElement {
	case focusTree:
		delegate(d.leftPanel)
	case focusEditor:
		delegate(d.sqlEditor)
	case focusResult:
		delegate(d.resultTable)
	default:
		delegate(d.leftPanel)
	}
}

// SetAppFocusHandler sets application focus handler
func (d *Database) SetAppFocusHandler(handler func()) {
	d.appFocusHandler = handler
}

// HideAllDialogs hides all sub dialogs
func (d *Database) HideAllDialogs() {
	if d.errorDialog.IsDisplay() {
		d.errorDialog.Hide()
	}
	if d.messageDialog.IsDisplay() {
		d.messageDialog.Hide()
	}
	if d.connDialog.IsDisplay() {
		d.connDialog.Hide()
	}
}

// SubDialogHasFocus returns whether or not sub dialog primitive has focus
func (d *Database) SubDialogHasFocus() bool {
	return d.errorDialog.HasFocus() || d.messageDialog.HasFocus() || d.connDialog.HasFocus()
}

// updateStatusBar updates the status bar
func (d *Database) updateStatusBar(message string) {
	highlightColor := style.GetColorHex(style.StatusInstalledColor)
	d.statusBar.SetText(fmt.Sprintf(" [%s]Status:[-] %s", highlightColor, message))
}

// handleConnect handles database connection
func (d *Database) handleConnect(driverType, dsn string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Close existing connection
	if d.driver != nil {
		d.driver.Close()
	}

	// Create new driver
	var driver db.Driver
	switch driverType {
	case "MySQL":
		driver = db.NewMySQL()
	case "PostgreSQL":
		driver = db.NewPostgres()
	default:
		d.showError("Invalid driver type: " + driverType)
		return
	}

	// Connect
	err := driver.Connect(dsn)
	if err != nil {
		d.showError(fmt.Sprintf("Connection failed: %v", err))
		return
	}

	d.driver = driver
	d.connected = true
	d.updateStatusBar(fmt.Sprintf("Connected to %s", driver.GetDriverName()))

	// Load databases
	d.loadDatabases()
}

// loadDatabases loads database list into tree
func (d *Database) loadDatabases() {
	if !d.connected || d.driver == nil {
		return
	}

	databases, err := d.driver.GetDatabases()
	if err != nil {
		d.showError(fmt.Sprintf("Failed to load databases: %v", err))
		return
	}

	// Build tree
	rootNode := tview.NewTreeNode(d.driver.GetDriverName())
	rootNode.SetColor(style.StatusInstalledColor)
	rootNode.SetExpanded(true)

	for _, dbName := range databases {
		dbNode := tview.NewTreeNode(dbName)
		dbNode.SetColor(style.FgColor)
		dbNode.SetReference(map[string]string{"type": "database", "name": dbName})
		dbNode.SetSelectable(true)
		rootNode.AddChild(dbNode)
	}

	d.leftPanel.SetRoot(rootNode)
	d.leftPanel.SetCurrentNode(rootNode)
}

// handleTreeSelection handles tree node selection
func (d *Database) handleTreeSelection(node *tview.TreeNode) {
	ref := node.GetReference()
	if ref == nil {
		return
	}

	data, ok := ref.(map[string]string)
	if !ok {
		return
	}

	switch data["type"] {
	case "database":
		d.loadTables(data["name"])
	case "table":
		// Could show table structure here
		d.updateStatusBar(fmt.Sprintf("Selected table: %s.%s", data["database"], data["name"]))
	}
}

// loadTables loads tables for a database
func (d *Database) loadTables(dbName string) {
	if !d.connected || d.driver == nil {
		return
	}

	d.currentDatabase = dbName
	d.updateStatusBar(fmt.Sprintf("Loading tables from %s...", dbName))

	tables, err := d.driver.GetTables(dbName)
	if err != nil {
		d.showError(fmt.Sprintf("Failed to load tables: %v", err))
		return
	}

	// Find database node and add tables
	root := d.leftPanel.GetRoot()
	for _, child := range root.GetChildren() {
		ref := child.GetReference()
		if data, ok := ref.(map[string]string); ok && data["name"] == dbName {
			// Clear existing children
			child.ClearChildren()

			// Add tables
			for _, tableName := range tables {
				tableNode := tview.NewTreeNode("  " + tableName)
				tableNode.SetColor(style.StatusSelectedColor)
				tableNode.SetReference(map[string]string{
					"type":     "table",
					"database": dbName,
					"name":     tableName,
				})
				tableNode.SetSelectable(true)
				child.AddChild(tableNode)
			}
			child.SetExpanded(true)
			break
		}
	}

	d.updateStatusBar(fmt.Sprintf("Loaded %d tables from %s", len(tables), dbName))
}

// executeQuery executes the SQL query
func (d *Database) executeQuery() {
	if !d.connected || d.driver == nil {
		d.showError("Not connected to database")
		return
	}

	query := strings.TrimSpace(d.sqlEditor.GetText())
	if query == "" {
		return
	}

	d.updateStatusBar("Executing query...")

	result, err := d.driver.ExecuteQuery(query)
	if err != nil {
		d.showError(fmt.Sprintf("Query error: %v", err))
		d.updateStatusBar("Query failed")
		return
	}

	// Display results
	d.displayResult(result)
	d.updateStatusBar(fmt.Sprintf("Query executed. Rows: %d", result.RowsAffected))
}

// displayResult displays query results in table
func (d *Database) displayResult(result *db.QueryResult) {
	d.resultTable.Clear()

	if len(result.Columns) == 0 {
		// DML/DDL result
		cell := tview.NewTableCell(fmt.Sprintf("Rows affected: %d", result.RowsAffected))
		cell.SetTextColor(style.StatusInstalledColor)
		d.resultTable.SetCell(0, 0, cell)
		return
	}

	// Add headers
	for i, col := range result.Columns {
		cell := tview.NewTableCell(strings.ToUpper(col))
		cell.SetExpansion(1)
		cell.SetBackgroundColor(style.PageHeaderBgColor)
		cell.SetTextColor(style.PageHeaderFgColor)
		cell.SetAlign(tview.AlignLeft)
		cell.SetSelectable(false)
		d.resultTable.SetCell(0, i, cell)
	}

	// Add rows
	for rowIdx, row := range result.Rows {
		for colIdx, value := range row {
			cell := tview.NewTableCell(value)
			cell.SetTextColor(style.FgColor)
			cell.SetAlign(tview.AlignLeft)
			d.resultTable.SetCell(rowIdx+1, colIdx, cell)
		}
	}

	d.resultTable.ScrollToBeginning()
}

// showError shows an error dialog
func (d *Database) showError(message string) {
	d.errorDialog.SetTitle("Error")
	d.errorDialog.SetText(message)
	d.errorDialog.Display()
}

// InputHandler returns the input handler for this primitive
func (d *Database) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		// Handle dialog input first if dialog has focus
		if d.SubDialogHasFocus() {
			// Pass event to the dialog that has focus
			if d.connDialog.HasFocus() {
				if handler := d.connDialog.InputHandler(); handler != nil {
					handler(event, setFocus)
				}
			} else if d.errorDialog.HasFocus() {
				if handler := d.errorDialog.InputHandler(); handler != nil {
					handler(event, setFocus)
				}
			} else if d.messageDialog.HasFocus() {
				if handler := d.messageDialog.InputHandler(); handler != nil {
					handler(event, setFocus)
				}
			}
			return
		}

		// Ctrl+N to open connection dialog
		if event.Key() == tcell.KeyCtrlN {
			d.connDialog.Display()
			setFocus(d.connDialog)
			return
		}

		// Ctrl+D to disconnect
		if event.Key() == tcell.KeyCtrlD && d.connected {
			d.disconnect()
			return
		}

		// Ctrl+R to execute query
		if event.Key() == tcell.KeyCtrlR {
			d.executeQuery()
			return
		}

		// Ctrl+Left/Right to switch panels
		if event.Key() == tcell.KeyLeft && event.Modifiers() == tcell.ModCtrl {
			d.focusedElement = (d.focusedElement - 1 + 3) % 3
			d.Focus(setFocus)
			return
		}
		if event.Key() == tcell.KeyRight && event.Modifiers() == tcell.ModCtrl {
			d.focusedElement = (d.focusedElement + 1) % 3
			d.Focus(setFocus)
			return
		}

		// Handle current focused element
		switch d.focusedElement {
		case focusTree:
			if d.leftPanel.HasFocus() {
				if handler := d.leftPanel.InputHandler(); handler != nil {
					handler(event, setFocus)
				}
			}
		case focusEditor:
			if d.sqlEditor.HasFocus() {
				if handler := d.sqlEditor.InputHandler(); handler != nil {
					handler(event, setFocus)
				}
			}
		case focusResult:
			if d.resultTable.HasFocus() {
				if handler := d.resultTable.InputHandler(); handler != nil {
					handler(event, setFocus)
				}
			}
		}
	})
}

// disconnect disconnects from database
func (d *Database) disconnect() {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.driver != nil {
		d.driver.Close()
		d.driver = nil
	}

	d.connected = false
	d.currentDatabase = ""

	rootNode := tview.NewTreeNode("Not Connected")
	rootNode.SetColor(style.StatusNotInstalledColor)
	d.leftPanel.SetRoot(rootNode)
	d.leftPanel.SetCurrentNode(rootNode)

	d.updateStatusBar("Disconnected. Press Ctrl+N to connect")
}

// Draw draws this primitive onto the screen
func (d *Database) Draw(screen tcell.Screen) {
	d.Box.DrawForSubclass(screen, d)
	x, y, width, height := d.GetInnerRect()

	d.mainFlex.SetRect(x, y, width, height)
	d.mainFlex.Draw(screen)

	// Draw dialogs
	if d.errorDialog.IsDisplay() {
		d.errorDialog.SetRect(x, y, width, height)
		d.errorDialog.Draw(screen)
	}
	if d.messageDialog.IsDisplay() {
		d.messageDialog.SetRect(x, y, width, height)
		d.messageDialog.Draw(screen)
	}
	if d.connDialog.IsDisplay() {
		d.connDialog.SetRect(x, y, width, height)
		d.connDialog.Draw(screen)
	}
}
