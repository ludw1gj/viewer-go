package main

import (
	"net/http"

	"log"

	"github.com/FriedPigeon/viewer-go/config"
	"github.com/FriedPigeon/viewer-go/db"
	"github.com/FriedPigeon/viewer-go/route"
	"github.com/FriedPigeon/viewer-go/session"
)

func main() {
	// load configuration, database, session, and routes
	c, err := config.Load("config.json")
	if err != nil {
		log.Fatalln(err)
	}
	err = db.Load(c)
	if err != nil {
		log.Fatalln(err)
	}
	session.Load(c)
	route.Load()

	// listen and serve
	http.ListenAndServe(":3000", nil)
}
