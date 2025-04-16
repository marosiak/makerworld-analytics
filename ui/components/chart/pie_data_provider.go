package chart

import (
	"makerworld-analytics/domain"
	"makerworld-analytics/echarts_wasm"
)

func (h *ChartsGridComponent) euroPerModelPieChart() echarts_wasm.ChartOption {
	chartData := echarts_wasm.PieData{}

	for designId, v := range h.Statistics.PointsPerDesign {
		if h.SelectedDesign != nil && h.SelectedDesign.ID != designId {
			// if filter is enabled - skip any other models
			continue
		}

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
