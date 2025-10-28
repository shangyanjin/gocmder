package database

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shangyanjin/gocmder/internal/ui/style"
	"github.com/shangyanjin/gocmder/internal/ui/utils"
)

const (
	connDialogWidth  = 80
	connDialogHeight = 20
)

// ConnectionDialog is a dialog for database connection
type ConnectionDialog struct {
	*tview.Box

	layout          *tview.Flex
	form            *tview.Form
	display         bool
	connectFunc     func(driverType, dsn string)
	appFocusHandler func()
	sessionName     string
	driverType      string
	host            string
	port            string
	username        string
	password        string
	database        string
}

// NewConnectionDialog creates a new connection dialog
func NewConnectionDialog(connectFunc func(driverType, dsn string)) *ConnectionDialog {
	bgColor := style.DialogBgColor

	dialog := &ConnectionDialog{
		Box:         tview.NewBox(),
		display:     false,
		connectFunc: connectFunc,
		sessionName: "Default",
		driverType:  "PostgreSQL",
		host:        "localhost",
		port:        "5432",
		username:    "postgres",
		password:    "postgres",
		database:    "postgres",
	}

	// Create form
	dialog.form = tview.NewForm()
	dialog.form.SetBackgroundColor(bgColor)
	dialog.form.SetButtonBackgroundColor(style.ButtonBgColor)
	dialog.form.SetFieldBackgroundColor(style.BgColor)
	dialog.form.SetLabelColor(style.FgColor)
	dialog.form.SetFieldTextColor(style.FgColor)

	// Add form fields - Session Name first, then DB Type
	dialog.form.AddInputField("Session Name", dialog.sessionName, 30, nil, func(text string) {
		dialog.sessionName = text
	})

	dialog.form.AddInputField("DB Type", dialog.driverType, 30, nil, func(text string) {
		dialog.driverType = text
	})

	dialog.form.AddInputField("Host", dialog.host, 30, nil, func(text string) {
		dialog.host = text
	})

	dialog.form.AddInputField("Port", dialog.port, 10, nil, func(text string) {
		dialog.port = text
	})

	dialog.form.AddInputField("Username", dialog.username, 30, nil, func(text string) {
		dialog.username = text
	})

	dialog.form.AddPasswordField("Password", "", 30, '*', func(text string) {
		dialog.password = text
	})

	dialog.form.AddInputField("Database", dialog.database, 30, nil, func(text string) {
		dialog.database = text
	})

	// Add buttons
	dialog.form.AddButton("Connect", func() {
		dialog.handleConnect()
	})

	dialog.form.AddButton("Cancel", func() {
		dialog.Hide()
		// Restore focus to parent page
		if dialog.appFocusHandler != nil {
			dialog.appFocusHandler()
		}
	})

	dialog.form.SetButtonsAlign(tview.AlignCenter)

	// Create keyboard shortcuts hint with consistent style
	highlightColor := style.GetColorHex(style.StatusInstalledColor)
	shortcutsHint := tview.NewTextView()
	shortcutsHint.SetBackgroundColor(bgColor)
	shortcutsHint.SetTextColor(style.FgColor)
	shortcutsHint.SetDynamicColors(true)
	shortcutsText := " [" + highlightColor + "]ALT+M[-] MySQL | [" + highlightColor + "]ALT+P[-] PgSQL | [" + highlightColor + "]ALT+L[-] SQLite | [" + highlightColor + "]ALT+S[-] Save | [" + highlightColor + "]ALT+C[-] Connect"
	shortcutsHint.SetText(shortcutsText)

	// Create layout - compact version without extra spacing
	dialog.layout = tview.NewFlex().SetDirection(tview.FlexRow)
	dialog.layout.AddItem(dialog.form, 0, 1, true)
	dialog.layout.AddItem(shortcutsHint, 1, 0, false)
	dialog.layout.SetBorder(true)
	dialog.layout.SetTitle(" Database Connection ")
	dialog.layout.SetTitleColor(style.FgColor)
	dialog.layout.SetBorderColor(style.DialogBorderColor)
	dialog.layout.SetBackgroundColor(bgColor)

	return dialog
}

// Display displays this primitive
func (d *ConnectionDialog) Display() {
	d.display = true
}

// IsDisplay returns true if primitive is shown
func (d *ConnectionDialog) IsDisplay() bool {
	return d.display
}

// Hide stops displaying this primitive
func (d *ConnectionDialog) Hide() {
	d.display = false
}

// HasFocus returns whether or not this primitive has focus
func (d *ConnectionDialog) HasFocus() bool {
	return d.display && (d.form.HasFocus() || d.Box.HasFocus())
}

// Focus is called when this primitive receives focus
func (d *ConnectionDialog) Focus(delegate func(p tview.Primitive)) {
	delegate(d.form)
}

