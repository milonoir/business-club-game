package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type WaitPanel struct {
	tv *tview.TextView
}

func NewWaitPanel(text string) *WaitPanel {
	p := &WaitPanel{
		tv: tview.NewTextView(),
	}

	p.tv.
		SetTextAlign(tview.AlignCenter).
		SetTextColor(tcell.ColorYellow).
		SetText(text)
	p.tv.
		SetBorderPadding(4, 0, 0, 0)

	return p
}

func (p *WaitPanel) GetTextView() *tview.TextView {
	return p.tv
}
