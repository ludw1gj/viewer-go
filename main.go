package main

import (
	"net/http"

	"github.com/FriedPigeon/viewer-go/db"
	"github.com/FriedPigeon/viewer-go/route"
)

func main() {
	db.Load()
	route.Load()
	http.ListenAndServe(":3000", nil)
}
