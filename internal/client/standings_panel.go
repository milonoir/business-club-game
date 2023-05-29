package client

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"
)

type StandingsPanel struct {
	g  *tview.Grid
	cp *CompanyProvider

	ptv *tview.TextView
	otv []*tview.TextView

	prices []int
}

func NewStandingsPanel(player string, opponents []string, cp *CompanyProvider) *StandingsPanel {
	p := &StandingsPanel{
		g:      tview.NewGrid(),
		cp:     cp,
		prices: make([]int, 4),
	}

	// Setup grid.
	p.g.
		SetColumns(15, 20, 20, 20, 20).
		SetRows(1, 1, 6).
		SetBorders(false).
		SetBorder(false).
		SetBorderPadding(0, 2, 1, 0)

	// "Header" column.
	tv := tview.NewTextView()
	tv.
		SetDynamicColors(true).
		SetTextAlign(tview.AlignRight).
		SetBorder(false)
	s := make([]string, 6)
	for i := 0; i < 4; i++ {
		name := cp.CompanyByIndex(i)
		s[i] = fmt.Sprintf("[%s]%s[white]: ", cp.ColorByCompany(name), name)
	}
	s[4] = fmt.Sprintf("[green]Cash[white]: ")
	s[5] = fmt.Sprintf("[white]Total value: ")
	tv.SetText(strings.Join(s, "\n"))
	p.g.AddItem(tv, 2, 0, 1, 1, 1, 1, false)

	// Header row.
	players := append([]string{player}, opponents...)
	for col := 1; col < 5; col++ {
		tv = createTextView()
		if col == 1 {
			tv.SetText(fmt.Sprintf("[green]%s", players[col-1]))
		} else {
			tv.SetText(fmt.Sprintf("[yellow]%s", players[col-1]))
		}
		p.g.AddItem(tv, 0, col, 1, 1, 1, 1, false)
	}

	// Separator line.
	tv = tview.NewTextView()
	tv.SetText(strings.Repeat("─", 95))
	p.g.AddItem(tv, 1, 0, 1, 5, 1, 1, false)

	// Player's breakdown.
	p.ptv = createTextView()
	p.g.AddItem(p.ptv, 2, 1, 1, 1, 1, 1, false)

	// Opponents' breakdown.
	p.otv = make([]*tview.TextView, 3)
	for i := range p.otv {
		p.otv[i] = createTextView()
		p.g.AddItem(p.otv[i], 2, i+2, 1, 1, 1, 1, false)
	}

	return p
}

func (p *StandingsPanel) GetGrid() *tview.Grid {
	return p.g
}

func (p *StandingsPanel) PlayerUpdate(n1, n2, n3, n4, cash int) {
	p.ptv.SetText(p.generateBreakdownString(n1, n2, n3, n4, cash, true, true))
}

func (p *StandingsPanel) OpponentUpdate(index, n1, n2, n3, n4, cash int, showTotal, showNumbers bool) {
	p.otv[index].SetText(p.generateBreakdownString(n1, n2, n3, n4, cash, showTotal, showNumbers))
}

func (p *StandingsPanel) generateBreakdownString(n1, n2, n3, n4, cash int, showTotal, showNumbers bool) string {
	sb := strings.Builder{}

	if showNumbers {
		sb.WriteString(fmt.Sprintf("[%s]%d\n", p.cp.ColorByCompanyIndex(0), n1))
		sb.WriteString(fmt.Sprintf("[%s]%d\n", p.cp.ColorByCompanyIndex(1), n2))
		sb.WriteString(fmt.Sprintf("[%s]%d\n", p.cp.ColorByCompanyIndex(2), n3))
		sb.WriteString(fmt.Sprintf("[%s]%d\n", p.cp.ColorByCompanyIndex(3), n4))
		sb.WriteString(fmt.Sprintf("[green]%d\n", cash))
	} else {
		sb.WriteString(fmt.Sprintf("[%s]%s\n", p.cp.ColorByCompanyIndex(0), strings.Repeat("■ ", n1)))
		sb.WriteString(fmt.Sprintf("[%s]%s\n", p.cp.ColorByCompanyIndex(1), strings.Repeat("■ ", n2)))
		sb.WriteString(fmt.Sprintf("[%s]%s\n", p.cp.ColorByCompanyIndex(2), strings.Repeat("■ ", n3)))
		sb.WriteString(fmt.Sprintf("[%s]%s\n", p.cp.ColorByCompanyIndex(3), strings.Repeat("■ ", n4)))
		sb.WriteString(fmt.Sprintf("[green]%s\n", strings.Repeat("$ ", cash)))
	}

	if showNumbers && showTotal {
		total := p.prices[0]*n1 + p.prices[1]*n2 + p.prices[2]*n3 + p.prices[3]*n4 + cash
		sb.WriteString(fmt.Sprintf("[white]%d\n", total))
	}

	return sb.String()
}

func (p *StandingsPanel) SetPrices(p1, p2, p3, p4 int) {
	p.prices = []int{p1, p2, p3, p4}
}

func createTextView() *tview.TextView {
	tv := tview.NewTextView()
	tv.
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetBorder(false)
	return tv
}
