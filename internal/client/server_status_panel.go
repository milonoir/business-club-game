package client

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"
)

type ServerStatusPanel struct {
	tv *tview.TextView

	connected bool
	authKey   string
	host      string
}

func NewServerStatus(host string) *ServerStatusPanel {
	p := &ServerStatusPanel{
		tv:   tview.NewTextView(),
		host: host,
	}

	p.tv.
		SetDynamicColors(true).
		SetTextAlign(tview.AlignRight).
		SetBorder(false)

	p.redraw()

	return p
}

func (p *ServerStatusPanel) GetTextView() *tview.TextView {
	return p.tv
}

func (p *ServerStatusPanel) SetConnection(isConnected bool) {
	p.connected = isConnected
	p.redraw()
}

func (p *ServerStatusPanel) SetAuthKey(authKey string) {
	p.authKey = authKey
	p.redraw()
}

func (p *ServerStatusPanel) redraw() {
	sb := strings.Builder{}

	if p.connected {
		sb.WriteString("[green]· ")
	} else {
		sb.WriteString("[red]· ")
	}
	sb.WriteString(fmt.Sprintf("[white]Server: [blue]%s   [white]Key: [blue]%s", p.host, p.authKey))

	p.tv.SetText(sb.String())
}
