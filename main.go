package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/mlctrez/imgtofactbp/components"
)

func main() {
	var static string
	flag.StringVar(&static, "static", "docs", "generate static files int at the provided path")
	flag.Parse()

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
	app.Route("/", &components.Index{})

	if static != "" {
		err := app.GenerateStaticWebsite(static, h)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	app.RunWhenOnBrowser()

	if err := http.ListenAndServe(":8989", h); err != nil {
		log.Fatal(err)
	}

}
