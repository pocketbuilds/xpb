package main

import (
	"log"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbuilds/xpb"
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
