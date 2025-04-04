package components

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"makerworld-analytics/domain"
	"time"
)

type Settings struct {
	StartDate       *time.Time
	EndDate         *time.Time
	MoneyMultiplier domain.MoneyMultiplier
}

func (s Settings) String() string {
	return fmt.Sprintf("[Settings] StartDate: %s, EndDate: %s, MoneyMultiplier: %s", s.StartDate, s.EndDate, s.MoneyMultiplier)
}

type SettingsComponent struct {
	app.Compo
	Statistics       *domain.Statistics
	OnSettingsChange func(settings Settings)
	Settings         Settings
}

func (h *SettingsComponent) Render() app.UI {

	return app.Div().Body(
		app.Div().Class("flex flex-row justify-between").Body(
			&TimeRangeComponent{
				OnChange: func(start, end *time.Time) {
					h.Settings.StartDate = start
					h.Settings.EndDate = end
					h.OnSettingsChange(h.Settings)
				},
			},
			&ButtonSwitchGroup{
				Buttons: []ButtonSwitch{
					{Text: "Bank payout"},
					{Text: "Vouchers"},
				},
				ActiveIndex: 1,
				OnChange: func(index int) {
					if index == 0 {
						h.Settings.MoneyMultiplier = domain.BankPayoutMultiplier
					} else {
						h.Settings.MoneyMultiplier = domain.VouchersMultiplier
					}
					h.OnSettingsChange(h.Settings)
				},
			},
		),
	)
}
