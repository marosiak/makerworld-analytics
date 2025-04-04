package components

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"makerworld-analytics/domain"
	"makerworld-analytics/echarts_wasm"
	"math"
	"sort"
	"time"
)

type ChartComponent struct {
	app.Compo
	Statistics *domain.Statistics
	startDate  *time.Time
	endDate    *time.Time
}

func roundFloat(val float32, precision uint) float32 {
	ratio := math.Pow(10, float64(precision))
	return float32(math.Round(float64(val)*ratio) / ratio)
}

func (h *ChartComponent) getTimeFormat() string {
	formatLayout := "2006-01-02"
	if h.startDate != nil {
		if h.endDate == nil {
			if time.Since(*h.startDate).Hours() < 72 {
				formatLayout = "2006-01-02 15:04"
			}
		} else {
			if h.endDate.Sub(*h.startDate).Hours() < 72 {
				formatLayout = "2006-01-02 15:04"
			}
		}
	}

	return formatLayout
}

func (h *ChartComponent) euroPerDayChartOption() echarts_wasm.ChartOption {
	chartData := echarts_wasm.NumericData{
		Values: []float32{},
	}

	xAxis := []echarts_wasm.XAxisOption{
		{Data: []string{}},
	}

	formatLayout := h.getTimeFormat()

	var dates []time.Time
	for date := range h.Statistics.PointsPerDay {
		dates = append(dates, date)
	}

	sort.Slice(dates, func(i, j int) bool {
		return dates[i].Before(dates[j])
	})

	for _, date := range dates {
		if h.startDate != nil && date.Before(*h.startDate) {
			continue
		}
		if h.endDate != nil && date.After(*h.endDate) {
			continue
		}
		pointsEarned := h.Statistics.PointsPerDay[date]
		xAxisData := date.Format(formatLayout)

		xAxis[0].Data = append(xAxis[0].Data, xAxisData)
		euroIncome := domain.Statistics{}.ToEuro(domain.VouchersMultiplier, pointsEarned)

		chartData.Values = append(chartData.Values, roundFloat(euroIncome, 2))
	}

	series := []echarts_wasm.SeriesOption{
		{
			Name: "Euro income",
			Type: "line",
			Data: chartData,
		},
	}
	return echarts_wasm.ChartOption{
		Color: []string{"#5470c6", "#91cc75"},
		Legend: echarts_wasm.LegendOption{
			Data:  []string{"Series 1"},
			Other: map[string]interface{}{},
		},
		Series: series,
		Title: echarts_wasm.TitleOption{
			More: map[string]interface{}{},
		},
		Toolbox: echarts_wasm.ToolboxOption{
			Show: true,
		},
		Tooltip: echarts_wasm.TooltipOption{
			Show: true,
		},
		XAxis: xAxis,
		YAxis: []echarts_wasm.YAxisOption{
			{Some: map[string]interface{}{}},
		},
	}
}

func (h *ChartComponent) Render() app.UI {
	return app.Div().Class("pt-12").Body(
		app.Div().Body(
			app.H1().Class("text-2xl opacity-70 mt-8 ml-2 mb-2 select-none").Text("Euro income"),
			&TimeRangeComponent{
				OnChange: func(start, end *time.Time) {
					h.startDate = start
					h.endDate = end
				},
			},
		),
		&echarts_wasm.EChartComp{
			ContainerID: "euro-per-day-chart",
			Option:      h.euroPerDayChartOption(),
		})
}
