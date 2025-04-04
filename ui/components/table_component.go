package components

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"makerworld-analytics/domain"
	"makerworld-analytics/state"
)

type TableComponent struct {
	app.Compo
	statistics      *domain.Statistics
	moneyMultiplier domain.MoneyMultiplier
}

func NewTable(statistics *domain.Statistics, moneyMultiplier domain.MoneyMultiplier) *TableComponent {
	return &TableComponent{statistics: statistics, moneyMultiplier: moneyMultiplier}
}

func (h *TableComponent) OnMount(ctx app.Context) {
	ctx.ObserveState(state.MoneyMultiplierKey, &h.moneyMultiplier)
}

func (h *TableComponent) renderTabView() app.UI {
	// TODO: Separate component in separate file
	bankPayoutClass := "btn btn-default btn-soft"
	if h.moneyMultiplier == domain.BankPayoutMultiplier {
		bankPayoutClass = " btn btn-secondary"
	}

	vouchersClass := "ml-2 btn btn-default btn-soft"
	if h.moneyMultiplier == domain.VouchersMultiplier {
		vouchersClass = "ml-2 btn btn-secondary"
	}

	return app.Div().Class("flex flex-row").Body(
		app.Button().Class(bankPayoutClass).Text("Bank payout").OnClick(func(ctx app.Context, e app.Event) {
			ctx.SetState(state.MoneyMultiplierKey, domain.BankPayoutMultiplier).Persist()
		}),
		app.Button().Class(vouchersClass).Text("Vouchers").OnClick(func(ctx app.Context, e app.Event) {
			ctx.SetState(state.MoneyMultiplierKey, domain.VouchersMultiplier).Persist()
		}),
	)
}

func (h *TableComponent) renderMoneyRow(value float32) app.HTMLTd {
	return app.Td().Text(fmt.Sprintf("€%.0f", h.statistics.ToEuro(h.moneyMultiplier, value)))
}

func (h *TableComponent) Render() app.UI {
	if h.statistics == nil {
		return app.H1().Text("Error, rendered general stats without valid stats")
	}
	return app.Div().ID("general-statistics").Body(
		app.H1().Class("text-2xl opacity-70 mt-8 ml-2 mb-2").Text("💰 Euro earned"),
		h.renderTabView(),
		app.Table().Class("table").Body(
			app.THead().Body(
				app.Tr().Body(
					app.Th().Text("Total money"),
					app.Th().Text("Boosts"),
					app.Th().Text("Design"),
					app.Th().Text("Other"),
				),
			),
			app.TBody().Body(
				app.Tr().Body(
					h.renderMoneyRow(h.statistics.TotalPoints),
					h.renderMoneyRow(h.statistics.PointsFromBoosts),
					h.renderMoneyRow(h.statistics.PointsFromDesign),
					h.renderMoneyRow(h.statistics.PointsOther),
				),
			),
		),
	)
}
