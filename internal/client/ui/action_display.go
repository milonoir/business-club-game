package ui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/milonoir/business-club-game/internal/game"
	"github.com/rivo/tview"
)

type ActionDisplay struct {
	tv *tview.TextView

	cp *CompanyProvider
}

func NewActionDisplay(cp *CompanyProvider) *ActionDisplay {
	a := &ActionDisplay{
		tv: tview.NewTextView(),
		cp: cp,
	}

	a.tv.
		SetDynamicColors(true).
		SetBorder(true).
		SetBorderColor(tcell.ColorRed).
		SetBorderPadding(0, 0, 1, 0).
		SetTitle(" Please wait for your turn ")

	return a
}

func (a *ActionDisplay) GetTextView() *tview.TextView {
	return a.tv
}

func (a *ActionDisplay) Update(cards []*game.Card) {
	a.tv.Clear()

	s := make([]string, 0, len(cards))
	for _, c := range cards {
		s = append(s, cardToString(a.cp, c))
	}

	a.tv.SetText(strings.Join(s, "\n"))
}
