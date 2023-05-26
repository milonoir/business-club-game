package client

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type VersionPanel struct {
	tv *tview.TextView
}

func NewVersionPanel() *VersionPanel {
	p := &VersionPanel{
		tv: tview.NewTextView(),
	}

	p.tv.
		SetTextAlign(tview.AlignLeft).
		SetTextColor(tcell.ColorYellow).
		SetBorder(false)

	return p
}

func (p *VersionPanel) GetTextView() *tview.TextView {
	return p.tv
}

func (p *VersionPanel) SetVersion(v string) {
	p.tv.SetText(fmt.Sprintf("The Business Club v%s", v))
}
