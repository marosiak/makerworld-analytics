package views

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"makerworld-analytics/domain"
	"makerworld-analytics/makerworld"
	"makerworld-analytics/ui/components"
)

type MainView struct {
	app.Compo
	statistics      *domain.Statistics
	updateAvailable bool
}

func (a *MainView) OnAppUpdate(ctx app.Context) {
	a.updateAvailable = ctx.AppUpdateAvailable() // Reports that an app update is available.
}

func (h *MainView) onJsonChange(ctx app.Context, e app.Event) {
	h.statistics = domain.NewStatistics(ctx.JSSrc().Get("value").String())
}

func (h *MainView) importMockedData(ctx app.Context, e app.Event) {
	h.statistics = domain.NewStatistics(makerworld.MockedRawJson)
}

func (h *MainView) Render() app.UI {

	return app.Div().Class("p-24").Body(
		app.If(h.updateAvailable, func() app.UI {
			return app.Button().
				Text("Update app!").
				Class("bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded").
				OnClick(h.onUpdateClick)
		}),
		app.H1().Text("Makerworld Analytics").Class("text-2xl font-bold mb-8"),
		app.Div().Class("flex flex-row mb-4").Body(

			app.P().Text("Paste the  "),
			app.A().Class("text-blue-400 pr-2 pl-2").Text("JSON from Makerworld").Href("https://makerworld.com/api/v1/point-service/point-bill/my?offset=0&limit=9999999999&filter=incomes").Target("_blank"),
			app.P().Text(" to get extended statistics."),
		),

		app.Textarea().Placeholder("Enter JSON from Makerworld").
			Class("shadow appearance-none border rounded w-full h-32 py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline").
			OnKeyDown(h.onJsonChange).
			OnPaste(h.onJsonChange).
			OnChange(h.onJsonChange).OnInput(h.onJsonChange),
		app.Button().Text("Import Maciej Rosiak data").Class("btn btn-soft btn-primary").OnClick(h.importMockedData),

		//components.NewGeneralStatistics(h.statistics),
		app.If(h.statistics != nil, func() app.UI {
			statsComponent := components.NewGeneralStatistics(h.statistics)
			return statsComponent
		}),
	)
}

func (a *MainView) onUpdateClick(ctx app.Context, e app.Event) {
	// Reloads the page to display the modifications.
	ctx.Reload()
}
