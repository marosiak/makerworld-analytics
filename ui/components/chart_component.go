package components

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"makerworld-analytics/domain"
	"makerworld-analytics/echarts_wasm"
	"math"
	"slices"
	"sort"
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
	}
}

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

func (h *ChartsGridComponent) accumulatedEuroPerModelStackedChart(designIDsWhitelist []domain.DesignID) echarts_wasm.ChartOption {
	var series []echarts_wasm.SeriesOption
	var xAxisDates []string
	var legendNames []string

	endDate := time.Now()
	if h.EndDate != nil {
		endDate = *h.EndDate
	}

	startDate := time.Now().Add(-24 * time.Hour * 31)
	if h.StartDate != nil {
		startDate = *h.StartDate
	}

	for date := startDate; date.Before(endDate); date = date.Add(24 * time.Hour) {
		dateWithoutTime := date.Format(h.getTimeFormat())
		if !slices.Contains(xAxisDates, dateWithoutTime) {
			xAxisDates = append(xAxisDates, dateWithoutTime)
		}
	}

	for designId, listOfPointsAssignments := range h.Statistics.PointsPerDesign {
		if !slices.Contains(designIDsWhitelist, designId) {
			continue
		}

		filtered := listOfPointsAssignments.FilterDate(&startDate, &endDate)
		filtered = filtered.SortByDate(true)
		design, designExists := h.Statistics.GetDesignByID(designId)

		if !designExists {
			fmt.Printf("Design with ID %d not found\n", designId)
			continue
		}

		chartData := echarts_wasm.NumericData{}
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

			chartData.Values[i] = roundFloat(float32(newValue), 2)
		}

		series = append(series, echarts_wasm.SeriesOption{
			Name:      design.Name,
			Type:      "line",
			Data:      chartData,
			Stack:     echarts_wasm.StackTypeTotal,
			AreaStyle: map[string]interface{}{},
			Emphasis:  &echarts_wasm.Emphasis{Focus: "series"},
		})
	}

	return echarts_wasm.ChartOption{
		Color:   []string{"#5470c6", "#91cc75", "#fac858", "#ee6666", "#73c0de", "#3ba272", "#fc8452", "#9a60b4", "#ea7ccc"},
		Series:  series,
		Toolbox: echarts_wasm.ToolboxOption{Show: false},
		Tooltip: echarts_wasm.TooltipOption{Show: false},
		Legend: echarts_wasm.LegendOption{
			Data: legendNames,
		},
		XAxis: []echarts_wasm.XAxisOption{
			{
				Data: xAxisDates,
				Other: map[string]interface{}{
					"type":       "category",
					"boundryGap": false,
				},
			},
		},
		YAxis: []echarts_wasm.YAxisOption{
			{
				Some: map[string]interface{}{
					"type": "value",
				},
			},
		},
		Title: echarts_wasm.TitleOption{
			More: map[string]interface{}{},
		},
	}
}

func (h *ChartsGridComponent) Render() app.UI {
	return app.Div().Class("flex flex-col w-full").Body(
		app.If(h.SelectedDesign != nil, func() app.UI {
			return &CardComponent{
				Body: []app.UI{
					app.H2().Class("text-xl").Text("Accumulated euro per model"),
					&echarts_wasm.EChartComp{
						ContainerID:   "euro-per-model-stacked-chart",
						Option:        h.accumulatedEuroPerModelStackedChart([]domain.DesignID{h.SelectedDesign.ID}),
						WidthCssValue: "77vw",
					},
				},
				Class: "",
			}
		}).Else(func() app.UI {
			return app.Div().Class("flex flex-row wrap w-full pt-1 justify-stretch").Body(
				&CardComponent{
					Body: []app.UI{
						app.H2().Class("text-xl").Text("Euro income"),
						&echarts_wasm.EChartComp{
							ContainerID:   "euro-per-day-chart",
							Option:        h.euroPerDayChartOption(),
							WidthCssValue: "37vw",
						},
					},
					Class: "",
				},
				app.Div().Class("w-4 h-1"),
				&CardComponent{
					Body: []app.UI{
						app.H2().Class("text-xl").Text("Euro per model"),
						&echarts_wasm.EChartComp{
							ContainerID:   "euro-per-model-pie-chart",
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
