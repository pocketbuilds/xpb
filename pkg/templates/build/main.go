package main

import (
	"log"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbuilds/xpb"
)

func main() {
	app := pocketbase.New()

	if err := xpb.Setup(app); err != nil {
		log.Fatal(err)
	}

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
