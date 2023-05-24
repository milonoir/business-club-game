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

	cp *CompanyProvider
}

func NewStockPricePanel(cp *CompanyProvider, startPrice int) *StockPricePanel {
	p := &StockPricePanel{
		tv: tview.NewTextView(),
		cp: cp,
	}

	p.tv.
		SetDynamicColors(true).
		SetBorder(true).
		SetTitle("Stock prices")

	p.Update(startPrice, startPrice, startPrice, startPrice)

	return p
}

func (p *StockPricePanel) GetTextView() *tview.TextView {
	return p.tv
}

func (p *StockPricePanel) Update(p1, p2, p3, p4 int) {
	sb := strings.Builder{}
	companies := p.cp.Companies()

	sb.WriteString(fmt.Sprintf("%14s %s\n", "", textScaleValues))
	sb.WriteString(p.line(p.cp.ColorByCompany(companies[0]), companies[0], p.scale(p1), p1))
	sb.WriteString(p.line(p.cp.ColorByCompany(companies[1]), companies[1], p.scale(p2), p2))
	sb.WriteString(p.line(p.cp.ColorByCompany(companies[2]), companies[2], p.scale(p3), p3))
	sb.WriteString(p.line(p.cp.ColorByCompany(companies[3]), companies[3], p.scale(p4), p4))

	p.tv.SetText(sb.String())
}

func (p *StockPricePanel) line(color, company, scale string, price int) string {
	return fmt.Sprintf("[%s]%14s[white] %s [%s]%3d\n", color, company, scale, color, price)
}

func (p *StockPricePanel) scale(price int) string {
	i := price / 10
	return "[grey]" + textScale[:i] + "[fuchsia]#[grey]" + textScale[i+1:]
}
