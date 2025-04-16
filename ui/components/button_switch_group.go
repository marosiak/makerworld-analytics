package components

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type ButtonSwitch struct {
	app.Compo
	Checked   bool
	Text      string
	OnChecked func(ctx app.Context, checked bool)
}

func (b *ButtonSwitch) Render() app.UI {
	class := "btn btn-default btn-soft"
	if b.Checked {
		class = " btn btn-secondary"
	}

	return app.Button().Class(class).Text(b.Text).OnClick(func(ctx app.Context, _ app.Event) {
		if b.OnChecked != nil {
			b.OnChecked(ctx, b.Checked)
		}
		ctx.Update()
	})
}

type ButtonSwitchGroup struct {
	app.Compo
	Buttons  []ButtonSwitch
	OnChange func(index int)
}

func (b *ButtonSwitchGroup) Render() app.UI {
	return app.Div().Class("flex flex-row").Body(
		app.Range(b.Buttons).Slice(func(i int) app.UI {
			containerClass := "mr-2"
			if i == len(b.Buttons)-1 {
				containerClass = ""
			}
			return app.Div().Class(containerClass).Body(
				&ButtonSwitch{
					Checked: b.Buttons[i].Checked,
					Text:    b.Buttons[i].Text,
					OnChecked: func(ctx app.Context, _ bool) {
						ctx.Update()
						if b.OnChange != nil {
							b.OnChange(i)
						}
					},
				})
		}),
	)
}
