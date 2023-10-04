package ui

import (
	"fmt"
	"strconv"

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

	header := fmt.Sprintf("%s [%s]%s stock", tradeType.AsString(), cp.ColorByIndex(company), cp.CompanyByIndex(company))

	t.form.
		AddTextView("", header, 30, 1, true, false).
		AddInputField(labelAmount, "", 10, PositiveIntegerValidator(maxAmount), nil).
		AddButton("Trade", t.callbackWrapper(cb))
	t.form.
		SetBorder(true).
		SetBorderPadding(1, 1, 1, 1)

	return t
}

func (t *TradeForm) GetForm() *tview.Form {
	return t.form
}

func (t *TradeForm) callbackWrapper(cb func(int)) func() {
	return func() {
		amount, _ := strconv.Atoi(t.form.GetFormItemByLabel(labelAmount).(*tview.InputField).GetText())
		cb(amount)
	}
}
