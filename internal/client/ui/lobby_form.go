package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/milonoir/business-club-game/internal/network"
	"github.com/rivo/tview"
)

type LobbyForm struct {
	form *tview.Form
}

func NewLobbyForm(startCb, quitCb func()) *LobbyForm {
	l := &LobbyForm{
		form: tview.NewForm(),
	}
	l.form.
		AddTextView("Players", "", 20, 5, true, false).
		AddButton("Start", startCb).
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
