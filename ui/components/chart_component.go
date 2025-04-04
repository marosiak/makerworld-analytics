package components

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"makerworld-analytics/domain"
	"makerworld-analytics/echarts_wasm"
)

type ChartComponent struct {
	app.Compo
	Statistics *domain.Statistics
}

func testChartOpt() echarts_wasm.ChartOption {
	return echarts_wasm.ChartOption{
		Color: []string{"#5470c6", "#91cc75"},
		Legend: echarts_wasm.LegendOption{
			Data:  []string{"Series 1"},
			Other: map[string]interface{}{},
		},
		Series: []echarts_wasm.SeriesOption{
			{
				Name: "Category A",
				Type: "bar",
				Data: []echarts_wasm.DataItem{
					{Value: 6},
					{Value: 208},
				},
				Label: map[string]interface{}{},
			},
		},
		Title: echarts_wasm.TitleOption{
			Text: "Chart",
			More: map[string]interface{}{},
		},
		Toolbox: echarts_wasm.ToolboxOption{
			Show: true,
		},
		Tooltip: echarts_wasm.TooltipOption{
			Show: true,
		},
		XAxis: []echarts_wasm.XAxisOption{
			{Data: []string{"A", "B"}},
		},
		YAxis: []echarts_wasm.YAxisOption{
			{Some: map[string]interface{}{}},
		},
	}
}

func (h *ChartComponent) statsToChart() echarts_wasm.ChartOption {
	//h.Statistics.PointsPerDay
	output := echarts_wasm.ChartOption{
		Color: []string{"#5470c6", "#91cc75"},
		Legend: echarts_wasm.LegendOption{
			Data:  []string{"Series 1"},
			Other: map[string]interface{}{},
		},
		Series: []echarts_wasm.SeriesOption{
			{
				Name:  "Category A",
				Type:  "bar",
				Data:  []echarts_wasm.DataItem{},
				Label: map[string]interface{}{},
			},
		},
		Title: echarts_wasm.TitleOption{
			Text: "Chart",
			More: map[string]interface{}{},
		},
		Toolbox: echarts_wasm.ToolboxOption{
			Show: true,
		},
		Tooltip: echarts_wasm.TooltipOption{
			Show: true,
		},
		XAxis: []echarts_wasm.XAxisOption{
			{Data: []string{}},
		},
		YAxis: []echarts_wasm.YAxisOption{
			{Some: map[string]interface{}{}},
		},
	}
	for i, v := range h.Statistics.PointsPerDay {
		output.Series[0].Data = append(output.Series[0].Data, echarts_wasm.DataItem{Value: h.Statistics.ToEuro(domain.VouchersMultiplier, v)})
		output.XAxis[0].Data = append(output.XAxis[0].Data, i)
	}
	return output
}

func (h *ChartComponent) Render() app.UI {

	return &echarts_wasm.EChartComp{
		ContainerID: "abc",
		Option:      h.statsToChart(),
	}
}
