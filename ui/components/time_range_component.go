package components

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"time"
)

type TimeRangeComponent struct {
	app.Compo
	OnChange func(start, end *time.Time)
}

func (c *TimeRangeComponent) inputChanged(ctx app.Context, e app.Event) {
	value := ctx.JSSrc().Get("value").String()

	var start, end *time.Time

	switch value {
	case "0":
	case "1":
		t := time.Now().AddDate(0, -1, 0)
		start = &t
	case "2":
		t := time.Now().AddDate(0, 0, -7)
		start = &t
	case "3":
		t := time.Now().Add(-1 * time.Hour * 48)
		start = &t
	case "4":
		t := time.Now().Add(-1 * time.Hour * 24)
		start = &t
	}

	c.OnChange(start, end)
}

func (c *TimeRangeComponent) Render() app.UI {
	return app.Div().Class("w-full max-w-xs").Body(
		app.Input().Type("range").OnChange(c.inputChanged).OnInput(c.inputChanged).Min(0).Max(4).Value(0).Class("range").Step(1),
		app.Div().Class("flex justify-between px-2.5 mt-2 text-xs select-none").Body(
			app.Span().Text("|"),
			app.Span().Text("|"),
			app.Span().Text("|"),
			app.Span().Text("|"),
			app.Span().Text("|"),
		),
		app.Div().Class("flex justify-between px-2.5 mt-2 text-xs select-none").Body(
			app.Span().Text("All time"),
			app.Span().Text("Last month"),
			app.Span().Text("Last week"),
			app.Span().Text("Last 48h"),
			app.Span().Text("Last 24h"),
		),
	)
}
