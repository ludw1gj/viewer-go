package main

import (
	"net/http"

	"github.com/FriedPigeon/viewer-go/route"
)

func main() {
	route.Load()
	http.ListenAndServe(":3000", nil)
}
