package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/ludw1gj/viewer-go/app/logic/config"
	"github.com/ludw1gj/viewer-go/app/logic/database"
	"github.com/ludw1gj/viewer-go/app/logic/session"
	"github.com/ludw1gj/viewer-go/app/router"
)

func main() {
	port := flag.Int("port", 3000, "Port number")
	dbFile := flag.String("dbFile", "data/viewer/viewer.db", "Database File")
	sessionConfigFile := flag.String("configFile", "data/viewer/session.json", "Session config json file")
	usersDirectory := flag.String("usersDir", "data/users", "Directory where user data will be stored")
	flag.Parse()

	config.SetUsersDirectory(*usersDirectory)

	// load database, session, and routes
	dbConn, err := database.CreateDBConn(*dbFile)
	if err != nil {
		log.Fatalln(err.Error())
	}
	sm, err := session.NewManager(*sessionConfigFile, dbConn)
	if err != nil {
		log.Fatalln(err.Error())
	}
	router.LoadRoutes(dbConn, sm)

	// listen and serve
	log.Printf("viewer-go listening on port %d...", *port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		log.Fatalln(err)
	}
}
