package main

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/mlctrez/imgtofactbp/pages"
)

func buildHandler() *app.Handler {
	h := &app.Handler{
		Title:     "image to factorio blueprint",
		Name:      "image to factorio blueprint",
		ShortName: "img to bp",
		Icon:      app.Icon{Large: "/web/logo-512.png", Default: "/web/logo-192.png"},
		Scripts: []string{
			"https://cdnjs.cloudflare.com/ajax/libs/material-components-web/13.0.0/material-components-web.min.js",
		},
		Styles: []string{
			"https://fonts.googleapis.com/icon?family=Material+Icons",
			"https://fonts.googleapis.com/css2?family=Roboto&display=swap",
			"https://cdnjs.cloudflare.com/ajax/libs/material-components-web/13.0.0/material-components-web.min.css",
			"/web/style.css",
		},
	}
	return h
}

func main() {
	app.Route("/", &pages.Index{})
	app.RunWhenOnBrowser()
	httpServer()
}
