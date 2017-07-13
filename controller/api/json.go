// This file contains structs used for json responses.

package api

// contentJSON is used for json response with generic key of "content".
type contentJSON struct {
	Content string `json:"content"`
}

// errorJSON is used for json response with the key of "error".
type errorJSON struct {
	Error string `json:"error"`
}
