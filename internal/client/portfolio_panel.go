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

	// Company names.
	c1, c2, c3, c4 string
}

func NewPortfolioPanel(c1, c2, c3, c4 string) *PortfolioPanel {
	p := &PortfolioPanel{
		tv: tview.NewTextView(),
		c1: c1,
		c2: c2,
		c3: c3,
		c4: c4,
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

	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("[blue]%s[white]: %d\n", p.c1, u.N1))
	sb.WriteString(fmt.Sprintf("[orange]%s[white]: %d\n", p.c2, u.N2))
	sb.WriteString(fmt.Sprintf("[yellow]%s[white]: %d\n", p.c3, u.N3))
	sb.WriteString(fmt.Sprintf("[red]%s[white]: %d\n", p.c4, u.N4))
	sb.WriteString(fmt.Sprintf("[green]Cash[white]: %d\n\n", u.Cash))
	sb.WriteString(fmt.Sprintf("[white]Total: %d\n", total))

	p.tv.SetText(sb.String())
}
