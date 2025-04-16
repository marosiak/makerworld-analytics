package views

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"makerworld-analytics/domain"
	"makerworld-analytics/makerworld"
	. "makerworld-analytics/ui/components"
	. "makerworld-analytics/ui/components/chart"
)

type MainView struct {
	app.Compo
	updateAvailable bool
	Statistics      *domain.Statistics
	Settings        domain.Settings
}

func (m *MainView) OnAppUpdate(ctx app.Context) {
	m.updateAvailable = ctx.AppUpdateAvailable() // Reports that an app update is available.
}

func (m *MainView) onJSONChange(ctx app.Context, e app.Event) {
	m.Statistics = domain.NewStatistics(ctx.JSSrc().Get("value").String())
}

func (m *MainView) importMockedData(ctx app.Context, e app.Event) {
	m.Statistics = domain.NewStatistics(makerworld.MockedRawJSON)
}

func (m *MainView) OnMount(ctx app.Context) {
	ctx.Dispatch(func(ctx app.Context) {
		m.Settings = domain.Settings{
			MoneyMultiplier: domain.VouchersMultiplier,
		}
	})
}

func (m *MainView) Render() app.UI {
	return app.Div().Class("p-24").Body(
		app.If(m.updateAvailable, func() app.UI {
			return app.Button().
				Text("Update app!").
				Class("bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded").
				OnClick(m.onUpdateClick)
		}),
		app.H1().Text("Makerworld Analytics").Class("text-2xl font-bold mb-8"),
		app.Div().Class("flex flex-row mb-4").Body(

			app.P().Text("Visit "),
			app.A().Class("text-blue-400 pr-1 pl-1").Text("this link").Href("https://makerworld.com/api/v1/point-service/point-bill/my?offset=0&limit=9999999999&filter=incomes").Target("_blank"),
			app.P().Text(" and copy all content"),
		),

		app.Textarea().Placeholder("Paste JSON from link above").
			Class("shadow appearance-none border rounded w-full m-0 py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline").
			OnKeyDown(m.onJSONChange).
			OnPaste(m.onJSONChange).
			OnChange(m.onJSONChange).OnInput(m.onJSONChange),
		app.Button().Text("Import test data").Class("btn btn-soft btn-primary mt-4").OnClick(m.importMockedData),
		app.P().Text("Basically my data, it may be old because I am not updating it often.").Class("text-xs opacity-40 mt-1"),

		app.If(m.Statistics != nil, func() app.UI {
			allModels := m.Statistics.PointsPerDate.FilterByDate(m.Settings.StartDate, m.Settings.EndDate)

			var incomeForPeriod float32
			var averageDaily float32
			if m.Settings.PublicationFilter == nil {
				incomeForPeriod = allModels.SumPointsChange()
				if len(allModels) >= 15 {
					averageDaily = allModels.AveragePointsPerDay()
				}
			} else {
				pointsPerDesignFilteredByDate := m.Statistics.PointsPerDesign[m.Settings.PublicationFilter.ID].FilterDate(m.Settings.StartDate, m.Settings.EndDate)
				incomeForPeriod = pointsPerDesignFilteredByDate.SumPointsChange()
				averageDaily = pointsPerDesignFilteredByDate.AveragePointsPerDay()
			}

			incomeForPeriod = m.Statistics.ToEuro(m.Settings.MoneyMultiplier, incomeForPeriod)
			averageDaily = m.Statistics.ToEuro(m.Settings.MoneyMultiplier, averageDaily)

			return app.Div().Class("mt-8 flex flex-col").Body(
				&SettingsComponent{
					Statistics: m.Statistics,
					Settings:   m.Settings,
					OnSettingsChange: func(settings domain.Settings) {
						m.Settings = settings
					},
				},
				app.Div().Class("flex flex-row mt-2").Body(
					m.renderStatisticsCardWidget("Euro income", fmt.Sprintf("+%.1f2€", incomeForPeriod)),
					app.If(averageDaily > 0, func() app.UI {
						return m.renderStatisticsCardWidget("Average daily", fmt.Sprintf("+%.1f2€", averageDaily))
					}).Else(func() app.UI {
						return m.renderStatisticsCardWidget("Average daily", "-")
					}),
				),
				&ChartsGridComponent{Statistics: m.Statistics, StartDate: m.Settings.StartDate, EndDate: m.Settings.EndDate, MoneyMultiplier: m.Settings.MoneyMultiplier, MinimumPointsThresholdForPieChart: 15, MinimumPointsThresholdForStackedChart: 0.1, SelectedDesign: m.Settings.PublicationFilter},
			)
		}))
}

func (m *MainView) renderStatisticsCardWidget(label, value string) *CardComponent {
	return &CardComponent{
		Body: []app.UI{
			app.H1().Class("text-xl opacity-40 mt-1 ml-2 select-none").Text(label),
			app.P().Class("text-3xl opacity-95 mt-1 select-none text-green-600").Textf(value),
		},
		Class: "flex flex-col mb-2 mt-2 flex-none ml-2",
	}
}

func (m *MainView) onUpdateClick(ctx app.Context, e app.Event) {
	ctx.Reload()
}
