package echarts

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type EChartComp struct {
	app.Compo
	ContainerID   string
	Option        ChartOption
	WidthCSSValue string
}

func (c *EChartComp) Render() app.UI {
	return app.Div().ID(c.ContainerID).Style("height", "600px").Style("width", c.WidthCSSValue)
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

	optionsMapped := c.Option.ToMap()
	app.Window().Get("console").Call("log", optionsMapped)
	chart.Call("setOption", optionsMapped)
}
