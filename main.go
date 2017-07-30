package main

import (
	"net/http"

	"log"
)

func main() {
	// listen and serve
	log.Println("viewer-go listening on port 3000...")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatalln(err)
	}
}
