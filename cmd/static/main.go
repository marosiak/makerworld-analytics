package main

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"log"
	"makerworld-analytics/ui/views"
)

func main() {
	// TODO fix DRY
	app.Route("/", func() app.Composer {
		return &views.MainView{}
	})

	app.RunWhenOnBrowser()

	err := app.GenerateStaticWebsite(".", &app.Handler{
		Name:        "Makerworld Analytics",
		Description: "Extended stats for makerworld",
		Resources:   app.GitHubPages("makerworld-analytics"),
		Scripts:     []string{"https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4"},
		Styles:      []string{"https://cdn.jsdelivr.net/npm/daisyui@5"},
	})

	if err != nil {
		log.Fatal(err)
	}
}
