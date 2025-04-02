package components

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"makerworld-analytics/domain"
	"makerworld-analytics/state"
)

type GeneralStatistics struct {
	app.Compo
	statistics      *domain.Statistics
	moneyMultiplier domain.MoneyMultiplier
}

func NewGeneralStatistics(statistics *domain.Statistics, moneyMultiplier domain.MoneyMultiplier) *GeneralStatistics {
	return &GeneralStatistics{statistics: statistics, moneyMultiplier: moneyMultiplier}
}
func (h *GeneralStatistics) OnNav(ctx app.Context) {
	fmt.Printf("OnNav()\n")
}

func (h *GeneralStatistics) OnMount(ctx app.Context) {
	fmt.Printf("OnMount()\n")
	ctx.Dispatch(func(ctx app.Context) {
		ctx.GetState(state.MoneyMultiplierKey, &h.moneyMultiplier)
	})

	println("GetState() multiplier = " + fmt.Sprintf("%f", h.moneyMultiplier))
	ctx.ObserveState(state.MoneyMultiplierKey, &h.moneyMultiplier).OnChange(func() {
		fmt.Println("Testowy print gdy zmienił się state = " + fmt.Sprintf("%f", h.moneyMultiplier))
		ctx.Dispatch(func(ctx app.Context) {
			h.moneyMultiplier = h.moneyMultiplier
		})
		//h.moneyMultiplier = h.moneyMultiplier
	})
}

func (h *GeneralStatistics) renderTabView() app.UI {
	// TODO: Separate component in separate file
	bankPayoutClass := "btn btn-default"
	if h.moneyMultiplier == domain.BankPayoutMultiplier {
		bankPayoutClass += "btn-soft"
	}

	vouchersClass := "ml-1 btn btn-default"
	if h.moneyMultiplier == domain.VouchersMultiplier {
		vouchersClass += "btn-soft"
	}
	return app.Div().Class("flex flex-row").Body(
		app.Button().Class(bankPayoutClass).Text("Bank payout").OnClick(func(ctx app.Context, e app.Event) {
			//h.moneyMultiplier = utils.ValueToPointer(domain.BankPayoutMultiplier)
			ctx.SetState(state.MoneyMultiplierKey, domain.BankPayoutMultiplier).Persist()
		}),
		app.Button().Class(vouchersClass).Text("Vouchers").OnClick(func(ctx app.Context, e app.Event) {
			//h.moneyMultiplier = utils.ValueToPointer(domain.VouchersMultiplier)
			ctx.SetState(state.MoneyMultiplierKey, domain.VouchersMultiplier).Persist()
		}),
	)
}

func (h *GeneralStatistics) Render() app.UI {

	if h.statistics == nil {
		return app.H1().Text("Error, rendered general stats without valid stats")
	}
	return app.Div().ID("general-statistics").Body(
		app.H1().Class("text-2xl opacity-70 mt-8 ml-2 mb-2").Text("All points"),
		app.Table().Class("table").Body(
			app.THead().Body(
				app.Tr().Body(
					app.Th().Text("Total points"),
					app.Th().Text("Boosts"),
					app.Th().Text("Design"),
					app.Th().Text("Other"),
				),
			),
			app.TBody().Body(
				app.Tr().Body(
					app.Td().Text(h.statistics.TotalPoints),
					app.Td().Text(h.statistics.PointsFromBoosts),
					app.Td().Text(h.statistics.PointsFromDesign),
					app.Td().Text(h.statistics.PointsOther),
				),
			),
		),
		app.H1().Class("text-2xl opacity-70 mt-8 ml-2 mb-2").Text(fmt.Sprintf("💰 Euro earned ~ debug = %f", h.moneyMultiplier)),
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
					app.Td().Text(h.statistics.ToEuro(h.moneyMultiplier, h.statistics.TotalPoints)),
					app.Td().Text(h.statistics.ToEuro(h.moneyMultiplier, h.statistics.PointsFromBoosts)),
					app.Td().Text(h.statistics.ToEuro(h.moneyMultiplier, h.statistics.PointsFromDesign)),
					app.Td().Text(h.statistics.ToEuro(h.moneyMultiplier, h.statistics.PointsOther)),
				),
			),
		),
	)
}
