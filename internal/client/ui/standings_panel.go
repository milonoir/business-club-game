package ui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/milonoir/business-club-game/internal/message"
	"github.com/rivo/tview"
)

type StandingsPanel struct {
	g *tview.Grid

	cntv *tview.TextView
	pntv *tview.TextView
	ontv []*tview.TextView

	ptv *tview.TextView
	otv []*tview.TextView

	cp *CompanyProvider
}

func NewStandingsPanel(cp *CompanyProvider) *StandingsPanel {
	p := &StandingsPanel{
		g:  tview.NewGrid(),
		cp: cp,
	}

	// Setup grid.
	p.g.
		SetColumns(15, 20, 20, 20, 20).
		SetRows(1, 1, 6).
		SetBorders(false).
		SetBorder(false).
		SetBorderPadding(0, 1, 1, 0)

	// "Header" column with company names.
	p.cntv = tview.NewTextView()
	p.cntv.
		SetDynamicColors(true).
		SetTextAlign(tview.AlignRight).
		SetBorder(false)
	p.g.AddItem(p.cntv, 2, 0, 1, 1, 1, 1, false)

	// Header row for player names.
	p.pntv = createTextView()
	p.g.AddItem(p.pntv, 0, 1, 1, 1, 1, 1, false)
	p.ontv = make([]*tview.TextView, 3)
	for i := 0; i < 3; i++ {
		p.ontv[i] = createTextView()
		p.g.AddItem(p.ontv[i], 0, i+2, 1, 1, 1, 1, false)
	}

	// Separator line.
	tv := tview.NewTextView().SetTextColor(tcell.ColorGrey)
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

func (p *StandingsPanel) Update(state *message.GameState) {
	p.refreshCompanyNames()

	opps := make([]string, 3)
	for i, o := range state.Opponents {
		opps[i] = o.Name
		p.opponentUpdate(i, state.StockPrices, o.Stocks, o.Cash, state.Ended, state.Ended)
	}

	p.playerUpdate(state.StockPrices, state.Player.Stocks, state.Player.Cash)
	p.setPlayerNames(state.Player.Name, opps)
}

func (p *StandingsPanel) setPlayerNames(player string, opponents []string) {
	p.pntv.SetText(fmt.Sprintf("[green]%s", player))
	for i := 0; i < 3; i++ {
		p.ontv[i].SetText(fmt.Sprintf("[yellow]%s", opponents[i]))
	}
}

func (p *StandingsPanel) refreshCompanyNames() {
	s := make([]string, 6)
	for i := 0; i < 4; i++ {
		name := p.cp.CompanyByIndex(i)
		s[i] = fmt.Sprintf("[%s]%s[white]: ", p.cp.ColorByIndex(i), name)
	}
	s[4] = fmt.Sprintf("[green]Cash[white]: ")
	s[5] = fmt.Sprintf("[white]Total value: ")
	p.cntv.SetText(strings.Join(s, "\n"))
}

func (p *StandingsPanel) playerUpdate(prices, stocks [4]int, cash int) {
	p.ptv.SetText(p.generateBreakdownString(prices, stocks, cash, true, true))
}

func (p *StandingsPanel) opponentUpdate(index int, prices, stocks [4]int, cash int, showTotal, showNumbers bool) {
	p.otv[index].SetText(p.generateBreakdownString(prices, stocks, cash, showTotal, showNumbers))
}

func (p *StandingsPanel) generateBreakdownString(prices, stocks [4]int, cash int, showTotal, showNumbers bool) string {
	sb := strings.Builder{}

	if showNumbers {
		for i := 0; i < 4; i++ {
			sb.WriteString(fmt.Sprintf("[%s]%d\n", p.cp.ColorByIndex(i), stocks[i]))
		}
		sb.WriteString(fmt.Sprintf("[green]%d\n", cash))
	} else {
		for i := 0; i < 4; i++ {
			s := strings.Repeat("♦", stocks[i])
			if s == "" {
				s = "-"
			}
			sb.WriteString(fmt.Sprintf("[%s]%s\n", p.cp.ColorByIndex(i), s))
		}
		s := strings.Repeat("$", cash)
		if s == "" {
			s = "-"
		}
		sb.WriteString(fmt.Sprintf("[green]%s\n", s))
	}

	if showNumbers && showTotal {
		total := prices[0]*stocks[0] + prices[1]*stocks[1] + prices[2]*stocks[2] + prices[3]*stocks[3] + cash
		sb.WriteString(fmt.Sprintf("[white]%d\n", total))
	}

	return sb.String()
}

func createTextView() *tview.TextView {
	tv := tview.NewTextView()
	tv.
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetBorder(false)
	return tv
}
