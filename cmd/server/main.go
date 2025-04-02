package main

import (
	"log"
	"makerworld-analytics/ui/views"
	"net/http"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// The main function is the entry point where the app is configured and started.
// It is executed in 2 different environments: A client (the web browser) and a
// server.
func main() {

	app.Route("/", func() app.Composer {
		return &views.MainView{}
	})

	app.RunWhenOnBrowser()
	http.Handle("/", &app.Handler{
		Name:        "Makerworld Analytics",
		Description: "Extended stats for makerworld",
		Scripts:     []string{"https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4"},
		Styles:      []string{"https://cdn.jsdelivr.net/npm/daisyui@5"},
	})

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
