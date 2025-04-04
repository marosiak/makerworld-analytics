package components

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"math/rand"
	"time"

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
	app.Window().
		Get("console").
		Call("log", "YOUR_ID")

	go func() {
		//<script type="text/javascript">
		//	"use strict";
		//let goecharts_RzqGdVyyqgiS = echarts.init(document.getElementById('RzqGdVyyqgiS'), "white", { renderer: "canvas" });
		//let option_RzqGdVyyqgiS = {"color":["#5470c6","#91cc75","#fac858","#ee6666","#73c0de","#3ba272","#fc8452","#9a60b4","#ea7ccc"],"legend":{},"series":[{"name":"Category A","type":"bar","data":[{"value":6},{"value":208},{"value":260},{"value":271},{"value":176},{"value":103},{"value":23}],"label":{}}],"title":{"text":"Chart"},"toolbox":{},"tooltip":{},"xAxis":[{"data":["A","B","C","D","E","F","G"]}],"yAxis":[{}]}
		//
		//goecharts_RzqGdVyyqgiS.setOption(option_RzqGdVyyqgiS);
		//</script>
		time.Sleep(time.Second * 2)
		app.Window().Get("console").Call("log", "Started go func()")

		doc := app.Window().Get("document")
		container := doc.Call("getElementById", "RzqGdVyyqgiS") // element <div> dla wykresu
		echarts := app.Window().Get("echarts")

		theChart := echarts.Call("init", container, "white",
			map[string]interface{}{"renderer": "canvas"})

		option := map[string]interface{}{
			"color": []interface{}{"#5470c6", "#91cc75", "#fac858", "#ee6666",
				"#73c0de", "#3ba272", "#fc8452", "#9a60b4", "#ea7ccc"},
			"legend": map[string]interface{}{}, // pusty obiekt
			"series": []interface{}{
				map[string]interface{}{
					"name": "Category A",
					"type": "bar",
					"data": []interface{}{ // tablica danych (każdy punkt jako obiekt)
						map[string]interface{}{"value": 6},
						map[string]interface{}{"value": 208},
						map[string]interface{}{"value": 260},
						map[string]interface{}{"value": 271},
						map[string]interface{}{"value": 176},
						map[string]interface{}{"value": 103},
						map[string]interface{}{"value": 23},
					},
					"label": map[string]interface{}{}, // pusty obiekt label
				},
			},
			"title":   map[string]interface{}{"text": "Chart"},
			"toolbox": map[string]interface{}{},
			"tooltip": map[string]interface{}{},
			"xAxis": []interface{}{
				map[string]interface{}{"data": []interface{}{"A", "B", "C", "D", "E", "F", "G"}},
			},
			"yAxis": []interface{}{
				map[string]interface{}{},
			},
		}
		theChart.Call("setOption", option)
	}()

	a := `
<div class="container">
    <div class="item" id="RzqGdVyyqgiS" style="width:900px;height:500px;"></div>
</div>
<style>
    .container {margin-top:30px; display: flex;justify-content: center;align-items: center;}
    .item {margin: auto;}
</style>`

	return app.Raw(a)

	//return app.Div().ID("chart-component").Body(
	//	app.H1().Text("Chart"),
	//	go_echarts_bridge.ComponentFromChart(chart),
	//)
}
