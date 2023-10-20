package ui

import (
	_ "embed"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

//go:embed congratulations.ascii
var congratsArt string

type CongratsPanel struct {
	tv *tview.TextView
}

func NewCongratsPanel() *CongratsPanel {
	p := &CongratsPanel{
		tv: tview.NewTextView(),
	}

	p.tv.
		SetTextAlign(tview.AlignCenter).
		SetTextColor(tcell.ColorGreen).
		SetText(congratsArt)
	p.tv.
		SetBorderPadding(1, 0, 0, 0)

	return p
}

func (p *CongratsPanel) GetTextView() *tview.TextView {
	return p.tv
}
