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

func NewServerStatus() *ServerStatusPanel {
	p := &ServerStatusPanel{
		tv: tview.NewTextView(),
	}

	p.tv.
		SetDynamicColors(true).
		SetTextAlign(tview.AlignRight).
		SetBorderPadding(0, 0, 0, 1).
		SetBorder(false)

	return p
}

func (p *ServerStatusPanel) GetTextView() *tview.TextView {
	return p.tv
}

func (p *ServerStatusPanel) SetHost(host string) {
	p.host = host
	p.redraw()
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

	cc := "red"
	if p.connected {
		cc = "green"
	} else {
	}
	sb.WriteString(fmt.Sprintf("[white]Server: [%s]%s   [white]Key: [blue]%s", cc, p.host, p.authKey))

	p.tv.SetText(sb.String())
}
