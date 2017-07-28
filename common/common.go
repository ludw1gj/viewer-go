// TODO: write doc
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

// TODO: write doc
func ValidateUser(r *http.Request) (user db.User, err error) {
	userId, err := session.GetUserIDFromSession(r)
	if err != nil {
		return
	}
	user, err = db.GetUser(userId)
	if err != nil {
		return
	}
	return
}

// TODO: write doc
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

// TODO: write doc
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
