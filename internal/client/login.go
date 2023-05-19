package client

import (
	"strconv"

	"github.com/rivo/tview"
)

const (
	// Form defaults
	defaultHost = "localhost"
	defaultPort = "8585"

	// Form labels
	labelUsername = "Username"
	labelHost     = "Host"
	labelPort     = "Port"
	labelAuthKey  = "Auth Key"
	labelTLS      = "TLS"
	labelLogin    = "Login"
	labelQuit     = "Quit"
)

// LoginData is the data returned from the login form.
type LoginData struct {
	Username string
	Host     string
	Port     uint16
	AuthKey  string
	TLS      bool
}

// Login is the login form.
type Login struct {
	form *tview.Form
}

// NewLogin creates a new login form.
func NewLogin(connectCb func(*LoginData), quitCb func()) *Login {
	l := &Login{
		form: tview.NewForm(),
	}

	l.form.
		AddInputField(labelUsername, "", 20, tview.InputFieldMaxLength(15), nil).
		AddInputField(labelHost, defaultHost, 20, tview.InputFieldMaxLength(15), nil).
		AddInputField(labelPort, defaultPort, 6, PositiveIntegerValidator(65535), nil).
		AddInputField(labelAuthKey, "", 20, tview.InputFieldMaxLength(15), nil).
		AddTextView("", "Provide auth key to\nreconnect if you got\ndisconnected.", 20, 3, true, false).
		AddCheckbox(labelTLS, false, nil).
		AddButton(labelLogin, l.Connect(connectCb)).
		AddButton(labelQuit, quitCb)
	l.form.
		SetBorderPadding(14, 1, 0, 1)

	return l
}

// GetForm returns the login form.
func (l *Login) GetForm() *tview.Form {
	return l.form
}

// Connect returns a callback that can be used to return validated login data.
func (l *Login) Connect(cb func(data *LoginData)) func() {
	return func() {
		user := l.form.GetFormItemByLabel(labelUsername).(*tview.InputField).GetText()
		if user == "" {
			return
		}
		host := l.form.GetFormItemByLabel(labelHost).(*tview.InputField).GetText()
		if host == "" {
			return
		}
		portStr := l.form.GetFormItemByLabel(labelPort).(*tview.InputField).GetText()
		port, err := strconv.ParseUint(portStr, 10, 16)
		if err != nil {
			return
		}

		cb(&LoginData{
			Username: user,
			Host:     host,
			Port:     uint16(port),
			AuthKey:  l.form.GetFormItemByLabel(labelAuthKey).(*tview.InputField).GetText(),
			TLS:      l.form.GetFormItemByLabel(labelTLS).(*tview.Checkbox).IsChecked(),
		})
	}
}
