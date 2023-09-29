package ui

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"
)

type TurnPanel struct {
	tv *tview.TextView
}

func NewTurnPanel() *TurnPanel {
	p := &TurnPanel{
		tv: tview.NewTextView(),
	}

	p.tv.
		SetDynamicColors(true).
		SetBorder(false)

	return p
}

func (p *TurnPanel) GetTextView() *tview.TextView {
	return p.tv
}

func (p *TurnPanel) Update(maxTurns, currentTurn int, playerOrder []string, currentPlayer int) {
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("[yellow]Turn: %d/%d\n\n", currentTurn, maxTurns))

	for i, name := range playerOrder {
		if i == currentPlayer {
			sb.WriteString(fmt.Sprintf("[red]» %s\n", name))
		} else {
			sb.WriteString(fmt.Sprintf("[white]  %s\n", name))
		}
	}

	if currentPlayer >= len(playerOrder) {
		sb.WriteString("[red]» BANK\n")
	} else {
		sb.WriteString("[purple]  BANK\n")
	}

	p.tv.SetText(sb.String())
}
