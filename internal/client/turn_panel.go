package client

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"
)

type TurnPanel struct {
	tv *tview.TextView

	maxTurns      int
	currentTurn   int
	currentPlayer int
	playerOrder   []string
}

func NewTurnPanel(max int) *TurnPanel {
	p := &TurnPanel{
		tv:       tview.NewTextView(),
		maxTurns: max,
	}

	p.tv.
		SetDynamicColors(true).
		SetBorder(true).
		SetTitle(fmt.Sprintf("Turn: 0/%d", p.maxTurns))

	return p
}

func (p *TurnPanel) GetTextView() *tview.TextView {
	return p.tv
}

func (p *TurnPanel) NewTurn(order []string) {
	if p.currentTurn >= p.maxTurns {
		return
	}
	p.currentTurn++

	p.tv.SetTitle(fmt.Sprintf("Turn: %d/%d", p.currentTurn, p.maxTurns))
	p.playerOrder = order
	p.currentPlayer = 0

	p.redraw()
}

func (p *TurnPanel) NextPlayer() {
	if p.currentPlayer >= len(p.playerOrder) {
		return
	}
	p.currentPlayer++

	p.redraw()
}

func (p *TurnPanel) redraw() {
	sb := strings.Builder{}

	for i, name := range p.playerOrder {
		if i == p.currentPlayer {
			sb.WriteString("[green]" + name + "[white]\n")
		} else {
			sb.WriteString(name + "\n")
		}
	}

	p.tv.SetText(sb.String())
}
