// Package common has functions used commonly by package api and frontend.
package common

import (
	"errors"
	"net/http"

	"bytes"
	"fmt"
	"reflect"

	"github.com/FriedPigeon/viewer-go/database"
	"github.com/FriedPigeon/viewer-go/session"
)

// ValidateUser checks if user's session is valid and then returns the user's information.
func ValidateUser(r *http.Request) (u database.User, err error) {
	userId, err := session.GetUserID(r)
	if err != nil {
		return u, err
	}
	u, err = database.GetUser(userId)
	if err != nil {
		return u, err
	}
	return u, nil
}

// ValidateAdmin checks if the user is valid and is admin.
func ValidateAdmin(r *http.Request) (u database.User, err error) {
	u, err = ValidateUser(r)
	if err != nil {
		return u, err
	}
	if !u.Admin {
		return u, errors.New("User is not an admin.")
	}
	return u, nil
}

// ValidateJsonInput checks if a passed struct object with json tags has no empty values.
func ValidateJsonInput(a interface{}) error {
	val := reflect.ValueOf(a)
	v := reflect.Indirect(val)

	for i := 0; i < v.Type().NumField(); i++ {
		switch v.Field(i).Type().Kind() {
		case reflect.String:
			if v.Field(i).String() == "" {
				return genJsonError(a)
			}
		case reflect.Int:
			if v.Field(i).Int() == 0 {
				return genJsonError(a)
			}
		case reflect.Float64:
			if v.Field(i).Float() == 0 {
				return genJsonError(a)
			}
		}
	}
	return nil
}

// genJsonError when passed a struct object with json tags generates an error which includes json structure (keys and
// key's types).
func genJsonError(a interface{}) error {
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
