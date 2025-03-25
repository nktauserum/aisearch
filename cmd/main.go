package main

import (
	"log"

	"github.com/nktauserum/aisearch/internal/app"
)

func main() {
	app := app.NewApplication(8081)
	if err := app.Run(); err != nil {
		log.Fatalf("error running app: %v", err)
	}
}
