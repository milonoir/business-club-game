package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/milonoir/business-club-game/internal/network"
	"github.com/rivo/tview"
)

const (
	labelReady    = "Ready!"
	labelNotReady = "Cancel"
)

type LobbyForm struct {
	form *tview.Form
}

func NewLobbyForm(readyCb func(bool), quitCb func()) *LobbyForm {
	l := &LobbyForm{
		form: tview.NewForm(),
	}
	l.form.
		AddTextView("Players", "", 20, 5, true, false).
		AddButton(labelReady, l.toggleReady(readyCb)).
		AddButton("Quit", quitCb)
	l.form.
		SetBorderPadding(14, 1, 0, 1)

	return l
}

func (l *LobbyForm) GetForm() *tview.Form {
	return l.form
}

func (l *LobbyForm) Update(state []network.Readiness) {
	// Sanity check.
	if state == nil {
		return
	}

	// Sort players by name.
	sort.Slice(state, func(i, j int) bool {
		return state[i].Name < state[j].Name
	})

	sb := strings.Builder{}
	for _, p := range state {
		c := "red"
		if p.Ready {
			c = "green"
		}
		sb.WriteString(fmt.Sprintf("[%s]%s\n", c, p.Name))
	}

	l.form.GetFormItem(0).(*tview.TextView).SetText(sb.String())
}

func (l *LobbyForm) toggleReady(readyCb func(bool)) func() {
	ready := false
	return func() {
		ready = !ready
		label := labelReady
		if ready {
			label = labelNotReady
		}
		l.form.GetButton(0).SetLabel(label)
		readyCb(ready)
	}
}
