package client

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type OpponentData struct {
	N1, N2, N3, N4, C int
}

type OpponentsUpdate struct {
	O1 OpponentData
	O2 OpponentData
	O3 OpponentData
}

type OpponentsPanel struct {
	t *tview.Table
}

func NewOpponentsPanel(o1, o2, o3 string) *OpponentsPanel {
	p := &OpponentsPanel{
		t: tview.NewTable(),
	}

	p.t.
		SetBorders(true)

	p.t.
		SetCell(0, 0, tview.NewTableCell(o1).SetTextColor(tcell.ColorYellow)).
		SetCell(0, 1, tview.NewTableCell(o2).SetTextColor(tcell.ColorYellow)).
		SetCell(0, 2, tview.NewTableCell(o3).SetTextColor(tcell.ColorYellow))

	p.Update(OpponentsUpdate{})

	return p
}

func (p *OpponentsPanel) GetTable() *tview.Table {
	return p.t
}

func (p *OpponentsPanel) Update(u OpponentsUpdate) {
	p.t.
		SetCell(1, 0, tview.NewTableCell(strings.Repeat("*", u.O1.N1)).SetTextColor(tcell.ColorBlue).SetAlign(tview.AlignCenter)).
		SetCell(2, 0, tview.NewTableCell(strings.Repeat("*", u.O1.N2)).SetTextColor(tcell.ColorOrange).SetAlign(tview.AlignCenter)).
		SetCell(3, 0, tview.NewTableCell(strings.Repeat("*", u.O1.N3)).SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter)).
		SetCell(4, 0, tview.NewTableCell(strings.Repeat("*", u.O1.N4)).SetTextColor(tcell.ColorRed).SetAlign(tview.AlignCenter)).
		SetCell(5, 0, tview.NewTableCell(strings.Repeat("$", u.O1.C)).SetTextColor(tcell.ColorGreen).SetAlign(tview.AlignCenter))

	p.t.
		SetCell(1, 1, tview.NewTableCell(strings.Repeat("*", u.O2.N1)).SetTextColor(tcell.ColorBlue).SetAlign(tview.AlignCenter)).
		SetCell(2, 1, tview.NewTableCell(strings.Repeat("*", u.O2.N2)).SetTextColor(tcell.ColorOrange).SetAlign(tview.AlignCenter)).
		SetCell(3, 1, tview.NewTableCell(strings.Repeat("*", u.O2.N3)).SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter)).
		SetCell(4, 1, tview.NewTableCell(strings.Repeat("*", u.O2.N4)).SetTextColor(tcell.ColorRed).SetAlign(tview.AlignCenter)).
		SetCell(5, 1, tview.NewTableCell(strings.Repeat("$", u.O2.C)).SetTextColor(tcell.ColorGreen).SetAlign(tview.AlignCenter))

	p.t.
		SetCell(1, 2, tview.NewTableCell(strings.Repeat("*", u.O3.N1)).SetTextColor(tcell.ColorBlue).SetAlign(tview.AlignCenter)).
		SetCell(2, 2, tview.NewTableCell(strings.Repeat("*", u.O3.N2)).SetTextColor(tcell.ColorOrange).SetAlign(tview.AlignCenter)).
		SetCell(3, 2, tview.NewTableCell(strings.Repeat("*", u.O3.N3)).SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter)).
		SetCell(4, 2, tview.NewTableCell(strings.Repeat("*", u.O3.N4)).SetTextColor(tcell.ColorRed).SetAlign(tview.AlignCenter)).
		SetCell(5, 2, tview.NewTableCell(strings.Repeat("$", u.O3.C)).SetTextColor(tcell.ColorGreen).SetAlign(tview.AlignCenter))
}
