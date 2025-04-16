package components

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type OptionData struct {
	Label string
	Value any
}
type SelectValue interface {
	string | int | float32 | float64
}

type SelectComponent[T SelectValue] struct {
	app.Compo
	ID             string
	OptionDataList []OptionData
	OnChange       func(ctx app.Context, value any)
	CurrentValue   T
}

func isEqual[T SelectValue](a, b T) bool {
	return a == b
}

func (s *SelectComponent[T]) isCurrentlySelected(value any) bool {
	if value == nil {
		return s.CurrentValue == *new(T) // Default zero value of T
	}

	castedValue, ok := value.(T)
	if !ok {
		return false
	}

	return isEqual(s.CurrentValue, castedValue)
}

func (s *SelectComponent[T]) Render() app.UI {
	return app.Select().Class("select select-md ").Body(
		app.Range(s.OptionDataList).Slice(func(i int) app.UI {
			value := s.OptionDataList[i].Value

			id := "option-id-"
			if value != nil {
				id += fmt.Sprintf("%s", value)
			} else {
				id += "nil"
			}

			isSelected := s.isCurrentlySelected(value)
			return app.Option().
				Text(s.OptionDataList[i].Label).
				Value(value).Selected(isSelected)
		}),
	).OnChange(s.onChange())
}

func (s *SelectComponent[T]) onChange() func(ctx app.Context, e app.Event) {
	return func(ctx app.Context, e app.Event) {
		if s.OnChange == nil {
			return
		}
		v := ctx.JSSrc().Get("value")

		var output any
		switch v.Type() {
		case app.TypeString:
			output = v.String()
		case app.TypeNumber:
			output = v.Int()
		case app.TypeBoolean:
			output = v.Bool()
		default:
			output = nil
		}

		fmt.Printf("emitting: %v", output)
		s.OnChange(ctx, output)
	}
}
