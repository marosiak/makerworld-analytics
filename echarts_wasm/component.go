package echarts_wasm

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type EChartComp struct {
	app.Compo
	ContainerID string
	Option      ChartOption
}

func (c *EChartComp) Render() app.UI {
	return app.Div().ID(c.ContainerID).Style("height", "600px").Style("width", "37vw")
}

func (c *EChartComp) OnMount(ctx app.Context) {
	c.OnUpdate(ctx)
}

func (c *EChartComp) OnUpdate(ctx app.Context) {
	if app.IsServer {
		return
	}
	cont := app.Window().
		Get("document").
		Call("getElementById", c.ContainerID)

	chart := app.Window().
		Get("echarts").
		Call("init", cont, "white",
			map[string]any{"renderer": "canvas"},
		)

	chart.Call("setOption", c.Option.ToMap())
}
