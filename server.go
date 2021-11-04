//go:build !wasm

package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

func httpServer() {
	var static string
	flag.StringVar(&static, "static", "", "generate static files int at the provided path")
	flag.Parse()
	h := buildHandler()
	if static != "" {
		h.Resources = app.GitHubPages("imgtofactbp")
		err := app.GenerateStaticWebsite(static, h)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	if err := http.ListenAndServe(":8000", h); err != nil {
		log.Println(err)
	}
}
