package main

import (
	"net/http"

	"github.com/FriedPigeon/viewer-go/route"
)

func main() {
	http.ListenAndServe(":3000", route.Load())
}
