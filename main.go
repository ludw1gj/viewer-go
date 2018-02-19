package main

import (
	"log"
	"net/http"

	"flag"

	"fmt"

	"github.com/robertjeffs/viewer-go/logic/database"
	"github.com/robertjeffs/viewer-go/logic/router"
	"github.com/robertjeffs/viewer-go/logic/session"
)

func main() {
	port := flag.Int("port", 3000, "Port number")
	dbFile := flag.String("dbFile", "data/viewer.db", "Database File")
	sessionConfigFile := flag.String("configFile", "data/session.json", "Session config json file")
	flag.Parse()

	// load database, session, and routes
	err := database.Load(*dbFile)
	if err != nil {
		log.Fatalln(err.Error())
	}
	err = session.Load(*sessionConfigFile)
	if err != nil {
		log.Fatalln(err.Error())
	}
	router.LoadRoutes()

	// listen and serve
	log.Printf("viewer-go listening on port %d...", *port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		log.Fatalln(err)
	}
}
