// Package common has functions used commonly by package api and frontend.
package common

import (
	"errors"
	"net/http"

	"bytes"
	"fmt"
	"reflect"

	"github.com/FriedPigeon/viewer-go/db"
	"github.com/FriedPigeon/viewer-go/session"
)

// ValidateUser checks if user's session is valid and then returns the user's information.
func ValidateUser(r *http.Request) (user db.User, err error) {
	userId, err := session.GetUserID(r)
	if err != nil {
		return
	}
	user, err = db.GetUser(userId)
	if err != nil {
		return
	}
	return
}

// genJSONError when passed a struct object with json tags generates an error which includes json structure (keys and
// key's types).
func genJSONError(a interface{}) error {
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

// ValidateJSONInput checks of a passed struct object with json tags has no empty values.
func ValidateJSONInput(a interface{}) error {
	val := reflect.ValueOf(a)
	v := reflect.Indirect(val)

	for i := 0; i < v.Type().NumField(); i++ {
		switch v.Field(i).Type().Kind() {
		case reflect.String:
			if v.Field(i).String() == "" {
				return genJSONError(a)
			}
		case reflect.Int:
			if v.Field(i).Int() == 0 {
				return genJSONError(a)
			}
		case reflect.Float64:
			if v.Field(i).Float() == 0 {
				return genJSONError(a)
			}
		}
	}
	return nil
}
