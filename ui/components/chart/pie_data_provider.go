package chart

import (
	"makerworld-analytics/domain"
	"makerworld-analytics/echarts"
)

func (h *ChartsGridComponent) euroPerModelPieChart() echarts.ChartOption {
	chartData := echarts.PieData{}

	for designID, v := range h.Statistics.PointsPerDesign {
		if h.SelectedDesign != nil && h.SelectedDesign.ID != designID {
			// if filter is enabled - skip any other models
			continue
		}

		filtered := v.FilterDate(h.StartDate, h.EndDate)
		pointsChange := filtered.SumPointsChange()

		design, designExists := h.Statistics.GetDesignByID(designID)
		if designExists {
			if h.MinimumPointsThresholdForPieChart > pointsChange {
				continue
			}
			chartData = append(chartData, echarts.PieDataItem{
				Name:  design.Name,
				Value: domain.Statistics{}.ToEuro(h.MoneyMultiplier, pointsChange),
			})
		}

	}
	series := []echarts.SeriesOption{
		{
			Name:     "Euro income",
			Type:     "pie",
			Data:     chartData,
			PadAngle: 15,
			ItemStyle: echarts.ItemStyle{
				BorderRadius: 10,
			},
			Radius: []string{"40%", "70%"},
		},
	}
	return echarts.ChartOption{
		Color:   []string{"#5470c6", "#91cc75", "#fac858", "#ee6666", "#73c0de", "#3ba272", "#fc8452", "#9a60b4", "#ea7ccc"},
		Series:  series,
		Toolbox: echarts.ToolboxOption{Show: false},
		Tooltip: echarts.TooltipOption{Show: false},
		Title: echarts.TitleOption{
			More: map[string]interface{}{},
		},
	}
}
