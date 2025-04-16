package components

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"makerworld-analytics/domain"
	"sort"
	"time"
)

type SettingsComponent struct {
	app.Compo
	Statistics       *domain.Statistics
	OnSettingsChange func(settings domain.Settings)
	Settings         domain.Settings
}

func (h *SettingsComponent) buildPublicationFilter() app.UI {
	optionDataList := []OptionData{
		{
			Label: "All designs",
			Value: nil,
		},
	}

	mapDesignIDtoIncome := make(map[domain.DesignID]float32)

	for designID, v := range h.Statistics.PointsPerDesign {
		filtered := v.FilterDate(h.Settings.StartDate, h.Settings.EndDate)
		pointsChange := filtered.SumPointsChange()

		mapDesignIDtoIncome[designID] = pointsChange
	}

	designsSorted := h.Statistics.AllPublishedDesigns
	sort.Slice(designsSorted, func(i, j int) bool {
		return mapDesignIDtoIncome[designsSorted[i].ID] > mapDesignIDtoIncome[designsSorted[j].ID]
	})

	for _, design := range designsSorted {
		pointsChange := mapDesignIDtoIncome[design.ID]

		designName := design.Name
		if len(designName) > 52 {
			designName = design.Name[:52] + "..."
		}

		euroEarned := domain.Statistics{}.ToEuro(h.Settings.MoneyMultiplier, pointsChange)
		label := fmt.Sprintf("%.2f€ - %s", euroEarned, designName)
		optionDataList = append(optionDataList, OptionData{
			Label: label,
			Value: fmt.Sprintf("%d", design.ID),
		})
	}

	var currentValue string
	if h.Settings.PublicationFilter != nil {
		currentValue = fmt.Sprintf("%d", h.Settings.PublicationFilter.ID)
	}

	return &SelectComponent[string]{
		OptionDataList: optionDataList,
		CurrentValue:   currentValue,
		OnChange: func(_ app.Context, selected any) {
			if selected != nil {
				selectedDesignID := selected.(string)
				for _, design := range h.Statistics.AllPublishedDesigns {
					if fmt.Sprintf("%d", design.ID) == selectedDesignID {
						h.Settings.PublicationFilter = &design
						h.OnSettingsChange(h.Settings)
						return
					}
				}
			}
			h.Settings.PublicationFilter = nil
			h.OnSettingsChange(h.Settings)
		},
	}
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
				h.buildPublicationFilter(),
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
