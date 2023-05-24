package client

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"
)

type PortfolioUpdate struct {
	P1, N1 int
	P2, N2 int
	P3, N3 int
	P4, N4 int
	Cash   int
}

type PortfolioPanel struct {
	tv *tview.TextView
	cp *CompanyProvider
}

func NewPortfolioPanel(cp *CompanyProvider) *PortfolioPanel {
	p := &PortfolioPanel{
		tv: tview.NewTextView(),
		cp: cp,
	}

	p.tv.
		SetDynamicColors(true).
		SetBorder(true).
		SetTitle("Portfolio")

	return p
}

func (p *PortfolioPanel) GetTextView() *tview.TextView {
	return p.tv
}

func (p *PortfolioPanel) Update(u PortfolioUpdate) {
	total := u.P1*u.N1 + u.P2*u.N2 + u.P3*u.N3 + u.P4*u.N4 + u.Cash
	c := p.cp.Companies()

	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("[%s]%12s[white]: %d\n", p.cp.ColorByCompany(c[0]), c[0], u.N1))
	sb.WriteString(fmt.Sprintf("[%s]%12s[white]: %d\n", p.cp.ColorByCompany(c[1]), c[1], u.N2))
	sb.WriteString(fmt.Sprintf("[%s]%12s[white]: %d\n", p.cp.ColorByCompany(c[2]), c[2], u.N3))
	sb.WriteString(fmt.Sprintf("[%s]%12s[white]: %d\n", p.cp.ColorByCompany(c[3]), c[3], u.N4))
	sb.WriteString(fmt.Sprintf("[green]%12s[white]: %d\n\n", "Cash", u.Cash))
	sb.WriteString(fmt.Sprintf("%12s: %d\n", "Total value", total))

	p.tv.SetText(sb.String())
}
