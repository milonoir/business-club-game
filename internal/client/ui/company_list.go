package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type CompanyList struct {
	l *tview.List

	cp *CompanyProvider
}

func NewCompanyList(cp *CompanyProvider, cb func(int)) *CompanyList {
	c := &CompanyList{
		l:  tview.NewList(),
		cp: cp,
	}

	c.l.
		SetWrapAround(false).
		ShowSecondaryText(false).
		SetHighlightFullLine(true).
		SetSelectedBackgroundColor(tcell.ColorDarkGray).
		SetBorderColor(tcell.ColorGreen).
		SetBorder(true).
		SetBorderPadding(0, 0, 1, 1).
		SetTitle(" Select a company ")

	for i, comp := range c.cp.Companies() {
		i := i
		comp := comp
		c.l.AddItem(fmt.Sprintf("[%s]%-12s", c.cp.ColorByIndex(i), comp), "", 0, func() {
			cb(i)
		})
	}

	return c
}

func (c *CompanyList) GetList() *tview.List {
	return c.l
}
