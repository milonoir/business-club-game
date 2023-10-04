package ui

import (
	"fmt"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/milonoir/business-club-game/internal/message"
	"github.com/rivo/tview"
)

const (
	labelAmount = "Amount"
)

type TradeForm struct {
	form *tview.Form
}

func NewTradeForm(cp *CompanyProvider, tradeType message.TradeType, company, maxAmount int, cb func(int)) *TradeForm {
	t := &TradeForm{
		form: tview.NewForm(),
	}

	header := fmt.Sprintf(" %s [%s]%s [white]stock ", tradeType.AsString(), cp.ColorByIndex(company), cp.CompanyByIndex(company))

	t.form.
		AddInputField(fmt.Sprintf("%s (max: %d):", labelAmount, maxAmount), "", 10, PositiveIntegerValidator(maxAmount), nil).
		AddButton(tradeType.AsString(), t.callbackWrapper(cb))
	t.form.
		SetTitle(header).
		SetBorder(true).
		SetBorderColor(tcell.ColorGreen).
		SetBorderPadding(1, 1, 1, 1)

	return t
}

func (t *TradeForm) GetForm() *tview.Form {
	return t.form
}

func (t *TradeForm) callbackWrapper(cb func(int)) func() {
	return func() {
		amount, _ := strconv.Atoi(t.form.GetFormItem(0).(*tview.InputField).GetText())
		cb(amount)
	}
}
