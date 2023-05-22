package client

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"
)

const (
	textConnected    = "[green]connected[white]"
	textDisconnected = "[red]disconnected[white]"
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
		SetBorder(true).
		SetTitle("Server")

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

	sb.WriteString(fmt.Sprintf("    Host: %s\n", p.host))
	sb.WriteString(fmt.Sprintf("Auth Key: [blue]%s[white]\n", p.authKey))

	status := textDisconnected
	if p.connected {
		status = textConnected
	}
	sb.WriteString(fmt.Sprintf("  Status: %s\n", status))

	p.tv.SetText(sb.String())
}
