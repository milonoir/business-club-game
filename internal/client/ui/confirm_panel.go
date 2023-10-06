package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ConfirmPanel struct {
	modal *tview.Modal
}

func NewConfirmPanel(text string, result chan bool) *ConfirmPanel {
	c := &ConfirmPanel{
		modal: tview.NewModal(),
	}

	c.modal.
		SetBackgroundColor(tcell.ColorBlue).
		SetText(text).
		AddButtons([]string{"Yes", "No"})

	c.modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Yes" {
			result <- true
		} else {
			result <- false
		}
	})

	return c
}

func (c *ConfirmPanel) GetModal() *tview.Modal {
	return c.modal
}
