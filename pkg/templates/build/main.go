package main

import (
	"log"

	"github.com/PocketBuilds/xpb"
	"github.com/pocketbase/pocketbase"
)

func main() {
	var app = pocketbase.New()

	if err := xpb.LoadConfig(app); err != nil {
		log.Fatal(err)
	}

	if err := xpb.InitPlugins(app); err != nil {
		log.Fatal(err)
	}

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
