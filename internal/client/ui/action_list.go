package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/milonoir/business-club-game/internal/game"
	"github.com/rivo/tview"
)

type ActionList struct {
	l *tview.List

	cp *CompanyProvider
}

func NewActionList(cp *CompanyProvider, cards []*game.Card, cb func(*game.Card)) *ActionList {
	a := &ActionList{
		l:  tview.NewList(),
		cp: cp,
	}

	a.l.
		SetWrapAround(false).
		ShowSecondaryText(false).
		SetHighlightFullLine(true).
		SetSelectedBackgroundColor(tcell.ColorDarkGrey).
		SetBorderColor(tcell.ColorGreen).
		SetBorder(true).
		SetBorderPadding(0, 0, 1, 1).
		SetTitle(" Select an action card ")

	for _, c := range cards {
		c := c
		a.l.AddItem(cardToString(a.cp, c), "", 0, func() {
			cb(c)
		})
	}

	return a
}

func (a *ActionList) GetList() *tview.List {
	return a.l
}
