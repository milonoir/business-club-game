package client

import (
	"github.com/rivo/tview"
)

type ActionType string

const (
	ActionChange = "change"
	ActionBuy    = "buy"
	ActionSell   = "sell"
)

type ActorType uint8

const (
	ActorPlayer ActorType = iota
	ActorBank
)

// Examples:
// [Player/Bank] action: [Company] [Mod] [Value] = [New price]
// [Player] deal: [Action - buy/sell] [Company] [Amount] @ [Price]

type LogItem struct {
	Actor      string
	ActorColor string
	Action     ActionType
	Company    string
}

type LogPanel struct {
	tv *tview.TextView

	cp   *CompanyProvider
	logs []*LogItem
}

func NewLogPanel() *LogPanel {
	p := &LogPanel{
		tv:   tview.NewTextView(),
		logs: make([]*LogItem, 0, 10),
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
