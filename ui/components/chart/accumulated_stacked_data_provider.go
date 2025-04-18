package chart

import (
	"fmt"
	"makerworld-analytics/domain"
	"makerworld-analytics/echarts"
	"slices"
	"time"
)

func (h *ChartsGridComponent) accumulatedEuroPerModelStackedChart(designIDsWhitelist []domain.DesignID) echarts.ChartOption {
	var series []echarts.SeriesOption
	var xAxisDates []string
	var legendNames []string

	endDate := time.Now()
	if h.EndDate != nil {
		endDate = *h.EndDate
	}

	startDate := h.getFirstOccurredPointDate(designIDsWhitelist)

	if h.StartDate != nil {
		startDate = *h.StartDate
	}

	for date := startDate; date.Before(endDate); date = date.Add(24 * time.Hour) {
		dateWithoutTime := date.Format(h.getTimeFormat())
		if !slices.Contains(xAxisDates, dateWithoutTime) {
			xAxisDates = append(xAxisDates, dateWithoutTime)
		}
	}

	for designID, listOfPointsAssignments := range h.Statistics.PointsPerDesign {
		if !slices.Contains(designIDsWhitelist, designID) {
			continue
		}

		filtered := listOfPointsAssignments.FilterDate(&startDate, &endDate)
		filtered = filtered.SortByDate(true)
		design, designExists := h.Statistics.GetDesignByID(designID)

		if !designExists {
			fmt.Printf("Design with ID %d not found\n", designID)
			continue
		}

		chartData := echarts.NumericData{}
		lastProcessedDate := ""
		for _, xAxisValue := range xAxisDates {

			found := false
			for _, pointsEntry := range filtered {
				if h.MinimumPointsThresholdForStackedChart > pointsEntry.PointChange {
					continue
				}

				dateWithoutTime := pointsEntry.CreateTime.Format(h.getTimeFormat())

				if dateWithoutTime == xAxisValue {
					euroIncomeFromEntry := roundFloat(domain.Statistics{}.ToEuro(h.MoneyMultiplier, pointsEntry.PointChange), 2)

					if lastProcessedDate == dateWithoutTime {
						// sum value
						chartData.Values[len(chartData.Values)-1] += euroIncomeFromEntry
					} else {
						// append value
						chartData.Values = append(chartData.Values, euroIncomeFromEntry)
					}
					found = true
				}
				lastProcessedDate = dateWithoutTime
			}
			if !found {
				chartData.Values = append(chartData.Values, 0)
			}
		}
		legendNames = append(legendNames, design.Name)

		// make it cumulative - if first day was +2e second day +3e, thirs day +1
		// it means that first day is equal to 2e, second day is equal to 5e, third day is equal to 6e
		accumulation := 0.0
		for i, v := range chartData.Values {
			if i == 0 {
				accumulation = float64(v)
				continue
			}
			newValue := float64(v) + accumulation
			accumulation += float64(v)

			f := roundFloat(float32(newValue), 2)
			chartData.Values[i] = f
		}

		series = append(series, echarts.SeriesOption{
			Name:      design.Name,
			Type:      "line",
			Data:      chartData,
			Stack:     echarts.StackTypeTotal,
			AreaStyle: map[string]interface{}{},
			Emphasis:  &echarts.Emphasis{Focus: "series"},
		})
	}

	return echarts.ChartOption{
		Color:   []string{"#5470c6", "#91cc75", "#fac858", "#ee6666", "#73c0de", "#3ba272", "#fc8452", "#9a60b4", "#ea7ccc"},
		Series:  series,
		Toolbox: echarts.ToolboxOption{Show: false},
		Tooltip: echarts.TooltipOption{
			Other: map[string]interface{}{
				"trigger": "axis",
			},
		},
		Legend: echarts.LegendOption{
			Data: legendNames,
		},
		XAxis: []echarts.XAxisOption{
			{
				Data: xAxisDates,
				Other: map[string]interface{}{
					"type":       "category",
					"boundryGap": false,
				},
			},
		},
		YAxis: []echarts.YAxisOption{
			{
				Some: map[string]interface{}{
					"type": "value",
				},
			},
		},
		Title: echarts.TitleOption{
			More: map[string]interface{}{},
		},
		DataZoom: []echarts.DataZoom{
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
