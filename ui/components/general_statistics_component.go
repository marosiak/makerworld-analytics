package components

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"makerworld-analytics/domain"
	"makerworld-analytics/state"
)

type GeneralStatistics struct {
	app.Compo
	Statistics      *domain.Statistics
	MoneyMultiplier domain.MoneyMultiplier
}

func NewGeneralStatistics(statistics *domain.Statistics) *GeneralStatistics {
	return &GeneralStatistics{Statistics: statistics, MoneyMultiplier: domain.VouchersMultiplier}
}

func (h *GeneralStatistics) OnMount(ctx app.Context) {
	fmt.Printf("OnMount()\n")
	println("OnMount, ale z innego printa")
	ctx.GetState(state.MoneyMultiplierKey, &h.MoneyMultiplier)
	ctx.ObserveState(state.MoneyMultiplierKey, &h.MoneyMultiplier).OnChange(func() {
		fmt.Println("Testowy print gdy zmienił się state = " + fmt.Sprintf("%f", h.MoneyMultiplier))
	})
	//abc := ""
	//ctx.ObserveState("greet-name", &abc).OnChange(func() {
	//	fmt.Println("Testowy print gdy zmienił się state greet name)")
	//})
	//fmt.Println("changed state greet name to = ", abc)
}

func (h *GeneralStatistics) renderTabView() app.UI {
	// TODO: Separate component in separate file
	bankPayoutClass := "btn btn-default"
	if h.MoneyMultiplier == domain.BankPayoutMultiplier {
		bankPayoutClass += "btn-soft"
	}

	vouchersClass := "ml-1 btn btn-default"
	if h.MoneyMultiplier == domain.VouchersMultiplier {
		vouchersClass += "btn-soft"
	}
	return app.Div().Class("flex flex-row").Body(
		app.Button().Class(bankPayoutClass).Text("Bank payout").OnClick(func(ctx app.Context, e app.Event) {
			//h.MoneyMultiplier = utils.ValueToPointer(domain.BankPayoutMultiplier)
			ctx.SetState(state.MoneyMultiplierKey, domain.BankPayoutMultiplier).Persist()
			state.HandleGreet(ctx, app.Action{Value: "test"})
		}),
		app.Button().Class(vouchersClass).Text("Vouchers").OnClick(func(ctx app.Context, e app.Event) {
			//h.MoneyMultiplier = utils.ValueToPointer(domain.VouchersMultiplier)
			ctx.SetState(state.MoneyMultiplierKey, domain.VouchersMultiplier).Persist()
			state.HandleGreet(ctx, app.Action{Value: "test"})
		}),
	)
}

func (h *GeneralStatistics) Render() app.UI {

	if h.Statistics == nil {
		return app.H1().Text("Error, rendered general stats without valid stats")
	}
	return app.Div().ID("general-Statistics").Body(
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
					app.Td().Text(h.Statistics.TotalPoints),
					app.Td().Text(h.Statistics.PointsFromBoosts),
					app.Td().Text(h.Statistics.PointsFromDesign),
					app.Td().Text(h.Statistics.PointsOther),
				),
			),
		),
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
					app.Td().Text(h.Statistics.ToEuro(h.MoneyMultiplier, h.Statistics.TotalPoints)),
					app.Td().Text(h.Statistics.ToEuro(h.MoneyMultiplier, h.Statistics.PointsFromBoosts)),
					app.Td().Text(h.Statistics.ToEuro(h.MoneyMultiplier, h.Statistics.PointsFromDesign)),
					app.Td().Text(h.Statistics.ToEuro(h.MoneyMultiplier, h.Statistics.PointsOther)),
				),
			),
		),
	)
}
