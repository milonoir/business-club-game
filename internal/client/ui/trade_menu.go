package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TradeOption int

const (
	Buy TradeOption = iota
	Sell
	EndTurn
)

type TradeMenu struct {
	l *tview.List
}

func NewTradeMenu(cb func(option TradeOption)) *TradeMenu {
	t := &TradeMenu{
		l: tview.NewList(),
	}

	t.l.
		SetWrapAround(false).
		ShowSecondaryText(false).
		SetHighlightFullLine(true).
		SetSelectedBackgroundColor(tcell.ColorDarkGray).
		SetBorderColor(tcell.ColorGreen).
		SetBorder(true).
		SetBorderPadding(0, 0, 1, 1).
		SetTitle(" Select an option ")

	options := []struct {
		val TradeOption
		str string
	}{
		{Buy, "Buy stock"},
		{Sell, "Sell stock"},
		{EndTurn, "End turn"},
	}

	for _, opt := range options {
		opt := opt
		t.l.AddItem(opt.str, "", 0, func() {
			cb(opt.val)
		})
	}

	return t
}

func (t *TradeMenu) GetList() *tview.List {
	return t.l
}
