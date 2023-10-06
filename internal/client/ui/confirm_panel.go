package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ConfirmPanel struct {
	form *tview.Form
}

func NewConfirmPanel(result chan bool) *ConfirmPanel {
	c := &ConfirmPanel{
		form: tview.NewForm(),
	}

	c.form.
		SetButtonsAlign(tview.AlignCenter).
		SetBorder(true).
		SetBorderColor(tcell.ColorGreen).
		SetBorderPadding(3, 1, 1, 1).
		SetTitle(" Are you sure? ")

	c.form.AddButton("Yes", func() {
		result <- true
	})
	c.form.AddButton("No", func() {
		result <- false
	})

	c.form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyDown, tcell.KeyRight:
			return tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone)
		case tcell.KeyUp, tcell.KeyLeft:
			return tcell.NewEventKey(tcell.KeyBacktab, 0, tcell.ModNone)
		}
		return event
	})

	return c
}

func (c *ConfirmPanel) GetForm() *tview.Form {
	return c.form
}
