package main

import (
	"net/http"

	"log"

	"github.com/FriedPigeon/viewer-go/db"
	"github.com/FriedPigeon/viewer-go/route"
)

func main() {
	err := db.Load()
	if err != nil {
		log.Fatalln(err)
	}

	route.Load()
	http.ListenAndServe(":3000", nil)
}
