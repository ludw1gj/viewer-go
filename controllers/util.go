package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path"
	"path/filepath"
	"reflect"

	"github.com/robertjeffs/viewer-go/logic/config"
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

// validateJSONInput checks if a passed struct object with json tags has no empty values.
func validateJSONInput(a interface{}) error {
	val := reflect.ValueOf(a)
	v := reflect.Indirect(val)

	// generateJSONError when passed a struct object with json tags generates an error which includes json structure -
	// keys and key's types.
	generateJSONError := func(a interface{}) error {
		val := reflect.ValueOf(a)
		v := reflect.Indirect(val)

		var buf bytes.Buffer
		fmt.Fprint(&buf, "invalid json: json must be {")

		for i := 0; i < v.Type().NumField(); i++ {
			fmt.Fprintf(&buf, `"%s": %s, `,
				v.Type().Field(i).Tag.Get("json"),
				v.Type().Field(i).Type)
		}
		buf.Truncate(len(buf.String()) - 2)
		fmt.Fprint(&buf, "}")
		return errors.New(buf.String())
	}

	for i := 0; i < v.Type().NumField(); i++ {
		switch v.Field(i).Type().Kind() {
		case reflect.String:
			if v.Field(i).String() == "" {
				return generateJSONError(a)
			}
		case reflect.Int:
			if v.Field(i).Int() == 0 {
				return generateJSONError(a)
			}
		case reflect.Float64:
			if v.Field(i).Float() == 0 {
				return generateJSONError(a)
			}
		}
	}
	return nil
}

func cleanPath(userRoot, folderPath string) string {
	return filepath.Join(config.GetUsersDirectory(), userRoot, filepath.FromSlash(path.Clean("/"+folderPath)))
}
