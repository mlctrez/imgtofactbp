package main

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/mlctrez/imgtofactbp/components"
)

func buildHandler() *app.Handler {
	h := &app.Handler{
		Title:     "factorio image to blueprint",
		Name:      "factorio image to blueprint",
		ShortName: "img to bp",
		Icon:      app.Icon{Large: "/web/logo-512.png", Default: "/web/logo-192.png"},
		Scripts: []string{
			"https://cdn.jsdelivr.net/npm/bootstrap@5.0.2/dist/js/bootstrap.bundle.min.js",
		},
		Styles: []string{
			"https://fonts.googleapis.com/css2?family=Roboto&display=swap",
			"https://cdn.jsdelivr.net/npm/bootstrap@5.0.2/dist/css/bootstrap.min.css",
			"/web/style.css",
		},
	}
	return h
}

func main() {
	app.Route("/", &components.Index{})
	app.RunWhenOnBrowser()
	httpServer()
}
