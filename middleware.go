package goreq

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

func RequestParse[T any](next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}

type structTags struct {
	In    string
	Label string
	Type  string
}

// parseTags finds all the tag components and parses
// them into an easier to consume format
func parseTags(val string) structTags {
	values := strings.Split(val, ",")

	tags := structTags{}

	for _, value := range values {
		parts := strings.Split(value, "=")
		if len(parts) != 2 {
			continue
		}
		if len(parts[1]) == 0 {
			continue
		}

		switch parts[0] {
		case "in":
			tags.In = parts[1]
		case "label":
			tags.Label = parts[1]
		case "type":
			tags.Type = parts[1]
		}
	}

	return tags
}

func parseParameters[T any](r *http.Request) (T, error) {
	var ret T

	// Used for finding struct tag definitions
	ty := reflect.TypeOf(ret)

	// Make sure we only try to parse into a struct-type. We are
	// unable to enforce this via the generic type so need to do it here.
	if ty.Kind() != reflect.Struct {
		return ret, fmt.Errorf("expected a struct, got a %s", ty)
	}

	// Used for assigning return values
	va := reflect.ValueOf(&ret)

	for i := 0; i < ty.NumField(); i++ {
		field := ty.Field(i)

		tag := parseTags(field.Tag.Get("goreq"))

		switch tag.In {
		case "header":
			headerVal := r.Header.Get(tag.Label)
			va.Elem().FieldByName(field.Name).SetString(headerVal)
		case "query":
			queryVal := r.URL.Query().Get(tag.Label)
			va.Elem().FieldByName(field.Name).SetString(queryVal)
		}
	}

	return ret, nil
}
