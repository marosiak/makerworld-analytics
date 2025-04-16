package chart

import (
	"makerworld-analytics/domain"
	"makerworld-analytics/echarts_wasm"
	"sort"
	"time"
)

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
		if h.SelectedDesign != nil && h.SelectedDesign.ID != 0 {
			// if filter is enabled - skip any other models
			continue
		}

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
		DataZoom: []echarts_wasm.DataZoom{
			{
				Type:  "inside",
				Start: 0,
				End:   100,
			},
			{
				Start: 0,
				End:   100,
			},
		},
	}
}
