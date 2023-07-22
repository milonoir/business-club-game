package ui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/milonoir/business-club-game/internal/game"
	"github.com/rivo/tview"
)

const (
	maxItem = 25
)

type HistoryItem interface {
	HistoryString(*CompanyProvider) string
}

type DealType uint8

const (
	DealBuy DealType = iota
	DealSell
)

type ActorType uint8

const (
	ActorPlayer ActorType = iota
	ActorBank
)

type ActionItem struct {
	ActorType ActorType
	Name      string
	Mod       *game.Modifier
	NewPrice  int
}

func (i *ActionItem) HistoryString(cp *CompanyProvider) string {
	sb := strings.Builder{}

	if i.ActorType == ActorPlayer {
		sb.WriteString(fmt.Sprintf("[yellow]%s ", i.Name))
	} else {
		sb.WriteString(fmt.Sprintf("[purple]BANK "))
	}

	company := cp.CompanyByIndex(i.Mod.Company)
	sb.WriteString(fmt.Sprintf("[white]action: [%s]%s ", cp.ColorByCompany(company), company))

	switch op := i.Mod.Mod.Op(); op {
	case "+":
		sb.WriteString(fmt.Sprintf("[green]+%d ", i.Mod.Mod.Value()))
	case "-":
		sb.WriteString(fmt.Sprintf("[red]-%d ", i.Mod.Mod.Value()))
	case "*":
		sb.WriteString(fmt.Sprintf("[yellow]*%d ", i.Mod.Mod.Value()))
	case "=":
		sb.WriteString(fmt.Sprintf("[blue]=%d ", i.Mod.Mod.Value()))
	}

	sb.WriteString(fmt.Sprintf("[white]--> %d", i.NewPrice))

	return sb.String()
}

type DealItem struct {
	Name         string
	Type         DealType
	CompanyIndex int
	Amount       int
	Price        int
}

func (i *DealItem) HistoryString(cp *CompanyProvider) string {
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("[yellow]%s [white]deal: ", i.Name))

	if i.Type == DealBuy {
		sb.WriteString(fmt.Sprintf("[green]buy "))
	} else {
		sb.WriteString(fmt.Sprintf("[red]sell "))
	}

	company := cp.CompanyByIndex(i.CompanyIndex)
	sb.WriteString(fmt.Sprintf("[%s]%s ", cp.ColorByCompany(company), company))
	sb.WriteString(fmt.Sprintf("[white]%d x %d = %d", i.Amount, i.Price, i.Amount*i.Price))

	return sb.String()
}

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

func (p *HistoryPanel) AddItem(li HistoryItem) {
	p.AddString(li.HistoryString(p.cp))
}

func (p *HistoryPanel) AddString(s string) {
	p.logs = append(p.logs, s)
	if len(p.logs) > maxItem {
		p.logs = p.logs[1:]
	}

	p.redraw()
}

func (p *HistoryPanel) redraw() {
	p.tv.SetText(strings.Join(p.logs, "\n"))
}
