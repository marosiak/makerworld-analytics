package components

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"makerworld-analytics/go_echarts_bridge"
	"math/rand"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func generateBarItems() []opts.BarData {
	items := make([]opts.BarData, 0)
	for i := 0; i < 7; i++ {
		items = append(items, opts.BarData{Value: rand.Intn(300)})
	}
	return items
}

type ChartComponent struct {
	app.Compo
}

func (h *ChartComponent) Render() app.UI {
	chart := charts.NewBar()
	chart.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title: "Chart",
	}))
	chart.SetXAxis([]string{"A", "B", "C", "D", "E", "F", "G"}).
		AddSeries("Category A", generateBarItems()).
		SetSeriesOptions(charts.WithLabelOpts(opts.Label{}))

	//bb := chart.RenderContent()
	//println(string(bb))

	return app.Div().ID("chart-component").Body(
		app.H1().Text("Chart"),
		go_echarts_bridge.ComponentFromChart(chart),
	)
}
