package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/robertjeffs/viewer-go/logic/config"
	"github.com/robertjeffs/viewer-go/logic/database"
	"github.com/robertjeffs/viewer-go/logic/router"
	"github.com/robertjeffs/viewer-go/logic/session"
)

func main() {
	port := flag.Int("port", 3000, "Port number")
	dbFile := flag.String("dbFile", "data/viewer/viewer.db", "Database File")
	sessionConfigFile := flag.String("configFile", "data/viewer/session.json", "Session config json file")
	usersDirectory := flag.String("usersDir", "data/users", "Directory where user data will be stored")
	flag.Parse()

	config.SetUsersDirectory(*usersDirectory)

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
