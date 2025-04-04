package components

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type ButtonSwitch struct {
	app.Compo
	checked   bool
	Text      string
	OnChecked func(ctx app.Context, checked bool)
}

func (b *ButtonSwitch) Render() app.UI {
	class := "btn btn-default btn-soft"
	if b.checked {
		class = " btn btn-secondary"
	}

	return app.Button().Class(class).Text("Bank payout").OnClick(func(ctx app.Context, e app.Event) {
		if b.OnChecked != nil {
			b.OnChecked(ctx, b.checked)
		}
		ctx.Update()
	})
}

type ButtonSwitchGroup struct {
	app.Compo
	Buttons     []ButtonSwitch
	OnChange    func(index int)
	ActiveIndex int
}

func (b *ButtonSwitchGroup) Render() app.UI {
	println("ActiveIndex = ", b.ActiveIndex)
	return app.Div().Class("flex flex-row").Body(
		app.Range(b.Buttons).Slice(func(i int) app.UI {
			return &ButtonSwitch{
				checked: i == b.ActiveIndex,
				Text:    b.Buttons[i].Text,
				OnChecked: func(ctx app.Context, checked bool) {
					b.ActiveIndex = i
					println("Setting index to = ", i)
					ctx.Update()
					if b.OnChange != nil {
						b.OnChange(i)
					}
				},
			}
		}),
	)
}
