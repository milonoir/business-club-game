package client

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"
)

// https://go.dev/play/p/lL2veiZOrxR

type GraphPanel struct {
	g    *tview.Grid
	tvs  []*tview.TextView
	cur  *tview.TextView
	data [][4]int
	cp   *CompanyProvider
}

func NewGraphPanel(cp *CompanyProvider) *GraphPanel {
	p := &GraphPanel{
		g:    tview.NewGrid(),
		tvs:  make([]*tview.TextView, 10),
		data: make([][4]int, 0, 10),
		cp:   cp,
	}

	// Setting up the grid.
	p.g.
		SetColumns(4, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9).
		SetRows(21, 1).
		SetBorderPadding(0, 0, 0, 1)

	// Y-axis.
	s := make([]string, 0, 20)
	for i := 400; i >= 0; i -= 20 {
		s = append(s, fmt.Sprintf("%3d ", i))
	}
	tv := tview.NewTextView()
	tv.SetDynamicColors(true).SetBorder(false)
	tv.SetText(strings.Join(s, "\n"))
	p.g.AddItem(tv, 0, 0, 1, 1, 1, 1, false)

	// Current prices - bottom row.
	p.cur = tview.NewTextView()
	p.cur.SetTextAlign(tview.AlignCenter).SetDynamicColors(true).SetBorder(false)
	p.g.AddItem(p.cur, 1, 0, 1, 11, 1, 1, false)

	// Initializing sub-graphs.
	for i := 0; i < 10; i++ {
		tv = tview.NewTextView()
		tv.SetDynamicColors(true).SetBorder(false)
		tv.SetText(p.emptyGraph())
		p.tvs[i] = tv
		p.g.AddItem(tv, 0, i+1, 1, 1, 1, 1, false)
	}

	return p
}

func (p *GraphPanel) GetGrid() *tview.Grid {
	return p.g
}

func (p *GraphPanel) Add(v1, v2, v3, v4 int) {
	d := [4]int{v1, v2, v3, v4}
	p.data = append(p.data, d)
	if len(p.data) > 10 {
		p.data = p.data[1:]
	}
	p.redraw()
}

func (p *GraphPanel) redraw() {
	for i, d := range p.data {
		rows := make([]string, 0, 20)
		for lvl := 400; lvl >= 0; lvl -= 20 {
			rows = append(rows, p.generateRow(d[0], d[1], d[2], d[3], lvl))
		}
		p.tvs[i].SetText(strings.Join(rows, "\n"))
	}

	current := p.data[len(p.data)-1]
	cs := make([]string, 4)
	for i, name := range p.cp.Companies() {
		cs[i] = fmt.Sprintf("[%s]%s: [white]%d", p.cp.ColorByCompany(name), name, current[i])
	}
	p.cur.SetText(strings.Join(cs, "   "))
}

func (p *GraphPanel) emptyGraph() string {
	return strings.Repeat("[grey]─────────\n", 21)
}

func (p *GraphPanel) generateRow(v1, v2, v3, v4, lvl int) string {
	sb := strings.Builder{}

	sb.WriteString("[grey]─")
	sb.WriteString(p.barOrEmpty(p.cp.ColorByCompanyIndex(0), v1, lvl))
	sb.WriteString(p.barOrEmpty(p.cp.ColorByCompanyIndex(1), v2, lvl))
	sb.WriteString(p.barOrEmpty(p.cp.ColorByCompanyIndex(2), v3, lvl))
	sb.WriteString(p.barOrEmpty(p.cp.ColorByCompanyIndex(3), v4, lvl))

	return sb.String()
}

func (p *GraphPanel) barOrEmpty(color string, v, lvl int) string {
	if lvl == 0 {
		if v > lvl {
			return fmt.Sprintf("[%s]▀[grey]─", color)
		}
		return "──"
	}

	if v > lvl {
		return fmt.Sprintf("[%s]█[grey]─", color)
	} else if v == lvl {
		return fmt.Sprintf("[%s]▄[grey]─", color)
	}
	return "──"
}
