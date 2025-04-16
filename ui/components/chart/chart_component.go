package chart

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"makerworld-analytics/domain"
	"makerworld-analytics/echarts_wasm"
	"makerworld-analytics/ui/components"
	"math"
	"time"
)

type ChartsGridComponent struct {
	app.Compo
	Statistics                            *domain.Statistics
	StartDate                             *time.Time
	EndDate                               *time.Time
	MoneyMultiplier                       domain.MoneyMultiplier
	MinimumPointsThresholdForPieChart     float32
	MinimumPointsThresholdForStackedChart float32

	SelectedDesign *domain.PublishedDesign
}

func roundFloat(val float32, precision uint) float32 {
	ratio := math.Pow(10, float64(precision))
	return float32(math.Round(float64(val)*ratio) / ratio)
}

func (h *ChartsGridComponent) getTimeFormat() string {
	formatLayout := "2006-01-02"
	if h.StartDate != nil {
		if h.EndDate == nil {
			if time.Since(*h.StartDate).Hours() < 72 {
				formatLayout = "2006-01-02 15:04"
			}
		} else {
			if h.EndDate.Sub(*h.StartDate).Hours() < 72 {
				formatLayout = "2006-01-02 15:04"
			}
		}
	}

	return formatLayout
}

func (h *ChartsGridComponent) Render() app.UI {
	return app.Div().Class("flex flex-col w-full").Body(
		app.If(h.SelectedDesign != nil, func() app.UI {
			accumulatedEuroPerModel := h.accumulatedEuroPerModelStackedChart([]domain.DesignID{h.SelectedDesign.ID})
			return app.If(len(accumulatedEuroPerModel.XAxis[0].Data) >= 7, func() app.UI {
				return &components.CardComponent{
					Body: []app.UI{
						app.H2().Class("text-xl").Text("Accumulated euro per model"),
						&echarts_wasm.EChartComp{
							ContainerID:   "euro-per-model-stacked-accumulatedEuroPerModel",
							Option:        accumulatedEuroPerModel,
							WidthCssValue: "77vw",
						},
					},
					Class: "",
				}
			}).Else(
				func() app.UI {
					return app.Div().Class("flex flex-col wrap w-full pt-1 justify-stretch").Body(
						app.H1().Class("text-xl opacity-75 mt-4 select-none").Text("Not enough data to render chart"),
						app.P().Class("text-md opacity-45 mt-1 select-none text-black").Textf("Try selecting different model, or date range"))
				})
		}).Else(func() app.UI {
			return app.Div().Class("flex flex-row wrap w-full pt-1 justify-stretch").Body(
				&components.CardComponent{
					Body: []app.UI{
						app.H2().Class("text-xl").Text("Euro income"),
						&echarts_wasm.EChartComp{
							ContainerID:   "euro-per-day-accumulatedEuroPerModel",
							Option:        h.euroPerDayChartOption(),
							WidthCssValue: "37vw",
						},
					},
					Class: "",
				},
				app.Div().Class("w-4 h-1"),
				&components.CardComponent{
					Body: []app.UI{
						app.H2().Class("text-xl").Text("Euro per model"),
						&echarts_wasm.EChartComp{
							ContainerID:   "euro-per-model-pie-accumulatedEuroPerModel",
							Option:        h.euroPerModelPieChart(),
							WidthCssValue: "37vw",
						},
					},
					Class: "",
				},
			)
		}),
	)
}
