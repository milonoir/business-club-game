package client

import (
	"fmt"
	"strings"

	"github.com/milonoir/business-club-game/internal/game"
	"github.com/rivo/tview"
)

type ActionList struct {
	l *tview.List

	companies []string
	cards     []*game.Card
}

func NewActionList(companies []string, cards []*game.Card) *ActionList {
	a := &ActionList{
		l: tview.NewList(),

		companies: companies,
		cards:     cards,
	}

	a.l.
		ShowSecondaryText(false).
		SetBorder(true).
		SetTitle("Actions")

	a.Update()

	return a
}

func (a *ActionList) GetList() *tview.List {
	return a.l
}

func (a *ActionList) Update() {
	a.dropAll()

	for i, c := range a.cards {
		a.l.AddItem(a.cardToString(c), "", rune('a'+i), func() {
			// TODO: send selected card to server; potentially use some callback func.
		})
	}
}

func (a *ActionList) cardToString(c *game.Card) string {
	sb := strings.Builder{}

	for _, mod := range c.Mods {
		var company string
		switch mod.Company {
		case 0:
			company = fmt.Sprintf("[blue]%-12s", a.companies[0])
		case 1:
			company = fmt.Sprintf("[orange]%-12s", a.companies[1])
		case 2:
			company = fmt.Sprintf("[yellow]%-12s", a.companies[2])
		case 3:
			company = fmt.Sprintf("[red]%-12s", a.companies[3])
		default:
			company = fmt.Sprintf("[fuchsia]%-12s", "???")
		}

		var modifier string
		switch op := mod.Mod.Op(); op {
		case "+":
			modifier = fmt.Sprintf("[green]%s %-3d", op, mod.Mod.Value())
		case "-":
			modifier = fmt.Sprintf("[red]%s %-3d", op, mod.Mod.Value())
		case "*":
			modifier = fmt.Sprintf("[yellow]%s %-3d", op, mod.Mod.Value())
		case "=":
			modifier = fmt.Sprintf("[blue]%s %-3d", op, mod.Mod.Value())
		}

		sb.WriteString(fmt.Sprintf("%s %s   ", company, modifier))
	}

	return sb.String()
}

func (a *ActionList) dropAll() {
	for a.l.GetItemCount() > 0 {
		a.l.RemoveItem(0)
	}
}
