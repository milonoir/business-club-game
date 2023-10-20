package ui

import (
	_ "embed"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

//go:embed bc_man.ascii
var logoArt string

type LogoPanel struct {
	tv *tview.TextView
}

func NewLogoPanel() *LogoPanel {
	p := &LogoPanel{
		tv: tview.NewTextView(),
	}

	p.tv.
		SetTextAlign(tview.AlignCenter).
		SetTextColor(tcell.ColorGrey).
		SetText(logoArt)
	p.tv.
		SetBorderPadding(1, 0, 0, 0)

	return p
}

func (p *LogoPanel) GetTextView() *tview.TextView {
	return p.tv
}
