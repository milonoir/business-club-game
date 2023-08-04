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
	form  *tview.Form
	ready bool
}

func NewLobbyForm(readyCb func(bool), leaveCb func()) *LobbyForm {
	l := &LobbyForm{
		form: tview.NewForm(),
	}
	l.form.
		AddTextView("Players:", "", 20, 5, true, false).
		AddButton(labelReady, l.toggleReady(readyCb)).
		AddButton("Leave", leaveCb)
	l.form.
		SetBorderPadding(20, 1, 0, 1)

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
	for i := 4 - len(state); i > 0; i-- {
		sb.WriteString("[grey]-- open slot --\n")
	}

	l.form.GetFormItem(0).(*tview.TextView).SetText(sb.String())
}

func (l *LobbyForm) toggleReady(readyCb func(bool)) func() {
	return func() {
		l.ready = !l.ready
		label := labelReady
		if l.ready {
			label = labelNotReady
		}
		l.form.GetButton(0).SetLabel(label)
		readyCb(l.ready)
	}
}

func (l *LobbyForm) Reset() {
	l.form.GetButton(0).SetLabel(labelReady)
	l.form.SetFocus(1)
	l.ready = false
}
