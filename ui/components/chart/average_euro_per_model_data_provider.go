package chart

import (
	"fmt"
	"makerworld-analytics/domain"
	"makerworld-analytics/echarts"
	"slices"
	"time"
)

func (h *ChartsGridComponent) getFirstOccurredPointDate(designIDsWhitelist []domain.DesignID) time.Time {
	startDate := time.Now()
	for designID, pointsAssigmentList := range h.Statistics.PointsPerDesign {
		if !slices.Contains(designIDsWhitelist, designID) {
			continue
		}

		filtered := pointsAssigmentList.SortByDate(true)

		if len(filtered) == 0 {
			continue
		}

		startDate = filtered[0].CreateTime
	}
	return startDate
}

func (h *ChartsGridComponent) averageEuroPerModelStackedChart(designIDsWhitelist []domain.DesignID, period domain.Period) echarts.ChartOption {
	var series []echarts.SeriesOption
	var xAxisDates []string
	var legendNames []string

	endDate := time.Now()
	if h.EndDate != nil {
		endDate = *h.EndDate
	}

	origStart := h.getFirstOccurredPointDate(designIDsWhitelist)
	if h.StartDate != nil {
		origStart = *h.StartDate
	}

	var chartStart time.Time
	var nextPeriod func(time.Time) time.Time

	switch period {
	case domain.PeriodWeek:
		chartStart = origStart.AddDate(0, 0, 7)
		nextPeriod = func(t time.Time) time.Time { return t.AddDate(0, 0, 7) }
	case domain.PeriodMonth:
		chartStart = origStart.AddDate(0, 1, 0)
		nextPeriod = func(t time.Time) time.Time { return t.AddDate(0, 1, 0) }
	default:
		chartStart = origStart.AddDate(0, 0, 7)
		nextPeriod = func(t time.Time) time.Time { return t.AddDate(0, 0, 7) }
	}

	for d := chartStart; d.Before(endDate); d = nextPeriod(d) {
		xAxisDates = append(xAxisDates, d.Format(h.getTimeFormat()))
	}

	for designID, pts := range h.Statistics.PointsPerDesign {
		if !slices.Contains(designIDsWhitelist, designID) {
			continue
		}

		filtered := pts.FilterDate(&origStart, &endDate).SortByDate(true)
		design, ok := h.Statistics.GetDesignByID(designID)
		if !ok {
			fmt.Printf("Design with ID %d not found\n", designID)
			continue
		}

		data := echarts.NumericData{}
		for _, ds := range xAxisDates {
			periodEnd, _ := time.Parse(h.getTimeFormat(), ds)
			var periodStart time.Time
			switch period {
			case domain.PeriodWeek:
				periodStart = periodEnd.AddDate(0, 0, -7)
			case domain.PeriodMonth:
				periodStart = periodEnd.AddDate(0, -1, 0)
			default:
				periodStart = periodEnd.AddDate(0, 0, -7)
			}
			daysCount := periodEnd.Sub(periodStart).Hours() / 24
			sum := 0.0
			for _, e := range filtered {
				if e.PointChange < h.MinimumPointsThresholdForStackedChart {
					continue
				}
				if e.CreateTime.After(periodStart) && !e.CreateTime.After(periodEnd) {
					sum += float64(roundFloat(domain.Statistics{}.ToEuro(h.MoneyMultiplier, e.PointChange), 2))
				}
			}
			avg := roundFloat(float32(sum/daysCount), 2)
			data.Values = append(data.Values, avg)
		}

		legendNames = append(legendNames, design.Name)
		series = append(series, echarts.SeriesOption{
			Name:      design.Name,
			Type:      "line",
			Data:      data,
			AreaStyle: map[string]interface{}{},
			Emphasis:  &echarts.Emphasis{Focus: "series"},
		})
	}

	return echarts.ChartOption{
		Color:   []string{"#91cc75", "#fac858", "#ee6666", "#73c0de", "#3ba272", "#fc8452", "#9a60b4", "#ea7ccc", "#5470c6"},
		Series:  series,
		Toolbox: echarts.ToolboxOption{Show: false},
		Tooltip: echarts.TooltipOption{Other: map[string]interface{}{"trigger": "axis"}},
		Legend:  echarts.LegendOption{Data: legendNames},
		XAxis: []echarts.XAxisOption{{
			Data:  xAxisDates,
			Other: map[string]interface{}{"type": "category", "boundaryGap": false},
		}},
		YAxis: []echarts.YAxisOption{{
			Some: map[string]interface{}{"type": "value"},
		}},
		Title: echarts.TitleOption{More: map[string]interface{}{}},
		DataZoom: []echarts.DataZoom{
			{Type: "inside", Start: 0, End: 100},
			{Start: 0, End: 100},
		},
	}
}
