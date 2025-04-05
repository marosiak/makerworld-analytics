package components

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"makerworld-analytics/domain"
	"time"
)

type SettingsComponent struct {
	app.Compo
	Statistics       *domain.Statistics
	OnSettingsChange func(settings domain.Settings)
	Settings         domain.Settings
}

func (h *SettingsComponent) Render() app.UI {
	return app.Div().Body(
		&CardComponent{
			Class: "flex flex-row justify-between",
			Body: []app.UI{
				&TimeRangeComponent{
					OnChange: func(start, end *time.Time) {
						h.Settings.StartDate = start
						h.Settings.EndDate = end
						h.OnSettingsChange(h.Settings)
					},
				},
				&ButtonSwitchGroup{
					Buttons: []ButtonSwitch{
						{Text: "Bank payout", Checked: h.Settings.MoneyMultiplier == domain.BankPayoutMultiplier},
						{Text: "Vouchers", Checked: h.Settings.MoneyMultiplier == domain.VouchersMultiplier},
					},
					OnChange: func(index int) {
						if index == 0 {
							h.Settings.MoneyMultiplier = domain.BankPayoutMultiplier
						} else {
							h.Settings.MoneyMultiplier = domain.VouchersMultiplier
						}
						h.OnSettingsChange(h.Settings)
					},
				},
			},
		},
	)
}
