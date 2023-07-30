package ui

import (
	"strconv"

	"github.com/rivo/tview"
)

const (
	// Form defaults
	defaultHost = "localhost"
	defaultPort = "8585"

	// Form labels
	labelUsername     = "Username"
	labelHost         = "Host"
	labelPort         = "Port"
	labelReconnectKey = "Reconnect key"
	labelTLS          = "TLS"
	labelLogin        = "Login"
	labelQuit         = "Quit"
)

// LoginData is the data returned by the LoginForm.
type LoginData struct {
	Username     string
	Host         string
	Port         uint16
	ReconnectKey string
	TLS          bool
}

// LoginForm wraps a tview.Form and adds some validation logic.
type LoginForm struct {
	form *tview.Form
}

// NewLoginForm creates a new LoginForm.
func NewLoginForm(loginCb func(*LoginData), quitCb func()) *LoginForm {
	l := &LoginForm{
		form: tview.NewForm(),
	}

	l.form.
		AddInputField(labelUsername, "", 20, tview.InputFieldMaxLength(15), nil).
		AddInputField(labelHost, defaultHost, 20, tview.InputFieldMaxLength(15), nil).
		AddInputField(labelPort, defaultPort, 6, PositiveIntegerValidator(65535), nil).
		AddInputField(labelReconnectKey, "", 20, tview.InputFieldMaxLength(15), nil).
		AddTextView("", "Provide key to\nreconnect if you got\ndisconnected.", 20, 3, true, false).
		AddCheckbox(labelTLS, false, nil).
		AddButton(labelLogin, l.formValidatorWrapper(loginCb)).
		AddButton(labelQuit, quitCb)
	l.form.
		SetBorderPadding(14, 1, 0, 1)

	return l
}

// GetForm returns the underlying tview.Form.
func (l *LoginForm) GetForm() *tview.Form {
	return l.form
}

// formValidatorWrapper returns a callback that can be used to return the validated form data as LoginData.
func (l *LoginForm) formValidatorWrapper(cb func(data *LoginData)) func() {
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
			Username:     user,
			Host:         host,
			Port:         uint16(port),
			ReconnectKey: l.form.GetFormItemByLabel(labelReconnectKey).(*tview.InputField).GetText(),
			TLS:          l.form.GetFormItemByLabel(labelTLS).(*tview.Checkbox).IsChecked(),
		})
	}
}

func (l *LoginForm) SetReconnectKey(key string) {
	l.form.GetFormItemByLabel(labelReconnectKey).(*tview.InputField).SetText(key)
}
