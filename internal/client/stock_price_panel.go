package client

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"
)

const (
	textScaleValues = "0        100       200       300       400"
	textScale       = "|....|....|....|....|....|....|....|....|"
)

type StockPricePanel struct {
	tv *tview.TextView

	// Company names.
	c1, c2, c3, c4 string
}

func NewStockPricePanel(c1, c2, c3, c4 string, startPrice int) *StockPricePanel {
	p := &StockPricePanel{
		tv: tview.NewTextView(),
		c1: c1,
		c2: c2,
		c3: c3,
		c4: c4,
	}

	p.tv.
		SetDynamicColors(true).
		SetBorder(true)

	p.Update(startPrice, startPrice, startPrice, startPrice)

	return p
}

func (p *StockPricePanel) GetTextView() *tview.TextView {
	return p.tv
}

func (p *StockPricePanel) Update(p1, p2, p3, p4 int) {
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("%14s %s\n", "", textScaleValues))
	sb.WriteString(fmt.Sprintf("[blue]%14s[white] %s [blue]%3d\n", p.c1, p.scale(p1), p1))
	sb.WriteString(fmt.Sprintf("[orange]%14s[white] %s [orange]%3d\n", p.c2, p.scale(p2), p2))
	sb.WriteString(fmt.Sprintf("[yellow]%14s[white] %s [yellow]%3d\n", p.c3, p.scale(p3), p3))
	sb.WriteString(fmt.Sprintf("[red]%14s[white] %s [red]%3d\n", p.c4, p.scale(p4), p4))

	p.tv.SetText(sb.String())
}

func (p *StockPricePanel) scale(price int) string {
	i := price / 10
	return "[grey]" + textScale[:i] + "[fuchsia]*[grey]" + textScale[i+1:]
}
