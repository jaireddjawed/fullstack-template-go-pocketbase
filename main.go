package main

import (
	"log"

	"github.com/jaireddjawed/fullstack-template-golang/internal/app"

	// Blank import registers all Go migrations with PocketBase's app
	// migrations list so they run on `serve` / `migrate up`.
	_ "github.com/jaireddjawed/fullstack-template-golang/migrations"
)

func main() {
	if err := app.New().Start(); err != nil {
		log.Fatal(err)
	}
}
