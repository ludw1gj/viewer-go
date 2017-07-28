package main

import (
	"net/http"

	"log"

	"github.com/FriedPigeon/viewer-go/db"
	"github.com/FriedPigeon/viewer-go/route"
	"github.com/FriedPigeon/viewer-go/session"
)

func main() {
	// load database
	err := db.Load()
	if err != nil {
		log.Fatalln(err)
	}

	// load session, and routes
	if err = session.Load("config.json"); err != nil {
		log.Fatalln(err)
	}
	route.Load()

	// listen and serve
	log.Println("viewer-go listening on port 3000...")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatalln(err)
	}
}
