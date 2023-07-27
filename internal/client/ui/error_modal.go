package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ErrorModal is a modal that displays an error message.
type ErrorModal struct {
	modal *tview.Modal
}

// NewErrorModal creates a new ErrorModal.
func NewErrorModal(err error) *ErrorModal {
	e := &ErrorModal{
		modal: tview.NewModal(),
	}

	e.modal.
		SetBackgroundColor(tcell.ColorRed).
		SetText(err.Error()).
		AddButtons([]string{"OK"})
	e.modal.SetTitle("Error")

	return e
}

// GetModal returns the underlying tview.Modal.
func (e *ErrorModal) GetModal() *tview.Modal {
	return e.modal
}

// SetHandler sets the done func for the modal.
func (e *ErrorModal) SetHandler(handler func(buttonIndex int, buttonLabel string)) *ErrorModal {
	e.modal.SetDoneFunc(handler)
	return e
}
