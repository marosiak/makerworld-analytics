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
	Statistics      *domain.Statistics
	StartDate       *time.Time
	EndDate         *time.Time
	MoneyMultiplier domain.MoneyMultiplier
}

func roundFloat(val float32, precision uint) float32 {
	ratio := math.Pow(10, float64(precision))
	return float32(math.Round(float64(val)*ratio) / ratio)
}

func (h *ChartComponent) getTimeFormat() string {
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

type componentDataPayload struct {
	chartOption                  echarts_wasm.ChartOption
	incomeFromScopedPeriodIncome float32
}

func (h *ChartComponent) retrieveComponentData() componentDataPayload {
	income := h.Statistics.PointsPerDate.FilterByDate(h.StartDate, h.EndDate).SumPointsChange()
	return componentDataPayload{
		chartOption:                  h.euroPerDayChartOption(),
		incomeFromScopedPeriodIncome: domain.Statistics{}.ToEuro(h.MoneyMultiplier, income),
	}
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
	for date := range h.Statistics.PointsPerDate {
		dates = append(dates, date)
	}

	sort.Slice(dates, func(i, j int) bool {
		return dates[i].Before(dates[j])
	})

	for _, date := range dates {
		if h.StartDate != nil && date.Before(*h.StartDate) {
			continue
		}
		if h.EndDate != nil && date.After(*h.EndDate) {
			continue
		}
		pointsEarned := h.Statistics.PointsPerDate[date]
		xAxisData := date.Format(formatLayout)

		xAxis[0].Data = append(xAxis[0].Data, xAxisData)
		euroIncome := domain.Statistics{}.ToEuro(h.MoneyMultiplier, pointsEarned)

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
	componentData := h.retrieveComponentData()

	return app.Div().Class("pt-1").Body(
		&echarts_wasm.EChartComp{
			ContainerID: "euro-per-day-chart",
			Option:      componentData.chartOption,
		})
}
