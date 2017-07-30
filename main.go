package main

import (
	"log"
	"net/http"

	"github.com/FriedPigeon/viewer-go/db"
	"github.com/FriedPigeon/viewer-go/route"
	"github.com/FriedPigeon/viewer-go/session"
)

func main() {
	// load db, session, and routes
	if err := db.Load("viewer.db"); err != nil {
		log.Fatalln(err.Error())
	}
	if err := session.Load("config.json"); err != nil {
		log.Fatalln(err.Error())
	}
	route.Load()

	// listen and serve
	log.Println("viewer-go listening on port 3000...")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatalln(err)
	}
}
