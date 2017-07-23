// This file contains structs used for json responses.

package api

import (
	"encoding/json"
	"net/http"
)

// sendErrorResponse writes an error json response to the client.
func sendErrorResponse(w http.ResponseWriter, errCode int, errMsg string) {
	type errorDesc struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	type errorJSON struct {
		Error errorDesc `json:"error"`
	}
	w.WriteHeader(errCode)
	json.NewEncoder(w).Encode(errorJSON{errorDesc{errCode, errMsg}})
}

// sendSuccessResponse writes a success json response to the client.
func sendSuccessResponse(w http.ResponseWriter, content string) {
	type dataDesc struct {
		Content string `json:"content"`
	}
	type dataJSON struct {
		Data dataDesc `json:"data"`
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dataJSON{dataDesc{Content: content}})
}
