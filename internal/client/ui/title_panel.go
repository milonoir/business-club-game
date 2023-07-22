package ui

import (
	_ "embed"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

//go:embed title.ascii
var titleArt string

type TitlePanel struct {
	tv *tview.TextView
}

func NewTitlePanel() *TitlePanel {
	p := &TitlePanel{
		tv: tview.NewTextView(),
	}

	p.tv.
		SetTextAlign(tview.AlignLeft).
		SetTextColor(tcell.ColorYellow).
		SetText(titleArt)
	p.tv.
		SetBorderPadding(6, 1, 10, 1)

	return p
}

func (p *TitlePanel) GetTextView() *tview.TextView {
	return p.tv
}
