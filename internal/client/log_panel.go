package client

import (
	"github.com/rivo/tview"
)

type LogItem struct {
	Who  string
	What string
}

type LogPanel struct {
	tv *tview.TextView

	cp   *CompanyProvider
	logs []*LogItem
}

func NewLogPanel() *LogPanel {
	p := &LogPanel{
		tv: tview.NewTextView(),
	}

	p.tv.
		SetDynamicColors(true).
		SetBorder(true).
		SetTitle("Log")

	return p
}

func (p *LogPanel) GetTextView() *tview.TextView {
	return p.tv
}

func (p *LogPanel) AddLogItem(li *LogItem) {
	p.logs = append(p.logs[1:], li)
}

func (p *LogPanel) redraw() {

}
