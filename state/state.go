package state

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

const (
	MoneyMultiplierKey = "money-multiplier"
)

func HandleGreet(ctx app.Context, a app.Action) {
	name, ok := a.Value.(string)
	if !ok {
		return
	}

	// Setting a state named "greet-name" with the name value.
	ctx.SetState("greet-name", name)
}
