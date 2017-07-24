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
	// load configuration file
	c, err := config.Load("config.json")
	if err != nil {
		log.Fatalln(err)
	}

	// load database
	err = db.Load()
	if err != nil {
		log.Fatalln(err)
	}

	// load session, and routes
	session.Load(c)
	route.Load()

	// listen and serve
	log.Println("viewer-go listening on port 3000...")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatalln(err)
	}
}
