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
	t  *tview.Table
	cp *CompanyProvider
}

func NewOpponentsPanel(opponents []string, cp *CompanyProvider) *OpponentsPanel {
	p := &OpponentsPanel{
		t:  tview.NewTable(),
		cp: cp,
	}

	p.t.
		SetBorder(true)

	for i, o := range opponents {
		p.t.SetCell(0, i, tview.NewTableCell(center(o, 15)).SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	}

	p.Update(OpponentsUpdate{})

	return p
}

func (p *OpponentsPanel) GetTable() *tview.Table {
	return p.t
}

func (p *OpponentsPanel) Update(u OpponentsUpdate) {
	colors := p.companyColors()

	p.t.
		SetCell(1, 0, tview.NewTableCell(strings.Repeat("#", u.O1.N1)).SetTextColor(colors[0]).SetAlign(tview.AlignCenter)).
		SetCell(2, 0, tview.NewTableCell(strings.Repeat("#", u.O1.N2)).SetTextColor(colors[1]).SetAlign(tview.AlignCenter)).
		SetCell(3, 0, tview.NewTableCell(strings.Repeat("#", u.O1.N3)).SetTextColor(colors[2]).SetAlign(tview.AlignCenter)).
		SetCell(4, 0, tview.NewTableCell(strings.Repeat("#", u.O1.N4)).SetTextColor(colors[3]).SetAlign(tview.AlignCenter)).
		SetCell(5, 0, tview.NewTableCell(strings.Repeat("$", u.O1.C)).SetTextColor(tcell.ColorGreen).SetAlign(tview.AlignCenter))

	p.t.
		SetCell(1, 1, tview.NewTableCell(strings.Repeat("#", u.O2.N1)).SetTextColor(colors[0]).SetAlign(tview.AlignCenter)).
		SetCell(2, 1, tview.NewTableCell(strings.Repeat("#", u.O2.N2)).SetTextColor(colors[1]).SetAlign(tview.AlignCenter)).
		SetCell(3, 1, tview.NewTableCell(strings.Repeat("#", u.O2.N3)).SetTextColor(colors[2]).SetAlign(tview.AlignCenter)).
		SetCell(4, 1, tview.NewTableCell(strings.Repeat("#", u.O2.N4)).SetTextColor(colors[3]).SetAlign(tview.AlignCenter)).
		SetCell(5, 1, tview.NewTableCell(strings.Repeat("$", u.O2.C)).SetTextColor(tcell.ColorGreen).SetAlign(tview.AlignCenter))

	p.t.
		SetCell(1, 2, tview.NewTableCell(strings.Repeat("#", u.O3.N1)).SetTextColor(colors[0]).SetAlign(tview.AlignCenter)).
		SetCell(2, 2, tview.NewTableCell(strings.Repeat("#", u.O3.N2)).SetTextColor(colors[1]).SetAlign(tview.AlignCenter)).
		SetCell(3, 2, tview.NewTableCell(strings.Repeat("#", u.O3.N3)).SetTextColor(colors[2]).SetAlign(tview.AlignCenter)).
		SetCell(4, 2, tview.NewTableCell(strings.Repeat("#", u.O3.N4)).SetTextColor(colors[3]).SetAlign(tview.AlignCenter)).
		SetCell(5, 2, tview.NewTableCell(strings.Repeat("$", u.O3.C)).SetTextColor(tcell.ColorGreen).SetAlign(tview.AlignCenter))
}

func (p *OpponentsPanel) companyColors() []tcell.Color {
	c := p.cp.Companies()
	colors := make([]tcell.Color, len(c))
	for i, name := range c {
		colors[i] = tcell.ColorNames[p.cp.ColorByCompany(name)]
	}
	return colors
}

func center(s string, width int) string {
	if len(s) >= width {
		return s[:width]
	}

	diff := width - len(s)
	left := diff / 2
	right := left

	if diff%2 != 0 {
		right += 1
	}

	return strings.Repeat(" ", left) + s + strings.Repeat(" ", right)
}
