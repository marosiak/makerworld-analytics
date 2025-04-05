package views

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"makerworld-analytics/domain"
	"makerworld-analytics/makerworld"
	. "makerworld-analytics/ui/components"
)

type MainView struct {
	app.Compo
	updateAvailable bool
	Statistics      *domain.Statistics
	Settings        domain.Settings
}

func (a *MainView) OnAppUpdate(ctx app.Context) {
	a.updateAvailable = ctx.AppUpdateAvailable() // Reports that an app update is available.
}

func (h *MainView) onJsonChange(ctx app.Context, e app.Event) {
	h.Statistics = domain.NewStatistics(ctx.JSSrc().Get("value").String())
}

func (h *MainView) importMockedData(ctx app.Context, e app.Event) {
	h.Statistics = domain.NewStatistics(makerworld.MockedRawJson)
}

func (h *MainView) OnMount(ctx app.Context) {
	ctx.Dispatch(func(ctx app.Context) {
		h.Settings = domain.Settings{
			MoneyMultiplier: domain.VouchersMultiplier,
		}
	})
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

		app.If(h.Statistics != nil, func() app.UI {
			filteredByDate := h.Statistics.PointsPerDate.FilterByDate(h.Settings.StartDate, h.Settings.EndDate)

			incomeForPeriod := domain.Statistics{}.ToEuro(h.Settings.MoneyMultiplier, filteredByDate.SumPointsChange())

			var averageDaily float32
			if len(filteredByDate) >= 15 {
				averageDaily = domain.Statistics{}.ToEuro(h.Settings.MoneyMultiplier, filteredByDate.AveragePointsChange())
			}

			return app.Div().Class("mt-8 flex flex-col").Body(
				&SettingsComponent{
					Statistics: h.Statistics,
					Settings:   h.Settings,
					OnSettingsChange: func(settings domain.Settings) {
						h.Settings = settings
					},
				},
				app.Div().Class("flex flex-row mt-2").Body(
					// TODO: implement DRY rule for these cards, could be provided from one function with 2 params
					&CardComponent{
						Body: []app.UI{
							app.H1().Class("text-xl opacity-40 mt-1 ml-2 select-none").Text("Euro income"),
							app.P().Class("text-3xl opacity-95 mt-1 select-none text-green-600").Textf("+%.1f2€", incomeForPeriod),
						},
						Class: "flex flex-col mb-2 mt-2 flex-none ml-2",
					},
					app.If(averageDaily > 0, func() app.UI {
						return &CardComponent{
							Body: []app.UI{
								app.H1().Class("text-xl opacity-40 mt-1 ml-2 select-none").Text("Average daily"),
								app.P().Class("text-3xl opacity-95 mt-1 select-none text-green-600").Textf("+%.1f2€", averageDaily),
							},
							Class: "flex flex-col mb-4 mt-2 flex-none ml-2",
						}
					}).Else(func() app.UI {
						return &CardComponent{
							Body: []app.UI{
								app.H1().Class("text-xl opacity-40 mt-1 ml-0 select-none").Text("Average daily"),
								app.P().Class("text-3xl opacity-45 mt-1 select-none text-black").Textf("-"),
							},
							Class: "flex flex-col mb-2 mt-2 flex-none ml-2",
						}
					}),
				),
				&ChartComponent{Statistics: h.Statistics, StartDate: h.Settings.StartDate, EndDate: h.Settings.EndDate, MoneyMultiplier: h.Settings.MoneyMultiplier},
			)
		}))
}

func (a *MainView) onUpdateClick(ctx app.Context, e app.Event) {
	ctx.Reload()
}
