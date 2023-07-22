package ui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/milonoir/business-club-game/internal/game"
	"github.com/rivo/tview"
)

type ActionList struct {
	l *tview.List

	cp *CompanyProvider
}

func NewActionList(cp *CompanyProvider) *ActionList {
	a := &ActionList{
		l: tview.NewList(),

		cp: cp,
	}

	a.l.
		ShowSecondaryText(false).
		SetHighlightFullLine(true).
		SetSelectedBackgroundColor(tcell.ColorDarkGrey).
		SetBorderColor(tcell.ColorGreen).
		SetBorder(true).
		SetBorderPadding(0, 0, 1, 1).
		SetTitle(" Actions ")

	return a
}

func (a *ActionList) GetList() *tview.List {
	return a.l
}

func (a *ActionList) Update(cards []*game.Card) {
	a.dropAll()

	for _, c := range cards {
		a.l.AddItem(a.cardToString(c), "", 0, func() {
			// TODO: send selected card to server; potentially use some callback func.
		})
	}
}

func (a *ActionList) cardToString(c *game.Card) string {
	sb := strings.Builder{}

	for _, mod := range c.Mods {
		company := fmt.Sprintf("[fuchsia]%-12s", "???")
		if mod.Company > -1 {
			name := a.cp.CompanyByIndex(mod.Company)
			company = fmt.Sprintf("[%s]%-12s", a.cp.ColorByCompany(name), name)
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
