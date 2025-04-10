package components

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"makerworld-analytics/domain"
	"makerworld-analytics/echarts_wasm"
	"math"
	"sort"
	"time"
)

type ChartsGridComponent struct {
	app.Compo
	Statistics                        *domain.Statistics
	StartDate                         *time.Time
	EndDate                           *time.Time
	MoneyMultiplier                   domain.MoneyMultiplier
	MinimumPointsThresholdForPieChart float32
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

func (h *ChartsGridComponent) euroPerDayChartOption() echarts_wasm.ChartOption {
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

	chartType := "line"
	if len(chartData.Values) > 7 {
		chartType = "bar"
	}
	series := []echarts_wasm.SeriesOption{
		{
			Name: "Euro income",
			Type: chartType,
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

func (h *ChartsGridComponent) euroPerModel() echarts_wasm.ChartOption {
	chartData := echarts_wasm.PieData{}

	for designId, v := range h.Statistics.PointsPerDesign {
		filtered := v.FilterDate(h.StartDate, h.EndDate)
		pointsChange := filtered.SumPointsChange()

		design, designExists := h.Statistics.GetDesignByID(designId)
		if designExists {
			if h.MinimumPointsThresholdForPieChart > pointsChange {
				continue
			}
			chartData = append(chartData, echarts_wasm.PieDataItem{
				Name:  design.Name,
				Value: domain.Statistics{}.ToEuro(h.MoneyMultiplier, pointsChange),
			})
		}

	}
	series := []echarts_wasm.SeriesOption{
		{
			Name:     "Euro income",
			Type:     "pie",
			Data:     chartData,
			PadAngle: 15,
			ItemStyle: echarts_wasm.ItemStyle{
				BorderRadius: 10,
			},
			Radius: []string{"40%", "70%"},
		},
	}
	return echarts_wasm.ChartOption{
		Color:   []string{"#5470c6", "#91cc75", "#fac858", "#ee6666", "#73c0de", "#3ba272", "#fc8452", "#9a60b4", "#ea7ccc"},
		Series:  series,
		Toolbox: echarts_wasm.ToolboxOption{Show: false},
		Tooltip: echarts_wasm.TooltipOption{Show: false},
		Title: echarts_wasm.TitleOption{
			More: map[string]interface{}{},
		},
	}
}

func (h *ChartsGridComponent) Render() app.UI {

	return app.Div().Class("flex flex-row wrap pt-1 justify-stretch").Body(
		&CardComponent{
			Body: []app.UI{
				app.H2().Class("text-xl").Text("Euro income"),
				&echarts_wasm.EChartComp{
					ContainerID: "euro-per-day-chart",
					Option:      h.euroPerDayChartOption(),
				},
			},
			Class: "",
		},
		app.Div().Class("w-4 h-1"),
		&CardComponent{
			Body: []app.UI{
				app.H2().Class("text-xl").Text("Euro per model"),
				&echarts_wasm.EChartComp{
					ContainerID: "euro-per-model-chart",
					Option:      h.euroPerModel(),
				},
			},
			Class: "",
		},
	)
}
