package main

import (
	"log"
	"pokapoka-viewer/pkg/ui"
)

func main() {
	app := ui.NewApp()
	if err := app.Run(); err != nil {
		log.Fatalf("Failed to start the app: %v", err)
	}
}
