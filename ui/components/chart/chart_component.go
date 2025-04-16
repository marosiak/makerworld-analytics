package chart

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"makerworld-analytics/domain"
	"makerworld-analytics/echarts"
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
			amountOfDataPoints := len(accumulatedEuroPerModel.XAxis[0].Data)

			validPeriodForAverage := domain.PeriodNone
			if amountOfDataPoints >= 7*3 {
				// at least 3 weeks of data
				validPeriodForAverage = domain.PeriodWeek
			}
			if amountOfDataPoints >= 31*3 {
				// at least 3 months of data
				validPeriodForAverage = domain.PeriodMonth
			}

			return app.IfSlice(amountOfDataPoints >= 7, func() []app.UI {
				title := "Average daily euro calculated from week"
				if validPeriodForAverage == domain.PeriodMonth {
					title = "Average daily euro calculated from month"
				}
				return []app.UI{
					buildChartCard(
						"Accumulated Euro per Model",
						"euro-per-model-stacked-accumulatedEuroPerModel",
						accumulatedEuroPerModel,
						"89vw",
					),
					app.Div().Class("w-1 h-4"),
					app.If(validPeriodForAverage != domain.PeriodNone, func() app.UI {
						averageEuroPerModel := h.averageEuroPerModelStackedChart([]domain.DesignID{h.SelectedDesign.ID}, validPeriodForAverage)
						return buildChartCard(
							title,
							"average-euro-per-model-chart",
							averageEuroPerModel,
							"89vw",
						)
					}),
				}
			}).Else(
				func() app.UI {
					return app.Div().Class("flex flex-col wrap w-full pt-1 justify-stretch").Body(
						app.H1().Class("text-xl opacity-75 mt-4 select-none").Text("Not enough data to render chart"),
						app.P().Class("text-md opacity-45 mt-1 select-none text-black").Textf("Try selecting different model, or date range"))
				})
		}).Else(func() app.UI {
			return app.Div().Class("flex flex-row wrap w-full pt-1 justify-stretch").Body(
				buildChartCard(
					"Euro income",
					"euro-per-day-accumulatedEuroPerModel",
					h.euroPerDayChartOption(),
					"38vw",
				),
				app.Div().Class("w-4 h-1"),
				buildChartCard(
					"Euro per model",
					"euro-per-model-pie-accumulatedEuroPerModel",
					h.euroPerModelPieChart(),
					"38vw",
				),
			)
		}),
	)
}

func buildChartCard(title, containerID string, option echarts.ChartOption, widthCSSValue string) *components.CardComponent {
	return &components.CardComponent{
		Body: []app.UI{
			app.H2().Class("text-xl").Text(title),
			&echarts.EChartComp{
				ContainerID:   containerID,
				Option:        option,
				WidthCSSValue: widthCSSValue,
			},
		},
		Class: "",
	}
}