// SetAppFocusHandler sets the app focus handler
func (d *ConnectionDialog) SetAppFocusHandler(handler func()) {
	d.appFocusHandler = handler
}

// setDatabasePreset sets database type and default values
func (d *ConnectionDialog) setDatabasePreset(dbType string) {
	d.driverType = dbType

	// Get form fields (Session Name is 0, DB Type is 1)
	dbTypeField := d.form.GetFormItem(1).(*tview.InputField)
	hostField := d.form.GetFormItem(2).(*tview.InputField)
	portField := d.form.GetFormItem(3).(*tview.InputField)
	userField := d.form.GetFormItem(4).(*tview.InputField)
	passField := d.form.GetFormItem(5).(*tview.InputField)
	dbField := d.form.GetFormItem(6).(*tview.InputField)

	// Set values based on database type
	switch dbType {
	case "PostgreSQL":
		d.host = "localhost"
		d.port = "5432"
		d.username = "postgres"
		d.password = "postgres"
		d.database = "postgres"
	case "MySQL":
		d.host = "localhost"
		d.port = "3306"
		d.username = "root"
		d.password = "root"
		d.database = ""
	case "SQLite":
		d.host = ""
		d.port = ""
		d.username = ""
		d.password = ""
		d.database = "./sqlite.db"
	}

	// Update form fields
	dbTypeField.SetText(d.driverType)
	hostField.SetText(d.host)
	portField.SetText(d.port)
	userField.SetText(d.username)
	passField.SetText(d.password)
	dbField.SetText(d.database)
}

// handleSave saves and connects without closing dialog (for Alt+S)
func (d *ConnectionDialog) handleSave() {
	var dsn string

	if d.driverType == "MySQL" {
		// MySQL DSN format: user:password@tcp(host:port)/database
		dsn = d.username + ":" + d.password + "@tcp(" + d.host + ":" + d.port + ")/"
		if d.database != "" {
			dsn += d.database
		}
	} else {
		// PostgreSQL DSN format: postgres://user:password@host:port/database
		dsn = "postgres://" + d.username + ":" + d.password + "@" + d.host + ":" + d.port + "/"
		if d.database != "" {
			dsn += d.database
		}
		dsn += "?sslmode=disable"
	}

	// Connect without closing dialog
	if d.connectFunc != nil {
		d.connectFunc(d.driverType, dsn)
	}
}

// handleConnect handles the connect button - connects and closes dialog
func (d *ConnectionDialog) handleConnect() {
	d.handleSave()
	d.Hide()

	// Restore focus to parent page
	if d.appFocusHandler != nil {
		d.appFocusHandler()
	}
}

// InputHandler returns input handler function for this primitive
func (d *ConnectionDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if event.Key() == utils.CloseDialogKey.Key {
			d.Hide()
			// Restore focus to parent page
			if d.appFocusHandler != nil {
				d.appFocusHandler()
			}
			return
		}

		// Handle Alt+Key shortcuts for quick database type selection
		if event.Key() == tcell.KeyRune && event.Modifiers()&tcell.ModAlt != 0 {
			switch event.Rune() {
			case 'p', 'P': // Alt+P for PostgreSQL
				d.setDatabasePreset("PostgreSQL")
				return
			case 'm', 'M': // Alt+M for MySQL
				d.setDatabasePreset("MySQL")
				return
			case 'l', 'L': // Alt+L for SQLite
				d.setDatabasePreset("SQLite")
				return
			case 's', 'S': // Alt+S for Save without closing
				d.handleSave()
				return
			case 'c', 'C': // Alt+C for Connect and close
				d.handleConnect()
				return
			}
		}

		if d.form.HasFocus() {
			if formHandler := d.form.InputHandler(); formHandler != nil {
				formHandler(event, setFocus)
				return
			}
		}
	})
}

// SetRect sets rects for this primitive
func (d *ConnectionDialog) SetRect(x, y, width, height int) {
	ws := (width - connDialogWidth) / 2
	hs := (height - connDialogHeight) / 2
	dy := y + hs
	bWidth := connDialogWidth
	bHeight := connDialogHeight

	if connDialogWidth > width {
		ws = 0
		bWidth = width - 1
	}

	if connDialogHeight >= height {
		dy = y + 1
		bHeight = height - 1
	}

	d.Box.SetRect(x+ws, dy, bWidth, bHeight)

	x, y, width, height = d.GetInnerRect()
	d.layout.SetRect(x, y, width, height)
}

// Draw draws this primitive onto the screen
func (d *ConnectionDialog) Draw(screen tcell.Screen) {
	if !d.display {
		return
	}

	d.DrawForSubclass(screen, d)
	d.layout.Draw(screen)
}
