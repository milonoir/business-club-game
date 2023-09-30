package ui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/milonoir/business-club-game/internal/message"
	"github.com/rivo/tview"
)

const (
	maxItem = 25
)

type HistoryPanel struct {
	tv *tview.TextView

	cp   *CompanyProvider
	logs []string
}

func NewHistoryPanel(cp *CompanyProvider) *HistoryPanel {
	p := &HistoryPanel{
		tv:   tview.NewTextView(),
		cp:   cp,
		logs: make([]string, 0, maxItem),
	}

	p.tv.
		SetDynamicColors(true).
		SetBorder(true).
		SetBorderColor(tcell.ColorGrey).
		//SetBorderStyle(tcell.Style{}.
		//	Foreground(tcell.ColorGrey)).
		SetTitle(" History ").
		SetBorderPadding(0, 0, 0, 1)

	return p
}

func (p *HistoryPanel) GetTextView() *tview.TextView {
	return p.tv
}

func (p *HistoryPanel) AddAction(a *message.Action) {
	sb := strings.Builder{}

	if a.ActorType == message.ActorPlayer {
		sb.WriteString(fmt.Sprintf("[yellow]%s ", a.Name))
	} else {
		sb.WriteString(fmt.Sprintf("[purple]BANK "))
	}

	company := p.cp.CompanyByIndex(a.Mod.Company)
	sb.WriteString(fmt.Sprintf("[white]action: [%s]%s ", p.cp.ColorByIndex(a.Mod.Company), company))

	switch op := a.Mod.Mod.Op(); op {
	case "+":
		sb.WriteString(fmt.Sprintf("[green]+%d ", a.Mod.Mod.Value()))
	case "-":
		sb.WriteString(fmt.Sprintf("[red]-%d ", a.Mod.Mod.Value()))
	case "*":
		sb.WriteString(fmt.Sprintf("[yellow]*%d ", a.Mod.Mod.Value()))
	case "=":
		sb.WriteString(fmt.Sprintf("[blue]=%d ", a.Mod.Mod.Value()))
	}

	sb.WriteString(fmt.Sprintf("[white]--> %d", a.NewPrice))

	p.addString(sb.String())
}

func (p *HistoryPanel) AddDeal(d *message.Deal) {
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("[yellow]%s [white]deal: ", d.Name))

	if d.Type == message.DealBuy {
		sb.WriteString(fmt.Sprintf("[green]buy "))
	} else {
		sb.WriteString(fmt.Sprintf("[red]sell "))
	}

	company := p.cp.CompanyByIndex(d.Company)
	sb.WriteString(fmt.Sprintf("[%s]%s ", p.cp.ColorByIndex(d.Company), company))
	sb.WriteString(fmt.Sprintf("[white]%d x %d = %d", d.Amount, d.Price, d.Amount*d.Price))

	p.addString(sb.String())
}

func (p *HistoryPanel) addString(s string) {
	p.logs = append(p.logs, s)
	if len(p.logs) > maxItem {
		p.logs = p.logs[1:]
	}

	p.redraw()
}

func (p *HistoryPanel) redraw() {
	p.tv.SetText(strings.Join(p.logs, "\n"))
}
